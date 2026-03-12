package oauth

import (
	"errors"
	"os"
)

type URLConfig struct {
	AdminURL   string
	PublicURL  string
	PrivateURL string
}

func LoadURLConfig() (*URLConfig, error) {
	adminURL := os.Getenv("OAUTH_ADMIN")
	publicURL := os.Getenv("OAUTH_PUBLIC")
	privateURL := os.Getenv("OAUTH_PRIVATE")

	// Fail fast if any critical value is missing
	if adminURL == "" || publicURL == "" || privateURL == "" {
		return nil, errors.New("OAuth URL configuration missing. Set OAUTH_ADMIN, OAUTH_PUBLIC, and OAUTH_PRIVATE")
	}
	return &URLConfig{
		AdminURL:   os.Getenv("OAUTH_ADMIN"),
		PublicURL:  os.Getenv("OAUTH_PUBLIC"),
		PrivateURL: os.Getenv("OAUTH_PRIVATE"),
	}, nil
}
