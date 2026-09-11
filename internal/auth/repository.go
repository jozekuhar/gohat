package auth

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

func (r *repository) CreateUser(
	ctx context.Context,
	tx pgx.Tx,
	params db.CreateUserParams,
) (db.User, error) {
	return r.GetQueries(tx).CreateUser(ctx, params)
}

func (r *repository) CreateAuthentication(
	ctx context.Context,
	tx pgx.Tx,
	params db.CreateAuthenticationParams,
) (db.Authentication, error) {
	return r.GetQueries(tx).CreateAuthentication(ctx, params)
}

func (r *repository) GetAuthenticationByEmail(
	ctx context.Context,
	params db.GetAuthenticationByEmailParams,
) (db.Authentication, error) {
	return r.GetQueries(nil).GetAuthenticationByEmail(ctx, params)
}

func (r *repository) GetAuthenticationByProvider(
	ctx context.Context,
	params db.GetAuthenticationByProviderParams,
) (db.Authentication, error) {
	return r.GetQueries(nil).GetAuthenticationByProvider(ctx, params)
}

func (r *repository) CreateSession(
	ctx context.Context,
	tx pgx.Tx,
	params db.CreateSessionParams,
) (db.Session, error) {
	return r.GetQueries(tx).CreateSession(ctx, params)
}

func (r *repository) GetSession(
	ctx context.Context,
	sessionID uuid.UUID,
) (db.Session, error) {
	return r.GetQueries(nil).GetSession(ctx, sessionID)
}

func (r *repository) DeleteSession(
	ctx context.Context,
	tx pgx.Tx,
	sessionID uuid.UUID,
) error {
	return r.GetQueries(tx).DeleteSession(ctx, sessionID)
}
