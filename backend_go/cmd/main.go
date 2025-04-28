package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-chi/chi/v5"

	"BackendGoLdap/auth"
	"BackendGoLdap/config"
)

// init function runs before main() to initialize environment variables
func init() {
	// Load .env
	viper.AutomaticEnv()
}

func main() {
	// Load typed application configuration
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	// Initialize zap logger with Kafka integration
	if err := config.InitLogger(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	// Ensure logger flushes all buffered logs before exit
	defer config.GetLogger().Sync()

	logger := config.GetLogger()

	// Log configuration values for debugging
	logger.Info("LDAP Host", zap.String("version", cfg.LDAPHost))

	// Add base context to logger
	logger = logger.With(
		zap.String("service", "ldap-api"),
		zap.String("env", "development"),
	)
	// Initialize Keycloak client
	kc := gocloak.NewClient(cfg.KeycloakBaseURL)
	// Configure TLS to skip certificate verification (since we use self-signed certs)
	kc.RestyClient().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// Build Chi router
	r := chi.NewRouter()

	// Public application root endpoint without authentication
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// Log incoming request details
		logger.Info("http request received",
			zap.String("path", r.URL.Path),
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
		)
		_, err = w.Write([]byte("Hello, from LDAP on golang!"))
		if err != nil {
			logger.Error("failed to write response", zap.Error(err))
		}
	})

	// Public Keycloak + LDAP login
	r.Post("/auth/login", auth.LoginHandlerByUID(kc, cfg))

	// Protected routes group
	r.Group(func(r chi.Router) {
		// Apply authentication middleware to all routes in this group
		r.Use(auth.AuthMiddleware(kc, cfg))
		//protected endpoint
		r.Get("/api/protected", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("🔒 your secret data"))
		})
	})

	// Start HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info("starting server",
		zap.String("address", addr),
		zap.Strings("kafka_brokers", cfg.KafkaBrokers),
		zap.String("kafka_topic", cfg.KafkaTopic),
	)
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal("server error", zap.Error(err))
	}
}
