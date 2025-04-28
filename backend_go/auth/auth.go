package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"BackendGoLdap/config"

	gocloak "github.com/Nerzal/gocloak/v13"
)

// Global context used for Keycloak operations
var ctx = context.Background()

// LoginHandlerByUID handles user authentication via Keycloak using username/password
// Takes gocloak client pointer and app config as dependencies
// Implements Resource Owner Password Credentials flow (ROPC)
func LoginHandlerByUID(ck *gocloak.GoCloak, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Define request payload structure
		var req struct {
			UID      string `json:"username"` // User's attribute username
			Password string `json:"password"` // User's password
		}
		// Decode JSON request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
func AuthMiddleware(ck *gocloak.GoCloak, cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract Authorization header
			auth := r.Header.Get("Authorization")
			// Verify Bearer token format
			if !strings.HasPrefix(auth, "Bearer ") {
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
func NewAuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	ck := gocloak.NewClient(cfg.KeycloakBaseURL)
	return AuthMiddleware(ck, cfg)
}
