package client

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Db struct {
	pool *pgxpool.Pool
}

func NewDbPool(ctx context.Context, dsn string) (*Db, error) {
	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(cctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}

	pctx, pcancel := context.WithTimeout(ctx, 5*time.Second)
	defer pcancel()

	if err := pool.Ping(pctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	return &Db{pool: pool}, nil
}

func (p *Db) Pool() *pgxpool.Pool { return p.pool }

func (p *Db) StdlibDB() *sql.DB {
	if p == nil || p.pool == nil {
		return nil
	}
	return stdlib.OpenDBFromPool(p.pool)
}

func (p *Db) Close() {
	if p != nil && p.pool != nil {
		p.pool.Close()
	}
}
