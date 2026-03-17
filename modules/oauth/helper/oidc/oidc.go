package oauth_helper_oidc

import (
	"context"

	oauth_domain "sso-bff/internal/domain/oauth"
	"time"
)

func NewVerifier(
	jwksURL,
	issuer,
	audience string,
	cacheTTL time.Duration,
) *oauth_domain.OauthOIDCVerifier {
	return &oauth_domain.OauthOIDCVerifier{
		JWKSURL:  jwksURL,
		Issuer:   issuer,
		Audience: audience,
		Client:   nil,
		CacheTTL: cacheTTL,
	}
}

func VerifyToken(
	ctx context.Context,
	verifier *oauth_domain.OauthOIDCVerifier,
	rawIDToken string,
) (map[string]interface{}, error) {
	return verifier.Verify(ctx, rawIDToken)
}
