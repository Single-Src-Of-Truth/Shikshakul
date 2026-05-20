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
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/infrastructure/email"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/infrastructure/queue"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/infrastructure/storage"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/routes"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/utility-service/internal/service"
	"github.com/Single-Src-Of-Truth/Shikshakul/lib/core-go/events"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if config.AppConfig.Environment == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx, cancelWorkers := context.WithCancel(context.Background())
	defer cancelWorkers()

	opt, err := redis.ParseURL(config.AppConfig.RedisUrl)
	if err != nil {
		logger.Fatal("Invalid Redis URL", zap.Error(err))
	}
	redisClient := redis.NewClient(opt)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	logger.Info("Connected to Redis successfully")

	s3Provider, err := storage.NewS3Provider(ctx, config.AppConfig.AWSRegion, config.AppConfig.DocBucketName, logger)
	if err != nil {
		logger.Fatal("Failed to initialize S3 provider", zap.Error(err))
	}
	if err := s3Provider.InitializeInfrastructure(ctx); err != nil {
		logger.Fatal("Failed to setup S3 infrastructure", zap.Error(err))
	}

	sesProvider, err := email.NewSESProvider(ctx, config.AppConfig.AWSRegion, config.AppConfig.SESSenderEmail, logger)
	if err != nil {
		logger.Fatal("Failed to initialize SES provider", zap.Error(err))
	}

	docService := service.NewDocumentService(s3Provider, logger)

	emailService := service.NewEmailService(sesProvider, 14, logger)

	streamConsumer := queue.NewStreamConsumer(redisClient, emailService, logger)
	streamConsumer.Start(ctx, events.StreamEmailHigh, "utility-email-group", "worker-1")
	streamConsumer.Start(ctx, events.StreamEmailLow, "utility-email-group", "worker-1")

	docController := controller.NewDocumentController(docService)
	engine := gin.Default()
	routes.Setup(engine, docController)

	srv := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: engine,
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

	cancelWorkers()

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	redisClient.Close()
	logger.Info("Utility Service exited properly")
}
