package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/config"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/cache"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/repository"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/crypto"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/fingerprint"
	"go.uber.org/zap"
)

type AuthService interface {
	Login(ctx context.Context, tenantID *string, identifier, password string, deviceInfo fingerprint.DeviceInfo) (string, *domain.User, string, error)
	Logout(ctx context.Context, rawToken string) error
	LogoutAll(ctx context.Context, userID string) error
	ChangePassword(ctx context.Context, userID, oldPass, newPass string) error
	ForgotPassword(ctx context.Context, tenantID *string, identifier string) (string, error)
	ResetPassword(ctx context.Context, tenantID *string, identifier, otp, newPass string) error
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	rbacRepo    repository.RBACRepository
	redisStore  *cache.SessionStore
	logger      *zap.Logger
}

func NewAuthService(ur repository.UserRepository, sr repository.SessionRepository, rr repository.RBACRepository, rs *cache.SessionStore, logger *zap.Logger) AuthService {
	return &authService{
		userRepo:    ur,
		sessionRepo: sr,
		rbacRepo:    rr,
		redisStore:  rs,
		logger:      logger,
	}
}

func (s *authService) Login(ctx context.Context, tenantID *string, identifier, password string, deviceInfo fingerprint.DeviceInfo) (string, *domain.User, string, error) {
	user, err := s.userRepo.FindByTenantAndIdentifier(ctx, tenantID, identifier)
	if err != nil {
		return "", nil, "", errors.New("invalid credentials")
	}

	if user.Status == domain.StatusSuspended {
		s.logger.Warn("Login blocked: User account is suspended", zap.String("user_id", user.ID.String()))
		return "", nil, "", errors.New("your account has been suspended by your administrator")
	}

	if user.Tenant != nil && !user.Tenant.IsActive {
		s.logger.Warn("Login blocked: Tenant is suspended", zap.String("tenant_id", user.TenantID.String()))
		return "", nil, "", errors.New("TENANT_SUSPENDED")
	}

	if !crypto.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, "", errors.New("invalid credentials")
	}

	maxSessions := config.AppConfig.MaxConcurrentSessions
	if maxSessions == 0 {
		maxSessions = 5
	}

	activeSessions, err := s.sessionRepo.GetActiveSessions(ctx, user.ID.String())
	if err == nil && len(activeSessions) >= maxSessions {
		overflow := len(activeSessions) - maxSessions + 1

		s.logger.Warn("Session limit reached. Revoking oldest sessions.", zap.Int("overflow", overflow))

		for i := range overflow {
			oldestSession := activeSessions[len(activeSessions)-1-i]

			s.sessionRepo.RevokeSessionByTokenHash(ctx, oldestSession.OpaqueTokenHash)
			s.redisStore.RevokeSession(ctx, oldestSession.OpaqueTokenHash)
		}
	}

	rawToken, err := crypto.GenerateOpaqueToken(crypto.PrefixSession)
	if err != nil {
		return "", nil, "", errors.New("failed to generate secure session")
	}
	tokenHash := crypto.HashToken(rawToken)

	session := &domain.UserSession{
		UserID:          user.ID,
		OpaqueTokenHash: tokenHash,
		IPAddress:       deviceInfo.IPAddress,
		Browser:         deviceInfo.Browser,
		OS:              deviceInfo.OS,
		DeviceType:      deviceInfo.DeviceType,
		Location:        "N/A",
		LastActiveAt:    time.Now(),
		ExpiresAt:       time.Now().Add(15 * 24 * time.Hour),
		IsActive:        true,
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		s.logger.Error("Failed to save session to DB", zap.Error(err))
		return "", nil, "", errors.New("internal server error during login")
	}

	var perms []string
	for _, r := range user.Roles {
		for _, p := range r.Permissions {
			perms = append(perms, p.ID)
		}
	}

	var tIDStr string
	if user.TenantID != nil {
		tIDStr = user.TenantID.String()
	}

	payload := cache.SessionPayload{
		UserID:      user.ID.String(),
		TenantID:    tIDStr,
		DeviceHash:  deviceInfo.Hash,
		Permissions: perms,
		CreatedAt:   time.Now().Unix(),
	}

	// TODO: If Redis fails, we should technically rollback Postgres, but for now, we just fail the login
	if err := s.redisStore.SaveSession(ctx, tokenHash, payload); err != nil {
		s.logger.Error("Failed to cache session", zap.Error(err))
		return "", nil, "", errors.New("internal server error during login")
	}

	redirectCommand := config.AppConfig.RedirectCommandAdmin

	if len(user.Roles) > 0 {
		primaryRole := user.Roles[0].Name

		switch primaryRole {
		case "Student":
			redirectCommand = config.AppConfig.RedirectCommandParent
		case "Teacher":
			redirectCommand = config.AppConfig.RedirectCommandTeacher
		}
	}

	s.logger.Info("User logged in successfully", zap.String("user_id", user.ID.String()))

	return rawToken, user, redirectCommand, nil
}

func (s *authService) Logout(ctx context.Context, rawToken string) error {
	tokenHash := crypto.HashToken(rawToken)

	_ = s.redisStore.RevokeSession(ctx, tokenHash)
	return s.sessionRepo.RevokeSessionByTokenHash(ctx, tokenHash)
}

func (s *authService) LogoutAll(ctx context.Context, userID string) error {
	activeSessions, err := s.sessionRepo.GetActiveSessions(ctx, userID)
	if err == nil {
		for _, session := range activeSessions {
			_ = s.redisStore.RevokeSession(ctx, session.OpaqueTokenHash)
		}
	}

	if err := s.sessionRepo.RevokeAllUserSessions(ctx, userID); err != nil {
		return err
	}

	s.logger.Info("All user sessions revoked successfully", zap.String("user_id", userID))
	return nil
}

func (s *authService) ChangePassword(ctx context.Context, userID, oldPass, newPass string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !crypto.CheckPasswordHash(oldPass, user.PasswordHash) {
		return errors.New("incorrect old password")
	}

	newHash, _ := crypto.HashPassword(newPass)
	return s.userRepo.UpdatePassword(ctx, userID, newHash)
}

func (s *authService) ForgotPassword(ctx context.Context, tenantID *string, identifier string) (string, error) {
	_, err := s.userRepo.FindByTenantAndIdentifier(ctx, tenantID, identifier)
	if err != nil {
		return "", nil
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	otp := fmt.Sprintf("%06d", n.Int64()+100000)

	err = s.redisStore.SaveOTP(ctx, identifier, otp)
	if err != nil {
		return "", err
	}

	s.logger.Info("FORGOT PASSWORD OTP GENERATED", zap.String("identifier", identifier), zap.String("otp", otp))
	return otp, nil
}

func (s *authService) ResetPassword(ctx context.Context, tenantID *string, identifier, otp, newPass string) error {
	if err := s.redisStore.VerifyAndBurnOTP(ctx, identifier, otp); err != nil {
		return err
	}

	user, err := s.userRepo.FindByTenantAndIdentifier(ctx, tenantID, identifier)
	if err != nil {
		return errors.New("user not found")
	}

	newHash, _ := crypto.HashPassword(newPass)
	return s.userRepo.UpdatePassword(ctx, user.ID.String(), newHash)
}
