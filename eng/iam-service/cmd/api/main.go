package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/bootstrap"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/config"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/database"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/logger"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/middleware"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/routes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()
	logger.InitLogger()
	defer logger.Sync()

	logger.Info("Starting IAM Service",
		zap.String("env", config.AppConfig.Environment),
		zap.String("port", config.AppConfig.Port),
	)

	deps := bootstrap.InitializeApp()

	if config.AppConfig.AutoMigrate {
		logger.Info("Running database migrations...")
		database.Migrate(logger.Log)
	}

	bootstrap.SeedDatabase(database.DB, logger.Log)

	if config.AppConfig.Environment != "DEV" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestLogger())

	routes.Setup(engine, deps)

	srv := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server startup failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Warn("Shutting down IAM service...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	database.Close(logger.Log)
	logger.Info("IAM Service exited properly")
}
