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
type M2MConfig struct {
	M2MID     string
	M2MSecret string
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

func LoadM2MConfig() (*M2MConfig, error) {
	m2mID := os.Getenv("M2M_CLIENT_ID")
	m2mSecret := os.Getenv("M2M_CLIENT_SECRET")

	// Fail fast if any critical value is missing
	if m2mID == "" || m2mSecret == "" {
		return nil, errors.New("M2M configuration missing. Set M2M_CLIENT_ID and M2M_CLIENT_SECRET")
	}

	return &M2MConfig{
		M2MID:     m2mID,
		M2MSecret: m2mSecret,
	}, nil
}
