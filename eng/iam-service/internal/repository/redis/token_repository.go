package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) *TokenRepository {
	return &TokenRepository{client: client}
}

// SetInviteToken stores: key="invite:{token}" -> value="{userID}"
func (r *TokenRepository) SetInviteToken(ctx context.Context, token string, userID string, duration time.Duration) error {
	key := fmt.Sprintf("invite:%s", token)
	return r.client.Set(ctx, key, userID, duration).Err()
}

// GetUserIDByToken checks if token exists and returns the associated UserID
func (r *TokenRepository) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("invite:%s", token)
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return result, nil
}

// DeleteToken removes the token so it cannot be re-used
func (r *TokenRepository) DeleteToken(ctx context.Context, token string) error {
	key := fmt.Sprintf("invite:%s", token)
	return r.client.Del(ctx, key).Err()
}
