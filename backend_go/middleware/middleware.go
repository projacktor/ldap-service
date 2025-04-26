package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Nerzal/gocloak/v13"
)

type AuthMiddleware struct {
	client *gocloak.GoCloak
	realm  string
}

func NewAuthMiddleware(client *gocloak.GoCloak, realm string) *AuthMiddleware {
	return &AuthMiddleware{
		client: client,
		realm:  realm,
	}
}

func (a *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		rptResult, err := a.client.RetrospectToken(ctx, tokenString, "", "", a.realm)
		if err != nil || !*rptResult.Active {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid, continue
		next.ServeHTTP(w, r)
	})
}
