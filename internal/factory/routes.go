package factory

import (
	"net/http"

	"sso-bff/internal/httpx"
	"sso-bff/internal/middleware"
	"sso-bff/modules/app"
	"sso-bff/modules/audit"
	identity_handler "sso-bff/modules/identity/handler"
	oauth_handler "sso-bff/modules/oauth/handler"
	permission_handler "sso-bff/modules/permission/handler"
)

const (
	healthPath      = "/health"
	discoveriesPath = "/discoveries"
)

type routeSpec struct {
	path    string
	methods []string
	handler http.HandlerFunc
}

type discoveryRoute struct {
	Path    string   `json:"path"`
	Methods []string `json:"methods"`
}

type httpHandlers struct {
	OAuth            *oauth_handler.OAuthHandler
	IdentityLogin    *identity_handler.IdentityLoginHandler
	IdentityAdmin    *identity_handler.IdentityAdminHandler
	IdentitySettings *identity_handler.IdentitySettingsHandler
	Permission       *permission_handler.PermissionHandler
	Audit            *audit.AuditHandler
	App              *app.AppHandler
}

type routeMiddleware struct {
	requireSession func(http.HandlerFunc) http.HandlerFunc
	requireAdmin   func(http.HandlerFunc) http.HandlerFunc
}

func newRouteMiddleware(auth *middleware.AuthMiddleware) routeMiddleware {
	guards := routeMiddleware{
		requireSession: identity,
		requireAdmin:   identity,
	}
	if auth == nil {
		return guards
	}
	guards.requireSession = auth.RequireSession
	guards.requireAdmin = auth.RequireAdmin
	return guards
}

func registerRoutes(
	mux *http.ServeMux,
	handlers httpHandlers,
	guards routeMiddleware,
) {
	routeSpecs := []routeSpec{
		{path: "/api/login", methods: []string{http.MethodGet}, handler: handlers.OAuth.Login},
		{path: "/api/callback", methods: []string{http.MethodGet}, handler: handlers.OAuth.Callback},
		{path: "/api/consent", methods: []string{http.MethodGet}, handler: handlers.OAuth.Consent},
		{path: "/api/launch", methods: []string{http.MethodGet}, handler: guards.requireSession(handlers.OAuth.LaunchApp)},
		{path: "/api/logout", methods: []string{http.MethodPost}, handler: verbHandler(http.MethodPost, handlers.OAuth.Logout)},
		{path: "/api/session", methods: []string{http.MethodGet}, handler: handlers.OAuth.Session},
		{path: "/api/session/refresh", methods: []string{http.MethodPost}, handler: verbHandler(http.MethodPost, handlers.OAuth.Refresh)},
		{path: "/api/identity/login", methods: []string{http.MethodPost}, handler: handlers.IdentityLogin.SubmitLogin},
		{path: "/api/identity/settings", methods: []string{http.MethodGet, http.MethodPost}, handler: guards.requireSession(identitySettingsHandler(handlers.IdentitySettings))},
		{path: "/api/admin/identities", methods: []string{http.MethodGet, http.MethodPost}, handler: guards.requireAdmin(identityAdminHandler(handlers.IdentityAdmin))},
		{path: "/api/permissions/tuple", methods: []string{http.MethodPost}, handler: guards.requireSession(handlers.Permission.WriteTuple)},
		{path: "/api/permissions/check", methods: []string{http.MethodGet}, handler: guards.requireSession(handlers.Permission.CheckTuple)},
		{path: "/api/permissions", methods: []string{http.MethodGet}, handler: guards.requireSession(handlers.Permission.ListTuples)},
		{path: "/api/audit/events", methods: []string{http.MethodGet}, handler: guards.requireSession(handlers.Audit.List)},
		{path: "/api/apps", methods: []string{http.MethodGet, http.MethodPost}, handler: guards.requireAdmin(appRootHandler(handlers.App))},
		{path: "/api/apps/", methods: []string{http.MethodGet, http.MethodPut, http.MethodDelete}, handler: guards.requireAdmin(appIDHandler(handlers.App))},
		{path: healthPath, methods: []string{http.MethodGet}, handler: healthCheckHandler()},
	}

	discovered := appendRouteSpecs(mux, routeSpecs)

	discoveryRoutes := append([]discoveryRoute(nil), discovered...)
	discoveryRoutes = append(discoveryRoutes, discoveryRoute{
		Path:    discoveriesPath,
		Methods: []string{http.MethodGet},
	})

	mux.HandleFunc(discoveriesPath, discoveriesHandler(discoveryRoutes))
}

func identity(handler http.HandlerFunc) http.HandlerFunc {
	return handler
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

func appendRouteSpecs(mux *http.ServeMux, specs []routeSpec) []discoveryRoute {
	routes := make([]discoveryRoute, 0, len(specs))
	for _, spec := range specs {
		mux.HandleFunc(spec.path, spec.handler)
		routes = append(routes, discoveryRoute{
			Path:    spec.path,
			Methods: copyMethods(spec.methods),
		})
	}
	return routes
}

func copyMethods(methods []string) []string {
	dst := make([]string, len(methods))
	copy(dst, methods)
	return dst
}

func healthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func discoveriesHandler(routes []discoveryRoute) http.HandlerFunc {
	captured := make([]discoveryRoute, len(routes))
	copy(captured, routes)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, map[string][]discoveryRoute{"routes": captured})
	}
}

func identitySettingsHandler(handler *identity_handler.IdentitySettingsHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.Start(w, r)
		case http.MethodPost:
			handler.Submit(w, r)
		default:
			methodNotAllowed(w)
		}
	}
}

func identityAdminHandler(handler *identity_handler.IdentityAdminHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.ListIdentities(w, r)
		case http.MethodPost:
			handler.CreateIdentity(w, r)
		default:
			methodNotAllowed(w)
		}
	}
}

func appRootHandler(handler *app.AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.List(w, r)
		case http.MethodPost:
			handler.Create(w, r)
		default:
			methodNotAllowed(w)
		}
	}
}

func appIDHandler(handler *app.AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.Get(w, r)
		case http.MethodPut:
			handler.Update(w, r)
		case http.MethodDelete:
			handler.Delete(w, r)
		default:
			methodNotAllowed(w)
		}
	}
}
