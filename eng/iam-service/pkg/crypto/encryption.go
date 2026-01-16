package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// Encrypt encrypts plainText using the 32-byte masterKey.
// Returns hex-encoded string: nonce + ciphertext
func Encrypt(plainText, masterKeyHex string) (string, error) {
	key, err := hex.DecodeString(masterKeyHex)
	if err != nil {
		return "", errors.New("invalid master key hex")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the hex-encoded cipherText using the masterKey.
func Decrypt(cipherTextHex, masterKeyHex string) (string, error) {
	key, err := hex.DecodeString(masterKeyHex)
	if err != nil {
		return "", errors.New("invalid master key hex")
	}

	data, err := hex.DecodeString(cipherTextHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("decryption failed: incorrect key or corrupted data")
	}

	return string(plainText), nil
}
