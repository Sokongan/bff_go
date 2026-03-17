package identity_helper

import (
	"errors"

	client "github.com/ory/kratos-client-go"
)

// CheckClient ensures the Kratos API client is configured.
func CheckClient(c *client.APIClient) error {
	if c == nil {
		return errors.New("Identity client not configured")
	}
	return nil
}
