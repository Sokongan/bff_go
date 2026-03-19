package oauth_helper_redirect

import (
	"encoding/json"
	"sso-bff/modules/oauth"
	"strings"
)

func EncodeRedirectPayload(appID, path string) (string, error) {
	payload := oauth.RedirectPayload{
		AppID: appID,
		Path:  path,
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}

func DecodeRedirectPayload(data string) (string, string, error) {
	if data == "" {
		return "", "", nil
	}

	var payload oauth.RedirectPayload

	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return "", data, nil
	}

	return payload.AppID, payload.Path, nil
}

func IsSafeRedirectPath(value string) bool {
	if value == "" {
		return false
	}
	if strings.HasPrefix(value, "//") || strings.Contains(value, "://") {
		return false
	}
	return strings.HasPrefix(value, "/")
}
