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

// InitOIDC initializes the OIDC provider and verifier
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


// LoginHandlerByUID - remains using gocloak ROPC flow
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
		json.NewEncoder(w).Encode(token)
	}
}

// AuthMiddleware using go-oidc verifier
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
			newCtx := context.WithValue(r.Context(), "idToken", idToken)
			next.ServeHTTP(w, r.WithContext(newCtx))
		})
	}
}

// GetUserDataFromToken fetches user claims from ID Token directly
func GetUserDataFromToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := config.GetLogger()

		idTokenValue := r.Context().Value("idToken")
		if idTokenValue == nil {
			logger.Error("missing idToken in context")
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		idToken := idTokenValue.(*oidc.IDToken)

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			logger.Error("failed to extract claims", zap.Error(err))
			http.Error(w, "invalid token claims", http.StatusInternalServerError)
			return
		}

		// Return claims as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(claims)
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