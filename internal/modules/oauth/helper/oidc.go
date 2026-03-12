package helper

import (
	"context"
	"sso-bff/internal/domain"
	"time"
)

func NewVerifier(jwksURL, issuer, audience string, cacheTTL time.Duration) *domain.OauthOIDCVerifier {
	return &domain.OauthOIDCVerifier{
		JWKSURL:  jwksURL,
		Issuer:   issuer,
		Audience: audience,
		Client:   nil,
		CacheTTL: cacheTTL,
	}
}

func VerifyToken(ctx context.Context, verifier *domain.OauthOIDCVerifier, rawIDToken string) (map[string]interface{}, error) {
	return verifier.Verify(ctx, rawIDToken)
}
