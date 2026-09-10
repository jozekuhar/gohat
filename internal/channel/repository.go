package channel

import (
	"context"
	"uuid"

	"mimokocke/internal/provider/db"

	"github.com/jackc/pgx/v5"
)

type repository struct {
	*db.BaseRepo
}

func NewRepository(baseRepo *db.BaseRepo) *repository {
	return &repository{
		BaseRepo: baseRepo,
	}
}

func (r *repository) ListChannels(ctx context.Context, orgID uuid.UUID) ([]Channel, error) {
	stmt := `
        SELECT * FROM channels
		WHERE organization_id = @organization_id
    `

	rows, err := r.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"organization_id": orgID,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Channel])
}
