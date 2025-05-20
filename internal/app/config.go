package app

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
	DB         DBConfig
	Auth       AuthConfig
}

type DBConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     int
	SSLMode  string
}

type AuthConfig struct {
	APIKey string
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	if config.Auth.APIKey == "" {
		log.Fatal("API key must be set in config or environment")
	}

	return &config
}
