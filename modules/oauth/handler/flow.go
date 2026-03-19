package oauth_handler

import (
	"errors"
	"log"
	"sso-bff/internal/domain"
	audit_domain "sso-bff/internal/domain/audit"
	"sso-bff/internal/httpx"
	"sso-bff/modules/app"
	"sso-bff/modules/audit"
	"sso-bff/modules/identity"
	"sso-bff/modules/oauth"
	oauth_factory_modules "sso-bff/modules/oauth/factory/modules"
	oauth_helper_redirect "sso-bff/modules/oauth/helper/redirect"
	"sso-bff/modules/permission"

	"net/http"

	"strings"
	"time"
)

type OAuthHandler struct {
	OAuth      *oauth_factory_modules.OAuthService
	Apps       *app.AppService
	Identity   identity.IdentityClient
	Audit      audit.AuditWriter
	Permission permission.TupleLister
	Cookies    httpx.CookieConfig
}

func NewOAuthHandler(
	svc *oauth_factory_modules.OAuthService,
	apps *app.AppService,
	identity identity.IdentityClient,
	audit audit.AuditWriter,
	perm permission.TupleLister,
	cookies httpx.CookieConfig,
) *OAuthHandler {

	return &OAuthHandler{
		OAuth:      svc,
		Apps:       apps,
		Identity:   identity,
		Audit:      audit,
		Permission: perm,
		Cookies:    cookies,
	}
}

func (h *OAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if h.OAuth == nil {
		http.Error(
			w,
			"oauth service unavailable",
			http.StatusInternalServerError,
		)
		return
	}

	appID := strings.TrimSpace(r.URL.Query().Get("app"))
	if appID == "" {
		dsn := strings.TrimSpace(r.URL.Query().Get("dsn"))
		if dsn == "" {
			dsn = strings.TrimSpace(r.Header.Get("Origin"))
		}
		if dsn == "" {
			dsn = strings.TrimSpace(r.Referer())
		}
		if dsn != "" && h.Apps != nil {
			if resolved, err := h.Apps.ResolveAppIDByDSN(r.Context(), dsn); err == nil {
				appID = resolved
			}
		}
		if appID == "" {
			httpx.WriteJSON(w, http.StatusBadRequest, map[string]string{
				"error":   "missing app or dsn",
				"example": "/api/login?dsn=http://sso-staging.doj.gov.ph&redirect=/dashboard",
			})
			return
		}
	}
	redirectPath := r.URL.Query().Get("redirect")
	if redirectPath != "" && !oauth_helper_redirect.
		IsSafeRedirectPath(redirectPath) {
		http.Error(w, "redirect not allowed", http.StatusBadRequest)
		return
	}

	redirectURL, err := h.OAuth.FlowService.Login(
		r.Context(),
		appID,
		redirectPath,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *OAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		http.Error(w, "missing code or state", http.StatusBadRequest)
		return
	}

	if h.OAuth == nil {
		http.Error(
			w,
			"oauth service unavailable",
			http.StatusInternalServerError,
		)
		return
	}

	token, err := h.OAuth.FlowService.Callback(r.Context(), code, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	kratosToken := ""
	if c, err := r.Cookie(httpx.KratosSessionTokenCookie); err == nil {
		kratosToken = c.Value
		http.SetCookie(w, &http.Cookie{
			Name:     httpx.KratosSessionTokenCookie,
			Value:    "",
			Path:     "/",
			Domain:   h.Cookies.Domain,
			HttpOnly: true,
			Secure:   httpx.Secure(h.Cookies, r),
			SameSite: httpx.SameSite(h.Cookies),
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
		})
	}

	sessionID, sessionTTL, err := h.OAuth.SessionService.CreateSession(
		r.Context(),
		token,
		kratosToken,
	)

	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, oauth.ErrInvalidIDToken) {
			status = http.StatusUnauthorized
		}
		log.Printf("oauth callback session error: %v", err)

		http.Error(w, err.Error(), status)
		return
	}

	redirectTarget, err := h.OAuth.RedirectService.RedirectForToken(
		r.Context(),
		state,
		token.AccessToken,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if h.Audit != nil {
		if sess, sessErr := h.OAuth.SessionService.GetSession(
			r.Context(),
			sessionID,
		); sessErr == nil {
			_ = h.Audit.Insert(r.Context(), audit_domain.AuditEvent{
				IdentityID: sess.Subject,
				EventType:  "login",
				IPAddress:  httpx.ClientIP(r),
				UserAgent:  r.UserAgent(),
			})
		}
	}

	payload := map[string]interface{}{
		"status":  "ok",
		"expiry":  token.Expiry,
		"session": "issued",
	}
	if redirectTarget != "" {
		payload["redirect"] = redirectTarget
	}

	httpx.SetSessionCookie(w, h.Cookies, sessionID, sessionTTL, r)
	httpx.WriteJSON(w, http.StatusOK, payload)
}

func (h *OAuthHandler) Consent(w http.ResponseWriter, r *http.Request) {
	consentChallenge := r.URL.Query().Get("consent_challenge")

	if consentChallenge == "" {
		http.Error(w, "missing consent_challenge", http.StatusBadRequest)
		return
	}

	if h.OAuth == nil {
		http.Error(
			w,
			"oauth service unavailable",
			http.StatusInternalServerError,
		)
		return
	}
	if h.Identity == nil {
		http.Error(
			w,
			"identity service unavailable",
			http.StatusInternalServerError,
		)
		return
	}

	identity, err := h.Identity.WhoAmI(r.Context(), r.Header.Get("Cookie"))
	if err != nil {
		if errors.Is(err, domain.ErrIdentityMisconfigured) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		identity = nil
	}

	redirectURL, err := h.OAuth.AcceptFlow.OauthConsent(
		r.Context(),
		consentChallenge,
		identity,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *OAuthHandler) LaunchApp(w http.ResponseWriter, r *http.Request) {
	if h == nil || h.OAuth == nil {
		http.Error(
			w,
			"oauth service unavailable",
			http.StatusInternalServerError,
		)
		return
	}

	sessionID := httpx.SessionIDFromRequest(r, h.Cookies)
	if sessionID == "" {
		http.Error(w, "missing session token", http.StatusUnauthorized)
		return
	}
	if _, err := h.OAuth.SessionService.GetSession(r.Context(), sessionID); err != nil {
		if errors.Is(err, oauth.ErrSessionNotFound) {
			http.Error(w, "session not found", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appID := strings.TrimSpace(r.URL.Query().Get("app"))
	if appID == "" {
		http.Error(w, "missing app", http.StatusBadRequest)
		return
	}

	redirectPath := strings.TrimSpace(r.URL.Query().Get("path"))
	if redirectPath != "" && !oauth_helper_redirect.
		IsSafeRedirectPath(redirectPath) {
		http.Error(w, "redirect not allowed", http.StatusBadRequest)
		return
	}

	if h.Apps == nil {
		http.Error(
			w,
			"app registry unavailable",
			http.StatusInternalServerError,
		)
		return
	}

	registry, err := h.Apps.ResolveRegistry(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appRedirect, ok := registry[appID]
	if !ok {
		http.Error(w, "unknown app", http.StatusBadRequest)
		return
	}

	target, err := oauth_helper_redirect.BuildAppRedirect(
		appRedirect,
		redirectPath,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, target, http.StatusFound)
}

func (h *OAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	if h == nil || h.OAuth == nil {
		http.Error(
			w,
			"oauth service unavailable",
			http.StatusInternalServerError,
		)
		return
	}

	sessionID := httpx.SessionIDFromRequest(r, h.Cookies)
	if sessionID != "" {
		session, sessErr := h.OAuth.SessionService.GetSession(
			r.Context(),
			sessionID,
		)
		if sessErr == nil && h != nil && h.Audit != nil {
			_ = h.Audit.Insert(r.Context(), audit_domain.AuditEvent{
				IdentityID: session.Subject,
				EventType:  "logout",
				IPAddress:  httpx.ClientIP(r),
				UserAgent:  r.UserAgent(),
			})
		}
		_ = h.OAuth.SessionService.DeleteSession(r.Context(), sessionID)
	}

	httpx.ClearSessionCookie(w, h.Cookies, r)
	w.WriteHeader(http.StatusNoContent)
}
