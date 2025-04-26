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
func LoginHandlerByEmail(ck *gocloak.GoCloak, cfg *config.Config) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Email    string `json:"email"`
            Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid payload", http.StatusBadRequest)
            return
        }

        // 1) Connect to LDAP
        l, err := ldap.DialURL(cfg.LDAPHost)
        if err != nil {
            http.Error(w, fmt.Sprintf("LDAP unreachable: %v", err), http.StatusServiceUnavailable)
            return
        }
        defer l.Close()

        // 2) Service bind
        if err := l.Bind(cfg.LDAPBindDN, cfg.LDAPBindPass); err != nil {
            http.Error(w, fmt.Sprintf("service bind failed: %v", err), http.StatusUnauthorized)
            return
        }

        // 3) Search by email
        filter := fmt.Sprintf("(mail=%s)", ldap.EscapeFilter(req.Email))
        search := ldap.NewSearchRequest(
            cfg.LDAPBaseDN,
            ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 1, 0, false,
            filter,
            []string{"dn", "uid"}, // request both DN and uid
            nil,
        )
        res, err := l.Search(search)
        if err != nil || len(res.Entries) == 0 {
            http.Error(w, "user not found", http.StatusUnauthorized)
            return
        }
        entry := res.Entries[0]
        userDN := entry.DN
        userUID := entry.GetAttributeValue("uid")
        if userUID == "" {
            http.Error(w, "LDAP entry missing uid", http.StatusInternalServerError)
            return
        }

        // 4) Bind as the user to verify password
        if err := l.Bind(userDN, req.Password); err != nil {
            http.Error(w, "invalid credentials", http.StatusUnauthorized)
            return
        }

        // 5) Resource Owner Password grant with Keycloak
        token, err := (*ck).Login(
            ctx,
            cfg.KeycloakClientID,
            cfg.KeycloakClientSecret,
            cfg.KeycloakRealm,
            userUID,       // Keycloak still gets the uid as username
            req.Password,  // same password
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