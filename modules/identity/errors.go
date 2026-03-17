package identity

import "errors"

var (
	ErrAdminMisconfigured    = errors.New("identity admin service misconfigured")
	ErrIdentityMisconfigured = errors.New("Identity service misconfigured")
	ErrOauthMisconfigured    = errors.New("Oauth service misconfigured")
	ErrMissingCookieHeader   = errors.New("missing cookie header")
)
