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

func (s *AcceptFlowService) OauthConsent(
	ctx context.Context,
	challenge string,
	identity *identity_domain.Identity,
) (string, error) {
	if s == nil || s.admin == nil {
		return "", oauth.ErrServiceMisconfigured
	}

	req, err := s.admin.GetConsentRequest(ctx, challenge)
	if err != nil {
		return "", fmt.Errorf("%w: %v", oauth.ErrHydraRequest, err)
	}

	if !oauth_helper_token.IsAllowedClient(req.ClientID, s.allowedClients) {
		return "", fmt.Errorf("%w: unauthorized client", oauth.ErrHydraRequest)
	}

	grantScopes := oauth_helper_token.FilterScopes(req.RequestedScope, s.allowedScopes)

	subject := req.Subject
	if identity != nil && identity.ID != "" {
		subject = identity.ID
	}

	claims := map[string]any{
		"sub": subject,
	}
	accessClaims := map[string]any{
		"sub": subject,
	}
	if identity != nil && len(identity.Traits) > 0 {
		accessClaims["traits"] = identity.Traits
	}

	redirect, err := s.admin.AcceptConsent(ctx, challenge, grantScopes, req.Audience, claims, accessClaims)
	if err != nil {
		return "", fmt.Errorf("%w: %v", oauth.ErrHydraRequest, err)
	}

	return redirect, nil
}
