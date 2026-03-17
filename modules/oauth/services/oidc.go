package service_oauth

import (
	"context"
	oauth_domain "sso-bff/internal/domain/oauth"
	"sso-bff/modules/oauth"
)

type OIDCService struct {
	verifier oauth.IDTokenVerifierPort
}

func NewOIDCService(verifier oauth.IDTokenVerifierPort) *OIDCService {
	return &OIDCService{verifier: verifier}
}

func (o *OIDCService) VerifyIDToken(
	ctx context.Context, rawIDToken string) (oauth_domain.IDTokenClaims, error) {
	return o.verifier.Verify(ctx, rawIDToken)
}
