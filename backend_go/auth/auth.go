package auth

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    
    "github.com/go-ldap/ldap/v3"
    gocloak "github.com/Nerzal/gocloak/v13"
    "BackendGoLdap/config"
)

var ctx = context.Background()

// LoginHandler returns a handler that binds against LDAP then mints a Keycloak token.
// Note we take *gocloak.GoCloak, not gocloak.GoCloak by value.
func LoginHandler(ck *gocloak.GoCloak, cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid payload", http.StatusBadRequest)
            return
        }

        // 1) LDAP bind as service account
        l, err := ldap.DialURL(cfg.LDAPHost)
        if err != nil {
            http.Error(w, "LDAP unreachable", http.StatusServiceUnavailable)
            return
        }
        defer l.Close()
        if err := l.Bind(cfg.LDAPBindDN, cfg.LDAPBindPass); err != nil {
            http.Error(w, "service bind failed", http.StatusUnauthorized)
            return
        }

        // 2) Look up user DN
        search := ldap.NewSearchRequest(
            cfg.LDAPBaseDN,
            ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 1, 0, false,
            fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(req.Username)),
            []string{"dn"}, nil,
        )
        res, err := l.Search(search)
        if err != nil || len(res.Entries) == 0 {
            http.Error(w, "user not found", http.StatusUnauthorized)
            return
        }
        userDN := res.Entries[0].DN

        // 3) Bind as user to verify password
        if err := l.Bind(userDN, req.Password); err != nil {
            http.Error(w, "invalid credentials", http.StatusUnauthorized)
            return
        }

        // 4) Resource-Owner-Password grant with Keycloak
        token, err := (*ck).Login(
            ctx,
            cfg.KeycloakClientID,
            cfg.KeycloakClientSecret,
            cfg.KeycloakRealm,
            req.Username,
            req.Password,
        )
        if err != nil {
            http.Error(w, "Keycloak login failed", http.StatusBadGateway)
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