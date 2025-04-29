package auth

import (
	"BackendGoLdap/config"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	gocloak "github.com/Nerzal/gocloak/v13"
	"github.com/coreos/go-oidc/v3/oidc"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

var ctx = context.Background()

var (
	provider    *oidc.Provider
	verifier    *oidc.IDTokenVerifier
	oauthConfig *oauth2.Config
)

// Struct for select only necessary user fields
type UserClaims struct {
    Email    string `json:"email"`
    Username string `json:"username"`
}

// InitOIDC initializes the OpenID provider and verifier
// Make configuration as in Keycloak
func InitOIDC() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	customHttpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	providerCtx := oidc.ClientContext(ctx, customHttpClient)

	provider, err = oidc.NewProvider(providerCtx, fmt.Sprintf("%s/realms/%s", cfg.KeycloakBaseURL, cfg.KeycloakRealm))
	if err != nil {
		return fmt.Errorf("failed to initialize OIDC provider: %w", err)
	}

	verifier = provider.Verifier(&oidc.Config{
		ClientID: cfg.KeycloakClientID,
	})

	oauthConfig = &oauth2.Config{
		ClientID:     cfg.KeycloakClientID,
		ClientSecret: cfg.KeycloakClientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return nil
}


// LoginHandlerByUID creates a login handler that takes a username and password
// as JSON payload and returns the Keycloak access token.
//
// The handler expects the following JSON payload:
//
// {
//     "username": string,
//     "password": string,
// }
//
// The handler returns a JSON response with the Keycloak access token:
//
// {
//     "access_token": string,
//     "expires_in":   int,
//     "refresh_expires_in": int,
//     "refresh_token": string,
//     "token_type":   string,
//     "id_token":     string,
//     "session_state": string,
//     "scope":        string,
// }
func LoginHandlerByUID(ck *gocloak.GoCloak) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := config.GetLogger()
		cfg, err := config.GetConfig()
		if err != nil {
			logger.Error("failed to get config", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		var req struct {
			UID      string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("failed to decode request body", zap.Error(err))
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		token, err := ck.Login(
			ctx,
			cfg.KeycloakClientID,
			cfg.KeycloakClientSecret,
			cfg.KeycloakRealm,
			req.UID,
			req.Password,
		)
		if err != nil {
			logger.Error("keycloak login failed", zap.Error(err))
			http.Error(w, fmt.Sprintf("Keycloak login failed: %v", err), http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(token)
		if err != nil {
			logger.Info("Error while encode token", zap.Error(err))
		}
	}
}

// AuthMiddleware
func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := config.GetLogger()

			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				logger.Error("missing bearer token")
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}

			rawIDToken := strings.TrimPrefix(authHeader, "Bearer ")

			idToken, err := verifier.Verify(ctx, rawIDToken)
			if err != nil {
				logger.Error("failed to verify ID Token", zap.Error(err))
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			// Store the idToken in request context for later handlers
			accessTokenKey := "accessToken"
			newCtx := context.WithValue(r.Context(), accessTokenKey, idToken)
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

// GetUserDataFromToken fetches user claims using ID Token
func GetUserDataFromToken() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        logger := config.GetLogger()

        accessTokenValue := r.Context().Value("accessToken")
        if accessTokenValue == nil {
            logger.Error("missing idToken in context")
            http.Error(w, "missing token", http.StatusUnauthorized)
            return
        }

        accessToken := accessTokenValue.(*oidc.IDToken)

		var allData map[string]interface{}

        if err := accessToken.Claims(&allData); err != nil {
            logger.Error("failed to extract claims", zap.Error(err))
            http.Error(w, "invalid token claims", http.StatusInternalServerError)
            return
        }

		username, ok := allData["preferred_username"].(string)
		if !ok {
			logger.Error("failed to extract preferred_username")
			http.Error(w, "invalid token claims", http.StatusInternalServerError)
			return
		}

		email, ok := allData["email"].(string)
		if !ok {
			logger.Error("failed to extract email")
			http.Error(w, "invalid token claims", http.StatusInternalServerError)
			return
		}

		resourceAccess, ok := allData["resource_access"].(map[string]interface{})
		if !ok {
			logger.Error("failed to extract resource_access")
			http.Error(w, "invalid token claims", http.StatusInternalServerError)
			return
		}

		account := resourceAccess["account"].(map[string]interface{})
		roles := account["roles"].([]interface{})
		
		var roleStrings []string
		for _, role := range roles {
			if roleStr, ok := role.(string); ok {
				roleStrings = append(roleStrings, roleStr)
			}
		}

		logger.Info("User claims", zap.String("username", username),
		 			zap.String("email", email),
					zap.Strings("roles", roleStrings))

        w.Header().Set("Content-Type", "application/json")
        err := json.NewEncoder(w).Encode(allData)
		if err != nil {
			logger.Error("Failed to pass user claims", zap.Error(err))
		}
    }
}

// NewTokenRefreshMiddleware creates a new token refresh middleware with default Keycloak client
func NewTokenRefreshMiddleware() func(http.Handler) http.Handler {
	cfg, err := config.GetConfig()
	if err != nil {
		config.GetLogger().Error("failed to get config", zap.Error(err))
		return nil
	}
	ck := gocloak.NewClient(cfg.KeycloakBaseURL)
	return TokenRefreshMiddleware(ck)
}

// NewRefreshTokenHandler creates a new refresh token handler with default Keycloak client
func NewRefreshTokenHandler() http.HandlerFunc {
	cfg, err := config.GetConfig()
	if err != nil {
		config.GetLogger().Error("failed to get config", zap.Error(err))
		return nil
	}
	ck := gocloak.NewClient(cfg.KeycloakBaseURL)
	return RefreshTokenHandler(ck)
}