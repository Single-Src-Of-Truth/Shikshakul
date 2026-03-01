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
	AppConfig = &Config{}

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.BindEnv("PORT")
	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("REDIS_URL")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("AUTO_MIGRATE")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Config file not found, relying on System Env Variables")
	}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatal("Unable to decode into struct:", err)
	}

	if AppConfig.Port == "" {
		AppConfig.Port = "9001"
	}
}
