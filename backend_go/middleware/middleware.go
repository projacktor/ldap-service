package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
)

// AuthMiddleware provides JWT authentication using Keycloak
// It validates Bearer tokens against a Keycloak server
type AuthMiddleware struct {
	client *gocloak.GoCloak // Keycloak client instance
	realm  string           // Keycloak realm to validate against
}

// NewAuthMiddleware creates a new AuthMiddleware instance
// Parameters:
//   - client: Initialized Keycloak client
//   - realm: Keycloak realm name for token validation
//
// Returns configured AuthMiddleware ready for use
func NewAuthMiddleware(client *gocloak.GoCloak, realm string) *AuthMiddleware {
	return &AuthMiddleware{
		client: client,
		realm:  realm,
	}
}

// Middleware returns an HTTP handler that validates JWT tokens
// Implements the standard middleware pattern for Go HTTP handlers
func (a *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		// Get raw token by removing "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Introspect token with Keycloak
		// Note: Empty clientID and clientSecret parameters mean public client
		rptResult, err := a.client.RetrospectToken(ctx, tokenString, "", "", a.realm)
		if err != nil || !*rptResult.Active {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid - proceed to next handler
		next.ServeHTTP(w, r)
	})
}
