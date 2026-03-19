package oauth_factory_modules

import (
	"sso-bff/modules/app"
	"sso-bff/modules/oauth"
	service_oauth "sso-bff/modules/oauth/services"
	"time"
)

type OAuthService struct {
	FlowService     *service_oauth.FlowService
	SessionService  *service_oauth.SessionService
	RedirectService *service_oauth.RedirectService
	OIDCService     *service_oauth.OIDCService
	AcceptFlow      *service_oauth.AcceptFlowService
}

func NewOAuthService(
	oauthClient oauth.OAuthClientPort,

	sessoinStore oauth.SessionStorePort,

	pkceStore oauth.PKCEStorePort,

	redirectStore oauth.RedirectStorePort,

	idVerifier oauth.IDTokenVerifierPort,

	adminPort oauth.OauthAdminPort,

	appRegistry *app.AppService,

	sessionTTL time.Duration,
	pkceTTL time.Duration,

	allowedClients map[string]struct{},
	allowedScopes map[string]struct{},

) *OAuthService {

	return &OAuthService{
		FlowService: service_oauth.NewFlowService(
			oauthClient,
			pkceStore,
			redirectStore,
			appRegistry,
			pkceTTL,
		),
		SessionService: service_oauth.NewSessionService(
			oauthClient,
			sessoinStore,
			idVerifier,
			sessionTTL,
		),
		RedirectService: service_oauth.NewRedirectService(
			redirectStore,
			appRegistry,
		),
		OIDCService: service_oauth.NewOIDCService(idVerifier),

		AcceptFlow: service_oauth.NewAcceptFlowService(
			adminPort,
			allowedClients,
			allowedScopes,
		),
	}
}
