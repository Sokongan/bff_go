package app

import (
	dbgen "sso-bff/internal/db/gen"
	app_domain "sso-bff/internal/domain/app"
)

func ToDomain(row dbgen.AppRegistry) app_domain.AppRegistry {
	return app_domain.AppRegistry{
		ID:           row.ID,
		DSN:          row.Dsn,
		RedirectPath: row.RedirectPath,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
