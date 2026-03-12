package oauth_gateway

import (
	"context"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type OAuthM2MGateway struct {
	cfg   *clientcredentials.Config
	token *oauth2.Token
	mu    sync.Mutex
}

func NewOAuthM2MGateway(cfg *clientcredentials.Config) *OAuthM2MGateway {
	return &OAuthM2MGateway{cfg: cfg}
}

func (g *OAuthM2MGateway) Token(
	ctx context.Context,
	scopes []string,
) (*oauth2.Token, error) {

	g.mu.Lock()
	defer g.mu.Unlock()

	if g.token != nil && g.token.Expiry.After(
		time.Now().Add(10*time.Second)) {
		return g.token, nil
	}

	cfg := *g.cfg
	if len(scopes) > 0 {
		cfg.Scopes = scopes
	}

	token, err := cfg.Token(ctx)

	if err != nil {
		return nil, err
	}

	g.token = token

	return token, nil
}
