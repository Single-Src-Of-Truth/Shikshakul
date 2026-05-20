package repository

import (
	"context"

	"github.com/Single-Src-Of-Truth/Shikshakul/eng/iam-service/internal/domain"
	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.UserSession) error
	RevokeSessionByTokenHash(ctx context.Context, tokenHash string) error
	RevokeAllUserSessions(ctx context.Context, userID string) error
	GetActiveSessions(ctx context.Context, userID string) ([]domain.UserSession, error)
	GetSessionsByRoleID(ctx context.Context, roleID string) ([]domain.UserSession, error)
}

type sessionRepo struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepo{db: db}
}

func (r *sessionRepo) CreateSession(ctx context.Context, session *domain.UserSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepo) RevokeSessionByTokenHash(ctx context.Context, tokenHash string) error {
	return r.db.WithContext(ctx).
		Model(&domain.UserSession{}).
		Where("opaque_token_hash = ?", tokenHash).
		Update("is_active", false).
		Delete(&domain.UserSession{}).Error
}

func (r *sessionRepo) RevokeAllUserSessions(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Model(&domain.UserSession{}).
		Where("user_id = ?", userID).
		Delete(&domain.UserSession{}).Error
}

func (r *sessionRepo) GetActiveSessions(ctx context.Context, userID string) ([]domain.UserSession, error) {
	var sessions []domain.UserSession
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ? AND expires_at > now()", userID, true).
		Order("last_active_at desc").
		Find(&sessions).Error
	return sessions, err
}

func (r *sessionRepo) GetSessionsByRoleID(ctx context.Context, roleID string) ([]domain.UserSession, error) {
	var sessions []domain.UserSession
	err := r.db.WithContext(ctx).
		Joins("JOIN users ON users.id = user_sessions.user_id").
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ? AND user_sessions.is_active = ?", roleID, true).
		Find(&sessions).Error
	return sessions, err
}
