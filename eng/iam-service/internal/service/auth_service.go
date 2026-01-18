package service

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/domain"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/internal/repository/postgres"
	"github.com/Modulix-IT/Shikshakul-WL/iam-service/pkg/crypto"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo    *postgres.UserRepository
	sessionRepo *postgres.SessionRepository
	keyManager  *KeyManager
}

func NewAuthService(u *postgres.UserRepository, s *postgres.SessionRepository, k *KeyManager) *AuthService {
	return &AuthService{userRepo: u, sessionRepo: s, keyManager: k}
}

func (s *AuthService) Login(ctx context.Context, req domain.LoginRequest, userAgent, ip string) (*domain.TokenResponse, error) {
	tenantID, ok := ctx.Value("tenant_id").(string)
	if !ok {
		return nil, errors.New("internal system error")
	}

	user, err := s.userRepo.GetByEmail(ctx, tenantID, req.Username)
	if err != nil || user == nil {
		slog.Warn("Login failed: User not found", "email", req.Username)
		return nil, errors.New("invalid credentials")
	}

	match, err := crypto.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		slog.Error("Login Failed: Crypto Verification Error", "error", err)
		return nil, err
	}
	if !match {
		slog.Warn("Login Failed: Password mismatch", "email", req.Username)
		return nil, errors.New("invalid credentials")
	}

	if user.Status != domain.UserStatusActive {
		slog.Warn("Login Failed: User not active", "status", user.Status)
		return nil, errors.New("account is not active")
	}

	kid, privKeyPEM, err := s.keyManager.GetSigningKey(ctx)
	if err != nil {
		slog.Error("Login Failed: Key Manager Error", "error", err)
		return nil, fmt.Errorf("key manager error: %w", err)
	}

	privKeyBytes, _ := hex.DecodeString(privKeyPEM)
	if len(privKeyBytes) != 64 {
		slog.Error("Login Failed: Invalid Private Key Length", "length", len(privKeyBytes))
		return nil, errors.New("invalid private key length")
	}

	signingKey := ed25519.PrivateKey(privKeyBytes)
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":       user.ID,
		"iss":       "shikshakul-iam",
		"aud":       "school-service",
		"iat":       now.Unix(),
		"exp":       now.Add(15 * time.Minute).Unix(),
		"tenant_id": user.TenantID,
		"role":      "teacher", // Future: Get from DB
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = kid

	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		slog.Error("JWT Signing Failed", "error", err)
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	refreshToken := generateRandomHex(32)
	if err := s.sessionRepo.CreateSession(ctx, user.ID, refreshToken, userAgent, ip); err != nil {
		slog.Error("Login Failed: Session Creation Error", "error", err)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	slog.Info("Login Success", "user_id", user.ID)

	return &domain.TokenResponse{
		AccessToken:  signedToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, oldRefreshToken string) (*domain.TokenResponse, error) {
	hash := sha256.Sum256([]byte(oldRefreshToken))
	oldTokenHash := hex.EncodeToString(hash[:])

	userID, expiresAt, err := s.sessionRepo.GetSessionByHash(ctx, oldTokenHash)
	if err != nil {
		slog.Warn("Refresh Failed: Session not found", "hash_prefix", oldTokenHash[:10])
		return nil, errors.New("invalid or expired refresh token")
	}

	if time.Now().After(expiresAt) {
		slog.Warn("Refresh Failed: Token expired", "user_id", userID)
		return nil, errors.New("invalid or expired refresh token")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("user not found")
	}

	if user.Status != domain.UserStatusActive {
		return nil, errors.New("user account is inactive")
	}

	kid, privKeyPEM, err := s.keyManager.GetSigningKey(ctx)
	if err != nil {
		return nil, err
	}
	privKeyBytes, _ := hex.DecodeString(privKeyPEM)
	signingKey := ed25519.PrivateKey(privKeyBytes)

	now := time.Now()
	claims := jwt.MapClaims{
		"sub":       user.ID,
		"iss":       "shikshakul-iam",
		"aud":       "school-service",
		"iat":       now.Unix(),
		"exp":       now.Add(15 * time.Minute).Unix(),
		"tenant_id": user.TenantID,
		"role":      "teacher", // TODO: Fetch
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = kid

	newAccessToken, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err
	}

	newRefreshToken := generateRandomHex(32)

	newHashRaw := sha256.Sum256([]byte(newRefreshToken))
	newTokenHash := hex.EncodeToString(newHashRaw[:])

	newExpiry := time.Now().Add(7 * 24 * time.Hour)
	if err := s.sessionRepo.RotateSessionToken(ctx, oldTokenHash, newTokenHash, newExpiry); err != nil {
		slog.Error("Refresh Failed: DB Update Error", "err", err)
		return nil, errors.New("failed to rotate token")
	}

	slog.Info("Token Refreshed Successfully", "user_id", user.ID)

	return &domain.TokenResponse{
		AccessToken:  newAccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		RefreshToken: newRefreshToken, // Client must replace the old one!
	}, nil
}

// Logout revokes the specific refresh token
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	if err := s.sessionRepo.RevokeSession(ctx, tokenHash); err != nil {
		slog.Error("Logout Failed: DB Error", "error", err)
		return errors.New("failed to revoke session")
	}

	slog.Info("Session Revoked Successfully")
	return nil
}

func generateRandomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
