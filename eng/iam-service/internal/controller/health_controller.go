package controller

import (
	"net/http"
	"time"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/cache"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/db"
	"github.com/Single-Src-Of-Truth/Shikshakul/lib/core-go/versioning"
	"github.com/gin-gonic/gin"
)

type HealthController struct {
	db    *db.Database
	redis *cache.SessionStore
}

func NewHealthController(database *db.Database, redisStore *cache.SessionStore) *HealthController {
	return &HealthController{db: database, redis: redisStore}
}

func (ctrl *HealthController) Ping(c *gin.Context) {
	ctx := c.Request.Context()
	status := gin.H{"postgres": "UP", "redis": "UP"}
	isHealthy := true

	if err := ctrl.db.Ping(ctx); err != nil {
		status["postgres"] = "DOWN"
		isHealthy = false
	}
	if err := ctrl.redis.Ping(ctx); err != nil {
		status["redis"] = "DOWN"
		isHealthy = false
	}

	response := gin.H{
		"success":    isHealthy,
		"timestamp":  time.Now().Format(time.RFC3339),
		"components": status,
	}

	if isHealthy {
		response["status"] = "UP"
		c.JSON(http.StatusOK, response)
	} else {
		response["status"] = "DOWN"
		c.JSON(http.StatusServiceUnavailable, response)
	}
}

func (ctrl *HealthController) GetVersions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"default": versioning.DefaultVersion,
		"history": versioning.Registry,
	})
}
