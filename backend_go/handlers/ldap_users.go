package handlers

import (
	"BackendGoLdap/config"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-ldap/ldap/v3"
	"go.uber.org/zap"
)

// LoginHandlerByUID handles user authentication via Keycloak using username/password
// Takes gocloak client pointer as dependency
// Implements Resource Owner Password Credentials flow (ROPC)
func GetUserData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := config.GetLogger()
		cfg, err := config.GetConfig()
		if err != nil {
			logger.Error("failed to get config", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		var req struct {
            Username    string `json:"username"`
        }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("invalid payload")
            http.Error(w, "invalid payload", http.StatusBadRequest)
            return
        }

        // 1) Connect to LDAP
        l, err := ldap.DialURL(cfg.LDAPHost)
        if err != nil {
			logger.Error("LDAP unreachable", zap.Error(err))
            http.Error(w, fmt.Sprintf("LDAP unreachable: %v", err), http.StatusServiceUnavailable)
            return
        }
        defer l.Close()

        // 2) Service bind
        if err := l.Bind(cfg.LDAPBindDN, cfg.LDAPBindPass); err != nil {
			logger.Error("service bind failed", zap.Error(err))
            http.Error(w, fmt.Sprintf("service bind failed: %v", err), http.StatusUnauthorized)
            return
        }

        // 3) Search by username
		filter := fmt.Sprintf("(uid=%s)", ldap.EscapeFilter(req.Username))
		search := ldap.NewSearchRequest(
			cfg.LDAPBaseDN,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 1, 0, false,
			filter,
			[]string{"dn", "uid"}, // request dn and uid
			nil,
		)
		res, err := l.Search(search)
		if err != nil || len(res.Entries) == 0 {
			logger.Error("user not found", zap.Error(err))
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}
		entry := res.Entries[0]
		userUID := entry.GetAttributeValue("uid")
		if userUID == "" {
			logger.Error("LDAP entry missing uid")
			http.Error(w, "LDAP entry missing uid", http.StatusInternalServerError)
			return
		}
		
		// Формируем JSON ответ с UID пользователя
		w.Header().Set("Content-Type", "application/json")
		response := struct {
			UID string `json:"uid"`
		}{
			UID: userUID,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("failed to encode response", zap.Error(err))
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}