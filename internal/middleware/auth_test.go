package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	permission_domain "sso-bff/internal/domain/permission"
	"sso-bff/internal/httpx"
	"sso-bff/modules/oauth"
)

type stubSessionResolver struct {
	session oauth.Session
	err     error
}

func (s stubSessionResolver) SubjectBySessionID(ctx context.Context, sessionID string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.session.Subject, nil
}

func (s stubSessionResolver) GetSession(ctx context.Context, sessionID string) (oauth.Session, error) {
	if s.err != nil {
		return oauth.Session{}, s.err
	}
	return s.session, nil
}

type stubPermissionChecker struct {
	allowed bool
	err     error
}

func (s stubPermissionChecker) CheckTuple(ctx context.Context, tuple permission_domain.RelationTuple) (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	return s.allowed, nil
}

func (s stubPermissionChecker) ListTuples(ctx context.Context, params permission_domain.ListTuplesParams) (permission_domain.ListTuplesResult, error) {
	return permission_domain.ListTuplesResult{}, nil
}

func TestRequireSessionLoadsSessionContext(t *testing.T) {
	mw := NewAuthMiddleware(
		stubSessionResolver{session: oauth.Session{Subject: "user-1", KratosToken: "kratos"}},
		nil,
		httpx.CookieConfig{Name: "session"},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/identity/settings", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "sid-1"})
	rec := httptest.NewRecorder()

	called := false
	mw.RequireSession(func(w http.ResponseWriter, r *http.Request) {
		called = true
		session, ok := SessionFromContext(r.Context())
		if !ok {
			t.Fatal("expected session in context")
		}
		if session.Subject != "user-1" {
			t.Fatalf("expected subject user-1, got %q", session.Subject)
		}

		subject, ok := SubjectFromContext(r.Context())
		if !ok || subject != "user-1" {
			t.Fatalf("expected subject context, got %q", subject)
		}

		w.WriteHeader(http.StatusNoContent)
	}).ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected wrapped handler to run")
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestRequireSessionRejectsMissingCookie(t *testing.T) {
	mw := NewAuthMiddleware(
		stubSessionResolver{session: oauth.Session{Subject: "user-1"}},
		nil,
		httpx.CookieConfig{Name: "session"},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/identity/settings", nil)
	rec := httptest.NewRecorder()

	mw.RequireSession(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("wrapped handler should not run")
	}).ServeHTTP(rec, req)

	assertJSONError(t, rec, http.StatusUnauthorized, "missing session token")
}

func TestRequireAdminRejectsNonAdmin(t *testing.T) {
	mw := NewAuthMiddleware(
		stubSessionResolver{session: oauth.Session{Subject: "user-1"}},
		stubPermissionChecker{allowed: false},
		httpx.CookieConfig{Name: "session"},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/apps", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "sid-1"})
	rec := httptest.NewRecorder()

	mw.RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("wrapped handler should not run")
	}).ServeHTTP(rec, req)

	assertJSONError(t, rec, http.StatusForbidden, "admin access required")
}

func TestRequireAdminPropagatesPermissionErrors(t *testing.T) {
	mw := NewAuthMiddleware(
		stubSessionResolver{session: oauth.Session{Subject: "user-1"}},
		stubPermissionChecker{err: errors.New("permission down")},
		httpx.CookieConfig{Name: "session"},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/apps", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "sid-1"})
	rec := httptest.NewRecorder()

	mw.RequireAdmin(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("wrapped handler should not run")
	}).ServeHTTP(rec, req)

	assertJSONError(t, rec, http.StatusBadGateway, "permission down")
}

func assertJSONError(t *testing.T, rec *httptest.ResponseRecorder, status int, message string) {
	t.Helper()

	if rec.Code != status {
		t.Fatalf("expected %d, got %d", status, rec.Code)
	}

	var payload map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload["error"] != message {
		t.Fatalf("expected error %q, got %q", message, payload["error"])
	}
}
