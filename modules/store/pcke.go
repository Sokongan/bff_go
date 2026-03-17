package store

import (
	"context"
	"sso-bff/modules/oauth"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisPKCEStore struct {
	Client *redis.Client
	Prefix string
}

func NewRedisPKCEStore(client *redis.Client) *RedisPKCEStore {
	return &RedisPKCEStore{
		Client: client,
		Prefix: "pkce:",
	}
}

func (s *RedisPKCEStore) SaveVerifier(
	ctx context.Context,
	state, verifier string,
	ttl time.Duration,
) error {
	if s == nil || s.Client == nil {
		return oauth.ErrStoreMisconfigured
	}
	return s.Client.Set(ctx, s.key(state), verifier, ttl).Err()
}

func (s *RedisPKCEStore) GetVerifier(
	ctx context.Context,
	state string,
) (string, error) {
	if s == nil || s.Client == nil {
		return "", oauth.ErrStoreMisconfigured
	}

	val, err := s.Client.Get(ctx, s.key(state)).Result()
	if err == redis.Nil {
		return "", oauth.ErrStateNotFound
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (s *RedisPKCEStore) DeleteVerifier(
	ctx context.Context,
	state string,
) error {
	if s == nil || s.Client == nil {
		return oauth.ErrStoreMisconfigured
	}
	return s.Client.Del(ctx, s.key(state)).Err()
}

func (s *RedisPKCEStore) key(state string) string {
	return s.Prefix + state
}
