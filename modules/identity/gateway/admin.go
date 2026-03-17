package identity_gateway

import (
	"context"
	"errors"
	"fmt"
	identity_domain "sso-bff/internal/domain/identity"
	identity_helper "sso-bff/modules/identity/helper"

	identity "github.com/ory/kratos-client-go"
)

type IdentityAdminGateway struct {
	Admin *identity.APIClient
}

func NewIdentityAdminGateway(admin *identity.APIClient) *IdentityAdminGateway {
	return &IdentityAdminGateway{Admin: admin}
}

func (g *IdentityAdminGateway) CreateIdentity(
	ctx context.Context,
	traits map[string]any,
	schemaID string,
) (*identity_domain.Identity, error) {

	if err := identity_helper.CheckClient(g.Admin); err != nil {
		return nil, err
	}

	if len(traits) == 0 {
		return nil, errors.New("traits required")
	}

	if schemaID == "" {
		schemaID = "default"
	}

	body := identity.CreateIdentityBody{
		SchemaId: schemaID,
		Traits:   traits,
	}

	created, _, err := g.Admin.IdentityAPI.
		CreateIdentity(ctx).
		CreateIdentityBody(body).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("client create identity: %w", err)
	}

	id := created.GetId()
	if id == "" {
		return nil, errors.New("client create identity returned empty id")
	}

	return &identity_domain.Identity{
		ID:     id,
		Traits: identity_helper.ExtractTraits(created.GetTraits()),
	}, nil
}

func (g *IdentityAdminGateway) UpdateIdentityTraits(
	ctx context.Context,
	id string,
	traits map[string]any,
) error {

	if err := identity_helper.CheckClient(g.Admin); err != nil {
		return err
	}

	if id == "" {
		return errors.New("identity id required")
	}

	body := identity.UpdateIdentityBody{
		Traits: traits,
	}

	_, _, err := g.Admin.IdentityAPI.
		UpdateIdentity(ctx, id).
		UpdateIdentityBody(body).
		Execute()

	if err != nil {
		return fmt.Errorf("Identity update identity: %w", err)
	}

	return nil
}

func (g *IdentityAdminGateway) GetIdentity(
	ctx context.Context,
	id string,
) (*identity_domain.Identity, error) {

	if err := identity_helper.CheckClient(g.Admin); err != nil {
		return nil, err
	}

	if id == "" {
		return nil, errors.New("identity id required")
	}

	ident, _, err := g.Admin.IdentityAPI.
		GetIdentity(ctx, id).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("kratos get identity: %w", err)
	}

	return &identity_domain.Identity{
		ID:     ident.GetId(),
		Traits: identity_helper.ExtractTraits(ident.GetTraits()),
	}, nil
}
func (g *IdentityAdminGateway) ListIdentities(
	ctx context.Context,
	params identity_domain.ListIdentitiesParams,
) ([]identity_domain.Identity, error) {

	if err := identity_helper.CheckClient(g.Admin); err != nil {
		return nil, err
	}

	req := g.Admin.IdentityAPI.ListIdentities(ctx)

	if params.PerPage > 0 {
		req = req.PerPage(params.PerPage)
	}
	if params.Page > 0 {
		req = req.Page(params.Page)
	}
	if params.PageSize > 0 {
		req = req.PageSize(params.PageSize)
	}
	if params.PageToken != "" {
		req = req.PageToken(params.PageToken)
	}
	if params.CredentialsIdentifier != "" {
		req = req.CredentialsIdentifier(params.CredentialsIdentifier)
	}

	idents, _, err := req.Execute()
	if err != nil {
		return nil, fmt.Errorf("kratos list identities: %w", err)
	}

	out := make([]identity_domain.Identity, 0, len(idents))

	for _, ident := range idents {
		out = append(out, identity_domain.Identity{
			ID:     ident.GetId(),
			Traits: identity_helper.ExtractTraits(ident.GetTraits()),
		})
	}

	return out, nil
}

func (g *IdentityAdminGateway) ListSessions(
	ctx context.Context,
	params identity_domain.ListSessionsParams,
) ([]identity_domain.IdentitySession, error) {

	if err := identity_helper.CheckClient(g.Admin); err != nil {
		return nil, err
	}

	req := g.Admin.IdentityAPI.ListSessions(ctx)
	if params.PageSize > 0 {
		req = req.PageSize(params.PageSize)
	}
	if params.PageToken != "" {
		req = req.PageToken(params.PageToken)
	}
	if params.Active != nil {
		req = req.Active(*params.Active)
	}

	sessions, _, err := req.Execute()
	if err != nil {
		return nil, fmt.Errorf("kratos list sessions: %w", err)
	}

	out := make([]identity_domain.IdentitySession, 0, len(sessions))
	for _, sess := range sessions {
		identityID := ""
		if sess.Identity != nil {
			identityID = sess.Identity.GetId()
		}
		if identityID == "" {
			continue
		}

		devices := make([]identity_domain.IdentitySessionDevice, 0)
		for _, device := range sess.GetDevices() {
			ip := device.GetIpAddress()
			ua := device.GetUserAgent()
			if ip == "" && ua == "" {
				continue
			}
			devices = append(devices, identity_domain.IdentitySessionDevice{
				IPAddress: ip,
				UserAgent: ua,
			})
		}

		out = append(out, identity_domain.IdentitySession{
			IdentityID: identityID,
			Devices:    devices,
		})
	}

	return out, nil
}
