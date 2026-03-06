package service

import (
	"context"
	"errors"
	"time"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/config"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/repository"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/crypto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OnboardingService interface {
	GenerateInvite(ctx context.Context, inviterID, tenantID, roleID string, identifier string) (string, string, error)
	AcceptInvite(ctx context.Context, rawToken, password, firstName, lastName string) (*domain.User, error)
	ExtendInvite(ctx context.Context, inviteID string) error
	CancelInvite(ctx context.Context, inviteID string) error
}

type onboardingService struct {
	repo   repository.OnboardingRepository
	logger *zap.Logger
}

func NewOnboardingService(repo repository.OnboardingRepository, logger *zap.Logger) OnboardingService {
	return &onboardingService{repo: repo, logger: logger}
}

func (s *onboardingService) GenerateInvite(ctx context.Context, inviterID, tenantID, roleID string, identifier string) (string, string, error) {
	iID, _ := uuid.Parse(inviterID)
	tID, _ := uuid.Parse(tenantID)
	rID, _ := uuid.Parse(roleID)

	rawToken, err := crypto.GenerateOpaqueToken(crypto.PrefixInvite)
	if err != nil {
		return "", "", errors.New("failed to generate invite token")
	}

	tokenHash := crypto.HashToken(rawToken)

	invite := &domain.Invitation{
		TenantID:   tID,
		InviterID:  iID,
		Identifier: identifier,
		RoleID:     rID,
		TokenHash:  tokenHash,
		Status:     domain.InvitePending,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // TODO: Replace this from config
	}

	if err := s.repo.CreateInvitation(ctx, invite); err != nil {
		s.logger.Error("Failed to save invitation", zap.Error(err))
		return "", "", errors.New("failed to create invite")
	}

	// TODO (Future): Send event to Utility Service to dispatch Email/SMS with `rawToken`
	s.logger.Info("Invite generated successfully", zap.String("identifier", identifier), zap.String("invite_id", invite.ID.String()))

	return invite.ID.String(), rawToken, nil
}

func (s *onboardingService) AcceptInvite(ctx context.Context, rawToken, password, firstName, lastName string) (*domain.User, error) {
	tokenHash := crypto.HashToken(rawToken)

	invite, err := s.repo.FindPendingInvite(ctx, tokenHash)
	if err != nil {
		s.logger.Warn("Invalid or expired invite token attempted", zap.Error(err))
		return nil, errors.New("invite link is invalid or has expired")
	}

	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to secure password")
	}

	user := &domain.User{
		TenantID:          &invite.TenantID,
		IdentifierType:    domain.IdentifierEmail,
		PrimaryIdentifier: invite.Identifier,
		PasswordHash:      hashedPassword,
		FirstName:         firstName,
		LastName:          lastName,
		Status:            domain.StatusActive,
	}

	if err := s.repo.ConsumeInviteAndCreateUser(ctx, invite, user); err != nil {
		s.logger.Error("Failed to consume invite and create user", zap.Error(err))
		return nil, errors.New("failed to finalize onboarding")
	}

	s.logger.Info("User onboarded successfully", zap.String("user_id", user.ID.String()))

	return user, nil
}

func (s *onboardingService) ExtendInvite(ctx context.Context, inviteID string) error {
	invite, err := s.repo.GetInvitationByID(ctx, inviteID)
	if err != nil {
		return errors.New("invitation not found")
	}

	if invite.Status != domain.InvitePending {
		return errors.New("only pending invitations can be extended")
	}

	hours := config.AppConfig.InviteExpiryHours
	if hours == 0 {
		hours = 168 // Fallback to 7 days if config is missing
	}

	newExpiry := time.Now().Add(time.Duration(hours) * time.Hour)

	if err := s.repo.UpdateInvitationExpiry(ctx, inviteID, newExpiry); err != nil {
		s.logger.Error("Failed to extend invitation", zap.Error(err))
		return errors.New("could not update expiration date")
	}

	// TODO: utity service to resend the email to `invite.Identifier` letting them know they have more time!
	s.logger.Info("Invitation extended", zap.String("invite_id", inviteID))
	return nil
}

func (s *onboardingService) CancelInvite(ctx context.Context, inviteID string) error {
	if err := s.repo.SoftDeleteInvitation(ctx, inviteID); err != nil {
		s.logger.Error("Failed to cancel invitation", zap.Error(err))
		return errors.New("could not cancel invitation")
	}

	s.logger.Info("Invitation cancelled", zap.String("invite_id", inviteID))
	return nil
}
