package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/config"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/repository"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/crypto"
	"github.com/Single-Src-Of-Truth/Shikshakul/lib/core-go/events"
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
	repo      repository.OnboardingRepository
	userRepo  repository.UserRepository
	rbacRepo  repository.RBACRepository
	publisher events.EventPublisher
	logger    *zap.Logger
}

func NewOnboardingService(
	repo repository.OnboardingRepository,
	userRepo repository.UserRepository,
	rbacRepo repository.RBACRepository,
	publisher events.EventPublisher,
	logger *zap.Logger,
) OnboardingService {
	return &onboardingService{
		repo:      repo,
		userRepo:  userRepo,
		rbacRepo:  rbacRepo,
		publisher: publisher,
		logger:    logger,
	}
}

func (s *onboardingService) GenerateInvite(ctx context.Context, inviterID, tenantID, roleID string, identifier string) (string, string, error) {
	existingUser, _ := s.userRepo.FindByTenantAndIdentifier(ctx, &tenantID, identifier)
	if existingUser != nil {
		s.logger.Warn("Attempted to invite an existing user", zap.String("identifier", identifier), zap.String("tenant_id", tenantID))
		return "", "", errors.New("user already exists in this workspace")
	}

	hasPending, _ := s.repo.HasPendingInvite(ctx, tenantID, identifier)
	if hasPending {
		s.logger.Warn("Attempted to duplicate an invite", zap.String("identifier", identifier), zap.String("tenant_id", tenantID))
		return "", "", errors.New("a pending invite already exists for this email")
	}

	iID, _ := uuid.Parse(inviterID)
	tID, _ := uuid.Parse(tenantID)
	rID, _ := uuid.Parse(roleID)

	role, err := s.rbacRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return "", "", errors.New("invalid role specified")
	}

	if role.IsSystem {
		userCount, err1 := s.rbacRepo.CountUsersByRole(ctx, roleID)
		inviteCount, err2 := s.repo.CountPendingInvitesByRole(ctx, roleID)

		if err1 != nil || err2 != nil {
			return "", "", errors.New("failed to verify system role constraints")
		}

		totalAssigned := int(userCount + inviteCount)
		limit := config.AppConfig.MaxUsersPerSystemRole

		if totalAssigned >= limit {
			s.logger.Warn("System role assignment limit reached", zap.String("role", role.Name), zap.Int("limit", limit))
			return "", "", errors.New("security constraint: maximum number of users for this system role has been reached")
		}
	}

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

	safeToken := url.QueryEscape(rawToken)
	inviteLink := fmt.Sprintf("%s/auth/accept-invite?token=%s", config.AppConfig.FrontendURL, safeToken)

	event := events.EmailEvent{
		EventID:   uuid.New().String(),
		Type:      "INVITE",
		Priority:  events.PriorityHigh,
		Source:    "iam-service",
		Recipient: identifier,
		Data: map[string]interface{}{
			"invite_link": inviteLink,
		},
		CreatedAt: time.Now(),
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if pubErr := s.publisher.PublishEmailEvent(bgCtx, event); pubErr != nil {
			s.logger.Error("Failed to publish invite email event", zap.Error(pubErr), zap.String("identifier", identifier))
		}
	}()

	s.logger.Info("Invite generated and email event published", zap.String("identifier", identifier), zap.String("invite_id", invite.ID.String()))

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
