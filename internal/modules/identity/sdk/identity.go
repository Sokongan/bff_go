package identity_sdk

import (
	client "github.com/ory/keto-client-go/v25"
)

type IdentitySDK struct {
	Public *client.APIClient
	Admin  *client.APIClient
}

func NewIdentitySDK(publicURL, adminURL string) *IdentitySDK {
	publicCfg := client.NewConfiguration()
	publicCfg.Servers = []client.ServerConfiguration{{URL: publicURL}}

	adminCfg := client.NewConfiguration()
	adminCfg.Servers = []client.ServerConfiguration{{URL: adminURL}}

	return &IdentitySDK{
		Public: client.NewAPIClient(publicCfg),
		Admin:  client.NewAPIClient(adminCfg),
	}
}
