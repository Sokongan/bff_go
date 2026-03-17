package audit

import (
	"net/http"
	"sso-bff/internal/httpx"
	"sso-bff/modules/oauth"
	"strconv"
)

type AuditHandler struct {
	Audit    *AuditService
	Sessions oauth.SubjectResolver
	Cookies  httpx.CookieConfig
}

func NewAuditHandler(
	auditSvc *AuditService,
	sessions oauth.SubjectResolver,
	cookies httpx.CookieConfig,
) *AuditHandler {
	return &AuditHandler{
		Audit:    auditSvc,
		Sessions: sessions,
		Cookies:  cookies,
	}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {

	if h == nil || h.Audit == nil || h.Sessions == nil {
		http.Error(
			w,
			"audit service unavailable",
			http.StatusInternalServerError)
		return
	}

	sessionID := httpx.SessionIDFromRequest(r, h.Cookies)
	if sessionID == "" {
		http.Error(w, "missing session token", http.StatusUnauthorized)
		return
	}

	subject, err := h.Sessions.SubjectBySessionID(
		r.Context(),
		sessionID,
	)

	if err != nil {
		if err == oauth.ErrSessionNotFound {
			http.Error(w, "session not found", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	limit := int32(50)
	if raw := r.URL.Query().Get("limit"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			if n > 200 {
				n = 200
			}
			limit = int32(n)
		}
	}

	events, err := h.Audit.ListRecent(r.Context(), subject, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{"events": events})
}
