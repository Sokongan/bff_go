package oauth

import (
	"context"

	oauth_domain "sso-bff/internal/domain/oauth"
	"time"

	"golang.org/x/oauth2"
)

type OAuthClientPort interface {
	AuthCodeURL(state, verifier string) string
	Exchange(
		ctx context.Context,
		code,
		verifier string,
	) (*oauth2.Token, error)
	Refresh(
		ctx context.Context,
		refreshToken string,
	) (*oauth2.Token, error)
}
type OauthAdminPort interface {
	GetConsentRequest(
		ctx context.Context,
		consentChallenge string,
	) (*oauth_domain.ConsentRequest, error)
	AcceptConsent(
		ctx context.Context,
		consentChallenge string,
		grantScope,
		grantAudience []string,
		idTokenClaims,
		accessTokenClaims map[string]any,
	) (string, error)
}

type OAuthM2MPort interface {
	Token(ctx context.Context, scopes []string) (*oauth2.Token, error)
}

type SessionStorePort interface {
	SaveSession(
		ctx context.Context,
		sessionID string,
		session Session,
		ttl time.Duration,
	) error
	GetSession(
		ctx context.Context,
		sessionID string,
	) (Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
}

type IDTokenVerifierPort interface {
	Verify(
		ctx context.Context,
		rawIDToken string,
	) (oauth_domain.IDTokenClaims, error)
}

type PKCEStorePort interface {
	SaveVerifier(
		ctx context.Context,
		state,
		verifier string,
		ttl time.Duration,
	) error
	GetVerifier(ctx context.Context, state string) (string, error)
	DeleteVerifier(ctx context.Context, state string) error
}

type RedirectStorePort interface {
	SaveRedirect(
		ctx context.Context,
		state,
		redirectURL string,
		ttl time.Duration,
	) error
	GetRedirect(ctx context.Context, state string) (string, error)
	DeleteRedirect(ctx context.Context, state string) error
}

type SubjectResolver interface {
	SubjectBySessionID(
		ctx context.Context,
		sessionID string,
	) (string, error)
}

type SessionReader interface {
	GetSession(
		ctx context.Context,
		sessionID string,
	) (Session, error)
}

type SessionResolver interface {
	SubjectBySessionID(
		ctx context.Context,
		sessionID string,
	) (string, error)

	GetSession(
		ctx context.Context,
		sessionID string,
	) (Session, error)
}
