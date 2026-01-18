package postgres

import (
	"context"
	"errors"

	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// GetByEmail checks if a user exists in a specific tenant (ignoring soft deleted ones)
func (r *UserRepository) GetByEmail(ctx context.Context, tenantID, email string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, email, status, password_hash, created_at
		FROM users
		WHERE tenant_id = $1 AND email = $2 AND deleted_at IS NULL
	`
	row := r.db.QueryRow(ctx, query, tenantID, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.TenantID, &user.Email, &user.Status, &user.PasswordHash, &user.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateInvitedUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (tenant_id, email, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, 'INVITED', NOW(), NOW())
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		user.TenantID,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt)

	return err
}

// ActivateUser updates the password and sets status to ACTIVE
func (r *UserRepository) ActivateUser(ctx context.Context, userID, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, status = 'ACTIVE', updated_at = NOW()
		WHERE id = $2
	`

	cmd, err := r.db.Exec(ctx, query, passwordHash, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `
		SELECT id, tenant_id, email, status, password_hash, created_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.Status, &user.PasswordHash, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
