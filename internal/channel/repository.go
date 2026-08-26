package channel

import (
	"context"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) ListChannels(ctx context.Context, orgID uuid.UUID) ([]Channel, error) {
	stmt := `
        SELECT * FROM channels
		WHERE organization_id = @organization_id
    `

	rows, err := r.pool.Query(ctx, stmt, pgx.NamedArgs{
		"organization_id": orgID,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Channel])
}
