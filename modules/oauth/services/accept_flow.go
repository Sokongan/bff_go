package service_oauth

import (
	"context"
	"fmt"
	identity_domain "sso-bff/internal/domain/identity"
	"sso-bff/modules/oauth"
	oauth_helper_token "sso-bff/modules/oauth/helper/token"
)

type AcceptFlowService struct {
	admin          oauth.OauthAdminPort
	allowedClients map[string]struct{}
	allowedScopes  map[string]struct{}
}

func NewAcceptFlowService(
	admin oauth.OauthAdminPort,
	allowedClients map[string]struct{},
	allowedScopes map[string]struct{},
) *AcceptFlowService {
	return &AcceptFlowService{
		admin:          admin,
		allowedClients: allowedClients,
		allowedScopes:  allowedScopes,
	}
}

func (a *AcceptFlowService) Consent(
	ctx context.Context,
	challenge string,
	identity *identity_domain.Identity,
) (string, error) {

	req, err := a.admin.GetConsentRequest(ctx, challenge)
	if err != nil {
		return "", fmt.Errorf("consent request: %w", err)
	}

	if !oauth_helper_token.IsAllowedClient(req.ClientID, a.allowedClients) {
		return "", fmt.Errorf("unauthorized client: %s", req.ClientID)
	}

	grantScopes := oauth_helper_token.FilterScopes(req.RequestedScope, a.allowedScopes)

	subject := req.Subject
	if identity != nil && identity.ID != "" {
		subject = identity.ID
	}

	idClaims := oauth_helper_token.BuildClaims(subject)
	accessClaims := oauth_helper_token.BuildAccessClaims(identity, subject)

	return a.admin.AcceptConsent(
		ctx,
		challenge,
		grantScopes,
		req.Audience,
		idClaims,
		accessClaims,
	)
}
