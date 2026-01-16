package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/repository/postgres"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/repository/redis"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/pkg/crypto"
)

type UserService struct {
	userRepo  *postgres.UserRepository
	tokenRepo *redis.TokenRepository
}

func NewUserService(userRepo *postgres.UserRepository, tokenRepo *redis.TokenRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

// InviteUser handles the logic of onboarding a new person
func (s *UserService) InviteUser(ctx context.Context, tenantID string, req domain.InviteRequest) (string, error) {
	// 1. Check if user already exists
	existing, err := s.userRepo.GetByEmail(ctx, tenantID, req.Email)
	if err != nil {
		return "", fmt.Errorf("database error checking user: %w", err)
	}
	if existing != nil {
		return "", errors.New("user with this email already exists in this school")
	}

	randomPass := make([]byte, 32)
	rand.Read(randomPass)
	placeholderHash, err := crypto.HashPassword(string(randomPass))
	if err != nil {
		return "", err
	}

	newUser := &domain.User{
		TenantID:     tenantID,
		Email:        req.Email,
		PasswordHash: placeholderHash,
		Status:       domain.UserStatusInvited,
	}

	if err := s.userRepo.CreateInvitedUser(ctx, newUser); err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	inviteToken := generateRandomToken(32)

	// Expires in 48 hours
	err = s.tokenRepo.SetInviteToken(ctx, inviteToken, newUser.ID, 48*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to save invite token: %w", err)
	}

	slog.Info("User Invited Successfully", "email", req.Email, "id", newUser.ID)
	return inviteToken, nil
}

func (s *UserService) AcceptInvite(ctx context.Context, token string, newPassword string) error {
	userID, err := s.tokenRepo.GetUserIDByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("redis error: %w", err)
	}
	if userID == "" {
		return errors.New("invalid or expired invite token")
	}

	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.ActivateUser(ctx, userID, hashedPassword); err != nil {
		return fmt.Errorf("failed to activate user: %w", err)
	}

	_ = s.tokenRepo.DeleteToken(ctx, token)

	slog.Info("User Activated Successfully", "user_id", userID)
	return nil
}

func generateRandomToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}
