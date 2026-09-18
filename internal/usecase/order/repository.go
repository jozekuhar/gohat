package orderuc

import (
	"context"
	"uuid"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/provider/db/sqlc"
)

type repository struct {
	*db.BaseRepo
}

func NewRepository(repo *db.BaseRepo) *repository {
	return &repository{
		repo,
	}
}

func (r *repository) ListOrdersWithDetails(
	ctx context.Context,
	orgID uuid.UUID,
) ([]sqlc.ListOrdersWithDetailsRow, error) {
	return r.GetQueries(nil).ListOrdersWithDetails(ctx, orgID)
}
