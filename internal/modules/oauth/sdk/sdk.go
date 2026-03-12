package oauth_sdk

import (
	oauth_types "sso-bff/internal/modules/oauth"

	oauth "github.com/ory/hydra-client-go/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type AuthorizationSDK struct {
	Browser  *oauth2.Config
	Internal *oauth2.Config
}

type OAuthSDK struct {
	Admin         *oauth.APIClient
	Authorization *AuthorizationSDK
	M2M           *clientcredentials.Config
}

func NewOAuthSDK(
	adminURL string,
	browser oauth_types.BrowserClient,
	internal oauth_types.InternalClient,
	m2m oauth_types.M2MClient,
) *OAuthSDK {

	auth := &AuthorizationSDK{
		Browser:  NewOAuthBrowserClient(browser),
		Internal: NewOAuthTokenClient(internal),
	}

	return &OAuthSDK{
		Admin:         NewOauthAdminClient(adminURL),
		Authorization: auth,
		M2M:           NewOAuthM2MClient(m2m),
	}
}
