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
	c, err := r.GetQueries(tx).CreateChannel(ctx, params)
	if err != nil {
		if db.IsUniqueConstraintViolation(err, "uq_channels_organization_id_name") {
			return db.Channel{}, db.ErrAlreadyExists
		}
		return db.Channel{}, err
	}
	return c, nil
}

func (r *repository) ListChannels(ctx context.Context, orgID uuid.UUID) ([]db.Channel, error) {
	return r.GetQueries(nil).ListChannels(ctx, orgID)
}

func (r *repository) GetChannel(
	ctx context.Context,
	params db.GetChannelParams,
) (db.Channel, error) {
	return r.GetQueries(nil).GetChannel(ctx, params)
}

func (r *repository) DeleteChannel(
	ctx context.Context,
	tx pgx.Tx,
	params db.DeleteChannelParams,
) error {
	return r.GetQueries(tx).DeleteChannel(ctx, params)
}
