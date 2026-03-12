package gateway

import (
	"context"
	"errors"
	"sso-bff/internal/domain"

	client "github.com/ory/hydra-client-go/v2"
)

type OauthAdminGateway struct {
	sdk *client.APIClient
}

func NewOauthAdminGateway(c *client.APIClient) *OauthAdminGateway {
	return &OauthAdminGateway{
		sdk: c,
	}
}

func (g *OauthAdminGateway) GetConsentRequest(
	ctx context.Context,
	challenge string,
) (*domain.ConsentRequest, error) {

	if g == nil || g.sdk == nil {
		return nil, errors.New("oauth admin gateway is not initialized")
	}

	res, _, err := g.sdk.OAuth2API.
		GetOAuth2ConsentRequest(ctx).
		ConsentChallenge(challenge).Execute()

	if err != nil {
		return nil, err
	}

	clientInfo := res.GetClient()

	if clientInfo.ClientId == nil || *clientInfo.ClientId == "" {
		return nil, errors.New("missing client id in consent request")
	}
	return &domain.ConsentRequest{
		Skip:           res.GetSkip(),
		ClientID:       *clientInfo.ClientId,
		RequestedScope: res.GetRequestedScope(),
		Audience:       res.GetRequestedAccessTokenAudience(),
		Subject:        res.GetSubject(),
	}, nil
}

func (g *OauthAdminGateway) AcceptConsent(
	ctx context.Context,
	challenge string,
	grantScope, grantAudience []string,
	idTokenClaims, accessTokenClaims map[string]any,
) (string, error) {
	if g == nil || g.sdk == nil {
		return "", errors.New("hydra client not configured")
	}

	accept := client.NewAcceptOAuth2ConsentRequest()
	accept.SetGrantScope(grantScope)

	if len(grantAudience) > 0 {
		accept.SetGrantAccessTokenAudience(grantAudience)
	}

	// session claims
	if len(idTokenClaims) > 0 || len(accessTokenClaims) > 0 {
		sess := client.NewAcceptOAuth2ConsentRequestSession()
		if len(idTokenClaims) > 0 {
			sess.SetIdToken(idTokenClaims)
		}
		if len(accessTokenClaims) > 0 {
			sess.SetAccessToken(accessTokenClaims)
		}
		accept.SetSession(*sess)
	}

	resp, _, err := g.sdk.OAuth2API.
		AcceptOAuth2ConsentRequest(ctx).
		ConsentChallenge(challenge).
		AcceptOAuth2ConsentRequest(*accept).
		Execute()
	if err != nil {
		return "", err
	}

	return resp.GetRedirectTo(), nil
}
func (g *OauthAdminGateway) AcceptLogin(ctx context.Context, loginChallenge, subject string) (string, error) {
	if g == nil || g.sdk == nil {
		return "", errors.New("client not configured")
	}
	if subject == "" {
		return "", errors.New("login subject missing")
	}

	body := client.NewAcceptOAuth2LoginRequest(subject)
	resp, _, err := g.sdk.OAuth2API.
		AcceptOAuth2LoginRequest(ctx).
		LoginChallenge(loginChallenge).
		AcceptOAuth2LoginRequest(*body).
		Execute()
	if err != nil {
		return "", err
	}
	return resp.GetRedirectTo(), nil
}
