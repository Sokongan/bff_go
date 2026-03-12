package services_config

import (
	"errors"
	"os"
)

type PermissionConfig struct {
	AdminURL  string
	PublicURL string
}

func LoadPermissionConfig() (*PermissionConfig, error) {

	adminURL := os.Getenv("PERMISSION_ADMIN")
	publicURL := os.Getenv("PERMISSION_PUBLIC")

	// Fail fast if any critical value is missing
	if adminURL == "" || publicURL == "" {
		return nil, errors.New("BFF client configuration missing. Set PERMISSION_ADMIN and PERMISSION_PUBLIC")
	}
	return &PermissionConfig{
		AdminURL:  adminURL,
		PublicURL: publicURL,
	}, nil
}
