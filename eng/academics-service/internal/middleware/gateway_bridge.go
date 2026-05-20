package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UUIDBridge() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantStr := c.GetString("tenant_id")
		userStr := c.GetString("user_id")

		tenantUUID, err1 := uuid.Parse(tenantStr)
		userUUID, err2 := uuid.Parse(userStr)

		if err1 != nil || err2 != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid UUID format in gateway headers",
			})
			return
		}

		c.Set("tenantID", tenantUUID)
		c.Set("userID", userUUID)

		c.Next()
	}
}
