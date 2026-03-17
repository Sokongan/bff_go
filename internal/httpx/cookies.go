package httpx

import (
	"net/http"
	"strings"
	"time"
)

type CookieConfig struct {
	Name     string
	Domain   string
	SameSite http.SameSite
	Secure   bool
}

func Name(cfg CookieConfig) string {
	if cfg.Name != "" {
		return cfg.Name
	}
	return "sso_session"
}

func Secure(cfg CookieConfig, r *http.Request) bool {
	if cfg.Secure {
		return true
	}
	return r != nil && r.TLS != nil
}

func SameSite(cfg CookieConfig) http.SameSite {
	if cfg.SameSite == 0 {
		return http.SameSiteLaxMode
	}
	return cfg.SameSite
}

func SetSessionCookie(w http.ResponseWriter, cfg CookieConfig, sessionID string, ttl time.Duration, r *http.Request) {
	c := &http.Cookie{
		Name:     Name(cfg),
		Value:    sessionID,
		Path:     "/",
		Domain:   cfg.Domain,
		HttpOnly: true,
		Secure:   Secure(cfg, r),
		SameSite: SameSite(cfg),
	}
	if ttl > 0 {
		c.MaxAge = int(ttl.Seconds())
		c.Expires = time.Now().Add(ttl)
	}
	http.SetCookie(w, c)
}

func ClearSessionCookie(w http.ResponseWriter, cfg CookieConfig, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     Name(cfg),
		Value:    "",
		Path:     "/",
		Domain:   cfg.Domain,
		HttpOnly: true,
		Secure:   Secure(cfg, r),
		SameSite: SameSite(cfg),
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

func SessionIDFromRequest(r *http.Request, cfg CookieConfig) string {
	if r == nil {
		return ""
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		lower := strings.ToLower(authHeader)
		if strings.HasPrefix(lower, "bearer ") {
			token := strings.TrimSpace(authHeader[len("Bearer "):])
			if token != "" {
				return token
			}
		}
	}

	if c, err := r.Cookie(Name(cfg)); err == nil && c.Value != "" {
		return c.Value
	}

	return ""
}

const KratosSessionTokenCookie = "kratos_session_token"
