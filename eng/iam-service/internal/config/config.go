package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port                  string `mapstructure:"PORT"`
	Environment           string `mapstructure:"ENVIRONMENT"`
	DBUrl                 string `mapstructure:"DATABASE_URL"`
	RedisUrl              string `mapstructure:"REDIS_URL"`
	AutoMigrate           bool   `mapstructure:"AUTO_MIGRATE"`
	InviteExpiryHours     int    `mapstructure:"INVITE_EXPIRY_HOURS"`
	MaxConcurrentSessions int    `mapstructure:"MAX_CONCURRENT_SESSIONS"`
}

var AppConfig *Config

func LoadConfig() {
	viper.SetConfigFile("config.env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config.env file found, relying on System Env Variables")
	}

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatal("Unable to decode into struct:", err)
	}

	if AppConfig.Port == "" {
		AppConfig.Port = "9000"
	}

	if AppConfig.InviteExpiryHours == 0 {
		AppConfig.InviteExpiryHours = 168
	}

	if AppConfig.MaxConcurrentSessions == 0 {
		AppConfig.MaxConcurrentSessions = 5
	}
}
