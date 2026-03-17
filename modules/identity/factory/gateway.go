package identity_factory

import (
	"errors"
	identity_gateway "sso-bff/modules/identity/gateway"
	identity_sdk "sso-bff/modules/identity/sdk"
)

type IdentityGateWay struct {
	Admin   *identity_gateway.IdentityAdminGateway
	Browser *identity_gateway.IdentityBrowserGateway
}

func NewIdentityGateways(sdk *identity_sdk.IdentitySDK) (*IdentityGateWay,
	error,
) {
	admin := identity_gateway.NewIdentityAdminGateway(sdk.Admin)
	if admin == nil {
		return nil, errors.New("failed to create identity admin gateway")
	}

	browser := identity_gateway.NewIdentityBrowserGateway(sdk.Public)
	if browser == nil {
		return nil, errors.New("failed to create identity admin gateway")
	}

	return &IdentityGateWay{
		Admin:   admin,
		Browser: browser,
	}, nil
}
