package db

import (
	"context"
	"fmt"
	"sso-bff/internal/db/client"
)

type Resources struct {
	Db    *client.Db
	Store *client.Store
}

func NewResources(ctx context.Context, DbDSN, storeAddr, storePassword string, storeDB int) (*Resources, error) {

	db, err := client.NewDbPool(ctx, DbDSN)
	if err != nil {
		return nil, fmt.Errorf("postgres connection: %w", err)
	}

	store, err := client.NewStoreClient(ctx, storeAddr, storePassword, storeDB)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("redis connection: %w", err)
	}

	return &Resources{
		Db:    db,
		Store: store,
	}, nil
}

func (r *Resources) Close() {
	if r == nil {
		return
	}
	if r.Store != nil {
		_ = r.Store.Close()
	}
	if r.Db != nil {
		r.Db.Close()
	}
}
