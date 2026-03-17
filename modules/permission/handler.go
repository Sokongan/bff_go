package permission

import (
	"encoding/json"
	"net/http"
	audit_domain "sso-bff/internal/domain/audit"
	permission_domain "sso-bff/internal/domain/permission"
	"sso-bff/internal/httpx"
	"sso-bff/modules/audit"
	"sso-bff/modules/oauth"
	"strconv"
	"strings"
)

type PermissionHandlerType struct {
	perm     *PermissionService
	sessions oauth.SubjectResolver
	audit    audit.AuditWriter
	cookies  httpx.CookieConfig
}

func NewPermissionHandler(
	perm *PermissionService,
	sessions oauth.SubjectResolver,
	audit audit.AuditWriter,
	cookies httpx.CookieConfig,
) *PermissionHandlerType {
	return &PermissionHandlerType{
		perm:     perm,
		sessions: sessions,
		audit:    audit,
		cookies:  cookies,
	}
}

func (h *PermissionHandlerType) WriteTuple(
	w http.ResponseWriter,
	r *http.Request,
) {

	if h == nil || h.perm == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "Permission service unavailable",
			})
		return
	}

	if h.sessions == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": "session service unavailable",
			})
		return
	}

	if r.Method != http.MethodPost {
		httpx.WriteJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{
				"error": "method not allowed",
			})
		return
	}

	sessionID := httpx.SessionIDFromRequest(r, h.cookies)

	if sessionID == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{
				"error": "missing session token",
			})
		return
	}

	subject, err := h.sessions.SubjectBySessionID(
		r.Context(),
		sessionID,
	)

	if err != nil {

		if err == oauth.ErrSessionNotFound {
			httpx.WriteJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{
					"error": "session not found",
				})
			return
		}

		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			})
		return
	}

	if strings.TrimSpace(subject) == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{
				"error": "missing session subject",
			})
		return
	}

	var req permission_domain.RelationTuple

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid request body",
			})
		return
	}
	if strings.TrimSpace(req.Namespace) == "" ||
		strings.TrimSpace(req.Object) == "" ||
		strings.TrimSpace(req.Relation) == "" ||
		strings.TrimSpace(req.SubjectID) == "" {

		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "missing tuple fields"})
		return
	}

	if err := h.perm.WriteTuple(
		r.Context(),
		permission_domain.RelationTuple(req),
	); err != nil {

		httpx.WriteJSON(
			w,
			http.StatusBadGateway,
			map[string]string{
				"error": err.Error(),
			})
		return
	}

	if h.audit != nil {
		_ = h.audit.Insert(r.Context(), audit_domain.AuditEvent{
			IdentityID: subject,
			EventType:  "keto_write",
			IPAddress:  httpx.ClientIP(r),
			UserAgent:  r.UserAgent(),
		})
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		map[string]string{"status": "ok"},
	)
}

func (h *PermissionHandlerType) CheckTuple(
	w http.ResponseWriter,
	r *http.Request,
) {
	if h == nil || h.perm == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "keto service unavailable"})
		return
	}

	if h.sessions == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "session service unavailable"})
		return
	}

	if r.Method != http.MethodGet {
		httpx.WriteJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "method not allowed"})
		return
	}

	sessionID := httpx.SessionIDFromRequest(r, h.cookies)

	if sessionID == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "missing session token"})
		return
	}

	subject, err := h.sessions.SubjectBySessionID(
		r.Context(),
		sessionID,
	)

	if err != nil {
		if err == oauth.ErrSessionNotFound {
			httpx.WriteJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{"error": "session not found"})
			return
		}
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()})
		return
	}
	if strings.TrimSpace(subject) == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "missing session subject"})
		return
	}

	req := permission_domain.RelationTuple{
		Namespace: r.URL.Query().Get("namespace"),
		Object:    r.URL.Query().Get("object"),
		Relation:  r.URL.Query().Get("relation"),
		SubjectID: r.URL.Query().Get("subject_id"),
	}
	if strings.TrimSpace(req.Namespace) == "" ||
		strings.TrimSpace(req.Object) == "" ||
		strings.TrimSpace(req.Relation) == "" ||
		strings.TrimSpace(req.SubjectID) == "" {

		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "missing tuple fields"})
		return
	}

	allowed, err := h.perm.CheckTuple(
		r.Context(),
		permission_domain.RelationTuple(req),
	)

	if err != nil {
		httpx.WriteJSON(
			w,
			http.StatusBadGateway,
			map[string]string{"error": err.Error()})
		return
	}

	if h.audit != nil {
		_ = h.audit.Insert(r.Context(), audit_domain.AuditEvent{
			IdentityID: subject,
			EventType:  "keto_check",
			IPAddress:  httpx.ClientIP(r),
			UserAgent:  r.UserAgent(),
		})
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		map[string]any{"allowed": allowed},
	)
}

func (h *PermissionHandlerType) ListTuples(
	w http.ResponseWriter,
	r *http.Request,
) {
	if h == nil || h.perm == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "keto service unavailable"})
		return
	}

	if h.sessions == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "session service unavailable"})
		return
	}

	if r.Method != http.MethodGet {
		httpx.WriteJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "method not allowed"})
		return
	}

	sessionID := httpx.SessionIDFromRequest(r, h.cookies)

	if sessionID == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "missing session token"})
		return
	}

	subject, err := h.sessions.SubjectBySessionID(
		r.Context(),
		sessionID,
	)

	if err != nil {
		if err == oauth.ErrSessionNotFound {
			httpx.WriteJSON(
				w,
				http.StatusUnauthorized,
				map[string]string{"error": "session not found"})
			return
		}
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()})
		return
	}
	if strings.TrimSpace(subject) == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "missing session subject"})
		return
	}

	q := r.URL.Query()
	params := permission_domain.ListTuplesParams{
		Namespace: strings.TrimSpace(q.Get("namespace")),
		Object:    strings.TrimSpace(q.Get("object")),
		Relation:  strings.TrimSpace(q.Get("relation")),
		SubjectID: strings.TrimSpace(q.Get("subject_id")),
		PageToken: strings.TrimSpace(q.Get("page_token")),
	}
	if params.Namespace == "" || params.Object == "" {
		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "namespace and object required"})
		return
	}
	if v := strings.TrimSpace(q.Get("page_size")); v != "" {
		if n, err := parseInt64(v); err == nil && n > 0 {
			params.PageSize = n
		} else if err != nil {
			httpx.WriteJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "invalid page_size"})
			return
		}
	}

	result, err := h.perm.ListTuples(r.Context(), params)
	if err != nil {
		httpx.WriteJSON(
			w,
			http.StatusBadGateway,
			map[string]string{"error": err.Error()})
		return
	}

	if h.audit != nil {
		_ = h.audit.Insert(r.Context(), audit_domain.AuditEvent{
			IdentityID: subject,
			EventType:  "keto_list",
			IPAddress:  httpx.ClientIP(r),
			UserAgent:  r.UserAgent(),
		})
	}

	httpx.WriteJSON(w, http.StatusOK, result)
}

func parseInt64(raw string) (int64, error) {
	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}
