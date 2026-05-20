package versioning

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func EnforceVersion() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientVersion := c.GetHeader("Skl-Service-Version")

		if clientVersion == "" {
			clientVersion = DefaultVersion
		}

		if !IsValid(clientVersion) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Unsupported or invalid Service version: " + clientVersion,
				"success": false,
			})
			return
		}

		c.Header("Skl-Service-Version", clientVersion)

		c.Set("api_version", clientVersion)

		c.Next()
	}
}
