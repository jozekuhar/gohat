package order

import (
	"context"
	"uuid"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/provider/db/sqlc"
)

type repository struct {
	*db.BaseRepo
}

func NewRepository(baseRepo *db.BaseRepo) *repository {
	return &repository{
		BaseRepo: baseRepo,
	}
}

func (r *repository) ListOrders(ctx context.Context, orgID uuid.UUID) ([]sqlc.Order, error) {
	return r.GetQueries(nil).ListOrders(ctx, orgID)
}
