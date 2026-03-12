package oauth_sdk

import oauth "github.com/ory/hydra-client-go/v2"

func NewOauthAdminClient(adminURL string) *oauth.APIClient {
	adminCfg := oauth.NewConfiguration()
	adminCfg.Servers = []oauth.ServerConfiguration{{URL: adminURL}}

	if adminURL != "" {
		adminCfg.Servers = []oauth.ServerConfiguration{
			{URL: adminURL},
		}
	}

	return oauth.NewAPIClient(adminCfg)
}
