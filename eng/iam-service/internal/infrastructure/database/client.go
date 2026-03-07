package database

import (
	"context"
	"time"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/config"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/cache"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	DB         *gorm.DB
	RedisStore *cache.SessionStore
)

func ConnectDB(logger *zap.Logger) {
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
	logger.Info("Connected to PostgreSQL", zap.String("env", config.AppConfig.Environment))
}

func ConnectRedis(logger *zap.Logger) {
	var err error
	RedisStore, err = cache.NewSessionStore(config.AppConfig.RedisUrl)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisStore.Ping(ctx); err != nil {
		logger.Fatal("Redis ping failed", zap.Error(err))
	}

	logger.Info("Connected to Redis")
}

func Migrate(logger *zap.Logger) {
	err := DB.AutoMigrate(
		&domain.Tenant{},
		&domain.User{},
		&domain.Role{},
		&domain.Permission{},
		&domain.Invitation{},
		&domain.UserSession{},
	)
	if err != nil {
		logger.Fatal("Migration Failed", zap.Error(err))
	}
	logger.Info("Database Migration Completed")
}

func Close(logger *zap.Logger) {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			_ = sqlDB.Close()
			logger.Info("PostgreSQL connection closed")
		}
	}
}
