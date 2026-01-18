package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env       string
	Port      string
	DBUrl     string
	RedisUrl  string
	MasterKey string
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println("No .env.local file found, relying on system env vars")
	}

	return &Config{
		Env:       getEnv("ENV", "production"),
		Port:      getEnv("PORT", "8080"),
		DBUrl:     getEnv("DB_URL", ""),
		RedisUrl:  getEnv("REDIS_URL", ""),
		MasterKey: getEnv("MASTER_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return fallback
}
