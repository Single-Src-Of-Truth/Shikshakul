package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"PORT"`
	DBUrl       string `mapstructure:"DATABASE_URL"`
	RedisUrl    string `mapstructure:"REDIS_URL"`
	JwtSecret   string `mapstructure:"JWT_SECRET"`
	Environment string `mapstructure:"ENVIRONMENT"`
	AutoMigrate bool   `mapstructure:"AUTO_MIGRATE"`
}

var AppConfig *Config

func LoadConfig() {
	viper.SetConfigFile("config.env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config.env file found, relying on System Env Variables")
		} else {
			log.Fatal("Error reading config file:", err)
		}
	}

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatal("Unable to decode into struct:", err)
	}

	if AppConfig.Port == "" {
		AppConfig.Port = "9001"
	}
}
