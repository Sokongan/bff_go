package oauth_helper_redirect

import (
	"fmt"
	"sso-bff/modules/oauth"
	"strings"
)

func BuildAppRedirect(
	app oauth.AppRedirect,
	path string,
) (string, error) {

	if path == "" {
		path = "/"
	}

	if len(app.AllowedPaths) > 0 {
		if _, ok := app.AllowedPaths[path]; !ok {
			return "", fmt.Errorf("redirect path not allowed")
		}
	}

	return strings.TrimRight(app.BaseURL, "/") + path, nil
}
