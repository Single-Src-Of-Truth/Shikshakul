package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/bootstrap"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/config"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/database"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/logger"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/middleware"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/routes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()

	logger.InitLogger()
	defer logger.Sync()

	logger.Info("Starting Academics Service",
		zap.String("env", config.AppConfig.Environment),
		zap.String("port", config.AppConfig.Port),
	)

	app := bootstrap.InitializeApp()

	if config.AppConfig.AutoMigrate {
		logger.Info("Running database migrations...")
		database.Migrate()
	}

	if config.AppConfig.Environment != "DEV" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	routes.RegisterRoutes(r, app)

	srv := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server startup failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Warn("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Cleaning up resources...")
	database.Close()

	logger.Info("Server exited properly")
}
