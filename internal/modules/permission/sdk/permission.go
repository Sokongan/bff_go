package sdk

import (
	client "github.com/ory/keto-client-go/v25"
)

type PermissionSDK struct {
	Admin  *client.APIClient
	Public *client.APIClient
}

func PermissionSDK(adminURL, publicURL string) *PermissionSDK {
	publicCfg := client.NewConfiguration()
	if publicURL != "" {
		publicCfg.Servers = []client.ServerConfiguration{{URL: publicURL}}
	}

	adminCfg := client.NewConfiguration()
	if adminURL != "" {
		adminCfg.Servers = []client.ServerConfiguration{{URL: adminURL}}
	}
	return &PermissionSDK{
		Admin:  client.NewAPIClient(adminCfg),
		Public: client.NewAPIClient(publicCfg),
	}
}
