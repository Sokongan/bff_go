package store

import (
	"github.com/redis/go-redis/v9"
)

type Stores struct {
	PKCE     *RedisPKCEStore
	Redirect *RedisRedirectStore
	Sessions *RedisSessionStore
}

func NewStore(rdb *redis.Client) *Stores {
	return &Stores{
		PKCE:     NewRedisPKCEStore(rdb),
		Redirect: NewRedisRedirectStore(rdb),
		Sessions: NewRedisSessionStore(rdb),
	}
}
