package app

import (
	"net/url"
	dbgen "sso-bff/internal/db/gen"
	app_domain "sso-bff/internal/domain/app"
	"strings"
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

func NormalizeBaseURL(raw string) string {
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return strings.TrimRight(strings.TrimSpace(raw), "/")
	}

	if parsed.Scheme == "" {
		return strings.TrimRight(strings.TrimSpace(raw), "/")
	}

	if parsed.Host != "" {
		return strings.TrimRight(parsed.Host, "/")
	}

	return strings.TrimRight(strings.TrimSpace(raw), "/")
}
