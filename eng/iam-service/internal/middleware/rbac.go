package middleware

import (
	"net/http"
	"slices"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/response"
	"github.com/gin-gonic/gin"
)

func RequirePermission(requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permsInterface, exists := c.Get("permissions")
		if !exists {
			response.Error(c, http.StatusForbidden, "missing permission context")
			c.Abort()
			return
		}

		permissions, ok := permsInterface.([]string)
		if !ok {
			response.Error(c, http.StatusInternalServerError, "invalid permission format in context")
			c.Abort()
			return
		}

		hasAccess := slices.Contains(permissions, requiredPerm)

		if !hasAccess {
			response.Error(c, http.StatusForbidden, "insufficient permissions: requires "+requiredPerm)
			c.Abort()
			return
		}

		c.Next()
	}
}
