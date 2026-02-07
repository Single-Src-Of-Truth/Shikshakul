package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/Modulix-IT/Shikshakul-Backend-MicroService/academics-service/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if tenantID, ok := claims["tenant_id"].(string); ok {
				c.Set("tenantID", uuid.MustParse(tenantID))
			}
			if roles, ok := claims["roles"].([]interface{}); ok {
				roleList := make([]string, len(roles))
				for i, r := range roles {
					roleList[i] = r.(string)
				}
				c.Set("roles", roleList)
			}
			if sub, ok := claims["sub"].(string); ok {
				c.Set("userID", uuid.MustParse(sub))
			}
		}

		c.Next()
	}
}

func RoleGuard(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles := c.GetStringSlice("roles")
		for _, uRole := range userRoles {
			if slices.Contains(allowedRoles, uRole) {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient Permissions"})
	}
}
