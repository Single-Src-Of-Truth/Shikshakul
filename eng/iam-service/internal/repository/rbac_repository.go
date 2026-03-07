package repository

import (
	"context"
	"errors"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RBACRepository interface {
	CreateRole(ctx context.Context, role *domain.Role, permissionIDs []string) error
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetAllPermissions(ctx context.Context) ([]domain.Permission, error)
	BulkCreatePermissions(ctx context.Context, permissions []domain.Permission) error
	UpdateRole(ctx context.Context, roleID string, updates map[string]interface{}, permissionIDs []string) error
	DeleteRoleSafely(ctx context.Context, roleID string) error
	GetRolesByTenant(ctx context.Context, tenantID *string) ([]domain.Role, error)
	GetRoleByID(ctx context.Context, roleID string) (*domain.Role, error)
	CountUsersByRole(ctx context.Context, roleID string) (int64, error)
}

type rbacRepo struct {
	db *gorm.DB
}

func NewRBACRepository(db *gorm.DB) RBACRepository {
	return &rbacRepo{db: db}
}

func (r *rbacRepo) CreateRole(ctx context.Context, role *domain.Role, permissionIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}

		var permissions []domain.Permission
		if err := tx.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}

		if len(permissions) != len(permissionIDs) {
			return errors.New("one or more provided permission_ids are invalid")
		}

		if err := tx.Model(role).Association("Permissions").Append(&permissions); err != nil {
			return err
		}
		return nil
	})
}

func (r *rbacRepo) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var permissionIDs []string

	query := `
		SELECT DISTINCT rp.permission_id
		FROM user_roles ur
		JOIN role_permissions rp ON ur.role_id = rp.role_id
		WHERE ur.user_id = ?
	`
	if err := r.db.WithContext(ctx).Raw(query, userID).Scan(&permissionIDs).Error; err != nil {
		return nil, err
	}

	return permissionIDs, nil
}

func (r *rbacRepo) GetAllPermissions(ctx context.Context) ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.WithContext(ctx).Order("id asc").Find(&permissions).Error
	return permissions, err
}

func (r *rbacRepo) BulkCreatePermissions(ctx context.Context, permissions []domain.Permission) error {
	return r.db.WithContext(ctx).Save(&permissions).Error
}

func (r *rbacRepo) UpdateRole(ctx context.Context, roleID string, updates map[string]interface{}, permissionIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(updates) > 0 {
			if err := tx.Model(&domain.Role{}).Where("id = ?", roleID).Updates(updates).Error; err != nil {
				return err
			}
		}

		if permissionIDs != nil {
			var permissions []domain.Permission
			if err := tx.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
				return err
			}

			if len(permissions) != len(permissionIDs) {
				return errors.New("one or more provided permission_ids are invalid")
			}

			role := domain.Role{Base: domain.Base{ID: uuid.MustParse(roleID)}}
			if err := tx.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *rbacRepo) DeleteRoleSafely(ctx context.Context, roleID string) error {
	var count int64
	r.db.WithContext(ctx).Table("user_roles").
		Joins("JOIN users ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ? AND users.deleted_at IS NULL", roleID).
		Count(&count)

	if count > 0 {
		return errors.New("conflict: cannot delete role while users are assigned to it")
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM user_roles WHERE role_id = ?", roleID).Error; err != nil {
			return err
		}

		role := domain.Role{Base: domain.Base{ID: uuid.MustParse(roleID)}}

		if err := tx.Model(&role).Association("Permissions").Clear(); err != nil {
			return err
		}

		return tx.Unscoped().Delete(&role).Error
	})
}

func (r *rbacRepo) GetRolesByTenant(ctx context.Context, tenantID *string) ([]domain.Role, error) {
	var roles []domain.Role
	query := r.db.WithContext(ctx).Preload("Permissions")

	if tenantID == nil {
		query = query.Where("tenant_id IS NULL AND is_system = ?", true)
	} else {
		query = query.Where("tenant_id = ? OR (tenant_id IS NULL AND is_system = ?)", *tenantID, true)
	}

	err := query.Find(&roles).Error
	return roles, err
}

func (r *rbacRepo) GetRoleByID(ctx context.Context, roleID string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).First(&role, "id = ?", roleID).Error
	return &role, err
}

func (r *rbacRepo) CountUsersByRole(ctx context.Context, roleID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("user_roles").
		Joins("JOIN users ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ? AND users.deleted_at IS NULL", roleID).
		Count(&count).Error
	return count, err
}
