package oauth

import (
	"net/http"
	"sso-bff/modules/permission"
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

type SessionPayload struct {
	Authenticated  bool                     `json:"authenticated"`
	Sub            string                   `json:"sub"`
	Exp            time.Time                `json:"exp"`
	Profile        map[string]any           `json:"profile,omitempty"`
	Roles          []permission.RolePayload `json:"roles,omitempty"`
	OrganizationID any                      `json:"organization_id,omitempty"`
	ProfileSource  string                   `json:"profile_source,omitempty"`
}
