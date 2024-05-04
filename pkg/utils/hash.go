// Package utils consists some useful methods
package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Sha256Hash method returns the Sha256 hash of provided data
func Sha256Hash(data []byte, key string) (string, error) {
	keyBytes := []byte(key)
	h := hmac.New(sha256.New, keyBytes)
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	hash := h.Sum(nil)
	return hex.EncodeToString(hash), nil
}
