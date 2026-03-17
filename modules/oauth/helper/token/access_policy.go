package oauth_helper_token

func IsAllowedClient(clientID string, allowed map[string]struct{}) bool {
	if len(allowed) == 0 {
		return true
	}

	_, ok := allowed[clientID]
	return ok
}

func FilterScopes(requested []string, allowed map[string]struct{}) []string {
	if len(allowed) == 0 {
		return requested
	}

	out := make([]string, 0, len(requested))

	for _, s := range requested {
		if _, ok := allowed[s]; ok {
			out = append(out, s)
		}
	}

	return out
}
