package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	audit_domain "sso-bff/internal/domain/audit"
	permission_domain "sso-bff/internal/domain/permission"
	"sso-bff/internal/httpx"
	"sso-bff/modules/audit"
	"sso-bff/modules/oauth"
	"sso-bff/modules/permission"

	"github.com/google/uuid"
)

var (
	appAdminTuple = permission_domain.RelationTuple{
		Namespace: "app",
		Object:    "sso-portal",
		Relation:  "admin",
	}
)

type AppHandler struct {
	Service    *AppService
	Sessions   oauth.SessionResolver
	Permission permission.PermissionChecker
	Audit      audit.AuditWriter
	Cookies    httpx.CookieConfig
}

func NewAppHandler(
	service *AppService,
	sessions oauth.SessionResolver,
	perm permission.PermissionChecker,
	auditWriter audit.AuditWriter,
	cookies httpx.CookieConfig,
) *AppHandler {
	return &AppHandler{
		Service:    service,
		Sessions:   sessions,
		Permission: perm,
		Audit:      auditWriter,
		Cookies:    cookies,
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
	subject, ok := h.requireAdmin(w, r)
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
	_, ok := h.requireAdmin(w, r)
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
	_, ok := h.requireAdmin(w, r)
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
	subject, ok := h.requireAdmin(w, r)
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
	subject, ok := h.requireAdmin(w, r)
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

func (h *AppHandler) requireAdmin(w http.ResponseWriter, r *http.Request) (string, bool) {
	if h == nil || h.Service == nil || h.Sessions == nil || h.Permission == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "app service unavailable"},
		)
		return "", false
	}

	sessionID := httpx.SessionIDFromRequest(r, h.Cookies)
	if sessionID == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "missing session token"},
		)
		return "", false
	}

	subject, err := h.Sessions.SubjectBySessionID(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, oauth.ErrSessionNotFound) {
			httpx.WriteJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{"error": "session not found"},
			)
			return "", false
		}
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()},
		)
		return "", false
	}
	if strings.TrimSpace(subject) == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "missing session subject"},
		)
		return "", false
	}

	tuple := appAdminTuple
	tuple.SubjectID = subject

	allowed, err := h.Permission.CheckTuple(r.Context(), tuple)
	if err != nil {
		httpx.WriteJSON(
			w,
			http.StatusBadGateway,
			map[string]string{"error": err.Error()},
		)
		return "", false
	}

	if !allowed {
		httpx.WriteJSON(
			w,
			http.StatusForbidden,
			map[string]string{"error": "admin access required"},
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
