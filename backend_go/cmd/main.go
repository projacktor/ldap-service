// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/go-chi/chi/v5"

	"BackendGoLdap/config"
    "BackendGoLdap/routes"
	"BackendGoLdap/logger"
)

func init() {
    // Load .env    
    viper.AutomaticEnv()
}

func main() {
    // Load typed config
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("config load error: %v", err)
    }

    // Initialize zap + Kafka logger
    logger, err := logger.NewKafkaLogger(cfg.KafkaBrokers, cfg.KafkaTopic)
    if err != nil {
        log.Fatalf("failed to init logger: %v", err)
    }
    defer logger.Sync()

    // Add base context to logger
    logger = logger.With(
        zap.String("service", "ldap-api"),
        zap.String("env", "development"),
    )

    // Build Chi router
    r := chi.NewRouter()

    routes.InitRoutes(r)

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
