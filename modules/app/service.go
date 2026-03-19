package app

import (
	"context"
	"errors"
	"fmt"
	app_domain "sso-bff/internal/domain/app"
	"sso-bff/modules/oauth"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrServiceMisconfigured = errors.New("app registry service misconfigured")

type AppService struct {
	repo AppRepository

	appCache    map[string]oauth.AppRedirect
	appLoadedAt time.Time
	appTTL      time.Duration
	appMu       sync.RWMutex
}

func NewService(repo AppRepository) *AppService {
	return &AppService{repo: repo}
}

func (s *AppService) Create(
	ctx context.Context,
	dsn,
	redirectPath string,
) (app_domain.AppRegistry, error) {

	if s == nil || s.repo == nil {
		return app_domain.AppRegistry{}, ErrServiceMisconfigured
	}

	if dsn == "" || redirectPath == "" {
		return app_domain.AppRegistry{},
			errors.New("dsn and redirect_path required")

	}
	return s.repo.Create(ctx, dsn, redirectPath)
}

func (s *AppService) Update(
	ctx context.Context,
	id uuid.UUID,
	dsn,
	redirectPath string,
) (app_domain.AppRegistry, error) {

	if s == nil || s.repo == nil {
		return app_domain.AppRegistry{}, ErrServiceMisconfigured
	}

	if id == uuid.Nil || dsn == "" || redirectPath == "" {
		return app_domain.AppRegistry{},
			errors.New("id, dsn and redirect_path required")
	}

	return s.repo.Update(ctx, id, dsn, redirectPath)
}

func (s *AppService) Get(
	ctx context.Context,
	id uuid.UUID,
) (app_domain.AppRegistry, error) {
	if s == nil || s.repo == nil {
		return app_domain.AppRegistry{}, ErrServiceMisconfigured
	}
	if id == uuid.Nil {
		return app_domain.AppRegistry{}, errors.New("id required")
	}
	return s.repo.Get(ctx, id)
}

func (s *AppService) List(
	ctx context.Context,
) ([]app_domain.AppRegistry, error) {

	if s == nil || s.repo == nil {
		return nil, ErrServiceMisconfigured
	}

	return s.repo.List(ctx)
}

func (s *AppService) Delete(ctx context.Context, id uuid.UUID) error {
	if s == nil || s.repo == nil {
		return ErrServiceMisconfigured
	}
	if id == uuid.Nil {
		return errors.New("id required")
	}
	return s.repo.Delete(ctx, id)
}

func (s *AppService) ResolveRegistry(
	ctx context.Context,
) (map[string]oauth.AppRedirect, error) {

	if s == nil || s.repo == nil {
		return nil, ErrServiceMisconfigured
	}

	s.appMu.RLock()
	if s.appCache != nil &&
		(s.appTTL <= 0 || time.Since(s.appLoadedAt) < s.appTTL) {

		out := s.appCache
		s.appMu.RUnlock()
		return out, nil
	}
	s.appMu.RUnlock()

	apps, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	registry := make(map[string]oauth.AppRedirect, len(apps))

	for _, app := range apps {

		if app.ID == uuid.Nil || strings.TrimSpace(app.DSN) == "" {
			continue
		}

		allowed := map[string]struct{}{}

		if p := strings.TrimSpace(app.RedirectPath); p != "" {
			allowed[p] = struct{}{}
		}

		registry[app.ID.String()] = oauth.AppRedirect{
			BaseURL:      strings.TrimSpace(app.DSN),
			AllowedPaths: allowed,
		}
	}

	s.appMu.Lock()
	s.appCache = registry
	s.appLoadedAt = time.Now()
	s.appMu.Unlock()

	return registry, nil
}

func (s *AppService) ResolveAppIDByDSN(
	ctx context.Context,
	dsn string,
) (string, error) {
	if s == nil {
		return "", ErrServiceMisconfigured
	}

	dsn = NormalizeBaseURL(dsn)
	if dsn == "" {
		return "", errors.New("dsn required")
	}

	if s.repo == nil {
		return "", ErrServiceMisconfigured
	}

	registry, err := s.ResolveRegistry(ctx)
	if err != nil {
		return "", err
	}

	for id, appRedirect := range registry {
		if NormalizeBaseURL(appRedirect.BaseURL) == dsn {
			return id, nil
		}
	}

	return "", fmt.Errorf("unknown app dsn: %s", dsn)
}
