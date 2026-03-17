package oauth

import "errors"

var (
	ErrStoreMisconfigured = errors.New("oauth store misconfigured")
	ErrStateNotFound      = errors.New("pkce state not found")

	ErrPKCEStore     = errors.New("pkce store error")
	ErrTokenExchange = errors.New("token exchange failed")
	ErrHydraRequest  = errors.New("hydra request failed")

	ErrServiceMisconfigured = errors.New("oauth service misconfigured")
	ErrSessionStore         = errors.New("session store error")
	ErrSessionNotFound      = errors.New("session not found")
	ErrInvalidIDToken       = errors.New("invalid id_token")
)
