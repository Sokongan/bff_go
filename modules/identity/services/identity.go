package identity_service

import (
	"context"
	"errors"
	identity_domain "sso-bff/internal/domain/identity"
	"sso-bff/modules/identity"
)

type GetIdentityService struct {
	client identity.IdentityClient
}

func NewGetIdentityService(client identity.IdentityClient) *GetIdentityService {
	return &GetIdentityService{client: client}
}

func (s *GetIdentityService) WhoAmI(ctx context.Context,
	cookieHeader string,
) (*identity_domain.Identity, error) {
	if s == nil || s.client == nil {
		return nil, identity.ErrIdentityMisconfigured
	}
	if cookieHeader == "" {
		return nil, identity.ErrMissingCookieHeader
	}
	return s.client.WhoAmI(ctx, cookieHeader)
}

func (s *GetIdentityService) WhoAmIWithSessionToken(ctx context.Context,
	sessionToken string,
) (*identity_domain.Identity, error) {
	if s == nil || s.client == nil {
		return nil, identity.ErrIdentityMisconfigured
	}
	if sessionToken == "" {
		return nil, errors.New("missing session token")
	}
	return s.client.WhoAmIWithSessionToken(ctx, sessionToken)
}
