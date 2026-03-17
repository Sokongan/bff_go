package service_oauth

import (
	"context"
	"fmt"
	"sso-bff/internal/lib"
	"sso-bff/modules/oauth"
	oauth_helper_redirect "sso-bff/modules/oauth/helper/redirect"

	"time"

	"golang.org/x/oauth2"
)

type FlowService struct {
	flow      oauth.OAuthClientPort
	pkce      oauth.PKCEStorePort
	redirects oauth.RedirectStorePort
	pkceTTL   time.Duration
}

func NewFlowService(
	flow oauth.OAuthClientPort,
	pkce oauth.PKCEStorePort,
	redirects oauth.RedirectStorePort,
	pkceTTL time.Duration,
) *FlowService {

	return &FlowService{
		flow:      flow,
		pkce:      pkce,
		redirects: redirects,
		pkceTTL:   pkceTTL,
	}
}

func (s *FlowService) Login(
	ctx context.Context,
	appID string,
	redirectPath string,
) (string, error) {

	state, err := lib.PKCEStateToken()
	if err != nil {
		return "", err
	}

	verifier, err := lib.PCKEVerifier()
	if err != nil {
		return "", err
	}

	if err := s.pkce.SaveVerifier(
		ctx,
		state,
		verifier,
		s.pkceTTL,
	); err != nil {

		return "", err
	}

	if s.redirects != nil && (appID != "" || redirectPath != "") {

		payload, err := oauth_helper_redirect.
			EncodeRedirectPayload(appID, redirectPath)
		if err != nil {
			return "", fmt.Errorf("redirect payload: %w", err)
		}

		if err := s.redirects.
			SaveRedirect(ctx, state, payload, s.pkceTTL); err != nil {
			return "", err
		}
	}

	return s.flow.AuthCodeURL(state, verifier), nil
}

func (s *FlowService) Callback(
	ctx context.Context,
	code string,
	state string,
) (*oauth2.Token, error) {

	verifier, err := s.pkce.GetVerifier(ctx, state)
	if err != nil {
		return nil, err
	}

	defer s.pkce.DeleteVerifier(ctx, state)

	return s.flow.Exchange(ctx, code, verifier)
}
