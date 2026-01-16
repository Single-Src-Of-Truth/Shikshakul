package service

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Modulix-IT/Shikshakul-WL/iam-service/pkg/crypto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KeyManager struct {
	db        *pgxpool.Pool
	masterKey string
}

type JSONWebKey struct {
	Kty string `json:"kty"` // Key Type (OKP)
	Crv string `json:"crv"` // Curve (Ed25519)
	Kid string `json:"kid"` // Key ID
	X   string `json:"x"`   // The Public Key (Base64URL encoded)
	Alg string `json:"alg"` // Algorithm (EdDSA)
	Use string `json:"use"` // Usage (sig)
}

type JWKS struct {
	Keys []JSONWebKey `json:"keys"`
}

func NewKeyManager(db *pgxpool.Pool, masterKey string) *KeyManager {
	if len(masterKey) != 64 {
		panic("MASTER_KEY must be exactly 64 hex characters (32 bytes)")
	}

	return &KeyManager{
		db:        db,
		masterKey: masterKey,
	}
}

// EnsureActiveKey checks if we have a valid key. If not, it creates one.
func (km *KeyManager) EnsureActiveKey(ctx context.Context) error {
	var exists bool
	err := km.db.QueryRow(ctx, "SELECT exists(SELECT 1 FROM jwk_keys WHERE status = 'ACTIVE')").Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	slog.Info("No active key found. Initializing new signing key...")
	return km.rotateKey(ctx)
}

func (km *KeyManager) RotateKey(ctx context.Context) error {
	return km.rotateKey(ctx)
}

func (km *KeyManager) rotateKey(ctx context.Context) error {
	privHex, pubHex, err := crypto.GenerateKeyPair()
	if err != nil {
		return err
	}

	encryptedPriv, err := crypto.Encrypt(privHex, km.masterKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %w", err)
	}

	kid := fmt.Sprintf("key_%d", time.Now().Unix())

	tx, err := km.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		UPDATE jwk_keys
		SET status = 'PASSIVE', updated_at = NOW()
		WHERE status = 'ACTIVE'
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO jwk_keys (kid, public_key, encrypted_private_key, algorithm, status, expires_at)
		VALUES ($1, $2, $3, 'EdDSA', 'ACTIVE', $4)
	`, kid, pubHex, encryptedPriv, time.Now().Add(24*time.Hour))
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	slog.Info("Rotated Signing Key", "kid", kid)
	return nil
}

func (km *KeyManager) GetSigningKey(ctx context.Context) (string, string, error) {
	var kid, encryptedPriv string

	err := km.db.QueryRow(ctx, `
		SELECT kid, encrypted_private_key
		FROM jwk_keys
		WHERE status = 'ACTIVE'
		LIMIT 1
	`).Scan(&kid, &encryptedPriv)

	if err == pgx.ErrNoRows {
		return "", "", errors.New("no active signing key found")
	} else if err != nil {
		return "", "", err
	}

	privKey, err := crypto.Decrypt(encryptedPriv, km.masterKey)
	if err != nil {
		slog.Error("Failed to decrypt signing key", "error", err)
		return "", "", errors.New("internal security error")
	}

	return kid, privKey, nil
}

// GetPublicJWKS fetches all ACTIVE public keys and formats them for the world to see
func (km *KeyManager) GetPublicJWKS(ctx context.Context) (*JWKS, error) {
	rows, err := km.db.Query(ctx, "SELECT kid, public_key FROM jwk_keys WHERE status = 'ACTIVE'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []JSONWebKey
	for rows.Next() {
		var kid, pubKeyHex string
		if err := rows.Scan(&kid, &pubKeyHex); err != nil {
			return nil, err
		}

		pubBytes, err := hex.DecodeString(pubKeyHex)
		if err != nil {
			slog.Error("Invalid public key in DB", "kid", kid)
			continue
		}

		pubBase64 := base64.RawURLEncoding.EncodeToString(pubBytes)

		keys = append(keys, JSONWebKey{
			Kty: "OKP",
			Crv: "Ed25519",
			Kid: kid,
			X:   pubBase64,
			Alg: "EdDSA",
			Use: "sig",
		})
	}

	return &JWKS{Keys: keys}, nil
}
