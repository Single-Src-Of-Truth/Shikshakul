package main

import (
	"log"

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
		log.Println("AutoMigrate is enabled. Running migrations...")
		database.Migrate()
	} else {
		log.Println("AutoMigrate is disabled. Skipping migrations.")
	}

	r := gin.Default()

	routes.RegisterRoutes(r, app)

	log.Printf("Server starting on port %s in %s mode", config.AppConfig.Port, config.AppConfig.Environment)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
