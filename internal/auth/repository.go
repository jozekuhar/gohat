package auth

import (
	"context"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) GetOrCreateUserWithGoogleID(
	ctx context.Context,
	id uuid.UUID,
	email, googleID string,
) (*User, error) {
	stmt := `
		INSERT INTO users(id, email, google_id)
		VALUES (@id, @email, @google_id)
		ON CONFLICT (google_id) 
		DO UPDATE
		SET email = EXCLUDED.email
		RETURNING *
    `

	rows, err := r.pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":        id,
		"email":     email,
		"google_id": googleID,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[User])
}

func (r *repository) CreateSession(
	ctx context.Context,
	sessionID uuid.UUID,
	userID uuid.UUID,
	expiresAt time.Time,
) (*Session, error) {
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

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Session])
}

func (r *repository) GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
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

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Session])
}

func (r *repository) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	stmt := `
        DELETE FROM sessions
		WHERE id = @id
    `

	_, err := r.pool.Exec(ctx, stmt, pgx.NamedArgs{
		"id": sessionID,
	})
	return err
}
