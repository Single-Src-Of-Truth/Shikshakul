package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Modulix-IT/Shikshakul-WL/iam-service/config"
	v1 "github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/api/v1"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/database"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/repository/postgres"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/repository/redis"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.LoadConfig()

	if cfg.MasterKey == "" {
		slog.Error("MASTER_KEY is missing. Cannot start IAM service.")
		os.Exit(1)
	}
	if len(cfg.MasterKey) != 64 {
		slog.Error("MASTER_KEY must be exactly 64 hex characters (32 bytes)")
		os.Exit(1)
	}

	db, err := database.NewPostgres(cfg.DBUrl)
	if err != nil {
		slog.Error("Failed to connect to Postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	rdb, err := database.NewRedis(cfg.RedisUrl)
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer rdb.Close()

	keyManager := service.NewKeyManager(db, cfg.MasterKey)
	ctx := context.Background()
	if err := keyManager.EnsureActiveKey(ctx); err != nil {
		slog.Error("Failed to initialize cryptographic keys", "error", err)
		os.Exit(1)
	}

	userRepo := postgres.NewUserRepository(db)
	tokenRepo := redis.NewTokenRepository(rdb)
	sessionRepo := postgres.NewSessionRepository(db)

	userService := service.NewUserService(userRepo, tokenRepo)
	authService := service.NewAuthService(userRepo, sessionRepo, keyManager)

	userHandler := v1.NewUserHandler(userService)
	authHandler := v1.NewAuthHandler(authService)
	jwksHandler := v1.NewJwksHandler(keyManager)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("IAM Service is Healthy"))
	})

	mux.HandleFunc("POST /api/v1/invite", userHandler.InviteUser)
	mux.HandleFunc("POST /api/v1/accept-invite", userHandler.AcceptInvite)
	mux.HandleFunc("POST /api/v1/oauth/token", authHandler.Login)
	mux.HandleFunc("POST /api/v1/oauth/refresh", authHandler.Refresh)
	mux.HandleFunc("POST /api/v1/oauth/revoke", authHandler.Revoke)
	mux.HandleFunc("GET /.well-known/jwks.json", jwksHandler.GetJWKS)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Starting IAM Service", "port", cfg.Port, "env", cfg.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	slog.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited properly")
}
