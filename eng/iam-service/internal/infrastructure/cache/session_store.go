package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	SlidingWindowTTL = 36 * time.Hour
	AbsoluteMaxTTL   = 15 * 24 * time.Hour
)

type SessionPayload struct {
	UserID      string   `json:"user_id"`
	TenantID    string   `json:"tenant_id"`
	DeviceHash  string   `json:"device_hash"`
	Permissions []string `json:"permissions"`
	CreatedAt   int64    `json:"created_at"`
}

type SessionStore struct {
	client *redis.Client
}

func NewSessionStore(redisURL string) (*SessionStore, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return &SessionStore{client: client}, nil
}

func (s *SessionStore) SaveSession(ctx context.Context, tokenHash string, payload SessionPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	redisKey := fmt.Sprintf("session:%s", tokenHash)
	return s.client.Set(ctx, redisKey, data, SlidingWindowTTL).Err()
}

func (s *SessionStore) VerifyAndRefresh(ctx context.Context, tokenHash string, currentDeviceHash string) (*SessionPayload, error) {
	redisKey := fmt.Sprintf("session:%s", tokenHash)

	data, err := s.client.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, errors.New("session expired or invalid")
	} else if err != nil {
		return nil, err
	}

	var payload SessionPayload
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return nil, err
	}

	if payload.DeviceHash != currentDeviceHash {
		s.RevokeSession(ctx, tokenHash)
		return nil, errors.New("security breach: device fingerprint mismatch")
	}

	createdAt := time.Unix(payload.CreatedAt, 0)
	if time.Since(createdAt) > AbsoluteMaxTTL {
		s.RevokeSession(ctx, tokenHash)
		return nil, errors.New("session reached absolute maximum lifetime (15 days). Please log in again")
	}

	ttl, _ := s.client.TTL(ctx, redisKey).Result()
	if ttl < (SlidingWindowTTL - 1*time.Hour) {
		s.client.Expire(ctx, redisKey, SlidingWindowTTL)
	}

	return &payload, nil
}

func (s *SessionStore) RevokeSession(ctx context.Context, tokenHash string) error {
	redisKey := fmt.Sprintf("session:%s", tokenHash)
	return s.client.Del(ctx, redisKey).Err()
}

func (s *SessionStore) SaveOTP(ctx context.Context, identifier string, otp string) error {
	key := fmt.Sprintf("pwd_reset:%s", identifier)
	return s.client.Set(ctx, key, otp, 10*time.Minute).Err()
}

func (s *SessionStore) VerifyAndBurnOTP(ctx context.Context, identifier string, providedOTP string) error {
	key := fmt.Sprintf("pwd_reset:%s", identifier)
	storedOTP, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil || storedOTP != providedOTP {
		return errors.New("invalid or expired OTP")
	}
	s.client.Del(ctx, key)
	return nil
}

func (s *SessionStore) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

func (s *SessionStore) SuspendTenant(ctx context.Context, tenantID string) error {
	key := fmt.Sprintf("tenant_suspended:%s", tenantID)
	return s.client.Set(ctx, key, "true", 0).Err()
}

func (s *SessionStore) UnsuspendTenant(ctx context.Context, tenantID string) error {
	key := fmt.Sprintf("tenant_suspended:%s", tenantID)
	return s.client.Del(ctx, key).Err()
}

func (s *SessionStore) IsTenantSuspended(ctx context.Context, tenantID string) bool {
	key := fmt.Sprintf("tenant_suspended:%s", tenantID)
	val, err := s.client.Get(ctx, key).Result()
	return err == nil && val == "true"
}
