package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// GenerateRandomID generates a 24-character hex string (12 random bytes)
func GenerateRandomID() (string, error) {
	bytes := make([]byte, 12) // 12 bytes = 24 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateRandomShortStr() (string, error) {
	bytes := make([]byte, 2) // 12 bytes = 24 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}
