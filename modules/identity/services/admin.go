package identity_service

import (
	"context"
	"errors"
	identity_domain "sso-bff/internal/domain/identity"
	"sso-bff/modules/identity"
)

type IdentityAdminService struct {
	admin identity.IdentityAdminClient
}

func NewIdentityAdminService(
	admin identity.IdentityAdminClient,
) *IdentityAdminService {
	return &IdentityAdminService{admin: admin}
}

func (s *IdentityAdminService) CreateIdentity(ctx context.Context,
	traits map[string]any,
	schemaID string,
) (*identity_domain.Identity, error) {
	if s == nil || s.admin == nil {
		return nil, identity.ErrAdminMisconfigured
	}
	if len(traits) == 0 {
		return nil, errors.New("traits required")
	}
	return s.admin.CreateIdentity(ctx, traits, schemaID)
}

func (s *IdentityAdminService) ListIdentities(ctx context.Context,
	params identity_domain.ListIdentitiesParams,
) ([]identity_domain.Identity, error) {
	if s == nil || s.admin == nil {
		return nil, identity.ErrAdminMisconfigured
	}
	return s.admin.ListIdentities(ctx, params)
}

func (s *IdentityAdminService) ListSessions(ctx context.Context,
	params identity_domain.ListSessionsParams,
) ([]identity_domain.IdentitySession, error) {
	if s == nil || s.admin == nil {
		return nil, identity.ErrAdminMisconfigured
	}
	return s.admin.ListSessions(ctx, params)
}
