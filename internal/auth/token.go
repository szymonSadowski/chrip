package auth

import (
	"fmt"
	"net/http"
)

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	// Check if the header starts with "Bearer "
	if len(authorization) < 7 || authorization[:7] != "Bearer " {
		return "", fmt.Errorf("invalid Authorization header format")
	}
	// Extract the token by removing the "Bearer " prefix
	token := authorization[7:]
	if token == "" {
		return "", fmt.Errorf("missing token in Authorization header")
	}
	// Return the token

	return token, nil
}
