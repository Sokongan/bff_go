package oauth_sdk

import (
	oauth_types "sso-bff/modules/oauth"

	"golang.org/x/oauth2/clientcredentials"
)

func NewOAuthM2MClient(cfg oauth_types.M2MClient) *clientcredentials.Config {

	return &clientcredentials.Config{
		ClientID:     cfg.M2MID,
		ClientSecret: cfg.M2MSecret,
		TokenURL:     cfg.TokenURL,
		Scopes:       cfg.Scopes,
	}
}
