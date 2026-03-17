package postgres

import (
	"context"
	"database/sql"
	dbgen "sso-bff/internal/db/gen"
	app_domain "sso-bff/internal/domain/app"
	"sso-bff/modules/app"

	"github.com/google/uuid"
)

type Repo struct {
	q *dbgen.Queries
}

func New(db *sql.DB) *Repo {
	return &Repo{q: dbgen.New(db)}
}

func (r *Repo) Create(
	ctx context.Context,
	dsn,
	redirectPath string,
) (app_domain.AppRegistry, error) {

	row, err := r.q.CreateAppRegistry(
		ctx,
		dbgen.CreateAppRegistryParams{
			Dsn:          dsn,
			RedirectPath: redirectPath,
		})

	if err != nil {
		return app_domain.AppRegistry{}, err
	}
	return app.ToDomain(row), nil
}

func (r *Repo) Update(
	ctx context.Context,
	id uuid.UUID,
	dsn,
	redirectPath string,
) (app_domain.AppRegistry, error) {

	row, err := r.q.UpdateAppRegistry(
		ctx,
		dbgen.UpdateAppRegistryParams{
			ID:           id,
			Dsn:          dsn,
			RedirectPath: redirectPath,
		})

	if err != nil {
		return app_domain.AppRegistry{}, err
	}
	return app.ToDomain(row), nil
}

func (r *Repo) Get(
	ctx context.Context,
	id uuid.UUID,
) (app_domain.AppRegistry, error) {

	row, err := r.q.GetAppRegistry(ctx, id)

	if err != nil {
		return app_domain.AppRegistry{}, err
	}

	return app.ToDomain(row), nil
}

func (r *Repo) List(ctx context.Context) (
	[]app_domain.AppRegistry, error,
) {

	rows, err := r.q.ListAppRegistries(ctx)

	if err != nil {
		return nil, err
	}

	out := make([]app_domain.AppRegistry, 0, len(rows))

	for _, row := range rows {
		out = append(out, app.ToDomain(row))
	}

	return out, nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteAppRegistry(ctx, id)
}
