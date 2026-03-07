package repository

import (
	"context"
	"time"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"gorm.io/gorm"
)

type OnboardingRepository interface {
	CreateInvitation(ctx context.Context, invite *domain.Invitation) error
	FindPendingInvite(ctx context.Context, tokenHash string) (*domain.Invitation, error)
	ConsumeInviteAndCreateUser(ctx context.Context, invite *domain.Invitation, user *domain.User) error
	GetInvitationByID(ctx context.Context, id string) (*domain.Invitation, error)
	UpdateInvitationExpiry(ctx context.Context, id string, newExpiry time.Time) error
	SoftDeleteInvitation(ctx context.Context, id string) error
	CountPendingInvitesByRole(ctx context.Context, roleID string) (int64, error)
}

type onboardingRepo struct {
	db *gorm.DB
}

func NewOnboardingRepository(db *gorm.DB) OnboardingRepository {
	return &onboardingRepo{db: db}
}

func (r *onboardingRepo) CreateInvitation(ctx context.Context, invite *domain.Invitation) error {
	return r.db.WithContext(ctx).Create(invite).Error
}

func (r *onboardingRepo) FindPendingInvite(ctx context.Context, tokenHash string) (*domain.Invitation, error) {
	var invite domain.Invitation
	err := r.db.WithContext(ctx).
		Where("token_hash = ? AND status = ? AND expires_at > ?", tokenHash, domain.InvitePending, time.Now()).
		First(&invite).Error
	return &invite, err
}

func (r *onboardingRepo) ConsumeInviteAndCreateUser(ctx context.Context, invite *domain.Invitation, user *domain.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		var role domain.Role
		if err := tx.First(&role, "id = ?", invite.RoleID).Error; err != nil {
			return err
		}

		if err := tx.Model(user).Association("Roles").Append(&role); err != nil {
			return err
		}

		if err := tx.Model(invite).Update("status", domain.InviteAccepted).Error; err != nil {
			return err
		}

		if err := tx.Delete(invite).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *onboardingRepo) GetInvitationByID(ctx context.Context, id string) (*domain.Invitation, error) {
	var invite domain.Invitation
	if err := r.db.WithContext(ctx).First(&invite, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &invite, nil
}

func (r *onboardingRepo) UpdateInvitationExpiry(ctx context.Context, id string, newExpiry time.Time) error {
	return r.db.WithContext(ctx).
		Model(&domain.Invitation{}).
		Where("id = ?", id).
		Update("expires_at", newExpiry).Error
}

func (r *onboardingRepo) SoftDeleteInvitation(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Invitation{}).Error
}

func (r *onboardingRepo) CountPendingInvitesByRole(ctx context.Context, roleID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Invitation{}).
		Where("role_id = ? AND status = ? AND expires_at > ?", roleID, domain.InvitePending, time.Now()).
		Count(&count).Error
	return count, err
}
