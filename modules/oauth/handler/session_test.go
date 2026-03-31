package oauth_handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	identity_domain "sso-bff/internal/domain/identity"
	"sso-bff/internal/httpx"
	"sso-bff/modules/oauth"
	oauth_factory_modules "sso-bff/modules/oauth/factory/modules"
	service_oauth "sso-bff/modules/oauth/services"
)

type stubSessionStore struct {
	session oauth.Session
}

func (s stubSessionStore) SaveSession(ctx context.Context, sessionID string, session oauth.Session, ttl time.Duration) error {
	return nil
}

func (s stubSessionStore) GetSession(ctx context.Context, sessionID string) (oauth.Session, error) {
	return s.session, nil
}

func (s stubSessionStore) DeleteSession(ctx context.Context, sessionID string) error {
	return nil
}

type stubIdentityClient struct {
	ident *identity_domain.Identity
	err   error
}

func (s stubIdentityClient) WhoAmI(ctx context.Context, cookieHeader string) (*identity_domain.Identity, error) {
	return s.ident, s.err
}

func (s stubIdentityClient) WhoAmIWithSessionToken(ctx context.Context, sessionToken string) (*identity_domain.Identity, error) {
	return s.ident, s.err
}

func TestSessionIncludesOrganizationIDFromKratosMetadataPublic(t *testing.T) {
	expiry := time.Date(2026, 3, 31, 10, 0, 0, 0, time.UTC)
	handler := &OAuthHandler{
		OAuth: &oauth_factory_modules.OAuthService{
			SessionService: service_oauth.NewSessionService(
				nil,
				stubSessionStore{
					session: oauth.Session{
						Subject: "user-123",
						Expiry:  expiry,
					},
				},
				nil,
				0,
			),
		},
		Identity: stubIdentityClient{
			ident: &identity_domain.Identity{
				ID: "user-123",
				Traits: map[string]any{
					"username": "paul.test",
					"name": map[string]any{
						"firstName": "Paul",
						"lastName":  "Test",
					},
				},
				MetadataPublic: map[string]any{
					"organization_id": "org-456",
				},
			},
		},
		Cookies: httpx.CookieConfig{},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/session", nil)
	req.Header.Set("Authorization", "Bearer session-1")
	rec := httptest.NewRecorder()

	handler.Session(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body["organization_id"] != "org-456" {
		t.Fatalf("expected organization_id org-456, got %#v", body["organization_id"])
	}
	profile, ok := body["profile"].(map[string]any)
	if !ok {
		t.Fatalf("expected profile object, got %#v", body["profile"])
	}
	if profile["username"] != "paul.test" {
		t.Fatalf("expected profile.username paul.test, got %#v", profile["username"])
	}
	if body["authenticated"] != true {
		t.Fatalf("expected authenticated true, got %#v", body["authenticated"])
	}
}
