package oauth

type OAuthConfig struct {
	URLs   *URLConfig
	Client *ClientConfig
	M2M    *M2MConfig
	Scopes *ClientScopesConfig
	OIDC   *OIDCConfig
	Cookie *CookieConfig
}

func LoadOAuthConfig() (*OAuthConfig, error) {
	urls, err := LoadURLConfig()
	if err != nil {
		return nil, err
	}

	client, err := LoadClientConfig() // capture both values
	if err != nil {
		return nil, err
	}

	m2m, err := LoadM2MConfig()
	if err != nil {
		return nil, err
	}

	scopes, err := LoadClientScopesConfig()
	if err != nil {
		return nil, err
	}

	oidc, err := LoadOIDCConfig()
	if err != nil {
		return nil, err
	}

	return &OAuthConfig{
		URLs:   urls,
		Client: client,
		M2M:    m2m,
		Scopes: scopes,
		OIDC:   oidc,
		Cookie: LoadCookieConfig(),
	}, nil
}
