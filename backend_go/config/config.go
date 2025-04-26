package config

import (
    "strings"
    "sync"

    "github.com/spf13/viper"
    "go.uber.org/zap"
)

var (
    cfg  *Config
    logg *zap.Logger
    once sync.Once
    logOnce sync.Once
)

// Config holds all the app configuration values
type Config struct {
    Host                string
    Port                int
    KafkaBrokers        []string
    KafkaTopic          string
    KafkaConsumerGroup  string

    // Keycloak
    KeycloakBaseURL     string
    KeycloakRealm       string
    KeycloakClientID    string
    KeycloakClientSecret string

    // LDAP
    LDAPHost       string
    LDAPBaseDN     string
    LDAPBindDN     string // service account
    LDAPBindPass   string
}

// Load reads .env and environment, applies defaults, and returns a Config
func Load() (*Config, error) {
    // 1) Tell viper to pick up all ENV variables
    viper.AutomaticEnv()

    // 2) Build the config struct
    cfg := &Config{
        Host:         viper.GetString("HOST"),
        Port:         viper.GetInt("PORT"),
        KafkaBrokers: strings.Split(viper.GetString("KAFKA_BROKERS"), ","),
        KafkaTopic:   viper.GetString("KAFKA_TOPIC"),
        KafkaConsumerGroup: viper.GetString("KAFKA_CONSUMER_GROUP"),
        
        KeycloakBaseURL:     viper.GetString("KEYCLOAK_BASE_URL"),
        KeycloakRealm:       viper.GetString("KEYCLOAK_REALM"),
        KeycloakClientID:    viper.GetString("KEYCLOAK_REST_API_CLIENT_ID"),
        KeycloakClientSecret:viper.GetString("KEYCLOAK_REST_API_CLIENT_SECRET"),

        LDAPHost:     viper.GetString("LDAP_HOST"),
        LDAPBaseDN:   strings.Trim(viper.GetString("LDAP_BASE_DN"), `"`),
        LDAPBindDN:   viper.GetString("LDAP_USER_DN"),
        LDAPBindPass: viper.GetString("LDAP_USER_PASSWORD"),
    }

    return cfg, nil
}

func GetConfig () (*Config, error) {
    var err error

    once.Do(func() {
        cfg, err = Load()
    })

    return cfg, err
}

func InitLogger() error {
    var err error
    logOnce.Do(func() {
        cfg, err := GetConfig()
        if err != nil {
            return 
        }
        logg, _ = NewKafkaLogger(cfg.KafkaBrokers, cfg.KafkaTopic) // log is global variable
    })

    return err
}

func GetLogger() *zap.Logger {
    return logg // log is global variable
}