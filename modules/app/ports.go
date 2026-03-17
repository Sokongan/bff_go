package app

import (
	"context"
	app_domain "sso-bff/internal/domain/app"

	"github.com/google/uuid"
)

type AppRepository interface {
	Create(
		ctx context.Context,
		dsn,
		redirectPath string,
	) (app_domain.AppRegistry, error)

	Update(
		ctx context.Context,
		id uuid.UUID,
		dsn,
		redirectPath string,
	) (app_domain.AppRegistry, error)

	Get(
		ctx context.Context,
		id uuid.UUID,
	) (app_domain.AppRegistry, error)

	List(ctx context.Context) ([]app_domain.AppRegistry, error)

	Delete(ctx context.Context, id uuid.UUID) error
}
