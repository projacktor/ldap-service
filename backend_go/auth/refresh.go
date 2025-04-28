package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"BackendGoLdap/config"

	gocloak "github.com/Nerzal/gocloak/v13"
	"go.uber.org/zap"
)

// RefreshTokenHandler handles token refresh requests
func RefreshTokenHandler(ck *gocloak.GoCloak) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := config.GetLogger()
		cfg, err := config.GetConfig()
		if err != nil {
			logger.Error("failed to get config", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		var req struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("failed to decode request body", zap.Error(err))
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		if req.RefreshToken == "" {
			logger.Error("missing refresh token")
			http.Error(w, "refresh token is required", http.StatusBadRequest)
			return
		}

		token, err := ck.RefreshToken(
			ctx,
			req.RefreshToken,
			cfg.KeycloakClientID,
			cfg.KeycloakClientSecret,
			cfg.KeycloakRealm,
		)
		if err != nil {
			logger.Error("failed to refresh token", zap.Error(err))
			http.Error(w, "failed to refresh token", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(token)
	}
}

// TokenRefreshMiddleware creates middleware for refreshing tokens
func TokenRefreshMiddleware(ck *gocloak.GoCloak) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := config.GetLogger()
			cfg, err := config.GetConfig()
			if err != nil {
				logger.Error("failed to get config", zap.Error(err))
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				logger.Error("missing bearer token")
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			tok := strings.TrimPrefix(auth, "Bearer ")
			active, err := ck.RetrospectToken(
				ctx,
				tok,
				cfg.KeycloakClientID,
				cfg.KeycloakClientSecret,
				cfg.KeycloakRealm,
			)

			if err != nil || active == nil || active.Active == nil || !*active.Active {
				// Try to refresh token if refresh token is present
				refreshToken := r.Header.Get("X-Refresh-Token")
				if refreshToken != "" {
					newToken, err := ck.RefreshToken(
						ctx,
						refreshToken,
						cfg.KeycloakClientID,
						cfg.KeycloakClientSecret,
						cfg.KeycloakRealm,
					)
					if err == nil {
						// Set new tokens in response headers
						w.Header().Set("X-New-Access-Token", newToken.AccessToken)
						w.Header().Set("X-New-Refresh-Token", newToken.RefreshToken)
						// Update request with new access token
						r.Header.Set("Authorization", "Bearer "+newToken.AccessToken)
					} else {
						logger.Error("invalid token and refresh failed", zap.Error(err))
						http.Error(w, "invalid token", http.StatusUnauthorized)
						return
					}
				} else {
					logger.Error("invalid token and no refresh token", zap.Error(err))
					http.Error(w, "invalid token", http.StatusUnauthorized)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
} 