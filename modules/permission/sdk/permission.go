package permission_sdk

import (
	client "github.com/ory/keto-client-go/v25"
)

type PermissionSDK struct {
	Admin  *client.APIClient
	Public *client.APIClient
}

func NewPermissionSDK(adminURL, publicURL string) *PermissionSDK {
	publicCfg := client.NewConfiguration()
	publicCfg.Servers = []client.ServerConfiguration{{URL: publicURL}}

	adminCfg := client.NewConfiguration()
	adminCfg.Servers = []client.ServerConfiguration{{URL: adminURL}}

	return &PermissionSDK{
		Admin:  client.NewAPIClient(adminCfg),
		Public: client.NewAPIClient(publicCfg),
	}
}
