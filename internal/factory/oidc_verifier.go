package factory

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"sso-bff/internal/config"
	oauth_domain "sso-bff/internal/domain/oauth"
	"sso-bff/modules/oauth"
	oauth_helper_oidc "sso-bff/modules/oauth/helper/oidc"
)

const oidcCacheTTL = 5 * time.Minute

type idTokenVerifier struct {
	inner *oauth_domain.OauthOIDCVerifier
}

func newIDTokenVerifier(cfg *config.Config) (oauth.IDTokenVerifierPort, error) {
	if cfg == nil || cfg.Oauth == nil || cfg.Oauth.OIDC == nil || cfg.Oauth.Client == nil {
		return nil, errors.New("oauth oidc configuration missing")
	}

	verifier := oauth_helper_oidc.NewVerifier(
		cfg.Oauth.OIDC.JWKS,
		cfg.Oauth.OIDC.Issuer,
		cfg.Oauth.Client.ClientID,
		oidcCacheTTL,
	)

	return &idTokenVerifier{inner: verifier}, nil
}

func (v *idTokenVerifier) Verify(ctx context.Context, rawIDToken string) (oauth_domain.IDTokenClaims, error) {
	if v == nil || v.inner == nil {
		return oauth_domain.IDTokenClaims{}, errors.New("id token verifier not configured")
	}
	claims, err := oauth_helper_oidc.VerifyToken(ctx, v.inner, rawIDToken)
	if err != nil {
		return oauth_domain.IDTokenClaims{}, err
	}

	subject, ok := claims["sub"].(string)
	if !ok || subject == "" {
		return oauth_domain.IDTokenClaims{}, errors.New("id token missing subject")
	}

	exp, err := parseExpiry(claims["exp"])
	if err != nil {
		return oauth_domain.IDTokenClaims{}, err
	}

	return oauth_domain.IDTokenClaims{
		Subject:   subject,
		ExpiresAt: exp,
	}, nil
}

func parseExpiry(raw any) (time.Time, error) {
	switch v := raw.(type) {
	case float64:
		return time.Unix(int64(v), 0), nil
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(i, 0), nil
	case int64:
		return time.Unix(v, 0), nil
	case int:
		return time.Unix(int64(v), 0), nil
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return time.Unix(i, 0), nil
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return time.Unix(int64(f), 0), nil
		}
	}
	return time.Time{}, errors.New("invalid exp claim")
}
