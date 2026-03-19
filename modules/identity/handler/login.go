package identity_handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sso-bff/internal/domain"
	audit_domain "sso-bff/internal/domain/audit"
	"sso-bff/internal/httpx"
	"sso-bff/modules/audit"
	"sso-bff/modules/identity"
	identity_factory_modules "sso-bff/modules/identity/factory/modules"
	"strings"
)

type HTTPError struct {
	Status int
	Body   []byte
	Err    error
}

func (e *HTTPError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "http error"
}

type IdentityLoginHandler struct {
	Flow    *identity_factory_modules.IdentityServices
	Audit   audit.AuditWriter
	Cookies httpx.CookieConfig
}

func NewIdentityFlowHandler(
	login *identity_factory_modules.IdentityServices,
	audit audit.AuditWriter,
	cookies httpx.CookieConfig,
) *IdentityLoginHandler {
	return &IdentityLoginHandler{
		Flow:    login,
		Audit:   audit,
		Cookies: cookies,
	}
}

func (h *IdentityLoginHandler) SubmitLogin(
	w http.ResponseWriter,
	r *http.Request,
) {
	if h == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "login handler unavailable"},
		)
		return
	}
	if h.Flow == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "login service unavailable"},
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

	var req identity.SubmitLoginRequest
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(
		contentType,
		"application/x-www-form-urlencoded") ||
		strings.HasPrefix(
			contentType,
			"multipart/form-data",
		) {
		if err := r.ParseForm(); err != nil {
			httpx.WriteJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "invalid request body"},
			)
			return
		}
		req.Identifier = r.FormValue("identifier")
		req.Password = r.FormValue("password")
		req.LoginChallenge = r.FormValue("login_challenge")
	} else {
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
		if err := dec.Decode(&struct{}{}); err != io.EOF {
			httpx.WriteJSON(
				w,
				http.StatusBadRequest,
				map[string]string{"error": "invalid request body"},
			)
			return
		}
	}

	redirectTo, identityID, sessionToken, err := h.
		Flow.Login.AuthenticatePassword(
		r.Context(),
		req.Identifier,
		req.Password,
		req.LoginChallenge,
	)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	if h.Audit != nil && identityID != "" {
		_ = h.Audit.Insert(r.Context(), audit_domain.AuditEvent{
			IdentityID: identityID,
			EventType:  "kratos_login",
			IPAddress:  httpx.ClientIP(r),
			UserAgent:  r.UserAgent(),
		})
	}

	if sessionToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     httpx.KratosSessionTokenCookie,
			Value:    sessionToken,
			Path:     "/",
			Domain:   h.Cookies.Domain,
			HttpOnly: true,
			Secure:   httpx.Secure(h.Cookies, r),
			SameSite: httpx.SameSite(h.Cookies),
		})
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		identity.SubmitLoginResponse{RedirectTo: redirectTo},
	)
}

func writeServiceError(w http.ResponseWriter, err error) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		if len(httpErr.Body) > 0 {
			if message, ok := kratosFlowMessage(httpErr.Body); ok {
				httpx.WriteJSON(w, httpErr.Status, map[string]string{"error": message})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(httpErr.Status)
			_, _ = w.Write(httpErr.Body)
			return
		}
		http.Error(w, err.Error(), httpErr.Status)
		return
	}

	switch {
	case errors.Is(err, domain.ErrMissingLoginChallenge),
		errors.Is(err, domain.ErrMissingLoginIdentifier),
		errors.Is(err, domain.ErrMissingLoginIdentity):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func kratosFlowMessage(body []byte) (string, bool) {
	var payload struct {
		UI struct {
			Messages []struct {
				Text string `json:"text"`
			} `json:"messages"`
		} `json:"ui"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", false
	}
	if len(payload.UI.Messages) == 0 {
		return "", false
	}
	return payload.UI.Messages[0].Text,
		payload.UI.Messages[0].Text != ""
}
