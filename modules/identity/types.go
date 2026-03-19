package identity

import identity_domain "sso-bff/internal/domain/identity"

type SubmitLoginRequest struct {
	Identifier     string `json:"identifier"`
	Password       string `json:"password"`
	LoginChallenge string `json:"login_challenge"`
}

type SubmitLoginResponse struct {
	RedirectTo string `json:"redirect_to"`
}

type IdentityRole struct {
	Object string `json:"object"`
	Role   string `json:"role"`
}
type IdentityWithRoles struct {
	ID             string                                  `json:"id"`
	Traits         map[string]any                          `json:"traits"`
	Roles          []IdentityRole                          `json:"roles,omitempty"`
	KratosSessions []identity_domain.IdentitySessionDevice `json:"kratos_sessions,omitempty"`
}
