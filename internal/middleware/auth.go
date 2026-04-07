package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	permission_domain "sso-bff/internal/domain/permission"
	"sso-bff/internal/httpx"
	"sso-bff/modules/oauth"
	"sso-bff/modules/permission"
)

type contextKey string

const (
	sessionContextKey contextKey = "auth.session"
	subjectContextKey contextKey = "auth.subject"
)

var adminTuple = permission_domain.RelationTuple{
	Namespace: "app",
	Object:    "sso-portal",
	Relation:  "admin",
}

type AuthMiddleware struct {
	sessions oauth.SessionResolver
	perm     permission.PermissionChecker
	cookies  httpx.CookieConfig
}

func NewAuthMiddleware(
	sessions oauth.SessionResolver,
	perm permission.PermissionChecker,
	cookies httpx.CookieConfig,
) *AuthMiddleware {
	return &AuthMiddleware{
		sessions: sessions,
		perm:     perm,
		cookies:  cookies,
	}
}

func (m *AuthMiddleware) RequireSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if m == nil || m.sessions == nil {
			httpx.WriteJSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "session service unavailable"},
			)
			return
		}

		sessionID := httpx.SessionIDFromRequest(r, m.cookies)
		if sessionID == "" {
			httpx.WriteJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{"error": "missing session token"},
			)
			return
		}

		session, err := m.sessions.GetSession(r.Context(), sessionID)
		if err != nil {
			if errors.Is(err, oauth.ErrSessionNotFound) {
				httpx.WriteJSON(
					w,
					http.StatusUnauthorized,
					map[string]string{"error": "session not found"},
				)
				return
			}

			httpx.WriteJSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": err.Error()},
			)
			return
		}

		subject := strings.TrimSpace(session.Subject)
		if subject == "" {
			httpx.WriteJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{"error": "missing session subject"},
			)
			return
		}

		ctx := context.WithValue(r.Context(), sessionContextKey, session)
		ctx = context.WithValue(ctx, subjectContextKey, subject)
		next(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return m.RequireSession(func(w http.ResponseWriter, r *http.Request) {
		if m == nil || m.perm == nil {
			httpx.WriteJSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "permission service unavailable"},
			)
			return
		}

		subject, ok := SubjectFromContext(r.Context())
		if !ok {
			httpx.WriteJSON(
				w,
				http.StatusInternalServerError,
				map[string]string{"error": "request subject unavailable"},
			)
			return
		}

		tuple := adminTuple
		tuple.SubjectID = subject

		allowed, err := m.perm.CheckTuple(r.Context(), tuple)
		if err != nil {
			httpx.WriteJSON(
				w,
				http.StatusBadGateway,
				map[string]string{"error": err.Error()},
			)
			return
		}

		if !allowed {
			httpx.WriteJSON(
				w,
				http.StatusForbidden,
				map[string]string{"error": "admin access required"},
			)
			return
		}

		next(w, r)
	})
}

func SubjectFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	subject, ok := ctx.Value(subjectContextKey).(string)
	return subject, ok && strings.TrimSpace(subject) != ""
}

func SessionFromContext(ctx context.Context) (oauth.Session, bool) {
	if ctx == nil {
		return oauth.Session{}, false
	}

	session, ok := ctx.Value(sessionContextKey).(oauth.Session)
	return session, ok
}
