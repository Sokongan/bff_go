package identity_handler

import (
	"encoding/json"
	"net/http"

	identity_domain "sso-bff/internal/domain/identity"
	permission_domain "sso-bff/internal/domain/permission"
	"sso-bff/internal/httpx"
	"sso-bff/modules/identity"

	identity_factory_modules "sso-bff/modules/identity/factory/modules"
	identity_helper "sso-bff/modules/identity/helper"
	"sso-bff/modules/oauth"
	"sso-bff/modules/permission"
	"strconv"
	"strings"
)

type IdentityAdminHandler struct {
	Service    *identity_factory_modules.IdentityServices
	Sessions   oauth.SessionResolver
	Permission permission.PermissionChecker
	Cookies    httpx.CookieConfig
}

func NewIdentityAdminHandler(
	admin *identity_factory_modules.IdentityServices,
	sessions oauth.SessionResolver,
	permission permission.PermissionChecker,
	cookies httpx.CookieConfig,
) *IdentityAdminHandler {

	return &IdentityAdminHandler{
		Service:    admin,
		Sessions:   sessions,
		Permission: permission,
		Cookies:    cookies,
	}
}

func (h *IdentityAdminHandler) CreateIdentity(
	w http.ResponseWriter,
	r *http.Request,
) {
	if h == nil || h.Service == nil ||
		h.Sessions == nil || h.Permission == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "identity admin unavailable"},
		)
		return
	}
	if r.Method != http.MethodPost {
		httpx.WriteJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "method not allowed"},
		)
		return
	}

	sessionID := httpx.SessionIDFromRequest(r, h.Cookies)
	if sessionID == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "missing session token"},
		)
		return
	}

	subject, err := h.Sessions.SubjectBySessionID(r.Context(), sessionID)
	if err != nil {
		if err == oauth.ErrSessionNotFound {
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
	if strings.TrimSpace(subject) == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "missing session subject"},
		)
		return
	}

	allowed, err := h.Permission.CheckTuple(
		r.Context(), permission_domain.RelationTuple{
			Namespace: "app",
			Object:    "sso-portal",
			Relation:  "admin",
			SubjectID: subject,
		})
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
			map[string]string{"error": "access required"},
		)
		return
	}

	var req struct {
		SchemaID string         `json:"schema_id"`
		Traits   map[string]any `json:"traits"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "invalid request body"},
		)
		return
	}
	if len(req.Traits) == 0 {
		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "traits required"},
		)
		return
	}

	ident, err := h.Service.Admin.CreateIdentity(
		r.Context(),
		req.Traits,
		req.SchemaID,
	)
	if err != nil {
		httpx.WriteJSON(
			w,
			http.StatusBadGateway,
			map[string]string{"error": err.Error()})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"id":     ident.ID,
		"traits": ident.Traits,
	})
}

func (h *IdentityAdminHandler) ListIdentities(
	w http.ResponseWriter,
	r *http.Request,
) {
	if h == nil || h.Service == nil ||
		h.Sessions == nil || h.Permission == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "identity admin unavailable"},
		)
		return
	}
	if r.Method != http.MethodGet {
		httpx.WriteJSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "method not allowed"},
		)
		return
	}

	sessionID := httpx.SessionIDFromRequest(r, h.Cookies)
	if sessionID == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "missing session token"},
		)
		return
	}
	subject, err := h.Sessions.SubjectBySessionID(r.Context(), sessionID)
	if err != nil {
		if err == oauth.ErrSessionNotFound {
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
	if strings.TrimSpace(subject) == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "missing session subject"},
		)
		return
	}

	allowed, err := h.Permission.CheckTuple(
		r.Context(), permission_domain.RelationTuple{
			Namespace: "app",
			Object:    "sso-portal",
			Relation:  "admin",
			SubjectID: subject,
		})

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

	q := r.URL.Query()
	params := identity_domain.ListIdentitiesParams{
		PageToken: strings.TrimSpace(q.Get("page_token")),
		CredentialsIdentifier: strings.
			TrimSpace(q.Get("credentials_identifier")),
	}
	tuplesNamespace := strings.TrimSpace(q.Get("tuples_namespace"))
	tuplesObject := strings.TrimSpace(q.Get("tuples_object"))
	tuplesRelation := strings.TrimSpace(q.Get("tuples_relation"))
	includeKratosSessions := true
	if v := strings.TrimSpace(
		q.Get("include_kratos_sessions")); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			includeKratosSessions = parsed
		} else {
			httpx.WriteJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "invalid include_kratos_sessions"})
			return
		}
	}
	activeSessions := true
	if v := strings.TrimSpace(q.Get("sessions_active")); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			activeSessions = parsed
		} else {
			httpx.WriteJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "invalid sessions_active"},
			)
			return
		}
	}
	if v := strings.TrimSpace(q.Get("page")); v != "" {
		if n, err := identity_helper.
			ParseInt64(v); err == nil && n > 0 {
			params.Page = n
		} else if err != nil {
			httpx.WriteJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "invalid page"},
			)
			return
		}
	}
	if v := strings.TrimSpace(q.Get("per_page")); v != "" {
		if n, err := identity_helper.
			ParseInt64(v); err == nil && n > 0 {
			params.PerPage = n
		} else if err != nil {
			httpx.WriteJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "invalid per_page"},
			)
			return
		}
	}
	if v := strings.TrimSpace(q.Get("page_size")); v != "" {
		if n, err := identity_helper.
			ParseInt64(v); err == nil && n > 0 {
			params.PageSize = n
		} else if err != nil {
			httpx.WriteJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "invalid page_size"},
			)
			return
		}
	}

	identities, err := h.Service.Admin.ListIdentities(
		r.Context(),
		params,
	)

	if err != nil {
		httpx.WriteJSON(
			w,
			http.StatusBadGateway,
			map[string]string{"error": err.Error()},
		)
		return
	}

	var bffSessionPayload map[string]any
	if h.Sessions != nil {
		if sess, err := h.Sessions.GetSession(r.Context(), sessionID); err == nil {
			bffSessionPayload = map[string]any{
				"subject": sess.Subject,
				"exp":     sess.Expiry,
			}
		}
	}

	kratosDevicesByIdentity := map[string][]identity_domain.
		IdentitySessionDevice{}

	if includeKratosSessions {
		params := identity_domain.ListSessionsParams{
			Active: &activeSessions,
		}
		sessions, err := h.Service.Admin.ListSessions(
			r.Context(),
			params,
		)

		if err != nil {
			httpx.WriteJSON(
				w,
				http.StatusBadGateway,
				map[string]string{"error": err.Error()},
			)
			return
		}
		for _, sess := range sessions {
			if sess.IdentityID == "" {
				continue
			}
			if len(sess.Devices) == 0 {
				continue
			}
			kratosDevicesByIdentity[sess.IdentityID] = append(
				kratosDevicesByIdentity[sess.IdentityID],
				sess.Devices...,
			)
		}
	}

	out := make([]identity.IdentityWithRoles, 0, len(identities))
	for _, ident := range identities {
		item := identity.IdentityWithRoles{
			ID:             ident.ID,
			Traits:         ident.Traits,
			KratosSessions: kratosDevicesByIdentity[ident.ID],
		}

		if tuplesNamespace != "" {
			tupleParams := permission_domain.ListTuplesParams{
				Namespace: tuplesNamespace,
				Object:    tuplesObject,
				Relation:  tuplesRelation,
				SubjectID: ident.ID,
			}
			tuples, err := h.Permission.ListTuples(
				r.Context(),
				tupleParams,
			)
			if err != nil {
				httpx.WriteJSON(
					w,
					http.StatusBadGateway,
					map[string]string{"error": err.Error()})
				return
			}
			if len(tuples.Tuples) > 0 {
				roles := make([]identity.IdentityRole,
					0,
					len(tuples.Tuples),
				)

				for _, t := range tuples.Tuples {
					if t.Relation == "" && t.Object == "" &&
						t.Namespace == "" {
						continue
					}
					roles = append(roles, identity.IdentityRole{
						Object: t.Object,
						Role:   t.Relation,
					})
				}
				item.Roles = roles
			}
		}

		out = append(out, item)
	}

	resp := map[string]any{"identities": out}
	if bffSessionPayload != nil {
		resp["bff_session"] = bffSessionPayload
	}
	httpx.WriteJSON(w, http.StatusOK, resp)
}
