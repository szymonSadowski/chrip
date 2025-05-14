package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// If there's an error reading random bytes, wrap it and return.
		return "", fmt.Errorf("failed to generate random bytes for refresh token: %w", err)
	}
	refreshToken := hex.EncodeToString(randomBytes)

	return refreshToken, nil
}
