package oauth_factory

import (
	"context"
	"testing"
	"time"

	app_domain "sso-bff/internal/domain/app"
	audit_domain "sso-bff/internal/domain/audit"
	identity_domain "sso-bff/internal/domain/identity"
	oauth_domain "sso-bff/internal/domain/oauth"
	permission_domain "sso-bff/internal/domain/permission"
	"sso-bff/internal/httpx"
	"sso-bff/modules/app"
	"sso-bff/modules/oauth"
	oauth_sdk "sso-bff/modules/oauth/sdk"

	"github.com/google/uuid"
)

type noopSessionStore struct{}

func (noopSessionStore) SaveSession(ctx context.Context, sessionID string, session oauth.Session, ttl time.Duration) error {
	return nil
}

func (noopSessionStore) GetSession(ctx context.Context, sessionID string) (oauth.Session, error) {
	return oauth.Session{}, nil
}

func (noopSessionStore) DeleteSession(ctx context.Context, sessionID string) error {
	return nil
}

type noopPKCEStore struct{}

func (noopPKCEStore) SaveVerifier(ctx context.Context, state, verifier string, ttl time.Duration) error {
	return nil
}

func (noopPKCEStore) GetVerifier(ctx context.Context, state string) (string, error) {
	return "", nil
}

func (noopPKCEStore) DeleteVerifier(ctx context.Context, state string) error {
	return nil
}

type noopRedirectStore struct{}

func (noopRedirectStore) SaveRedirect(ctx context.Context, state, redirectURL string, ttl time.Duration) error {
	return nil
}

func (noopRedirectStore) GetRedirect(ctx context.Context, state string) (string, error) {
	return "", nil
}

func (noopRedirectStore) DeleteRedirect(ctx context.Context, state string) error {
	return nil
}

type noopIDTokenVerifier struct{}

func (noopIDTokenVerifier) Verify(ctx context.Context, rawIDToken string) (oauth_domain.IDTokenClaims, error) {
	return oauth_domain.IDTokenClaims{}, nil
}

type noopAuditWriter struct{}

func (noopAuditWriter) Insert(ctx context.Context, e audit_domain.AuditEvent) error {
	return nil
}

type noopIdentityClient struct{}

func (noopIdentityClient) WhoAmI(ctx context.Context, cookieHeader string) (*identity_domain.Identity, error) {
	return &identity_domain.Identity{}, nil
}

func (noopIdentityClient) WhoAmIWithSessionToken(ctx context.Context, sessionToken string) (*identity_domain.Identity, error) {
	return &identity_domain.Identity{}, nil
}

type noopTupleLister struct{}

func (noopTupleLister) ListTuples(ctx context.Context, params permission_domain.ListTuplesParams) (permission_domain.ListTuplesResult, error) {
	return permission_domain.ListTuplesResult{}, nil
}

type noopAppRepo struct{}

func (noopAppRepo) Create(ctx context.Context, dsn, redirectPath string) (app_domain.AppRegistry, error) {
	return app_domain.AppRegistry{}, nil
}

func (noopAppRepo) Update(ctx context.Context, id uuid.UUID, dsn, redirectPath string) (app_domain.AppRegistry, error) {
	return app_domain.AppRegistry{}, nil
}

func (noopAppRepo) Get(ctx context.Context, id uuid.UUID) (app_domain.AppRegistry, error) {
	return app_domain.AppRegistry{}, nil
}

func (noopAppRepo) List(ctx context.Context) ([]app_domain.AppRegistry, error) {
	return nil, nil
}

func (noopAppRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func TestNewOauthModule_Success(t *testing.T) {
	sdk := oauth_sdk.NewOAuthSDK(
		"http://admin",
		oauth.BrowserClient{
			BrowserPublicURL: "http://public",
			ClientID:         "cid",
			ClientSecret:     "secret",
			RedirectURL:      "http://public/oauth2/callback",
			Scopes:           []string{"openid"},
		},
		oauth.InternalClient{
			TokenURL:     "http://internal",
			ClientID:     "cid",
			ClientSecret: "secret",
			RedirectURL:  "http://internal/oauth2/callback",
			Scopes:       []string{"openid"},
		},
	)

	module, err := NewOauthModule(
		sdk,
		noopSessionStore{},
		noopPKCEStore{},
		noopRedirectStore{},
		noopIDTokenVerifier{},
		app.NewService(noopAppRepo{}),
		time.Minute,
		time.Minute,
		map[string]struct{}{},
		map[string]struct{}{},
		noopIdentityClient{},
		noopAuditWriter{},
		noopTupleLister{},
		httpx.CookieConfig{},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if module == nil || module.Service == nil || module.Handler == nil {
		t.Fatalf("expected module with service and handler")
	}
}

func TestNewOauthModule_MissingSDK(t *testing.T) {
	_, err := NewOauthModule(
		nil,
		noopSessionStore{},
		noopPKCEStore{},
		noopRedirectStore{},
		noopIDTokenVerifier{},
		app.NewService(noopAppRepo{}),
		time.Minute,
		time.Minute,
		map[string]struct{}{},
		map[string]struct{}{},
		noopIdentityClient{},
		noopAuditWriter{},
		noopTupleLister{},
		httpx.CookieConfig{},
	)
	if err == nil {
		t.Fatalf("expected error when sdk is nil")
	}
}
