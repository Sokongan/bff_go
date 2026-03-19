package identity_helper

import (
	"errors"
	"strconv"

	client "github.com/ory/kratos-client-go"
)

// CheckClient ensures the Kratos API client is configured.
func CheckClient(c *client.APIClient) error {
	if c == nil {
		return errors.New("Identity client not configured")
	}
	return nil
}

func ExtractTraits(v any) map[string]any {
	if traits, ok := v.(map[string]any); ok {
		return traits
	}
	return map[string]any{}
}

func ParseInt64(raw string) (int64, error) {
	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}
