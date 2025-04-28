package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"BackendGoLdap/config"

	gocloak "github.com/Nerzal/gocloak/v13"
	"go.uber.org/zap"
)

// Global context used for Keycloak operations
var ctx = context.Background()

// LoginHandlerByUID handles user authentication via Keycloak using username/password
// Takes gocloak client pointer as dependency
// Implements Resource Owner Password Credentials flow (ROPC)
func LoginHandlerByUID(ck *gocloak.GoCloak) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := config.GetLogger()
		cfg, err := config.GetConfig()
		if err != nil {
			logger.Error("failed to get config", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// Define request payload structure
		var req struct {
			UID string `json:"username"` // User's attribute username
			Password string `json:"password"` // User's password
		}
		// Decode JSON request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("failed to decode request body", zap.Error(err))
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		// Authenticate with Keycloak using ROPC flow
		token, err := ck.Login(
			ctx,
			cfg.KeycloakClientID,     // Client identifier
			cfg.KeycloakClientSecret, // Client secret
			cfg.KeycloakRealm,        // Keycloak realm
			req.UID,                  // Keycloak gets the uid attribute as username
			req.Password,             // same password
		)
		if err != nil {
			logger.Error("keycloak login failed", zap.Error(err))
			// Return 502 if Keycloak communication fails
			http.Error(w, fmt.Sprintf("Keycloak login failed: %v", err), http.StatusBadGateway)
			return
		}
		// Return JWT tokens to client
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(token)
	}
}

// AuthMiddleware creates middleware for protecting routes with JWT validation
// Verifies token active status using Keycloak's token introspection endpoint
func AuthMiddleware(ck *gocloak.GoCloak) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := config.GetLogger()
			cfg, err := config.GetConfig()
			if err != nil {
				logger.Error("failed to get config", zap.Error(err))
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			// Extract Authorization header
			auth := r.Header.Get("Authorization")
			// Verify Bearer token format
			if !strings.HasPrefix(auth, "Bearer ") {
				logger.Error("missing bearer token")
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}
			// Extract raw token
			tok := strings.TrimPrefix(auth, "Bearer ")
			// Introspect token with Keycloak
			active, err := (*ck).RetrospectToken(
				ctx,
				tok,
				cfg.KeycloakClientID,
				cfg.KeycloakClientSecret,
				cfg.KeycloakRealm,
			)
			// Validate token status
			if err != nil || active == nil || active.Active == nil || !*active.Active {
				logger.Error("invalid token", zap.Error(err))
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			// Token is valid - proceed to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// NewAuthMiddleware creates a new auth middleware with default Keycloak client
// Convenience function that initializes the Keycloak client internally
func NewAuthMiddleware() func(http.Handler) http.Handler {
	cfg, err := config.GetConfig()
	if err != nil {
		config.GetLogger().Error("failed to get config", zap.Error(err))
		return nil
	}
	ck := gocloak.NewClient(cfg.KeycloakBaseURL)
	return AuthMiddleware(ck)
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
