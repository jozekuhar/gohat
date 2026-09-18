package db

import (
	"mimokocke/internal/provider/db/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseRepo struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewBaseRepo(pool *pgxpool.Pool) *BaseRepo {
	return &BaseRepo{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *BaseRepo) Pool() *pgxpool.Pool {
	return r.pool
}

func (r *BaseRepo) GetQueries(tx pgx.Tx) *sqlc.Queries {
	if tx != nil {
		return r.queries.WithTx(tx)
	}
	return r.queries
}
