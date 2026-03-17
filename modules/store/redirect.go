package store

import (
	"context"
	"fmt"
	"sso-bff/modules/oauth"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRedirectStore struct {
	Client *redis.Client
	Prefix string
}

func NewRedisRedirectStore(client *redis.Client) *RedisRedirectStore {
	return &RedisRedirectStore{
		Client: client,
		Prefix: "oauth_redirect:",
	}
}

func (s *RedisRedirectStore) SaveRedirect(
	ctx context.Context,
	state,
	redirectURL string,
	ttl time.Duration,
) error {
	if s == nil || s.Client == nil {
		return fmt.Errorf("%w: redis client not configured",
			oauth.ErrPKCEStore)
	}
	return s.Client.Set(ctx, s.key(state), redirectURL, ttl).Err()
}

func (s *RedisRedirectStore) GetRedirect(
	ctx context.Context,
	state string,
) (string, error) {
	if s == nil || s.Client == nil {
		return "", fmt.Errorf("%w: redis client not configured",
			oauth.ErrPKCEStore)
	}
	val, err := s.Client.Get(ctx, s.key(state)).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (s *RedisRedirectStore) DeleteRedirect(
	ctx context.Context,
	state string,
) error {
	if s == nil || s.Client == nil {
		return fmt.Errorf("%w: redis client not configured",
			oauth.ErrPKCEStore)
	}
	return s.Client.Del(ctx, s.key(state)).Err()
}

func (s *RedisRedirectStore) key(state string) string {
	return s.Prefix + state
}
