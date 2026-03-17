package oauth_helper_token

import (
	identity_domain "sso-bff/internal/domain/identity"
)

func BuildClaims(subject string) map[string]any {
	return map[string]any{
		"sub": subject,
	}
}

func BuildAccessClaims(identity *identity_domain.Identity, subject string) map[string]any {
	claims := map[string]any{
		"sub": subject,
	}

	if identity != nil && len(identity.Traits) > 0 {
		claims["traits"] = identity.Traits
	}

	return claims
}
