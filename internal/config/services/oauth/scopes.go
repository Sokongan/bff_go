package oauth

import (
	"errors"
	"os"
	"strings"
)

type ClientScopesConfig struct {
	BFFScopes     []string
	AllowedClient map[string]struct{}
	AllowedScope  map[string]struct{}
}

func parseCSVSet(v string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, s := range strings.Split(v, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			out[s] = struct{}{}
		}
	}
	return out
}

func LoadClientScopesConfig() (*ClientScopesConfig, error) {
	// Hard-coded default scopes
	bffScopes := []string{"openid", "offline"}
	allowedClients := os.Getenv("ALLOWED_CLIENT_IDS")
	allowedScopes := os.Getenv("ALLOWED_SCOPES")

	// Fail fast if any critical value is missing
	if allowedClients == "" || allowedScopes == "" {
		return nil, errors.New("oauth scope configuration missing. Set ALLOWED_CLIENT_IDS and ALLOWED_SCOPES")
	}
	return &ClientScopesConfig{
		BFFScopes:     bffScopes,
		AllowedClient: parseCSVSet(allowedClients), // still configurable per environment
		AllowedScope:  parseCSVSet(allowedScopes),  // optional override
	}, nil
}
