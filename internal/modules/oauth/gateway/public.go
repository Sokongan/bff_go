package oauth_gateway

import (
	"context"
	"errors"
	"sso-bff/internal/lib"

	"golang.org/x/oauth2"
)

type OAuthAuthorizationGateway struct {
	Browser  *oauth2.Config
	Internal *oauth2.Config
}

func NewOAuthAuthorizationGateway(browser, internal *oauth2.Config) (*OAuthAuthorizationGateway, error) {
	if browser == nil || internal == nil {
		return nil, errors.New("oauth config is nil")
	}
	return &OAuthAuthorizationGateway{
		Browser:  browser,
		Internal: internal,
	}, nil
}

func (g *OAuthAuthorizationGateway) AuthCodeURL(
	state,
	codeVerifier string,
) string {
	opts := []oauth2.AuthCodeOption{}
	if codeVerifier != "" {
		opts = append(opts,
			oauth2.SetAuthURLParam("code_challenge", lib.PKCEChallenge(codeVerifier)),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		)
	}
	return g.Browser.AuthCodeURL(state, opts...)
}

func (g *OAuthAuthorizationGateway) Exchange(
	ctx context.Context,
	code,
	codeVerifier string,
) (*oauth2.Token, error) {
	opts := []oauth2.AuthCodeOption{}
	if codeVerifier != "" {
		opts = append(opts, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	}
	return g.Internal.Exchange(ctx, code, opts...)
}

func (g *OAuthAuthorizationGateway) Refresh(
	ctx context.Context,
	refreshToken string,
) (*oauth2.Token, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token missing")
	}
	token := &oauth2.Token{RefreshToken: refreshToken}
	return g.Internal.TokenSource(ctx, token).Token()
}
