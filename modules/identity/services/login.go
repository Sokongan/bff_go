package identity_service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sso-bff/internal/domain"
	"sso-bff/modules/identity"

	client "github.com/ory/kratos-client-go"
)

type HTTPError struct {
	Status int
	Body   []byte
	Err    error
}

func (e *HTTPError) Error() string {
	if e == nil || e.Err == nil {
		return "http error"
	}
	return e.Err.Error()
}

type IdentityLoginService struct {
	ident identity.IdentityLoginClient
	oauth identity.OauthLoginAdmin
}

func NewIdentityLoginService(ident identity.IdentityLoginClient,
	oauth identity.OauthLoginAdmin,
) *IdentityLoginService {

	return &IdentityLoginService{
		ident: ident,
		oauth: oauth,
	}
}

func (s *IdentityLoginService) AuthenticatePassword(ctx context.Context, identifier, password, loginChallenge string) (string, string, string, error) {
	if s == nil || s.ident == nil {
		return "", "", "", identity.ErrIdentityMisconfigured
	}
	if s.oauth == nil {
		return "", "", "", identity.ErrOauthMisconfigured
	}
	if identifier == "" || password == "" {
		return "", "", "", domain.ErrMissingLoginIdentifier
	}
	if loginChallenge == "" {
		return "", "", "", domain.ErrMissingLoginChallenge
	}

	flow, resp, err := s.ident.CreateNativeLoginFlow(ctx)
	if err != nil {
		return "", "", "", wrapHTTPError(resp, err)
	}

	method :=
		client.NewUpdateLoginFlowWithPasswordMethod(identifier,
			"password",
			password,
		)

	body :=
		client.UpdateLoginFlowWithPasswordMethodAsUpdateLoginFlowBody(method)

	login, resp, err := s.ident.UpdateLoginFlow(ctx,
		flow.GetId(),
		body,
	)

	if err != nil {
		return "", "", "", wrapHTTPError(resp, err)
	}

	var session *client.Session
	if sess, ok := login.GetSessionOk(); ok {
		session = sess
	}

	identityID := identityIDFromSession(session)
	if identityID == "" {
		return "", "", "", domain.ErrMissingLoginIdentity
	}

	redirectTo, err := s.oauth.AcceptLogin(ctx, loginChallenge, identityID)
	if err != nil {
		return "", "", "", fmt.Errorf("accept hydra login: %w", err)
	}

	sessionToken := ""
	if token, ok := login.GetSessionTokenOk(); ok {
		sessionToken = *token
	}

	return redirectTo, identityID, sessionToken, nil
}

func identityIDFromSession(session *client.Session) string {
	if session == nil || session.Identity == nil {
		return ""
	}
	return session.Identity.GetId()
}

func wrapHTTPError(resp *http.Response, err error) error {
	if err == nil {
		return nil
	}
	status := http.StatusBadGateway
	if resp != nil {
		status = resp.StatusCode
	}
	var apiErr *client.GenericOpenAPIError
	if errors.As(err, &apiErr) {
		return &HTTPError{Status: status, Body: apiErr.Body(), Err: err}
	}
	return &HTTPError{Status: status, Err: err}
}
