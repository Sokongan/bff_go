package factory

import (
	"encoding/json"
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
	registerRoutes(mux, handlers, newRouteMiddleware(nil))

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
		{http.MethodGet, healthPath},
		{http.MethodGet, discoveriesPath},
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

func TestRegisterRoutes_DiscoveriesPayload(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux, stubHandlers(), newRouteMiddleware(nil))

	req := httptest.NewRequest(http.MethodGet, discoveriesPath, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload discoveryResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode discovery payload: %v", err)
	}

	if len(payload.Routes) == 0 {
		t.Fatal("expected discovery payload to include routes")
	}

	if !hasRoute(payload.Routes, "/api/login") {
		t.Fatal("discovery payload missing /api/login")
	}
	if !hasRoute(payload.Routes, healthPath) {
		t.Fatalf("discovery payload missing %s", healthPath)
	}

	methods, ok := findRouteMethods(payload.Routes, discoveriesPath)
	if !ok {
		t.Fatalf("discovery payload missing %s", discoveriesPath)
	}
	if !containsMethod(methods, http.MethodGet) {
		t.Fatalf("%s should declare GET", discoveriesPath)
	}
}

type discoveryResponse struct {
	Routes []discoveryRouteSpec `json:"routes"`
}

type discoveryRouteSpec struct {
	Path    string   `json:"path"`
	Methods []string `json:"methods"`
}

func hasRoute(routes []discoveryRouteSpec, path string) bool {
	for _, route := range routes {
		if route.Path == path {
			return true
		}
	}
	return false
}

func findRouteMethods(routes []discoveryRouteSpec, path string) ([]string, bool) {
	for _, route := range routes {
		if route.Path == path {
			return route.Methods, true
		}
	}
	return nil, false
}

func containsMethod(methods []string, method string) bool {
	for _, m := range methods {
		if m == method {
			return true
		}
	}
	return false
}
