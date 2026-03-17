package oauth_helper_token

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func ClientIDFromAccessToken(accessToken string) (string, error) {
	parts := strings.Split(accessToken, ".")
	if len(parts) < 2 {
		return "", errors.New("invalid access token")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("decode access token: %w", err)
	}
	var claims struct {
		ClientID string `json:"client_id"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", fmt.Errorf("parse access token: %w", err)
	}
	if claims.ClientID == "" {
		return "", errors.New("client_id missing in access token")
	}
	return claims.ClientID, nil
}
