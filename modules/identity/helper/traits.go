package identity_helper

func ExtractTraits(v any) map[string]any {
	if traits, ok := v.(map[string]any); ok {
		return traits
	}
	return map[string]any{}
}
