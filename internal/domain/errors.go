package domain

import "errors"

var (
	ErrMissingLoginChallenge  = errors.New("missing login_challenge")
	ErrMissingLoginIdentity   = errors.New("missing identity in login response")
	ErrMissingLoginIdentifier = errors.New("missing identifier or password")

	ErrIdentityMisconfigured = errors.New("identity service misconfigured")
)
