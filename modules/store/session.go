package store

import (
	"context"
	"encoding/json"
	"sso-bff/modules/oauth"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSessionStore struct {
	Client *redis.Client
	Prefix string
}

func NewRedisSessionStore(client *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{
		Client: client,
		Prefix: "oauth_session:",
	}
}

func (s *RedisSessionStore) SaveSession(
	ctx context.Context,
	sessionID string,
	session oauth.Session,
	ttl time.Duration,
) error {
	if s == nil || s.Client == nil {
		return oauth.ErrStoreMisconfigured
	}
	payload, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return s.Client.Set(ctx, s.key(sessionID), payload, ttl).Err()
}

func (s *RedisSessionStore) GetSession(
	ctx context.Context,
	sessionID string,
) (oauth.Session, error) {
	if s == nil || s.Client == nil {
		return oauth.Session{}, oauth.ErrStoreMisconfigured
	}
	val, err := s.Client.Get(ctx, s.key(sessionID)).Result()
	if err == redis.Nil {
		return oauth.Session{}, oauth.ErrSessionNotFound
	}
	if err != nil {
		return oauth.Session{}, err
	}

	var session oauth.Session
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		return oauth.Session{}, err
	}
	return session, nil
}

func (s *RedisSessionStore) DeleteSession(
	ctx context.Context,
	sessionID string,
) error {
	if s == nil || s.Client == nil {
		return oauth.ErrStoreMisconfigured
	}
	return s.Client.Del(ctx, s.key(sessionID)).Err()
}

func (s *RedisSessionStore) key(sessionID string) string {
	return s.Prefix + sessionID
}
