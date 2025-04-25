package config

import (
    "strings"

    "github.com/spf13/viper"
)

// Config holds all the app configuration values
type Config struct {
    Host                string
    Port                int
    KafkaBrokers        []string
    KafkaTopic          string
    KafkaConsumerGroup  string
    DBHost              string
    DBPort              int
    DBUser              string
    DBPass              string
    DBName              string
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
        DBHost:       viper.GetString("DB_HOST"),
        DBPort:       viper.GetInt("DB_PORT"),
        DBUser:       viper.GetString("DB_USER"),
        DBPass:       viper.GetString("DB_PASS"),
        DBName:       viper.GetString("DB_NAME"),
    }

    return cfg, nil
}
