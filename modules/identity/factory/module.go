package identity_factory

import (
	"fmt"

	"sso-bff/internal/httpx"
	"sso-bff/modules/audit"
	identity_factory_modules "sso-bff/modules/identity/factory/modules"
	identity_handler "sso-bff/modules/identity/handler"
	identity_sdk "sso-bff/modules/identity/sdk"
	"sso-bff/modules/oauth"
	oauth_gateway "sso-bff/modules/oauth/gateway"
	"sso-bff/modules/permission"
)

type Module struct {
	Services *identity_factory_modules.IdentityServices
	Login    *identity_handler.IdentityLoginHandler
	Admin    *identity_handler.IdentityAdminHandler
	Settings *identity_handler.IdentitySettingsHandler
}

func NewModule(
	sdk *identity_sdk.IdentitySDK,
	oauthAdmin *oauth_gateway.OauthAdminGateway,
	sessions oauth.SessionResolver,
	perm permission.PermissionChecker,
	auditWriter audit.AuditWriter,
	cookies httpx.CookieConfig,
) (*Module, error) {

	if sdk == nil {
		return nil, fmt.Errorf("identity sdk is nil")
	}

	gws, err := identity_factory_modules.NewIdentityGateways(sdk)
	if err != nil {
		return nil, err
	}

	if oauthAdmin == nil {
		return nil, fmt.Errorf("oauth admin gateway is nil")
	}

	services := identity_factory_modules.NewIdentityServices(
		gws.Admin,
		gws.Browser,
		oauthAdmin,
	)

	loginHandler := identity_handler.NewIdentityFlowHandler(
		services,
		auditWriter,
		cookies,
	)

	adminHandler := identity_handler.NewIdentityAdminHandler(
		services,
		sessions,
		perm,
		cookies,
	)

	settingsHandler := identity_handler.NewIdentitySettingsHandler(
		services.Settings,
		sessions,
		cookies,
	)

	return &Module{
		Services: services,
		Login:    loginHandler,
		Admin:    adminHandler,
		Settings: settingsHandler,
	}, nil
}
