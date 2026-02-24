package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/config"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/controller"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/infrastructure/storage"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if config.AppConfig.Environment == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx := context.Background()

	s3Provider, err := storage.NewS3Provider(ctx, config.AppConfig.AWSRegion, config.AppConfig.DocBucketName, logger)
	if err != nil {
		logger.Fatal("Failed to initialize S3 provider", zap.Error(err))
	}

	if err := s3Provider.InitializeInfrastructure(ctx); err != nil {
		logger.Fatal("Failed to setup S3 infrastructure", zap.Error(err))
	}

	docService := service.NewDocumentService(s3Provider, logger)
	docController := controller.NewDocumentController(docService)

	router := gin.Default()
	docController.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: router,
	}

	go func() {
		logger.Info("Utility Service started", zap.String("port", config.AppConfig.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Graceful Shutdown Logic
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Warn("Shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Utility Service exited properly")
}
