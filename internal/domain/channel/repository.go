package channel

import (
	"context"
	"uuid"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/provider/db/sqlc"

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
	params sqlc.CreateChannelParams,
) (sqlc.Channel, error) {
	c, err := r.GetQueries(tx).CreateChannel(ctx, params)
	if err != nil {
		if db.IsUniqueConstraintViolation(err, "uq_channels_organization_id_name") {
			return sqlc.Channel{}, db.ErrAlreadyExists
		}
		return sqlc.Channel{}, err
	}
	return c, nil
}

func (r *repository) ListChannels(ctx context.Context, orgID uuid.UUID) ([]sqlc.Channel, error) {
	return r.GetQueries(nil).ListChannels(ctx, orgID)
}

func (r *repository) GetChannel(
	ctx context.Context,
	params sqlc.GetChannelParams,
) (sqlc.Channel, error) {
	return r.GetQueries(nil).GetChannel(ctx, params)
}

func (r *repository) DeleteChannel(
	ctx context.Context,
	tx pgx.Tx,
	params sqlc.DeleteChannelParams,
) error {
	return r.GetQueries(tx).DeleteChannel(ctx, params)
}
