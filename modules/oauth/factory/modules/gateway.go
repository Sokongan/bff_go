package oauth_factory_modules

import (
	"errors"
	oauth_gateway "sso-bff/modules/oauth/gateway"
	oauth_sdk "sso-bff/modules/oauth/sdk"
)

type OauthGateways struct {
	Authorization *oauth_gateway.OAuthAuthorizationGateway
	Admin         *oauth_gateway.OauthAdminGateway
	M2M           *oauth_gateway.OAuthM2MGateway
}

func NewOauthGateways(sdk *oauth_sdk.OAuthSDK) (*OauthGateways, error) {
	authGW, _ := oauth_gateway.NewOAuthAuthorizationGateway(
		sdk.Authorization.Browser,
		sdk.Authorization.Internal,
	)

	if authGW == nil {
		return nil, errors.New("failed to create authorization gateway")
	}
	adminGW := oauth_gateway.NewOauthAdminGateway(sdk.Admin)

	if adminGW == nil {
		return nil, errors.New("failed to create admin gateway")
	}

	m2mGW := oauth_gateway.NewOAuthM2MGateway(sdk.M2M)

	if m2mGW == nil {
		return nil, errors.New("failed to create M2M gateway")
	}

	return &OauthGateways{
		Authorization: authGW,
		Admin:         adminGW,
		M2M:           m2mGW,
	}, nil
}
