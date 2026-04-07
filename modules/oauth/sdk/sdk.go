package oauth_sdk

import (
	oauth_types "sso-bff/modules/oauth"

	oauth "github.com/ory/hydra-client-go/v2"
	"golang.org/x/oauth2"
)

type AuthorizationSDK struct {
	Browser  *oauth2.Config
	Internal *oauth2.Config
}

type OAuthSDK struct {
	Admin         *oauth.APIClient
	Authorization *AuthorizationSDK
}

func NewOAuthSDK(
	adminURL string,
	browser oauth_types.BrowserClient,
	internal oauth_types.InternalClient,
) *OAuthSDK {

	auth := &AuthorizationSDK{
		Browser:  NewOAuthBrowserClient(browser),
		Internal: NewOAuthTokenClient(internal),
	}

	return &OAuthSDK{
		Admin:         NewOauthAdminClient(adminURL),
		Authorization: auth,
	}
}
