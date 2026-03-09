package middleware

import (
	"net/http"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/config"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/cache"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/crypto"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/fingerprint"
	"github.com/gin-gonic/gin"
)

const SessionCookieName = "skl_session"

func RequireAuth(redisStore *cache.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(SessionCookieName)
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized access", "success": false})
			return
		}

		tokenHash := crypto.HashToken(token)
		deviceInfo := fingerprint.Extract(c.Request)

		payload, err := redisStore.VerifyAndRefresh(c.Request.Context(), tokenHash, deviceInfo.Hash)
		if err != nil {
			c.SetSameSite(http.SameSiteLaxMode)
			c.SetCookie("skl_session", "", -1, "/", config.AppConfig.CookieDomain, true, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "success": false})
			return
		}

		if payload.TenantID != "" {
			if redisStore.IsTenantSuspended(c.Request.Context(), payload.TenantID) {
				_ = redisStore.RevokeSession(c.Request.Context(), tokenHash)

				c.SetSameSite(http.SameSiteLaxMode)
				c.SetCookie("skl_session", "", -1, "/", config.AppConfig.CookieDomain, true, true)

				c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{
					"error":   "Service suspended. Please contact administration to resume access.",
					"success": false,
				})
				return
			}
		}

		c.Set("user_id", payload.UserID)
		c.Set("tenant_id", payload.TenantID)
		c.Set("permissions", payload.Permissions)

		c.Next()
	}
}
