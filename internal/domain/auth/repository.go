package auth

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

func (r *repository) CreateUser(
	ctx context.Context,
	tx pgx.Tx,
	params sqlc.CreateUserParams,
) (sqlc.User, error) {
	return r.GetQueries(tx).CreateUser(ctx, params)
}

func (r *repository) CreateAuthentication(
	ctx context.Context,
	tx pgx.Tx,
	params sqlc.CreateAuthenticationParams,
) (sqlc.Authentication, error) {
	return r.GetQueries(tx).CreateAuthentication(ctx, params)
}

func (r *repository) GetAuthenticationByEmail(
	ctx context.Context,
	params sqlc.GetAuthenticationByEmailParams,
) (sqlc.Authentication, error) {
	return r.GetQueries(nil).GetAuthenticationByEmail(ctx, params)
}

func (r *repository) GetAuthenticationByProvider(
	ctx context.Context,
	params sqlc.GetAuthenticationByProviderParams,
) (sqlc.Authentication, error) {
	return r.GetQueries(nil).GetAuthenticationByProvider(ctx, params)
}

func (r *repository) CreateSession(
	ctx context.Context,
	tx pgx.Tx,
	params sqlc.CreateSessionParams,
) (sqlc.Session, error) {
	return r.GetQueries(tx).CreateSession(ctx, params)
}

func (r *repository) GetSession(
	ctx context.Context,
	sessionID uuid.UUID,
) (sqlc.GetSessionRow, error) {
	return r.GetQueries(nil).GetSession(ctx, sessionID)
}

func (r *repository) DeleteSession(
	ctx context.Context,
	tx pgx.Tx,
	sessionID uuid.UUID,
) error {
	return r.GetQueries(tx).DeleteSession(ctx, sessionID)
}
