package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseRepo struct {
	pool    *pgxpool.Pool
	queries *Queries
}

func NewBaseRepo(pool *pgxpool.Pool) *BaseRepo {
	return &BaseRepo{
		pool:    pool,
		queries: New(pool),
	}
}

func (r *BaseRepo) Pool() *pgxpool.Pool {
	return r.pool
}

func (r *BaseRepo) GetQueries(tx pgx.Tx) *Queries {
	if tx != nil {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func LoadPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("creating pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pinging db: %w", err)
	}

	return pool, err
}
