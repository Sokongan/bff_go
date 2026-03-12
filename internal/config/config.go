package config

import (
	db_config "sso-bff/internal/config/db"
	services_config "sso-bff/internal/config/services"
	"sso-bff/internal/config/services/oauth"
)

type Config struct {
	DB         *db_config.DbConfig
	Store      *db_config.StoreConfig
	Oauth      *oauth.OAuthConfig
	Identity   *services_config.IdentityConfig
	Permission *services_config.PermissionConfig
}

func LoadConfig() (*Config, error) {
	db, err := db_config.LoadDBConfig()
	if err != nil {
		return nil, err
	}
	store, err := db_config.LoadRedisConfig()
	if err != nil {
		return nil, err
	}

	oauth, err := oauth.LoadOAuthConfig()
	if err != nil {
		return nil, err
	}
	identity, err := services_config.LoadIdentityConfig()
	if err != nil {
		return nil, err
	}
	permission, err := services_config.LoadPermissionConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		DB:         db,
		Store:      store,
		Oauth:      oauth,
		Identity:   identity,
		Permission: permission,
	}, nil
}
