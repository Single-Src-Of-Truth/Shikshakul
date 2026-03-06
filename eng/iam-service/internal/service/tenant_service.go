package service

import (
	"context"
	"errors"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/cache"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/repository"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/dto"
	"go.uber.org/zap"
)

type TenantService interface {
	CreateTenant(ctx context.Context, req dto.CreateTenantRequest) (*dto.TenantResponse, error)
	ListTenants(ctx context.Context) ([]dto.TenantResponse, error)
	GetTenant(ctx context.Context, id string) (*dto.TenantResponse, error)
	UpdateTenant(ctx context.Context, id string, req dto.UpdateTenantRequest) error
	UpdateUserStatus(ctx context.Context, userID string, status string) error
	DeleteUser(ctx context.Context, userID string) error
}

type tenantService struct {
	repo        repository.TenantRepository
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	redisStore  *cache.SessionStore
	logger      *zap.Logger
}

func NewTenantService(repo repository.TenantRepository, userRepo repository.UserRepository, sessionRepo repository.SessionRepository, redisStore *cache.SessionStore, logger *zap.Logger) TenantService {
	return &tenantService{
		repo:        repo,
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		redisStore:  redisStore,
		logger:      logger,
	}
}

func (s *tenantService) CreateTenant(ctx context.Context, req dto.CreateTenantRequest) (*dto.TenantResponse, error) {
	authConfig := req.AuthConfig
	if authConfig == "" {
		authConfig = "{}"
	}

	tenant := &domain.Tenant{
		Name:       req.Name,
		Domain:     req.Domain,
		AuthConfig: authConfig,
		IsActive:   true,
	}

	if err := s.repo.CreateTenant(ctx, tenant); err != nil {
		s.logger.Error("Failed to create tenant", zap.Error(err))
		return nil, err
	}

	return &dto.TenantResponse{
		ID:       tenant.ID.String(),
		Name:     tenant.Name,
		Domain:   tenant.Domain,
		IsActive: tenant.IsActive,
	}, nil
}

func (s *tenantService) ListTenants(ctx context.Context) ([]dto.TenantResponse, error) {
	tenants, err := s.repo.GetAllTenants(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.TenantResponse
	for _, t := range tenants {
		response = append(response, dto.TenantResponse{
			ID:       t.ID.String(),
			Name:     t.Name,
			Domain:   t.Domain,
			IsActive: t.IsActive,
		})
	}
	return response, nil
}

func (s *tenantService) GetTenant(ctx context.Context, id string) (*dto.TenantResponse, error) {
	t, err := s.repo.GetTenantByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.TenantResponse{ID: t.ID.String(), Name: t.Name, Domain: t.Domain, IsActive: t.IsActive}, nil
}

func (s *tenantService) UpdateTenant(ctx context.Context, id string, req dto.UpdateTenantRequest) error {
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Domain != nil {
		updates["domain"] = *req.Domain
	}
	if req.AuthConfig != nil {
		updates["auth_config"] = *req.AuthConfig
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive

		if !*req.IsActive {
			s.redisStore.SuspendTenant(ctx, id)
			s.logger.Warn("Tenant suspended system-wide", zap.String("tenant_id", id))
		} else {
			s.redisStore.UnsuspendTenant(ctx, id)
			s.logger.Info("Tenant unsuspended", zap.String("tenant_id", id))
		}
	}

	if len(updates) == 0 {
		return nil
	}
	return s.repo.UpdateTenant(ctx, id, updates)
}

func (s *tenantService) UpdateUserStatus(ctx context.Context, userID string, status string) error {
	userStatus := domain.UserStatus(status)

	if err := s.userRepo.UpdateStatus(ctx, userID, userStatus); err != nil {
		s.logger.Error("Failed to update user status", zap.Error(err))
		return errors.New("failed to update user status")
	}

	if userStatus == domain.StatusSuspended {
		if err := s.sessionRepo.RevokeAllUserSessions(ctx, userID); err != nil {
			s.logger.Error("Failed to revoke sessions for suspended user", zap.Error(err))
		}
		s.logger.Warn("User suspended system-wide", zap.String("user_id", userID))
	}

	return nil
}

func (s *tenantService) DeleteUser(ctx context.Context, userID string) error {
	if err := s.userRepo.SoftDelete(ctx, userID); err != nil {
		s.logger.Error("Failed to soft delete user", zap.Error(err))
		return errors.New("failed to delete user account")
	}

	if err := s.sessionRepo.RevokeAllUserSessions(ctx, userID); err != nil {
		s.logger.Error("Failed to revoke sessions for deleted user", zap.Error(err))
	}

	s.logger.Info("User account deleted", zap.String("user_id", userID))
	return nil
}
