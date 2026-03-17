package identity_gateway

import (
	"context"
	"errors"
	"net/http"
	identity_domain "sso-bff/internal/domain/identity"
	identity_helper "sso-bff/modules/identity/helper"

	client "github.com/ory/kratos-client-go"
)

type IdentityBrowserGateway struct {
	Browser *client.APIClient
}

func NewIdentityBrowserGateway(
	Browser *client.APIClient,
) *IdentityBrowserGateway {

	return &IdentityBrowserGateway{Browser: Browser}
}

func (g *IdentityBrowserGateway) WhoAmI(
	ctx context.Context,
	cookieHeader string,
) (*identity_domain.Identity, error) {
	if err := identity_helper.CheckClient(g.Browser); err != nil {
		return nil, err
	}

	session, _, err := g.Browser.FrontendAPI.
		ToSession(ctx).
		Cookie(cookieHeader).
		Execute()
	if err != nil {
		return nil, err
	}

	identity := session.GetIdentity()
	id := identity.GetId()
	if id == "" {
		return nil, errors.New("Identity whoami missing id")
	}

	return &identity_domain.Identity{
		ID:     id,
		Traits: identity_helper.ExtractTraits(identity.GetTraits()),
	}, nil
}

func (g *IdentityBrowserGateway) WhoAmIWithSessionToken(
	ctx context.Context,
	sessionToken string,
) (*identity_domain.Identity, error) {

	if err := identity_helper.CheckClient(g.Browser); err != nil {
		return nil, err
	}
	if sessionToken == "" {
		return nil, errors.New("Identity session token required")
	}

	session, _, err := g.Browser.FrontendAPI.
		ToSession(ctx).
		XSessionToken(sessionToken).
		Execute()
	if err != nil {
		return nil, err
	}

	identity := session.GetIdentity()
	id := identity.GetId()
	if id == "" {
		return nil, errors.New("Identity whoami missing identity id")
	}

	return &identity_domain.Identity{
		ID:     id,
		Traits: identity_helper.ExtractTraits(identity.GetTraits()),
	}, nil
}

func (g *IdentityBrowserGateway) CreateNativeLoginFlow(
	ctx context.Context,
) (*client.LoginFlow, *http.Response, error) {

	if err := identity_helper.CheckClient(g.Browser); err != nil {
		return nil, nil, err
	}

	return g.Browser.FrontendAPI.CreateNativeLoginFlow(ctx).Execute()
}

func (g *IdentityBrowserGateway) CreateBrowserSettingsFlow(
	ctx context.Context,
	cookieHeader string,
) (*client.SettingsFlow, *http.Response, error) {

	if err := identity_helper.CheckClient(g.Browser); err != nil {
		return nil, nil, err
	}

	return g.Browser.FrontendAPI.
		CreateBrowserSettingsFlow(ctx).
		Cookie(cookieHeader).
		Execute()
}

func (g *IdentityBrowserGateway) CreateNativeSettingsFlow(
	ctx context.Context,
	sessionToken string,
) (*client.SettingsFlow, *http.Response, error) {
	if err := identity_helper.CheckClient(g.Browser); err != nil {
		return nil, nil, err
	}

	return g.Browser.FrontendAPI.
		CreateNativeSettingsFlow(ctx).
		XSessionToken(sessionToken).
		Execute()
}

func (g *IdentityBrowserGateway) UpdateLoginFlow(
	ctx context.Context,
	flowID string,
	body client.UpdateLoginFlowBody,
) (*client.SuccessfulNativeLogin, *http.Response, error) {
	if err := identity_helper.CheckClient(g.Browser); err != nil {
		return nil, nil, err
	}

	return g.Browser.FrontendAPI.
		UpdateLoginFlow(ctx).
		Flow(flowID).
		UpdateLoginFlowBody(body).
		Execute()
}

func (g *IdentityBrowserGateway) UpdateSettingsFlow(
	ctx context.Context,
	flowID string,
	body client.UpdateSettingsFlowBody,
	cookieHeader string,
) (*client.SettingsFlow, *http.Response, error) {

	if err := identity_helper.CheckClient(g.Browser); err != nil {
		return nil, nil, err
	}

	return g.Browser.FrontendAPI.
		UpdateSettingsFlow(ctx).
		Flow(flowID).
		UpdateSettingsFlowBody(body).
		Cookie(cookieHeader).
		Execute()
}

func (g *IdentityBrowserGateway) UpdateSettingsFlowWithSessionToken(
	ctx context.Context,
	flowID string,
	body client.UpdateSettingsFlowBody,
	sessionToken string,
) (*client.SettingsFlow, *http.Response, error) {

	if err := identity_helper.CheckClient(g.Browser); err != nil {
		return nil, nil, err
	}

	return g.Browser.FrontendAPI.
		UpdateSettingsFlow(ctx).
		Flow(flowID).
		UpdateSettingsFlowBody(body).
		XSessionToken(sessionToken).
		Execute()
}

func (g *IdentityBrowserGateway) ToSession(
	ctx context.Context,
	sessionToken string,
) (*client.Session, *http.Response, error) {
	if err := identity_helper.CheckClient(g.Browser); err != nil {
		return nil, nil, err
	}

	return g.Browser.FrontendAPI.
		ToSession(ctx).
		XSessionToken(sessionToken).
		Execute()
}
