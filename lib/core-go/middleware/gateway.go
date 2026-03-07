package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireGatewayAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("Skl-User-Id")

		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Missing gateway authorization headers. Direct access forbidden.",
			})
			return
		}

		permsStr := c.GetHeader("Skl-Permissions")
		var permissions []string
		if permsStr != "" {
			permissions = strings.Split(permsStr, ",")
		}

		c.Set("user_id", userID)
		c.Set("tenant_id", c.GetHeader("Skl-Tenant-Id"))
		c.Set("permissions", permissions)

		c.Next()
	}
}

func RequireGatewayPermission(requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permsInterface, exists := c.Get("permissions")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Missing permission context",
			})
			return
		}

		permissions, ok := permsInterface.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Invalid permission format in context",
			})
			return
		}

		if !slices.Contains(permissions, requiredPerm) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Insufficient permissions: requires " + requiredPerm,
			})
			return
		}

		c.Next()
	}
}
