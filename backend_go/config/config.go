package config

import (
    "strings"
    "sync"

    "github.com/spf13/viper"
    "go.uber.org/zap"

    "BackendGoLdap/logger"
)

var (
    cfg  *Config
    log *zap.Logger
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

    BaseUrl             string
	Realm               string
	RestApiClientId     string
	RestApiClientSecret string
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

        BaseUrl:             viper.GetString("KEYCLOAK_BASE_URL"),
		Realm:               viper.GetString("KEYCLOAK_REALM"),
		RestApiClientId:     viper.GetString("KEYCLOAK_REST_API_CLIENT_ID"),
		RestApiClientSecret: viper.GetString("KEYCLOAK_REST_API_CLIENT_SECRET"),
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
        log, err = logger.NewKafkaLogger(cfg.KafkaBrokers, cfg.KafkaTopic) // log is global variable
    })

    return err
}

func GetLogger() *zap.Logger {
    return log // log is global variable
}