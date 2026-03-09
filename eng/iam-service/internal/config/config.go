package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port                   string `mapstructure:"PORT"`
	Environment            string `mapstructure:"ENVIRONMENT"`
	DBUrl                  string `mapstructure:"DATABASE_URL"`
	RedisUrl               string `mapstructure:"REDIS_URL"`
	AutoMigrate            bool   `mapstructure:"AUTO_MIGRATE"`
	InviteExpiryHours      int    `mapstructure:"INVITE_EXPIRY_HOURS"`
	MaxConcurrentSessions  int    `mapstructure:"MAX_CONCURRENT_SESSIONS"`
	MaxSystemRoles         int    `mapstructure:"MAX_SYSTEM_ROLES"`
	MaxUsersPerSystemRole  int    `mapstructure:"MAX_USERS_PER_SYSTEM_ROLE"`
	RedirectCommandAdmin   string `mapstructure:"REDIRECT_COMMAND_ADMIN"`
	RedirectCommandTeacher string `mapstructure:"REDIRECT_COMMAND_TEACHER"`
	RedirectCommandParent  string `mapstructure:"REDIRECT_COMMAND_PARENT"`
	FrontendURL            string `mapstructure:"FRONTEND_URL"`
	CookieDomain           string `mapstructure:"COOKIE_DOMAIN"`
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{}

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.BindEnv("PORT")
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("DATABASE_URL")
	viper.BindEnv("REDIS_URL")
	viper.BindEnv("AUTO_MIGRATE")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Config file not found, relying on System Env Variables")
	}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatal("Unable to decode into struct:", err)
	}

	if AppConfig.Port == "" {
		AppConfig.Port = "9000"
	}
	if AppConfig.Environment == "" {
		AppConfig.Environment = "DEV"
	}
	if AppConfig.InviteExpiryHours == 0 {
		AppConfig.InviteExpiryHours = 168
	}
	if AppConfig.MaxConcurrentSessions == 0 {
		AppConfig.MaxConcurrentSessions = 5
	}
	if AppConfig.MaxSystemRoles == 0 {
		AppConfig.MaxSystemRoles = 2
	}
	if AppConfig.MaxUsersPerSystemRole == 0 {
		AppConfig.MaxUsersPerSystemRole = 1
	}
	if AppConfig.RedirectCommandAdmin == "" {
		AppConfig.RedirectCommandAdmin = "NAV_PORTAL_ADMIN"
	}
	if AppConfig.RedirectCommandTeacher == "" {
		AppConfig.RedirectCommandTeacher = "NAV_PORTAL_TEACHER"
	}
	if AppConfig.RedirectCommandParent == "" {
		AppConfig.RedirectCommandParent = "NAV_PORTAL_PARENT"
	}
	if AppConfig.FrontendURL == "" {
		AppConfig.FrontendURL = "https://apps.shikshakul.com"
	}
	if AppConfig.CookieDomain == "" {
		if AppConfig.Environment == "PROD" || AppConfig.Environment == "STAGING" {
			AppConfig.CookieDomain = ".shikshakul.com"
		} else {
			AppConfig.CookieDomain = "localhost"
		}
	}
}
