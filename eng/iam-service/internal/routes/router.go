package routes

import (
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/controller"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/cache"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/middleware"
	"github.com/Single-Src-Of-Truth/Shikshakul/lib/core-go/versioning"
	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	RedisStore  *cache.SessionStore
	HealthCtrl  *controller.HealthController
	AuthCtrl    *controller.AuthController
	ProfileCtrl *controller.ProfileController
	TenantCtrl  *controller.TenantController
	RBACCtrl    *controller.RBACController
	InviteCtrl  *controller.InviteController
}

func Setup(engine *gin.Engine, deps RouterDependencies) {
	engine.Use(versioning.EnforceVersion())

	// PUBLIC ROUTES
	engine.GET("/ping", deps.HealthCtrl.Ping)
	engine.GET("/", deps.HealthCtrl.Ping)
	engine.GET("/versions", deps.HealthCtrl.GetVersions)

	// Auth Public
	engine.POST("/login", deps.AuthCtrl.Login)
	engine.POST("/invites/accept", deps.InviteCtrl.AcceptInvite)
	engine.POST("/password/forgot", deps.AuthCtrl.ForgotPassword)
	engine.POST("/password/reset", deps.AuthCtrl.ResetPassword)

	// PROTECTED ROUTES (Require Authentication)
	protected := engine.Group("")
	protected.Use(middleware.RequireAuth(deps.RedisStore))
	{
		// Profile & Sessions (No specific RBAC needed, applies to "self")
		protected.POST("/logout", deps.AuthCtrl.Logout)
		protected.GET("/verify", deps.AuthCtrl.Verify)
		protected.GET("/profile/me", deps.ProfileCtrl.GetMyProfile)
		protected.GET("/profile/sessions", deps.ProfileCtrl.GetMySessions)
		protected.PATCH("/profile/me", deps.ProfileCtrl.UpdateMyProfile)
		protected.POST("/sessions/logout-all", deps.AuthCtrl.LogoutAll)
		protected.POST("/profile/password/change", deps.AuthCtrl.ChangePassword)

		// TENANT MANAGEMENT
		tenants := protected.Group("/tenants")
		{
			tenants.POST("", middleware.RequirePermission("iam:tenants:create"), deps.TenantCtrl.CreateTenant)
			tenants.GET("", middleware.RequirePermission("iam:tenants:read"), deps.TenantCtrl.ListTenants)
			tenants.GET("/:id", middleware.RequirePermission("iam:tenants:read"), deps.TenantCtrl.GetTenant)
			tenants.PATCH("/:id", middleware.RequirePermission("iam:tenants:update"), deps.TenantCtrl.UpdateTenant)
		}

		// USER / STAFF MANAGEMENT
		users := protected.Group("/users")
		{
			users.PATCH("/:userId/status", middleware.RequirePermission("iam:users:execute"), deps.TenantCtrl.UpdateUserStatus)
			users.PATCH("/:userId/role", middleware.RequirePermission("iam:users:execute"), deps.TenantCtrl.ChangeUserRole)
			users.DELETE("/:userId", middleware.RequirePermission("iam:users:delete"), deps.TenantCtrl.DeleteUser)
		}

		// RBAC: ROLES & PERMISSIONS
		roles := protected.Group("/roles")
		{
			roles.GET("", middleware.RequirePermission("iam:roles:read"), deps.RBACCtrl.ListRoles)
			roles.POST("", middleware.RequirePermission("iam:roles:create"), deps.RBACCtrl.CreateRole)
			roles.PATCH("/:id", middleware.RequirePermission("iam:roles:update"), deps.RBACCtrl.UpdateRole)
			roles.DELETE("/:id", middleware.RequirePermission("iam:roles:delete"), deps.RBACCtrl.DeleteRole)
		}

		protected.GET("/permissions", middleware.RequirePermission("iam:permissions:read"), deps.RBACCtrl.ListPermissions)
		protected.POST("/permissions/bulk", middleware.RequirePermission("iam:permissions:create"), deps.RBACCtrl.BulkCreatePermissions)

		// ONBOARDING & INVITATIONS
		invites := protected.Group("/invites")
		{
			invites.POST("/generate", middleware.RequirePermission("iam:invites:execute"), deps.InviteCtrl.GenerateInvite)
			invites.PATCH("/:id/extend", middleware.RequirePermission("iam:invites:execute"), deps.InviteCtrl.ExtendInvite)
			invites.DELETE("/:id", middleware.RequirePermission("iam:invites:delete"), deps.InviteCtrl.CancelInvite)
		}
	}
}
