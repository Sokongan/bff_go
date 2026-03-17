package oauth_helper_session

import (
	oauth_domain "sso-bff/internal/domain/oauth"
	"sso-bff/modules/oauth"

	"golang.org/x/oauth2"
)

func BuildSession(
	token *oauth2.Token,
	idToken string,
	claims *oauth_domain.IDTokenClaims,
	kratosToken string,
) oauth.Session {

	return oauth.Session{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		IDToken:      idToken,
		Expiry:       claims.ExpiresAt,
		Subject:      claims.Subject,
		KratosToken:  kratosToken,
	}
}
