package middleware

import (
	"time"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		tenantID := ""
		if t, exists := c.Get("tenantID"); exists {
			if id, ok := t.(string); ok {
				tenantID = id
			}
		}

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("tenant_id", tenantID),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error("Request Error", append(fields, zap.String("error", e))...)
			}
		} else {
			if status >= 500 {
				logger.Error("Server Error", fields...)
			} else if status >= 400 {
				logger.Warn("Client Error", fields...)
			} else {
				logger.Info("Request Processed", fields...)
			}
		}
	}
}
