package controller

import (
	"net/http"
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/database"
	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (ctrl *HealthController) HealthCheck(c *gin.Context) {
	postgresStatus := "DOWN"
	if database.CheckPostgres() {
		postgresStatus = "UP"
	}

	redisStatus := "DOWN"
	if database.CheckRedis() {
		redisStatus = "UP"
	}

	status := "UP"
	httpCode := http.StatusOK

	if postgresStatus == "DOWN" || redisStatus == "DOWN" {
		status = "DOWN"
		httpCode = http.StatusServiceUnavailable
	}

	c.JSON(httpCode, gin.H{
		"service":   "Academics Service(acads-pa.client1.shikshakul.com)",
		"status":    status,
		"timestamp": time.Now().Format(time.RFC3339),
		"dependencies": gin.H{
			"postgresql": postgresStatus,
			"redis":      redisStatus,
		},
	})
}
