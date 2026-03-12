package oauth

import (
	"errors"
	"fmt"
	"os"
)

type OIDCConfig struct {
	Issuer string
	JWKS   string
}

func LoadOIDCConfig() (*OIDCConfig, error) {
	issuer := os.Getenv("OIDC_ISSUER")

	if issuer == "" {
		return nil, errors.New("OIDC issuer configuration missing. Set OIDC_ISSUER")
	}

	return &OIDCConfig{
		Issuer: issuer,
		JWKS:   fmt.Sprintf("%s/.well-known/jwks.json", issuer),
	}, nil
}
