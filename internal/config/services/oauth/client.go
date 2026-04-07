package oauth

import (
	"errors"
	"os"
)

type ClientConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func LoadClientConfig() (*ClientConfig, error) {
	clientID := os.Getenv("BFF_CLIENT_ID")
	clientSecret := os.Getenv("BFF_CLIENT_SECRET")
	redirectURL := os.Getenv("BFF_REDIRECT_URL")

	// Fail fast if any critical value is missing
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, errors.New("BFF client configuration missing. Set BFF_CLIENT_ID, BFF_CLIENT_SECRET, and BFF_REDIRECT_URL")
	}

	return &ClientConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}, nil
}
