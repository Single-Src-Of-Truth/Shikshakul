package service

import (
	"context"
	"errors"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/cache"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/repository"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/crypto"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/dto"
	"go.uber.org/zap"
)

type ProfileService interface {
	GetMyProfile(ctx context.Context, userID string) (*dto.UserProfileResponse, error)
	GetMySessions(ctx context.Context, userID string, currentRawToken string) ([]dto.SessionResponse, error)
	UpdateMyProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) error
	DeactivateMyAccount(ctx context.Context, userID string, currentRawToken string) error
}

type profileService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	redisStore  *cache.SessionStore
	logger      *zap.Logger
}

func NewProfileService(ur repository.UserRepository, sr repository.SessionRepository, rs *cache.SessionStore, logger *zap.Logger) ProfileService {
	return &profileService{userRepo: ur, sessionRepo: sr, redisStore: rs, logger: logger}
}

func (s *profileService) GetMyProfile(ctx context.Context, userID string) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to fetch user profile", zap.Error(err))
		return nil, errors.New("profile not found")
	}

	var tID *string
	if user.TenantID != nil {
		tStr := user.TenantID.String()
		tID = &tStr
	}

	var roles []string
	var perms []string
	for _, r := range user.Roles {
		roles = append(roles, r.Name)
		for _, p := range r.Permissions {
			perms = append(perms, p.ID)
		}
	}

	return &dto.UserProfileResponse{
		ID:             user.ID.String(),
		TenantID:       tID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Identifier:     user.PrimaryIdentifier,
		IdentifierType: string(user.IdentifierType),
		Status:         string(user.Status),
		Roles:          roles,
		Permissions:    perms,
	}, nil
}

func (s *profileService) GetMySessions(ctx context.Context, userID string, currentRawToken string) ([]dto.SessionResponse, error) {
	currentTokenHash := crypto.HashToken(currentRawToken)

	sessions, err := s.sessionRepo.GetActiveSessions(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to fetch sessions", zap.Error(err))
		return nil, err
	}

	var responses []dto.SessionResponse
	for _, sess := range sessions {
		responses = append(responses, dto.SessionResponse{
			IPAddress:    sess.IPAddress,
			Browser:      sess.Browser,
			OS:           sess.OS,
			DeviceType:   sess.DeviceType,
			Location:     sess.Location,
			LastActiveAt: sess.LastActiveAt.Format("2006-01-02 15:04:05"),
			IsCurrent:    sess.OpaqueTokenHash == currentTokenHash,
		})
	}

	return responses, nil
}

func (s *profileService) UpdateMyProfile(ctx context.Context, userID string, req dto.UpdateProfileRequest) error {
	if err := s.userRepo.UpdateProfile(ctx, userID, req.FirstName, req.LastName); err != nil {
		s.logger.Error("Failed to update profile", zap.Error(err))
		return err
	}
	return nil
}

func (s *profileService) DeactivateMyAccount(ctx context.Context, userID string, currentRawToken string) error {
	if err := s.userRepo.SoftDelete(ctx, userID); err != nil {
		s.logger.Error("Failed to deactivate account", zap.Error(err))
		return err
	}

	if err := s.sessionRepo.RevokeAllUserSessions(ctx, userID); err != nil {
		s.logger.Error("Failed to revoke postgres sessions", zap.Error(err))
	}

	if currentRawToken != "" {
		_ = s.redisStore.RevokeSession(ctx, currentRawToken)
	}

	s.logger.Info("User account deactivated", zap.String("user_id", userID))
	return nil
}
