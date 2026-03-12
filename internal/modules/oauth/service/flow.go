package oauth_service_flow

import (
	"context"
	"encoding/json"
	"fmt"
	"sso-bff/internal/lib"
	"sso-bff/internal/modules/oauth"
	"time"

	"golang.org/x/oauth2"
)

type FlowService struct {
	flow      oauth.OAuthClientPort
	pkce      oauth.PKCEStorePort
	redirects oauth.RedirectStorePort
	pkceTTL   time.Duration
	sessions  oauth.SessionStorePort
}

func NewFlowService(
	flow oauth.OAuthClientPort,
	pkce oauth.PKCEStorePort,
	redirects oauth.RedirectStorePort,
	pkceTTL time.Duration,
	sessions oauth.SessionStorePort,
) *FlowService {
	return &FlowService{
		flow:      flow,
		pkce:      pkce,
		redirects: redirects,
		pkceTTL:   pkceTTL,
		sessions:  sessions,
	}
}

func (f *FlowService) Login(
	ctx context.Context,
	appID,
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

	if err := f.pkce.SaveVerifier(
		ctx,
		state,
		verifier,
		f.pkceTTL,
	); err != nil {
		return "", err
	}

	if (redirectPath != "" || appID != "") && s.redirects != nil {
		payload := oauth.RedirectPayload{
			AppID: appID,
			Path:  redirectPath,
		}
		raw, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("encode redirect: %w", err)
		}
		if err := s.redirects.SaveRedirect(ctx, state, string(raw), s.pkceTTL); err != nil {
			return "", fmt.Errorf("%w: %w", oauth.ErrPKCEStore, err)
		}
	}

	return s.flow.AuthCodeURL(state, verifier), nil
}

func (f *FlowService) Callback(ctx context.Context, code, state string) (*oauth2.Token, error) {
	verifier, err := f.pkce.GetVerifier(ctx, state)
	if err != nil {
		return nil, err
	}
	defer f.pkce.DeleteVerifier(ctx, state)

	return f.flow.Exchange(ctx, code, verifier)
}
