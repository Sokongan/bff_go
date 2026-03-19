package factory

import (
	"net/http"

	"sso-bff/internal/httpx"
	"sso-bff/modules/app"
	"sso-bff/modules/audit"
	identity_handler "sso-bff/modules/identity/handler"
	oauth_handler "sso-bff/modules/oauth/handler"
	permission_handler "sso-bff/modules/permission/handler"
)

type httpHandlers struct {
	OAuth            *oauth_handler.OAuthHandler
	IdentityLogin    *identity_handler.IdentityLoginHandler
	IdentityAdmin    *identity_handler.IdentityAdminHandler
	IdentitySettings *identity_handler.IdentitySettingsHandler
	Permission       *permission_handler.PermissionHandler
	Audit            *audit.AuditHandler
	App              *app.AppHandler
}

func registerRoutes(mux *http.ServeMux, handlers httpHandlers) {
	mux.HandleFunc("/api/login", handlers.OAuth.Login)
	mux.HandleFunc("/api/callback", handlers.OAuth.Callback)
	mux.HandleFunc("/api/consent", handlers.OAuth.Consent)
	mux.HandleFunc("/api/launch", handlers.OAuth.LaunchApp)
	mux.HandleFunc("/api/logout", verbHandler(http.MethodPost, handlers.OAuth.Logout))
	mux.HandleFunc("/api/session", handlers.OAuth.Session)
	mux.HandleFunc("/api/session/refresh", verbHandler(http.MethodPost, handlers.OAuth.Refresh))

	mux.HandleFunc("/api/identity/login", handlers.IdentityLogin.SubmitLogin)
	mux.HandleFunc("/api/identity/settings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.IdentitySettings.Start(w, r)
		case http.MethodPost:
			handlers.IdentitySettings.Submit(w, r)
		default:
			methodNotAllowed(w)
		}
	})

	mux.HandleFunc("/api/admin/identities", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.IdentityAdmin.ListIdentities(w, r)
		case http.MethodPost:
			handlers.IdentityAdmin.CreateIdentity(w, r)
		default:
			methodNotAllowed(w)
		}
	})

	mux.HandleFunc("/api/permissions/tuple", handlers.Permission.WriteTuple)
	mux.HandleFunc("/api/permissions/check", handlers.Permission.CheckTuple)
	mux.HandleFunc("/api/permissions", handlers.Permission.ListTuples)
	mux.HandleFunc("/api/audit/events", handlers.Audit.List)

	mux.HandleFunc("/api/apps", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.App.List(w, r)
		case http.MethodPost:
			handlers.App.Create(w, r)
		default:
			methodNotAllowed(w)
		}
	})
	mux.HandleFunc("/api/apps/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.App.Get(w, r)
		case http.MethodPut:
			handlers.App.Update(w, r)
		case http.MethodDelete:
			handlers.App.Delete(w, r)
		default:
			methodNotAllowed(w)
		}
	})
}

func verbHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			methodNotAllowed(w)
			return
		}
		handler(w, r)
	}
}

func methodNotAllowed(w http.ResponseWriter) {
	httpx.WriteJSON(
		w,
		http.StatusMethodNotAllowed,
		map[string]string{"error": "method not allowed"},
	)
}
