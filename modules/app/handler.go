package app

import (
	"encoding/json"
	"net/http"
	"strings"

	audit_domain "sso-bff/internal/domain/audit"
	"sso-bff/internal/httpx"
	"sso-bff/internal/middleware"
	"sso-bff/modules/audit"

	"github.com/google/uuid"
)

type AppHandler struct {
	Service *AppService
	Audit   audit.AuditWriter
}

func NewAppHandler(
	service *AppService,
	auditWriter audit.AuditWriter,
) *AppHandler {
	return &AppHandler{
		Service: service,
		Audit:   auditWriter,
	}
}

func (h *AppHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpx.WriteJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "method not allowed"},
		)
		return
	}
	subject, ok := requestSubject(w, r)
	if !ok {
		return
	}

	var payload struct {
		DSN          string `json:"dsn"`
		RedirectPath string `json:"redirect_path"`
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "invalid request body"},
		)
		return
	}
	if strings.TrimSpace(payload.DSN) == "" || strings.TrimSpace(payload.RedirectPath) == "" {
		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "dsn and redirect_path required"},
		)
		return
	}

	app, err := h.Service.Create(
		r.Context(),
		payload.DSN,
		payload.RedirectPath,
	)
	if err != nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()},
		)
		return
	}

	h.logAudit(r, subject, "app_create")
	httpx.WriteJSON(w, http.StatusCreated, app)
}

func (h *AppHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "method not allowed"},
		)
		return
	}
	_, ok := requestSubject(w, r)
	if !ok {
		return
	}

	apps, err := h.Service.List(r.Context())
	if err != nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()},
		)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, apps)
}

func (h *AppHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpx.WriteJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "method not allowed"},
		)
		return
	}
	_, ok := requestSubject(w, r)
	if !ok {
		return
	}

	id, err := h.parseID(r.URL.Path)
	if err != nil {
		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "invalid app id"},
		)
		return
	}

	app, err := h.Service.Get(r.Context(), id)
	if err != nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()},
		)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, app)
}

func (h *AppHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		httpx.WriteJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "method not allowed"},
		)
		return
	}
	subject, ok := requestSubject(w, r)
	if !ok {
		return
	}

	id, err := h.parseID(r.URL.Path)
	if err != nil {
		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "invalid app id"},
		)
		return
	}

	var payload struct {
		DSN          string `json:"dsn"`
		RedirectPath string `json:"redirect_path"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&payload); err != nil {
		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "invalid request body"},
		)
		return
	}
	if strings.TrimSpace(payload.DSN) == "" || strings.TrimSpace(payload.RedirectPath) == "" {
		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "dsn and redirect_path required"},
		)
		return
	}

	app, err := h.Service.Update(
		r.Context(),
		id,
		payload.DSN,
		payload.RedirectPath,
	)
	if err != nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()},
		)
		return
	}

	h.logAudit(r, subject, "app_update")
	httpx.WriteJSON(w, http.StatusOK, app)
}

func (h *AppHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		httpx.WriteJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "method not allowed"},
		)
		return
	}
	subject, ok := requestSubject(w, r)
	if !ok {
		return
	}

	id, err := h.parseID(r.URL.Path)
	if err != nil {
		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "invalid app id"},
		)
		return
	}

	if err := h.Service.Delete(r.Context(), id); err != nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()},
		)
		return
	}

	h.logAudit(r, subject, "app_delete")
	w.WriteHeader(http.StatusNoContent)
}

func requestSubject(w http.ResponseWriter, r *http.Request) (string, bool) {
	if r == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "request unavailable"},
		)
		return "", false
	}

	subject, ok := middleware.SubjectFromContext(r.Context())
	if !ok {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "request subject unavailable"},
		)
		return "", false
	}

	return subject, true
}

func (h *AppHandler) parseID(path string) (uuid.UUID, error) {
	id := strings.TrimPrefix(path, "/api/apps/")
	if idx := strings.Index(id, "/"); idx >= 0 {
		id = id[:idx]
	}
	return uuid.Parse(strings.TrimSpace(id))
}

func (h *AppHandler) logAudit(r *http.Request, subject, eventType string) {
	if h == nil || h.Audit == nil || subject == "" {
		return
	}
	_ = h.Audit.Insert(r.Context(), audit_domain.AuditEvent{
		IdentityID: subject,
		EventType:  eventType,
		IPAddress:  httpx.ClientIP(r),
		UserAgent:  r.UserAgent(),
	})
}
