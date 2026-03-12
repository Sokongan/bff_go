package domain

import (
	"crypto/rsa"
	"net/http"
	"sync"
	"time"
)

type OauthOIDCVerifier struct {
	JWKSURL   string
	Issuer    string
	Audience  string
	Client    *http.Client
	CacheTTL  time.Duration
	Nonce     string
	mu        sync.Mutex
	cachedAt  time.Time
	cachedKey map[string]*rsa.PublicKey
}
