package database

import (
	"context"
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/config"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

func ConnectDB() {
	var logLevel gormLogger.LogLevel
	if config.AppConfig.Environment == "DEV" {
		logLevel = gormLogger.Info
	} else {
		logLevel = gormLogger.Error
	}

	var err error
	DB, err = gorm.Open(postgres.Open(config.AppConfig.DBUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
	})

	if err != nil {
		logger.Fatal("Failed to connect to Database", zap.Error(err))
	}
	logger.Info("Connected to PostgreSQL", zap.String("log_level", config.AppConfig.Environment))
}

func ConnectRedis() {
	opt, err := redis.ParseURL(config.AppConfig.RedisUrl)
	if err != nil {
		logger.Fatal("Invalid Redis URL", zap.Error(err))
	}

	RedisClient = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Connected to Redis")
}

func CheckPostgres() bool {
	if DB == nil {
		return false
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return false
	}
	return true
}

func CheckRedis() bool {
	if RedisClient == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		return false
	}
	return true
}

func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			logger.Error("Failed to get generic database object", zap.Error(err))
		} else {
			if err := sqlDB.Close(); err != nil {
				logger.Error("Error closing PostgreSQL connection", zap.Error(err))
			} else {
				logger.Info("PostgreSQL connection closed")
			}
		}
	}

	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			logger.Error("Error closing Redis connection", zap.Error(err))
		} else {
			logger.Info("Redis connection closed")
		}
	}
}
