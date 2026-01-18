package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
)

// GenerateKeyPair creates a new Ed25519 pair.
// Returns privateKey (hex), publicKey (hex), error
func GenerateKeyPair() (string, string, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	return hex.EncodeToString(priv), hex.EncodeToString(pub), nil
}
