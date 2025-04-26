package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"
    "github.com/Nerzal/gocloak/v13"

	"BackendGoLdap/config"
    "BackendGoLdap/auth"
)

func init() {
    // Load .env    
    viper.AutomaticEnv()
}

func main() {
    // Load typed config
    cfg, err := config.GetConfig()
    if err != nil {
        log.Fatalf("config load error: %v", err)
    }

    // Initialize zap + Kafka logger
    if err := config.InitLogger(); err != nil {
        log.Fatalf("failed to init logger: %v", err)
    }
    defer config.GetLogger().Sync()

    logger := config.GetLogger()
    
    logger.Info("LDAP Host", zap.String("version", cfg.LDAPHost))

    // Add base context to logger
    logger = logger.With(
        zap.String("service", "ldap-api"),
        zap.String("env", "development"),
    )

    kc := gocloak.NewClient(cfg.KeycloakBaseURL)

    // Build Chi router
    r := chi.NewRouter()

    // Application endpoint
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

    // Public LDAP→Keycloak login
    r.Post("/auth/login", auth.LoginHandler(kc, cfg))

    // Protected routes
    r.Group(func(r chi.Router) {
        r.Use(auth.AuthMiddleware(kc, cfg))
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
