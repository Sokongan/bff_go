package identity

import (
	"context"
	"net/http"
	identity_domain "sso-bff/internal/domain/identity"

	identity "github.com/ory/kratos-client-go"
)

type IdentityLoginClient interface {
	CreateNativeLoginFlow(ctx context.Context) (*identity.LoginFlow,
		*http.Response, error,
	)
	UpdateLoginFlow(ctx context.Context,
		flowID string,
		body identity.UpdateLoginFlowBody,
	) (*identity.SuccessfulNativeLogin, *http.Response, error)
}

type IdentityClient interface {
	WhoAmI(ctx context.Context, cookieHeader string) (*identity_domain.Identity, error)
	WhoAmIWithSessionToken(ctx context.Context, sessionToken string) (*identity_domain.Identity, error)
}

type IdentitySettingsClient interface {
	CreateBrowserSettingsFlow(ctx context.Context,
		cookieHeader string,
	) (*identity.SettingsFlow, *http.Response, error)
	UpdateSettingsFlow(ctx context.Context,
		flowID string,
		body identity.UpdateSettingsFlowBody,
		cookieHeader string,
	) (*identity.SettingsFlow, *http.Response, error)
	CreateNativeSettingsFlow(ctx context.Context,
		sessionToken string,
	) (*identity.SettingsFlow, *http.Response, error)
	UpdateSettingsFlowWithSessionToken(ctx context.Context,
		flowID string,
		body identity.UpdateSettingsFlowBody,
		sessionToken string,
	) (*identity.SettingsFlow, *http.Response, error)
}

type IdentityAdminClient interface {
	CreateIdentity(ctx context.Context,
		traits map[string]any,
		schemaID string,
	) (*identity_domain.Identity, error)
	ListIdentities(ctx context.Context,
		params identity_domain.ListIdentitiesParams,
	) ([]identity_domain.Identity, error)
	ListSessions(ctx context.Context,
		params identity_domain.ListSessionsParams,
	) ([]identity_domain.IdentitySession, error)
}

type OauthLoginAdmin interface {
	AcceptLogin(ctx context.Context,
		loginChallenge,
		subject string,
	) (string, error)
}
