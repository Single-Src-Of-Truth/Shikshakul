package service

import (
	"context"
	"errors"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/infrastructure/cache"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/repository"
	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/pkg/dto"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RBACService interface {
	CreateCustomRole(ctx context.Context, tenantID *string, name, description string, permissionIDs []string) (*domain.Role, error)
	FetchUserPermissionSet(ctx context.Context, userID string) ([]string, error)
	GetAllPermissions(ctx context.Context) ([]domain.Permission, error)
	BulkSeedPermissions(ctx context.Context, req dto.BulkCreatePermissionsRequest) error
	UpdateCustomRole(ctx context.Context, roleID string, req dto.UpdateRoleRequest) error
	DeleteCustomRole(ctx context.Context, roleID string) error
	ListRoles(ctx context.Context, tenantID *string) ([]dto.RoleResponse, error)
}

type rbacService struct {
	repo        repository.RBACRepository
	sessionRepo repository.SessionRepository
	redisStore  *cache.SessionStore
	logger      *zap.Logger
}

func NewRBACService(repo repository.RBACRepository, sessionRepo repository.SessionRepository, redisStore *cache.SessionStore, logger *zap.Logger) RBACService {
	return &rbacService{repo: repo, sessionRepo: sessionRepo, redisStore: redisStore, logger: logger}
}

func (s *rbacService) CreateCustomRole(ctx context.Context, tenantID *string, name, description string, permissionIDs []string) (*domain.Role, error) {
	var tID *uuid.UUID
	if tenantID != nil {
		parsed, err := uuid.Parse(*tenantID)
		if err != nil {
			return nil, errors.New("invalid tenant ID")
		}
		tID = &parsed
	}

	role := &domain.Role{
		TenantID:    tID,
		Name:        name,
		Description: description,
		IsSystem:    false,
	}

	if err := s.repo.CreateRole(ctx, role, permissionIDs); err != nil {
		s.logger.Error("Failed to create custom role", zap.Error(err))
		return nil, errors.New("could not create role")
	}

	return role, nil
}

func (s *rbacService) FetchUserPermissionSet(ctx context.Context, userID string) ([]string, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}
	return s.repo.GetUserPermissions(ctx, uid)
}

func (s *rbacService) GetAllPermissions(ctx context.Context) ([]domain.Permission, error) {
	return s.repo.GetAllPermissions(ctx)
}

func (s *rbacService) BulkSeedPermissions(ctx context.Context, req dto.BulkCreatePermissionsRequest) error {
	var perms []domain.Permission
	for _, p := range req.Permissions {
		perms = append(perms, domain.Permission{ID: p.ID, Description: p.Description})
	}
	return s.repo.BulkCreatePermissions(ctx, perms)
}

func (s *rbacService) UpdateCustomRole(ctx context.Context, roleID string, req dto.UpdateRoleRequest) error {
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	err := s.repo.UpdateRole(ctx, roleID, updates, req.PermissionIDs)
	if err != nil {
		return err
	}

	if req.PermissionIDs != nil {
		sessions, _ := s.sessionRepo.GetSessionsByRoleID(ctx, roleID)
		for _, sess := range sessions {
			s.sessionRepo.RevokeSessionByTokenHash(ctx, sess.OpaqueTokenHash)
			s.redisStore.RevokeSession(ctx, sess.OpaqueTokenHash)
		}
		s.logger.Info("Role updated. Invalidated sessions for immediate propagation.", zap.Int("sessions_killed", len(sessions)))
	}

	return nil
}

func (s *rbacService) DeleteCustomRole(ctx context.Context, roleID string) error {
	return s.repo.DeleteRoleSafely(ctx, roleID)
}

func (s *rbacService) ListRoles(ctx context.Context, tenantID *string) ([]dto.RoleResponse, error) {
	roles, err := s.repo.GetRolesByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var response []dto.RoleResponse
	for _, r := range roles {
		var perms []string
		for _, p := range r.Permissions {
			perms = append(perms, p.ID)
		}
		response = append(response, dto.RoleResponse{
			ID:          r.ID.String(),
			Name:        r.Name,
			Description: r.Description,
			IsSystem:    r.IsSystem,
			Permissions: perms,
		})
	}
	return response, nil
}
