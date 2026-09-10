package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// type DBTX interface {
// 	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
// 	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
// 	QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row
// }

type BaseRepo struct {
	Pool    *pgxpool.Pool
	Queries *Queries
}

func NewBaseRepo(pool *pgxpool.Pool) *BaseRepo {
	return &BaseRepo{
		Pool:    pool,
		Queries: New(pool),
	}
}

func (r *BaseRepo) DB(tx pgx.Tx) DBTX {
	if tx != nil {
		return tx
	}
	return r.Pool
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

// TODO(jozekuhar): this is tx for sqlc
