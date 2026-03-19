package oauth_handler

import (
	"errors"
	"net/http"
	permission_domain "sso-bff/internal/domain/permission"
	"sso-bff/internal/httpx"
	"sso-bff/modules/oauth"
	"sso-bff/modules/permission"
	"strconv"
	"strings"
)

func (h *OAuthHandler) Session(
	w http.ResponseWriter,
	r *http.Request,
) {

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

		httpx.WriteJSON(
			w,
			http.StatusOK,
			map[string]interface{}{"authenticated": false},
		)

		return
	}

	session, err := h.OAuth.SessionService.GetSession(
		r.Context(),
		sessionID,
	)

	if err != nil {
		if errors.Is(err, oauth.ErrSessionNotFound) {
			httpx.WriteJSON(
				w,
				http.StatusOK,
				map[string]interface{}{"authenticated": false},
			)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type rolePayload struct {
		Object string `json:"object"`
		Role   string `json:"role"`
	}

	payload := oauth.SessionPayload{
		Authenticated: true,
		Sub:           session.Subject,
		Exp:           session.Expiry,
	}

	if h.Identity != nil {
		profileSource := "none"

		ident, err := h.Identity.WhoAmI(
			r.Context(),
			r.Header.Get("Cookie"),
		)

		if err == nil && ident != nil && len(ident.Traits) > 0 {
			profileSource = "cookie"
		} else {
			if session.KratosToken != "" {
				if alt, altErr := h.Identity.WhoAmIWithSessionToken(
					r.Context(),
					session.KratosToken,
				); altErr == nil {
					ident = alt
					if ident != nil && len(ident.Traits) > 0 {
						profileSource = "token"
					}
				}
			}
		}
		if ident != nil && len(ident.Traits) > 0 {
			profile := map[string]any{}
			if nameRaw, ok := ident.Traits["name"].(map[string]any); ok {
				name := map[string]string{}
				if v, ok := nameRaw["firstName"].(string); ok && v != "" {
					name["first_name"] = v
				}
				if v, ok := nameRaw["lastName"].(string); ok && v != "" {
					name["last_name"] = v
				}
				if len(name) > 0 {
					profile["name"] = name
				}
			}
			if len(profile) > 0 {
				payload.Profile = profile
			}
		}
		if profileSource != "none" {
			payload.ProfileSource = profileSource
		}
	}

	includeTuples := false
	if v := strings.TrimSpace(
		r.URL.Query().Get("include_tuples")); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			includeTuples = parsed
		} else {
			httpx.WriteJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "invalid include_tuples"},
			)
			return
		}
	}
	if includeTuples && h.Permission != nil {
		tupleParams := permission_domain.ListTuplesParams{
			Namespace: strings.TrimSpace(
				r.URL.Query().Get("tuples_namespace"),
			),
			Object: strings.TrimSpace(
				r.URL.Query().Get("tuples_object"),
			),
			Relation: strings.TrimSpace(
				r.URL.Query().Get("tuples_relation"),
			),
			SubjectID: session.Subject,
		}
		if tupleParams.Namespace == "" {
			httpx.WriteJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "tuples_namespace required"},
			)
			return
		}
		tuples, err := h.Permission.ListTuples(r.Context(), tupleParams)
		if err != nil {
			httpx.WriteJSON(
				w,
				http.StatusBadGateway,
				map[string]string{"error": err.Error()})
			return
		}
		if len(tuples.Tuples) > 0 {
			roles := make([]permission.RolePayload, 0, len(tuples.Tuples))
			for _, t := range tuples.Tuples {
				if t.Object == "" && t.Relation == "" {
					continue
				}
				roles = append(roles, permission.RolePayload{
					Object: t.Object,
					Role:   t.Relation,
				})
			}
			payload.Roles = roles
		}
	}
	httpx.WriteJSON(w, http.StatusOK, payload)
}

func (h *OAuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
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

	ttl, err := h.OAuth.SessionService.RefreshSession(
		r.Context(),
		sessionID,
	)

	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, oauth.ErrSessionNotFound),
			errors.Is(err, oauth.ErrInvalidIDToken):
			status = http.StatusUnauthorized
		}
		http.Error(w, err.Error(), status)
		return
	}

	if ttl > 0 {
		httpx.SetSessionCookie(w, h.Cookies, sessionID, ttl, r)
	}
	w.WriteHeader(http.StatusNoContent)
}
