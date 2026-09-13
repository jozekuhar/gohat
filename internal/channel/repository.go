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

func (r *repository) CreateChannel(
	ctx context.Context,
	tx pgx.Tx,
	params db.CreateChannelParams,
) (db.Channel, error) {
	return r.GetQueries(tx).CreateChannel(ctx, params)
}

func (r *repository) ListChannels(ctx context.Context, orgID uuid.UUID) ([]db.Channel, error) {
	return r.GetQueries(nil).ListChannels(ctx, orgID)
}
