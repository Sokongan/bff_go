// domain/oidc.go
package oauth_domain

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type OauthOIDCVerifier struct {
	JWKSURL  string
	Issuer   string
	Audience string
	Client   *http.Client
	CacheTTL time.Duration
	Nonce    string

	mu        sync.Mutex
	cachedAt  time.Time
	cachedKey map[string]*rsa.PublicKey
}

type IDTokenClaims struct {
	Subject   string
	ExpiresAt time.Time
}

// Domain method handles JWKS fetch + caching internally
func (v *OauthOIDCVerifier) GetKeys(ctx context.Context) (map[string]*rsa.PublicKey, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.CacheTTL > 0 && time.Since(v.cachedAt) < v.CacheTTL && len(v.cachedKey) > 0 {
		return v.cachedKey, nil
	}

	if v.Client == nil {
		v.Client = http.DefaultClient
	}
	if v.JWKSURL == "" {
		return nil, errors.New("jwks url missing")
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, v.JWKSURL, nil)
	resp, err := v.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("jwks http status: %d", resp.StatusCode)
	}

	var doc struct {
		Keys []struct {
			Kid string `json:"kid"`
			Kty string `json:"kty"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, err
	}

	keys := make(map[string]*rsa.PublicKey)
	for _, k := range doc.Keys {
		n, err := base64.RawURLEncoding.DecodeString(k.N)
		if err != nil {
			continue
		}
		e, err := base64.RawURLEncoding.DecodeString(k.E)
		if err != nil {
			continue
		}
		key := &rsa.PublicKey{
			N: new(big.Int).SetBytes(n),
			E: int(new(big.Int).SetBytes(e).Int64()),
		}
		keys[k.Kid] = key
	}

	if len(keys) == 0 {
		return nil, errors.New("no jwks keys parsed")
	}

	v.cachedKey = keys
	v.cachedAt = time.Now()
	return keys, nil
}

// Domain method: verify token
func (v *OauthOIDCVerifier) Verify(ctx context.Context, rawIDToken string) (map[string]interface{}, error) {
	if rawIDToken == "" {
		return nil, errors.New("id_token missing")
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))
	claims := jwt.MapClaims{}

	token, err := parser.ParseWithClaims(rawIDToken, claims, v.keyFunc(ctx))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token invalid")
	}
	return claims, nil
}

func (v *OauthOIDCVerifier) keyFunc(ctx context.Context) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		kid, _ := token.Header["kid"].(string)
		if kid == "" {
			return nil, errors.New("kid missing")
		}
		keys, err := v.GetKeys(ctx)
		if err != nil {
			return nil, err
		}
		key := keys[kid]
		if key == nil {
			return nil, fmt.Errorf("key not found: %s", kid)
		}
		return key, nil
	}
}
