package factory

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sso-bff/internal/config"
	oauthcfg "sso-bff/internal/config/services/oauth"
	"sso-bff/internal/db"
	"sso-bff/internal/db/client"
	"sso-bff/internal/httpx"
	"sso-bff/modules"
	"sso-bff/modules/app"
	"sso-bff/modules/audit"
	identity_handler "sso-bff/modules/identity/handler"
	oauth_handler "sso-bff/modules/oauth/handler"
	permission_handler "sso-bff/modules/permission/handler"
)

func TestRegisterRoutesRegistersEndpoints(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux, stubHandlers(), newRouteMiddleware(nil))

	paths := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/login"},
		{http.MethodGet, "/api/callback"},
		{http.MethodGet, "/api/logout"},
		{http.MethodGet, "/api/session"},
		{http.MethodGet, "/api/identity/login"},
		{http.MethodGet, "/api/identity/settings"},
		{http.MethodGet, "/api/admin/identities"},
		{http.MethodGet, "/api/permissions"},
		{http.MethodGet, "/api/audit/events"},
		{http.MethodGet, "/api/apps"},
		{http.MethodGet, "/api/apps/123"},
	}

	for _, p := range paths {
		req := httptest.NewRequest(p.method, p.path, nil)
		handler, pattern := mux.Handler(req)
		if pattern == "" {
			t.Fatalf("expected route %s to be registered", p.path)
		}
		if handler == nil {
			t.Fatalf("handler missing for route %s", p.path)
		}
	}
}

func TestVerbHandlerRejectsWrongMethod(t *testing.T) {
	called := false
	handler := verbHandler(http.MethodPost, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/session/refresh", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
	if called {
		t.Fatal("handler should not be invoked when method is wrong")
	}
}

func TestCookieConfigFrom(t *testing.T) {
	cfg := &oauthcfg.CookieConfig{
		Name:     "session",
		Domain:   "example.com",
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	}
	got := cookieConfigFrom(cfg)
	if got.Name != cfg.Name || got.Domain != cfg.Domain || got.SameSite != cfg.SameSite || got.Secure != cfg.Secure {
		t.Fatalf("unexpected cookie config: %#v", got)
	}

	if empty := cookieConfigFrom(nil); empty != (httpx.CookieConfig{}) {
		t.Fatalf("expected zero value for nil config, got %#v", empty)
	}
}

func TestNewHandlersRequiresConfig(t *testing.T) {
	_, err := NewHandlers(nil, nil, &modules.SDKs{})
	if err == nil {
		t.Fatal("expected error when config is nil")
	}
}

func TestNewHandlersRequiresResources(t *testing.T) {
	cfg := &config.Config{Oauth: minimalOAuthConfig()}
	resources := &db.Resources{Store: &client.Store{}}
	_, err := NewHandlers(cfg, resources, &modules.SDKs{})
	if err == nil {
		t.Fatal("expected error when sql resources are unavailable")
	}
}

func TestNewHandlersRequiresSDKs(t *testing.T) {
	cfg := &config.Config{Oauth: minimalOAuthConfig()}
	resources := stubResources()
	_, err := NewHandlers(cfg, resources, nil)
	if err == nil {
		t.Fatal("expected error when sdk bundle is nil")
	}
}

func TestNewHandlersRequiresOauthConfig(t *testing.T) {
	cfg := &config.Config{Oauth: &oauthcfg.OAuthConfig{
		Cookie: &oauthcfg.CookieConfig{},
		Client: &oauthcfg.ClientConfig{},
		OIDC:   &oauthcfg.OIDCConfig{Issuer: "https://issuer", JWKS: "https://issuer/.well-known/jwks.json"},
	}}
	resources := stubResources()
	_, err := NewHandlers(cfg, resources, &modules.SDKs{})
	if err == nil {
		t.Fatal("expected error when oauth config is incomplete")
	}
}

func stubHandlers() httpHandlers {
	return httpHandlers{
		OAuth:            &oauth_handler.OAuthHandler{},
		IdentityLogin:    &identity_handler.IdentityLoginHandler{},
		IdentityAdmin:    &identity_handler.IdentityAdminHandler{},
		IdentitySettings: &identity_handler.IdentitySettingsHandler{},
		Permission:       &permission_handler.PermissionHandler{},
		Audit:            &audit.AuditHandler{},
		App:              &app.AppHandler{},
	}
}

func stubResources() *db.Resources {
	return &db.Resources{
		Db:    &client.Db{},
		Store: &client.Store{},
	}
}

func minimalOAuthConfig() *oauthcfg.OAuthConfig {
	return &oauthcfg.OAuthConfig{
		Cookie: &oauthcfg.CookieConfig{Name: "session"},
		Client: &oauthcfg.ClientConfig{
			ClientID:     "cid",
			ClientSecret: "secret",
			RedirectURL:  "https://example.com/oauth2/callback",
		},
		Scopes: &oauthcfg.ClientScopesConfig{
			AllowedClient: map[string]struct{}{"cid": {}},
			AllowedScope:  map[string]struct{}{"scope": {}},
			BFFScopes:     []string{"openid"},
		},
		OIDC: &oauthcfg.OIDCConfig{
			Issuer: "https://issuer",
			JWKS:   "https://issuer/.well-known/jwks.json",
		},
	}
}
