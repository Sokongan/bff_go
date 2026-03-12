package services_config

import (
	"errors"
	"os"
)

type IdentityConfig struct {
	PublicURL string
	AdminURL  string
}

func LoadIdentityConfig() (*IdentityConfig, error) {
	adminURL := os.Getenv("IDENTITY_ADMIN")
	publicURL := os.Getenv("IDENTITY_PUBLIC")

	// Fail fast if any critical value is missing
	if adminURL == "" || publicURL == "" {
		return nil, errors.New("BFF client configuration missing. Set IDENTITY_ADMIN and IDENTITY_PUBLIC")
	}
	return &IdentityConfig{
		PublicURL: os.Getenv("IDENTITY_PUBLIC"),
		AdminURL:  os.Getenv("IDENTITY_ADMIN"),
	}, nil
}
