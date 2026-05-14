package repository

import (
	"context"
	"errors"

	"gohat/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotExist = errors.New("user not exist")

type User struct {
	pool *pgxpool.Pool
}

func NewUser(pool *pgxpool.Pool) *User {
	return &User{
		pool: pool,
	}
}

func (r *User) GetOrCreateUserWithGoogleID(
	ctx context.Context,
	email, googleID string,
) (*model.User, error) {
	stmt := `
		INSERT INTO users(email, google_id)
		VALUES (@email, @google_id)
		ON CONFLICT (google_id) 
		DO UPDATE
		SET email = EXCLUDED.email
		RETURNING *
    `

	rows, err := r.pool.Query(ctx, stmt, pgx.NamedArgs{
		"email":     email,
		"google_id": googleID,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.User])
}
