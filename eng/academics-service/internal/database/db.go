package database

import (
	"context"
	"log"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

func ConnectDB() {
	var err error
	dsn := config.AppConfig.DBUrl

	dbConfig := &gorm.Config{}
	if config.AppConfig.Environment == "DEV" {
		dbConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	DB, err = gorm.Open(postgres.Open(dsn), dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	log.Println("Connected to PostgreSQL")

	opt, err := redis.ParseURL(config.AppConfig.RedisUrl)
	if err != nil {
		log.Fatal("Invalid Redis URL:", err)
	}

	RedisClient = redis.NewClient(opt)

	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis")
}
