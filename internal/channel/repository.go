package channel

import (
	"context"

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
