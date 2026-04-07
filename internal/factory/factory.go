package factory

import (
	"context"
	"errors"
	"net/http"
	"time"

	"sso-bff/internal/config"
	oauthcfg "sso-bff/internal/config/services/oauth"
	"sso-bff/internal/db"
	permission_domain "sso-bff/internal/domain/permission"
	"sso-bff/internal/httpx"
	"sso-bff/internal/middleware"
	"sso-bff/modules"
	"sso-bff/modules/app"
	app_adapter "sso-bff/modules/app/adapter"
	"sso-bff/modules/audit"
	audit_adapter "sso-bff/modules/audit/adapter"
	identity_factory "sso-bff/modules/identity/factory"
	identity_factory_modules "sso-bff/modules/identity/factory/modules"
	oauth_factory "sso-bff/modules/oauth/factory"
	oauth_gateway "sso-bff/modules/oauth/gateway"
	"sso-bff/modules/permission"
	permission_factory "sso-bff/modules/permission/factory"
	"sso-bff/modules/store"
)

const (
	pkceTTL = 10 * time.Minute
)

type Module struct {
	Handler http.Handler
}

func NewHandlers(cfg *config.Config, resources *db.Resources, sdks *modules.SDKs) (*Module, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}
	if resources == nil || resources.Db == nil || resources.Store == nil {
		return nil, errors.New("resources unavailable")
	}
	if sdks == nil {
		return nil, errors.New("sdk bundle is nil")
	}
	if cfg.Oauth == nil || cfg.Oauth.Cookie == nil || cfg.Oauth.Scopes == nil ||
		cfg.Oauth.Client == nil || cfg.Oauth.OIDC == nil {
		return nil, errors.New("oauth configuration incomplete")
	}

	sqlDB := resources.Db.StdlibDB()
	if sqlDB == nil {
		return nil, errors.New("sql db is nil")
	}

	if resources.Store.Client() == nil {
		return nil, errors.New("redis client unavailable")
	}

	stores := store.NewStore(resources.Store.Client())
	cookieCfg := cookieConfigFrom(cfg.Oauth.Cookie)

	appService := app.NewService(app_adapter.New(sqlDB))
	auditService := audit.NewAuditService(audit_adapter.New(sqlDB))

	identityGateways, err := identity_factory_modules.NewIdentityGateways(sdks.Identity)
	if err != nil {
		return nil, err
	}

	oauthAdmin := oauth_gateway.NewOauthAdminGateway(sdks.OAuth.Admin)
	identityServices := identity_factory_modules.NewIdentityServices(
		identityGateways.Admin,
		identityGateways.Browser,
		oauthAdmin,
	)

	idVerifier, err := newIDTokenVerifier(cfg)
	if err != nil {
		return nil, err
	}

	tupleProxy := &tupleListerProxy{}

	oauthModule, err := oauth_factory.NewOauthModule(
		sdks.OAuth,
		stores.Sessions,
		stores.PKCE,
		stores.Redirect,
		idVerifier,
		appService,
		cfg.Oauth.Cookie.TTL,
		pkceTTL,
		cfg.Oauth.Scopes.AllowedClient,
		cfg.Oauth.Scopes.AllowedScope,
		identityServices.GetIdentity,
		auditService,
		tupleProxy,
		cookieCfg,
	)
	if err != nil {
		return nil, err
	}

	sessionService := oauthModule.Service.SessionService
	if sessionService == nil {
		return nil, errors.New("session service unavailable")
	}

	permissionModule, err := permission_factory.NewPermissionModule(
		sdks.Permission,
		auditService,
	)
	if err != nil {
		return nil, err
	}
	tupleProxy.target = permissionModule.Service

	identityModule, err := identity_factory.NewModule(
		sdks.Identity,
		oauthAdmin,
		permissionModule.Service,
		auditService,
		cookieCfg,
	)
	if err != nil {
		return nil, err
	}

	auditHandler := audit.NewAuditHandler(
		auditService,
	)

	appHandler := app.NewAppHandler(
		appService,
		auditService,
	)

	authMiddleware := middleware.NewAuthMiddleware(
		sessionService,
		permissionModule.Service,
		cookieCfg,
	)

	mux := http.NewServeMux()
	registerRoutes(mux, httpHandlers{
		OAuth:            oauthModule.Handler,
		IdentityLogin:    identityModule.Login,
		IdentityAdmin:    identityModule.Admin,
		IdentitySettings: identityModule.Settings,
		Permission:       permissionModule.Handler,
		Audit:            auditHandler,
		App:              appHandler,
	}, newRouteMiddleware(authMiddleware))

	return &Module{Handler: mux}, nil
}

func cookieConfigFrom(cfg *oauthcfg.CookieConfig) httpx.CookieConfig {
	if cfg == nil {
		return httpx.CookieConfig{}
	}
	return httpx.CookieConfig{
		Name:     cfg.Name,
		Domain:   cfg.Domain,
		SameSite: cfg.SameSite,
		Secure:   cfg.Secure,
	}
}

type tupleListerProxy struct {
	target permission.TupleLister
}

func (p *tupleListerProxy) ListTuples(ctx context.Context, params permission_domain.ListTuplesParams) (permission_domain.ListTuplesResult, error) {
	if p == nil || p.target == nil {
		return permission_domain.ListTuplesResult{}, permission.ErrMisconfig
	}
	return p.target.ListTuples(ctx, params)
}
