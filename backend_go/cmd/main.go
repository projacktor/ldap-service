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
	"github.com/go-chi/cors"

	"BackendGoLdap/auth"
	"BackendGoLdap/config"
)

func init() {
	viper.AutomaticEnv()
}

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	if err := config.InitLogger(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer config.GetLogger().Sync()

	logger := config.GetLogger()

	logger.Info("LDAP Host", zap.String("version", cfg.LDAPHost))

	logger = logger.With(
		zap.String("service", "ldap-api"),
		zap.String("env", "development"),
	)

	// Initialize Keycloak gocloak client
	kc := gocloak.NewClient(cfg.KeycloakBaseURL)
	kc.RestyClient().
		SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	// Initialize OIDC provider and verifier
	if err := auth.InitOIDC(); err != nil {
		log.Fatalf("failed to initialize OIDC: %v", err)
	}

	r := chi.NewRouter()

	// CORS settings
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
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

	// Public login with username/password
	r.Post("/auth/login", auth.LoginHandlerByUID(kc))

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware())

		// Get users from LDAP db through keycloak
		r.Get("/users", auth.GetUserDataFromToken())
	})

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
