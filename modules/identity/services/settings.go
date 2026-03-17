package identity_service

import (
	"context"
	"errors"
	"sso-bff/modules/identity"

	client "github.com/ory/kratos-client-go"
)

var ErrSettingsMisconfigured = errors.New("identity settings service misconfigured")

type IdentitySettingsService struct {
	ident identity.IdentitySettingsClient
}

func NewIdentitySettingsService(
	ident identity.IdentitySettingsClient) *IdentitySettingsService {
	return &IdentitySettingsService{ident: ident}
}

func (s *IdentitySettingsService) CreateBrowserSettingsFlow(
	ctx context.Context,
	cookieHeader string,
) (*client.SettingsFlow, error) {

	if s == nil || s.ident == nil {
		return nil, ErrSettingsMisconfigured
	}
	if cookieHeader == "" {
		return nil, errors.New("missing session token")
	}

	flow, resp, err := s.ident.CreateNativeSettingsFlow(
		ctx,
		cookieHeader,
	)

	if err != nil {
		return nil, wrapHTTPError(resp, err)
	}
	return flow, nil
}

func (s *IdentitySettingsService) SubmitSettingsFlow(
	ctx context.Context,
	flowID string,
	body client.UpdateSettingsFlowBody,
	cookieHeader string,
) (*client.SettingsFlow, error) {

	if s == nil || s.ident == nil {
		return nil, ErrSettingsMisconfigured
	}

	if flowID == "" {
		return nil, errors.New("missing flow id")
	}

	if cookieHeader == "" {
		return nil, errors.New("missing session token")
	}

	flow, resp, err := s.ident.UpdateSettingsFlowWithSessionToken(
		ctx,
		flowID,
		body,
		cookieHeader,
	)

	if err != nil {
		return nil, wrapHTTPError(resp, err)
	}

	return flow, nil
}
