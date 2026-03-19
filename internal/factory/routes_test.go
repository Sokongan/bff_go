package factory

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sso-bff/modules/app"
	"sso-bff/modules/audit"
	identity_handler "sso-bff/modules/identity/handler"
	oauth_handler "sso-bff/modules/oauth/handler"
	permission_handler "sso-bff/modules/permission/handler"
)

func TestRegisterRoutes_AttachesHandlers(t *testing.T) {
	mux := http.NewServeMux()
	handlers := httpHandlers{
		OAuth:            &oauth_handler.OAuthHandler{},
		IdentityLogin:    &identity_handler.IdentityLoginHandler{},
		IdentityAdmin:    &identity_handler.IdentityAdminHandler{},
		IdentitySettings: &identity_handler.IdentitySettingsHandler{},
		Permission:       &permission_handler.PermissionHandler{},
		Audit:            &audit.AuditHandler{},
		App:              &app.AppHandler{},
	}
	registerRoutes(mux, handlers)

	cases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/login"},
		{http.MethodGet, "/api/callback"},
		{http.MethodPost, "/api/logout"},
		{http.MethodPost, "/api/session/refresh"},
		{http.MethodPost, "/api/identity/login"},
		{http.MethodGet, "/api/identity/settings"},
		{http.MethodPost, "/api/identity/settings"},
		{http.MethodGet, "/api/admin/identities"},
		{http.MethodPost, "/api/admin/identities"},
		{http.MethodPost, "/api/permissions/tuple"},
		{http.MethodGet, "/api/permissions/check"},
		{http.MethodGet, "/api/permissions"},
		{http.MethodGet, "/api/audit/events"},
		{http.MethodGet, "/api/apps"},
		{http.MethodPost, "/api/apps"},
		{http.MethodGet, "/api/apps/123"},
		{http.MethodPut, "/api/apps/123"},
		{http.MethodDelete, "/api/apps/123"},
	}

	for _, tt := range cases {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		handler, pattern := mux.Handler(req)
		if handler == nil {
			t.Fatalf("handler missing for %s %s", tt.method, tt.path)
		}
		if pattern == "" {
			t.Fatalf("no pattern registered for %s %s", tt.method, tt.path)
		}
	}
}
