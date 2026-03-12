package oauth

import (
	"net/http"
	"time"
)

type CookieConfig struct {
	Name     string
	TTL      time.Duration
	Domain   string
	Secure   bool
	SameSite http.SameSite
}

func LoadCookieConfig() *CookieConfig {
	return &CookieConfig{
		Name:     "sso_session",
		TTL:      2 * time.Hour,
		Domain:   ".doj.gov.ph",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
