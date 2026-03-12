package oauth_sdk

import (
	oauth_types "sso-bff/internal/modules/oauth"

	"golang.org/x/oauth2"
)

func NewOAuthBrowserClient(browser oauth_types.BrowserClient) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     browser.ClientID,
		ClientSecret: browser.ClientSecret,
		RedirectURL:  browser.RedirectURL,
		Scopes:       browser.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL: browser.BrowserPublicURL + "/oauth2/auth",
		},
	}
}
