package oauth_gateway

import (
	"context"
	"sso-bff/internal/domain"
	oauth_types "sso-bff/internal/modules/oauth"
	"time"

	"golang.org/x/oauth2"
)

type OAuthClient interface {
	AuthCodeURL(state, verifier string) string
	Exchange(ctx context.Context, code, verifier string) (*oauth2.Token, error)
	Refresh(ctx context.Context, refreshToken string) (*oauth2.Token, error)
}
type OauthAdmin interface {
	GetConsentRequest(ctx context.Context, consentChallenge string) (*domain.ConsentRequest, error)
	AcceptConsent(ctx context.Context, consentChallenge string, grantScope, grantAudience []string, idTokenClaims, accessTokenClaims map[string]any) (string, error)
}

type SessionStore interface {
	SaveSession(ctx context.Context, sessionID string, session oauth_types.Session, ttl time.Duration) error
	GetSession(ctx context.Context, sessionID string) (oauth_types.Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
}

type IDTokenVerifier interface {
	Verify(ctx context.Context, rawIDToken string) (IDTokenClaims, error)
}

type PKCEStore interface {
	SaveVerifier(ctx context.Context, state, verifier string, ttl time.Duration) error
	GetVerifier(ctx context.Context, state string) (string, error)
	DeleteVerifier(ctx context.Context, state string) error
}

type RedirectStore interface {
	SaveRedirect(ctx context.Context, state, redirectURL string, ttl time.Duration) error
	GetRedirect(ctx context.Context, state string) (string, error)
	DeleteRedirect(ctx context.Context, state string) error
}

type OAuthM2M interface {
	Token(ctx context.Context, scopes []string) (*oauth2.Token, error)
}
