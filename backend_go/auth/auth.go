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

var ctx = context.Background()

// Note we take *gocloak.GoCloak, not gocloak.GoCloak by value.
func LoginHandlerByEmail(ck *gocloak.GoCloak, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UID      string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		userUID := req.UID
		// Resource Owner Password grant with Keycloak
		token, err := (*ck).Login(
			ctx,
			cfg.KeycloakClientID,
			cfg.KeycloakClientSecret,
			cfg.KeycloakRealm,
			userUID,      // Keycloak still gets the uid as username
			req.Password, // same password
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Keycloak login failed: %v", err), http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(token)
	}
}

// AuthMiddleware protects routes by introspecting the Bearer token.
// Fix: active.Active is *bool, so we must dereference it.
func AuthMiddleware(ck *gocloak.GoCloak, cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}
			tok := strings.TrimPrefix(auth, "Bearer ")

			active, err := (*ck).RetrospectToken(
				ctx,
				tok,
				cfg.KeycloakClientID,
				cfg.KeycloakClientSecret,
				cfg.KeycloakRealm,
			)
			if err != nil || active == nil || active.Active == nil || !*active.Active {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func NewAuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	ck := gocloak.NewClient(cfg.KeycloakBaseURL)
	return AuthMiddleware(ck, cfg)
}
