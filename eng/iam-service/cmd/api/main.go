package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/bootstrap"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/config"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/controller"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/cache"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/db"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/middleware"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/repository"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if config.AppConfig.Environment == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	}

	database, err := db.Connect(config.AppConfig.DBUrl)
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}

	redisStore, err := cache.NewSessionStore(config.AppConfig.RedisUrl)
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	if config.AppConfig.AutoMigrate {
		logger.Info("Running Database Auto-Migrations...")
		err = database.AutoMigrate(
			&domain.Tenant{},
			&domain.User{},
			&domain.Role{},
			&domain.Permission{},
			&domain.Invitation{},
			&domain.UserSession{},
		)
		if err != nil {
			logger.Fatal("Failed to run migrations", zap.Error(err))
		}
	}

	bootstrap.SeedDatabase(database.DB, logger)

	userRepo := repository.NewUserRepository(database.DB)
	sessionRepo := repository.NewSessionRepository(database.DB)
	rbacRepo := repository.NewRBACRepository(database.DB)
	tenantRepo := repository.NewTenantRepository(database.DB)
	onboardingRepo := repository.NewOnboardingRepository(database.DB)

	authSvc := service.NewAuthService(userRepo, sessionRepo, rbacRepo, redisStore, logger)
	profileSvc := service.NewProfileService(userRepo, sessionRepo, redisStore, logger)
	tenantSvc := service.NewTenantService(tenantRepo, userRepo, sessionRepo, redisStore, logger)
	rbacSvc := service.NewRBACService(rbacRepo, sessionRepo, redisStore, logger)
	onboardingSvc := service.NewOnboardingService(onboardingRepo, logger)

	healthCtrl := controller.NewHealthController(database, redisStore)
	authCtrl := controller.NewAuthController(authSvc)
	profileCtrl := controller.NewProfileController(profileSvc)
	tenantCtrl := controller.NewTenantController(tenantSvc)
	rbacCtrl := controller.NewRBACController(rbacSvc)
	inviteCtrl := controller.NewInviteController(onboardingSvc)

	router := gin.Default()

	router.GET("/ping", healthCtrl.Ping)
	router.GET("/", healthCtrl.Ping)

	api := router.Group("/api/v1/iam")
	{
		api.POST("/login", authCtrl.Login)
		api.POST("/invites/accept", inviteCtrl.AcceptInvite)
		api.POST("/password/forgot", authCtrl.ForgotPassword)
		api.POST("/password/reset", authCtrl.ResetPassword)

		protected := api.Group("")
		protected.Use(middleware.RequireAuth(redisStore))
		{
			protected.POST("/logout", authCtrl.Logout)
			protected.GET("/verify", authCtrl.Verify)

			protected.GET("/profile/me", profileCtrl.GetMyProfile)
			protected.GET("/profile/sessions", profileCtrl.GetMySessions)
			protected.PATCH("/profile/me", profileCtrl.UpdateMyProfile)
			protected.POST("/sessions/logout-all", authCtrl.LogoutAll)
			protected.POST("/profile/password/change", authCtrl.ChangePassword)

			tenants := protected.Group("/tenants")
			tenants.Use(middleware.RequirePermission("iam:tenants:write"))
			{
				tenants.POST("", tenantCtrl.CreateTenant)
				tenants.GET("", tenantCtrl.ListTenants)
				tenants.GET("/:id", tenantCtrl.GetTenant)
				tenants.PATCH("/:id", tenantCtrl.UpdateTenant)
			}

			staffGroup := protected.Group("/users")
			staffGroup.Use(middleware.RequirePermission("iam:roles:write"))
			{
				staffGroup.PATCH("/:userId/status", tenantCtrl.UpdateUserStatus)
				staffGroup.DELETE("/:userId", tenantCtrl.DeleteUser)
			}

			rbacGroup := protected.Group("")
			rbacGroup.Use(middleware.RequirePermission("iam:roles:write"))
			{
				rbacGroup.GET("/roles", rbacCtrl.ListRoles)
				rbacGroup.GET("/permissions", rbacCtrl.ListPermissions)
				rbacGroup.POST("/roles", rbacCtrl.CreateRole)
				rbacGroup.POST("/invites/generate", inviteCtrl.GenerateInvite)
				rbacGroup.PATCH("/invites/:id/extend", inviteCtrl.ExtendInvite)
				rbacGroup.DELETE("/invites/:id", inviteCtrl.CancelInvite)
				rbacGroup.POST("/permissions/bulk", rbacCtrl.BulkCreatePermissions)
				rbacGroup.PATCH("/roles/:id", rbacCtrl.UpdateRole)
				rbacGroup.DELETE("/roles/:id", rbacCtrl.DeleteRole)
			}
		}
	}

	srv := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: router,
	}

	go func() {
		logger.Info("IAM Service started", zap.String("port", config.AppConfig.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Warn("Shutting down IAM service...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("IAM Service exited properly")
}
