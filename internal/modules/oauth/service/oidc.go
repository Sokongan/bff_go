package oauth_service_oidc

import (
	"context"
	"sso-bff/internal/domain"
	"sso-bff/internal/oauth/gateway"
)

type OIDCService struct {
	verifier *gateway.OIDCVerifierGateway
}

func NewOIDCService(verifier *gateway.OIDCVerifierGateway) *OIDCService {
	return &OIDCService{verifier: verifier}
}

func (o *OIDCService) VerifyIDToken(ctx context.Context, rawIDToken string) (domain.IDTokenClaims, error) {
	return o.verifier.Verify(ctx, rawIDToken)
}
