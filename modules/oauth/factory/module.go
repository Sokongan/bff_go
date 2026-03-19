package oauth_factory

import (
	"fmt"
	"time"

	"sso-bff/internal/httpx"
	"sso-bff/modules/app"
	"sso-bff/modules/audit"
	"sso-bff/modules/identity"
	"sso-bff/modules/oauth"
	oauth_factory_modules "sso-bff/modules/oauth/factory/modules"
	oauth_handler "sso-bff/modules/oauth/handler"
	oauth_sdk "sso-bff/modules/oauth/sdk"
	"sso-bff/modules/permission"
)

// Module wires sdk -> gateways -> services -> handler for OAuth.
type Module struct {
	Service *oauth_factory_modules.OAuthService
	Handler *oauth_handler.OAuthHandler
}

// NewModule builds the OAuth stack and returns the service/handler bundle.
func NewOauthModule(
	sdk *oauth_sdk.OAuthSDK,
	sessionStore oauth.SessionStorePort,
	pkceStore oauth.PKCEStorePort,
	redirectStore oauth.RedirectStorePort,
	idVerifier oauth.IDTokenVerifierPort,
	appService *app.AppService,
	sessionTTL time.Duration,
	pkceTTL time.Duration,
	allowedClients map[string]struct{},
	allowedScopes map[string]struct{},
	identityClient identity.IdentityClient,
	auditWriter audit.AuditWriter,
	perm permission.TupleLister,
	cookies httpx.CookieConfig,
) (*Module, error) {

	if sdk == nil {
		return nil, fmt.Errorf("oauth sdk is nil")
	}

	gws, err := oauth_factory_modules.NewOauthGateways(sdk)
	if err != nil {
		return nil, err
	}

	svc := oauth_factory_modules.NewOAuthService(
		gws.Authorization,
		sessionStore,
		pkceStore,
		redirectStore,
		idVerifier,
		gws.Admin,
		appService,
		sessionTTL,
		pkceTTL,
		allowedClients,
		allowedScopes,
	)

	h := oauth_handler.NewOAuthHandler(
		svc,
		appService,
		identityClient,
		auditWriter,
		perm,
		cookies,
	)

	return &Module{
		Service: svc,
		Handler: h,
	}, nil
}
