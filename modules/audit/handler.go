package audit

import (
	"net/http"
	"sso-bff/internal/httpx"
	"sso-bff/internal/middleware"
	"strconv"
)

type AuditHandler struct {
	Audit *AuditService
}

func NewAuditHandler(auditSvc *AuditService) *AuditHandler {
	return &AuditHandler{
		Audit: auditSvc,
	}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {

	if h == nil || h.Audit == nil {
		http.Error(
			w,
			"audit service unavailable",
			http.StatusInternalServerError)
		return
	}

	subject, ok := middleware.SubjectFromContext(r.Context())
	if !ok {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "request subject unavailable"},
		)
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
