package identity_factory

import (
	identity_gateway "sso-bff/modules/identity/gateway"
	identity_service "sso-bff/modules/identity/services"
	oauth_gateway "sso-bff/modules/oauth/gateway"
)

type IdentityServices struct {
	Admin       *identity_service.IdentityAdminService
	Login       *identity_service.IdentityLoginService
	GetIdentity *identity_service.GetIdentityService
	Settings    *identity_service.IdentitySettingsService
}

func NewIdentityServices(
	admin *identity_gateway.IdentityAdminGateway,
	browser *identity_gateway.IdentityBrowserGateway,
	oauthAdmin *oauth_gateway.OauthAdminGateway,
) *IdentityServices {

	adminSvc := identity_service.NewIdentityAdminService(admin)

	loginSvc := identity_service.NewIdentityLoginService(
		browser,
		oauthAdmin,
	)

	getIdentitySvc := identity_service.NewGetIdentityService(
		browser,
	)

	settingsSvc := identity_service.NewIdentitySettingsService(
		browser,
	)

	return &IdentityServices{
		Admin:       adminSvc,
		Login:       loginSvc,
		GetIdentity: getIdentitySvc,
		Settings:    settingsSvc,
	}
}
