package identity_handler

import (
	"encoding/json"
	"io"
	"net/http"
	"sso-bff/internal/httpx"
	"sso-bff/internal/middleware"
	identity_service "sso-bff/modules/identity/services"

	kratos "github.com/ory/kratos-client-go"
)

type IdentitySettingsHandler struct {
	Settings *identity_service.IdentitySettingsService
}

func NewIdentitySettingsHandler(settings *identity_service.IdentitySettingsService) *IdentitySettingsHandler {
	return &IdentitySettingsHandler{
		Settings: settings,
	}
}

func (h *IdentitySettingsHandler) Start(
	w http.ResponseWriter,
	r *http.Request,
) {
	if h == nil || h.Settings == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "settings service unavailable"},
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

	session, ok := middleware.SessionFromContext(r.Context())
	if !ok {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "request session unavailable"},
		)
		return
	}

	if session.KratosToken == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "missing kratos session"},
		)
		return
	}

	flow, err := h.Settings.CreateBrowserSettingsFlow(
		r.Context(),
		session.KratosToken,
	)

	if err != nil {
		writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, flow)
}

func (h *IdentitySettingsHandler) Submit(
	w http.ResponseWriter,
	r *http.Request,
) {
	if h == nil || h.Settings == nil {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "settings service unavailable"},
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

	flowID := r.URL.Query().Get("flow")
	if flowID == "" {
		httpx.WriteJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "missing flow"},
		)
		return
	}

	var body kratos.UpdateSettingsFlowBody
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&body); err != nil {
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

	session, ok := middleware.SessionFromContext(r.Context())
	if !ok {
		httpx.WriteJSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": "request session unavailable"},
		)
		return
	}

	if session.KratosToken == "" {
		httpx.WriteJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "missing kratos session"},
		)
		return
	}

	flow, err := h.Settings.SubmitSettingsFlow(
		r.Context(),
		flowID,
		body,
		session.KratosToken,
	)

	if err != nil {
		writeServiceError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, flow)
}
