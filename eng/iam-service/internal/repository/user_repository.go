package repository

import (
	"context"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByTenantAndIdentifier(ctx context.Context, tenantID *string, identifier string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	UpdatePassword(ctx context.Context, id string, newPasswordHash string) error
	UpdateProfile(ctx context.Context, id, firstName, lastName string) error
	SoftDelete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status domain.UserStatus) error
	UpdateUserRole(ctx context.Context, userID string, roleID string) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) FindByTenantAndIdentifier(ctx context.Context, tenantID *string, identifier string) (*domain.User, error) {
	var user domain.User

	query := r.db.WithContext(ctx).
		Preload("Tenant").
		Preload("Roles").
		Preload("Roles.Permissions").
		Where("primary_identifier = ?", identifier)

	if tenantID == nil {
		query = query.Where("tenant_id IS NULL")
	} else {
		query = query.Where("tenant_id = ?", *tenantID)
	}

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Preload("Roles").Preload("Roles.Permissions").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, id string, newPasswordHash string) error {
	return r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Update("password_hash", newPasswordHash).Error
}

func (r *userRepo) UpdateProfile(ctx context.Context, id, firstName, lastName string) error {
	return r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"first_name": firstName,
			"last_name":  lastName,
		}).Error
}

func (r *userRepo) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("status", domain.StatusSuspended).Delete(&domain.User{}).Error
}

func (r *userRepo) UpdateStatus(ctx context.Context, id string, status domain.UserStatus) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *userRepo) UpdateUserRole(ctx context.Context, userID string, roleID string) error {
	var role domain.Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", roleID).Error; err != nil {
		return err
	}

	user := domain.User{Base: domain.Base{ID: uuid.MustParse(userID)}}

	return r.db.WithContext(ctx).Model(&user).Association("Roles").Replace(&role)
}
