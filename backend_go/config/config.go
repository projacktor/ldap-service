package config

import (
	"strings"
	"sync"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Package-level variables with sync.Once for thread-safe initialization
var (
	cfg     *Config     // Singleton config instance
	logg    *zap.Logger // Singleton logger instance
	once    sync.Once   // Ensures config loads only once
	logOnce sync.Once   // Ensures logger initializes only once
)

// Config holds all application configuration parameters
// Struct is organized by functional groups with clear field naming
type Config struct {
	// Server configuration
	Host string
	Port int

	// Kafka configuration
	KafkaBrokers       []string
	KafkaTopic         string
	KafkaConsumerGroup string

	// Keycloak authentication
	KeycloakBaseURL      string
	KeycloakRealm        string
	KeycloakClientID     string
	KeycloakClientSecret string

	// LDAP directory services
	LDAPHost     string
	LDAPBaseDN   string
	LDAPBindDN   string // Service account for LDAP binds
	LDAPBindPass string // Sensitive - should be handled carefully
}

// Load reads configuration from environment variables and .env file
// Returns initialized Config struct or error
// Note: Viper handles automatic env variable binding
func Load() (*Config, error) {
	// Enable automatic environment variable loading
	viper.AutomaticEnv()

	// Build and populate config struct
	cfg := &Config{
		Host: viper.GetString("HOST"),
		Port: viper.GetInt("PORT"),
		// Kafka brokers are comma-separated in env
		KafkaBrokers:       strings.Split(viper.GetString("KAFKA_BROKERS"), ","),
		KafkaTopic:         viper.GetString("KAFKA_TOPIC"),
		KafkaConsumerGroup: viper.GetString("KAFKA_CONSUMER_GROUP"),

		// Keycloak OAuth2 settings
		KeycloakBaseURL:      viper.GetString("KEYCLOAK_BASE_URL"),
		KeycloakRealm:        viper.GetString("KEYCLOAK_REALM"),
		KeycloakClientID:     viper.GetString("KEYCLOAK_REST_API_CLIENT_ID"),
		KeycloakClientSecret: viper.GetString("KEYCLOAK_REST_API_CLIENT_SECRET"),

		// LDAP connection parameters
		LDAPHost: viper.GetString("LDAP_HOST"),
		// Trim quotes that might come from env files
		LDAPBaseDN:   strings.Trim(viper.GetString("LDAP_BASE_DN"), `"`),
		LDAPBindDN:   strings.Trim(viper.GetString("LDAP_USER_DN"), `"`),
		LDAPBindPass: viper.GetString("LDAP_USER_PASSWORD"),
	}

	return cfg, nil
}

// GetConfig provides thread-safe singleton access to configuration
// Uses sync.Once to ensure initialization happens exactly once
func GetConfig() (*Config, error) {
	var err error

	once.Do(func() {
		cfg, err = Load()
	})

	return cfg, err
}

// InitLogger initializes the global zap logger with Kafka integration
// Also uses sync.Once pattern for safe initialization
func InitLogger() error {
	var err error
	logOnce.Do(func() {
		cfg, err := GetConfig()
		if err != nil {
			return
		}
		// Initialize Kafka logger (errors ignored in this example)
		logg, _ = NewKafkaLogger(cfg.KafkaBrokers, cfg.KafkaTopic)
	})

	return err
}

// GetLogger provides access to the global logger instance
func GetLogger() *zap.Logger {
	return logg
}
