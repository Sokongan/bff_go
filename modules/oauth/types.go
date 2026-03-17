package oauth

import (
	"net/http"
	"time"
)

type BrowserClient struct {
	BrowserPublicURL string
	ClientID         string
	ClientSecret     string
	RedirectURL      string
	Scopes           []string
}

type InternalClient struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

type M2MClient struct {
	TokenURL  string
	M2MID     string
	M2MSecret string
	Scopes    []string
}

type Session struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	IDToken      string    `json:"id_token"`
	Expiry       time.Time `json:"expiry"`
	Subject      string    `json:"subject"`
	KratosToken  string    `json:"kratos_token"`
}

type CookieConfig struct {
	Name     string
	Domain   string
	SameSite http.SameSite
	Secure   bool
}

type AppRedirect struct {
	BaseURL      string
	AllowedPaths map[string]struct{}
}

type RedirectPayload struct {
	AppID string `json:"app_id"`
	Path  string `json:"path"`
}
