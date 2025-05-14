package auth

import (
	"fmt"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("missing Authorization header")
	}
	// Check if the header starts with "ApiKey "
	const prefix = "ApiKey "
	if len(authorization) < len(prefix) || authorization[:len(prefix)] != prefix {
		return "", fmt.Errorf("invalid Authorization header format")
	}
	token := authorization[len(prefix):]
	if token == "" {
		return "", fmt.Errorf("missing token in Authorization header")
	}
	// Return the token

	return token, nil
}
