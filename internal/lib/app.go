package lib

import (
	"sso-bff/internal/config"
	"sso-bff/internal/db"

	internal "sso-bff/modules"
)

type App struct {
	Config    *config.Config
	Resources *db.Resources
	SDK       *internal.SDKs
}
