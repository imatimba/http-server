package auth

import (
	"fmt"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	apiKeyPrefix := "ApiKey "

	if len(apiKey) <= len(apiKeyPrefix) || apiKey[:len(apiKeyPrefix)] != apiKeyPrefix {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return apiKey[len(apiKeyPrefix):], nil
}

func GetAuthToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return authHeader[len(bearerPrefix):], nil
}
