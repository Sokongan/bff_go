package oauth_sdk

import (
	oauth_types "sso-bff/internal/modules/oauth"

	"golang.org/x/oauth2"
)

func NewOAuthTokenClient(cfg oauth_types.InternalClient) *oauth2.Config {

	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       cfg.Scopes,
		Endpoint: oauth2.Endpoint{
			TokenURL: cfg.TokenURL + "/oauth2/token",
		},
	}
}
