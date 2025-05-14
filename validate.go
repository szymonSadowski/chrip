package main

import (
	"net/http"
	"regexp"
	"strings"
)

func (cfg *apiConfig) ValidateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", http.ErrBodyNotAllowed
	}
	cleaned, err := cleanBody(body)
	if err != nil {
		return "", err
	}
	return cleaned, nil
}

func cleanBody(body string) (string, error) {
	forbiddenWords := []string{"kerfuffle", "sharbert", "fornax"}

	pattern := "(?i)(" + strings.Join(forbiddenWords, "|") + ")"
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	cleanedBody := re.ReplaceAllString(body, "****")

	return cleanedBody, nil
}
