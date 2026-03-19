package lib

import (
	"crypto/rand"
	"encoding/base64"
)

func randomToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func PCKEVerifier() (string, error) {
	return randomToken(32)
}

func PKCEStateToken() (string, error) {
	return randomToken(32)
}
