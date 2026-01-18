package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

// CreateSession stores the hash of the refresh token
func (r *SessionRepository) CreateSession(ctx context.Context, userID, refreshToken string, userAgent string, ip string) error {
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	query := `
		INSERT INTO sessions (user_id, refresh_token_hash, client_ip, user_agent, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err := r.db.Exec(ctx, query, userID, tokenHash, ip, userAgent, expiresAt)
	return err
}

// GetSessionByHash finds the session looking for the hashed token
// We return the UserID and ExpiresAt so we can validate it
func (r *SessionRepository) GetSessionByHash(ctx context.Context, tokenHash string) (string, time.Time, error) {
	var userID string
	var expiresAt time.Time

	query := `
		SELECT user_id, expires_at
		FROM sessions
		WHERE refresh_token_hash = $1
	`
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(&userID, &expiresAt)
	if err != nil {
		return "", time.Time{}, err
	}

	return userID, expiresAt, nil
}

// RotateSessionToken updates the old hash with a new one and extends expiry
func (r *SessionRepository) RotateSessionToken(ctx context.Context, oldHash, newHash string, newExpiry time.Time) error {
	query := `
		UPDATE sessions
		SET refresh_token_hash = $1, expires_at = $2, updated_at = NOW()
		WHERE refresh_token_hash = $3
	`
	cmd, err := r.db.Exec(ctx, query, newHash, newExpiry, oldHash)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("session not found or already rotated")
	}

	return nil
}

// RevokeSession deletes the session associated with the token hash
func (r *SessionRepository) RevokeSession(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM sessions WHERE refresh_token_hash = $1`
	_, err := r.db.Exec(ctx, query, tokenHash)
	return err
}
