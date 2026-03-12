package lib

import (
	"sso-bff/internal"
	"sso-bff/internal/config"
	"sso-bff/internal/db"
)

type App struct {
	Config    *config.Config
	Resources *db.Resources
	SDK       *internal.SDKs
}
