package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

const (
	PrefixSession = "skl_sess_"
	PrefixInvite  = "skl_inv_"
)

func GenerateOpaqueToken(prefix string) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	randomString := base64.RawURLEncoding.EncodeToString(bytes)

	return fmt.Sprintf("%s%s", prefix, randomString), nil
}

func HashToken(rawToken string) string {
	hash := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(hash[:])
}
