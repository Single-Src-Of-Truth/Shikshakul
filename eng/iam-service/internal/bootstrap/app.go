package bootstrap

import (
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/controller"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/database"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/logger"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/repository"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/routes"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/service"
	"github.com/Single-Src-Of-Truth/Shikshakul/lib/core-go/events"
)

func InitializeApp() routes.RouterDependencies {
	database.ConnectDB(logger.Log)
	database.ConnectRedis(logger.Log)

	db := database.DB
	redisStore := database.RedisStore

	rawRedisClient := redisStore.GetClient()
	eventPublisher := events.NewRedisPublisher(rawRedisClient, logger.Log)

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	rbacRepo := repository.NewRBACRepository(db)
	tenantRepo := repository.NewTenantRepository(db)
	onboardingRepo := repository.NewOnboardingRepository(db)

	authSvc := service.NewAuthService(userRepo, sessionRepo, rbacRepo, redisStore, eventPublisher, logger.Log)
	profileSvc := service.NewProfileService(userRepo, sessionRepo, redisStore, logger.Log)
	tenantSvc := service.NewTenantService(tenantRepo, userRepo, sessionRepo, rbacRepo, onboardingRepo, redisStore, logger.Log)
	rbacSvc := service.NewRBACService(rbacRepo, sessionRepo, redisStore, logger.Log)
	onboardingSvc := service.NewOnboardingService(onboardingRepo, userRepo, rbacRepo, eventPublisher, logger.Log)

	healthCtrl := controller.NewHealthController(db, redisStore)
	authCtrl := controller.NewAuthController(authSvc)
	profileCtrl := controller.NewProfileController(profileSvc)
	tenantCtrl := controller.NewTenantController(tenantSvc)
	rbacCtrl := controller.NewRBACController(rbacSvc)
	inviteCtrl := controller.NewInviteController(onboardingSvc)

	return routes.RouterDependencies{
		RedisStore:  redisStore,
		HealthCtrl:  healthCtrl,
		AuthCtrl:    authCtrl,
		ProfileCtrl: profileCtrl,
		TenantCtrl:  tenantCtrl,
		RBACCtrl:    rbacCtrl,
		InviteCtrl:  inviteCtrl,
	}
}
