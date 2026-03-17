package modules

import (
	"sso-bff/internal/config"
	identity_sdk "sso-bff/modules/identity/sdk"
	"sso-bff/modules/oauth"
	oauth_sdk "sso-bff/modules/oauth/sdk"
	permission_sdk "sso-bff/modules/permission/sdk"
)

type SDKs struct {
	OAuth      *oauth_sdk.OAuthSDK
	Identity   *identity_sdk.IdentitySDK
	Permission *permission_sdk.PermissionSDK
}

func NewSDKs(cfg *config.Config) *SDKs {

	browser := oauth.BrowserClient{
		BrowserPublicURL: cfg.Oauth.URLs.PublicURL,
		ClientID:         cfg.Oauth.Client.ClientID,
		ClientSecret:     cfg.Oauth.Client.ClientSecret,
		RedirectURL:      cfg.Oauth.Client.RedirectURL,
		Scopes:           cfg.Oauth.Scopes.BFFScopes,
	}

	internal := oauth.InternalClient{
		TokenURL:     cfg.Oauth.URLs.PrivateURL,
		ClientID:     cfg.Oauth.Client.ClientID,
		ClientSecret: cfg.Oauth.Client.ClientSecret,
		RedirectURL:  cfg.Oauth.Client.RedirectURL,
		Scopes:       cfg.Oauth.Scopes.BFFScopes,
	}

	m2m := oauth.M2MClient{
		TokenURL:  cfg.Oauth.URLs.PublicURL,
		M2MID:     cfg.Oauth.M2M.M2MID,
		M2MSecret: cfg.Oauth.M2M.M2MSecret,
		Scopes:    cfg.Oauth.Scopes.M2MScopes,
	}

	return &SDKs{
		OAuth: oauth_sdk.NewOAuthSDK(
			cfg.Oauth.URLs.AdminURL,
			browser,
			internal,
			m2m,
		),

		Identity: identity_sdk.NewIdentitySDK(
			cfg.Identity.PublicURL,
			cfg.Identity.AdminURL,
		),

		Permission: permission_sdk.NewPermissionSDK(
			cfg.Permission.AdminURL,
			cfg.Permission.PublicURL,
		),
	}
}
