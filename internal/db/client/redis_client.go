package client

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func NewStoreClient(ctx context.Context, addr, password string, db int) (*Store, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	pctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := rdb.Ping(pctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &Store{client: rdb}, nil
}

func (r *Store) Client() *redis.Client { return r.client }

func (r *Store) Close() error {
	if r == nil || r.client == nil {
		return nil
	}
	return r.client.Close()
}
