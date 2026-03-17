package service_oauth

import (
	"context"
	"fmt"

	"sso-bff/modules/app"
	"sso-bff/modules/oauth"

	oauth_helper_redirect "sso-bff/modules/oauth/helper/redirect"
	oauth_helper_token "sso-bff/modules/oauth/helper/token"
)

type RedirectService struct {
	redirects oauth.RedirectStorePort
	registry  *app.AppService
}

func NewRedirectService(
	redirects oauth.RedirectStorePort,
	registry *app.AppService,
) *RedirectService {
	return &RedirectService{
		redirects: redirects,
		registry:  registry,
	}
}

func (s *RedirectService) ConsumeRedirect(
	ctx context.Context,
	state string,
) (string, string, error) {
	if s == nil || s.redirects == nil {
		return "", "", nil
	}

	redirect, err := s.redirects.GetRedirect(ctx, state)
	if err != nil {
		return "", "", err
	}

	if err := s.redirects.DeleteRedirect(ctx, state); err != nil {
		return "", "", err
	}

	if redirect == "" {
		return "", "", nil
	}

	appID, redirectPath, err := oauth_helper_redirect.
		DecodeRedirectPayload(redirect)
	if err != nil {
		return "", "", err
	}

	return appID, redirectPath, nil
}

func (s *RedirectService) RedirectForToken(
	ctx context.Context,
	state, accessToken string,
) (string, error) {
	if s == nil || s.redirects == nil || s.registry == nil {
		return "", oauth.ErrServiceMisconfigured
	}

	appID, redirectPath, err := s.ConsumeRedirect(ctx, state)
	if err != nil {
		return "", err
	}

	if redirectPath == "" {
		redirectPath = "/"
	}

	if appID == "" {
		clientID, err := oauth_helper_token.
			ClientIDFromAccessToken(accessToken)

		if err != nil {
			return "", err
		}
		appID = clientID
	}

	registry, err := s.registry.ResolveRegistry(ctx)
	if err != nil {
		return "", err
	}
	if len(registry) == 0 {
		return "", oauth.ErrServiceMisconfigured
	}

	app, ok := registry[appID]
	if !ok {
		return "", fmt.Errorf("unknown app_id: %s", appID)
	}

	return oauth_helper_redirect.BuildAppRedirect(app, redirectPath)
}
