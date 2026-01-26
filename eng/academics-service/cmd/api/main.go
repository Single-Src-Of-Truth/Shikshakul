package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/bootstrap"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/config"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/database"
	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	app := bootstrap.InitializeApp()

	if config.AppConfig.AutoMigrate {
		log.Println("Running migrations...")
		database.Migrate()
	}

	r := gin.Default()
	routes.RegisterRoutes(r, app)

	srv := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", config.AppConfig.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}
