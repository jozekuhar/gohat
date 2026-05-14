package repository

import (
	"context"
	"time"

	"gohat/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Auth struct {
	pool *pgxpool.Pool
}

func NewAuth(pool *pgxpool.Pool) *Auth {
	return &Auth{
		pool: pool,
	}
}

func (r *Auth) CreateSession(
	ctx context.Context,
	sessionID uuid.UUID,
	userID int64,
	expiresAt time.Time,
) (*model.Session, error) {
	stmt := `
        INSERT INTO sessions (id, user_id, expires_at)
		VALUES (@id, @user_id, @expires_at)
		RETURNING *
    `

	rows, err := r.pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":         sessionID,
		"user_id":    userID,
		"expires_at": expiresAt,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.Session])
}

func (r *Auth) GetSession(ctx context.Context, sessionID uuid.UUID) (*model.Session, error) {
	stmt := `
        SELECT * FROM sessions
		WHERE id = @session_id
    `

	rows, err := r.pool.Query(ctx, stmt, pgx.NamedArgs{
		"session_id": sessionID,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.Session])
}

func (r *Auth) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	stmt := `
        DELETE FROM sessions
		WHERE id = @id
    `

	_, err := r.pool.Exec(ctx, stmt, pgx.NamedArgs{
		"id": sessionID,
	})
	return err
}
