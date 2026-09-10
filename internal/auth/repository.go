package auth

import (
	"context"
	"errors"
	"fmt"
	"uuid"

	"mimokocke/internal/model"
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

func (r *repository) CreateUser(ctx context.Context, tx pgx.Tx, u model.User) (model.User, error) {
	stmt := `
		INSERT INTO users (id, email) 
		VALUES (@id, @email) 
		RETURNING *
	`

	rows, err := r.DB(tx).Query(ctx, stmt, pgx.NamedArgs{
		"id":    u.ID,
		"email": u.Email,
	})
	if err != nil {
		// TODO(jozekuhar): errors for users, kot je index err user already exists
		return model.User{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])
}

func (r *repository) CreateAuthentication(
	ctx context.Context,
	tx pgx.Tx,
	a model.Authentication,
) (model.Authentication, error) {
	stmt := `
        INSERT INTO authentications (id, user_id, provider, provider_id, password_hash)
		VALUES (@id, @user_id, @provider, @provider_id, @password_hash)
		RETURNING *
    `

	rows, err := r.DB(tx).Query(ctx, stmt, pgx.NamedArgs{
		"id":            a.ID,
		"user_id":       a.UserID,
		"provider":      a.Provider,
		"provider_id":   a.ProviderID,
		"password_hash": a.PasswordHash,
	})
	if err != nil {
		// TODO(jozekuhar): error authentication for user_id and provider_id already exists
		return model.Authentication{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Authentication])
}

func (r *repository) GetAuthenticationByEmail(
	ctx context.Context,
	email string,
	provider model.AuthProvider,
) (model.Authentication, error) {
	stmt := `
		SELECT a.*
		FROM authentications a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE u.email = @email
		  AND a.provider = @provider
    `

	rows, err := r.DB(nil).Query(ctx, stmt, pgx.NamedArgs{
		"email":    email,
		"provider": provider,
	})
	if err != nil {
		return model.Authentication{}, err
	}
	defer rows.Close()

	authentication, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.Authentication],
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Authentication{}, fmt.Errorf("authentication: %w", db.ErrNotFound)
	}

	return authentication, err
}

func (r *repository) GetAuthenticationByProvider(
	ctx context.Context,
	provider model.AuthProvider,
	providerID string,
) (model.Authentication, error) {
	stmt := `
		SELECT *
		FROM authentications
		WHERE provider = @provider
		  AND provider_id = @provider_id
    `

	rows, err := r.DB(nil).Query(ctx, stmt, pgx.NamedArgs{
		"provider":    provider,
		"provider_id": providerID,
	})
	if err != nil {
		return model.Authentication{}, err
	}
	defer rows.Close()

	authentication, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[model.Authentication],
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Authentication{}, fmt.Errorf("authentication: %w", db.ErrNotFound)
	}

	return authentication, err
}

func (r *repository) CreateSession(
	ctx context.Context,
	tx pgx.Tx,
	s model.Session,
) (model.Session, error) {
	stmt := `
        INSERT INTO sessions (id, user_id, expires_at)
		VALUES (@id, @user_id, @expires_at)
		RETURNING *
    `

	rows, err := r.DB(tx).Query(ctx, stmt, pgx.NamedArgs{
		"id":         s.ID,
		"user_id":    s.UserID,
		"expires_at": s.ExpiresAt,
	})
	if err != nil {
		// TODO(jozekuhar): session already exists for user?
		return model.Session{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Session])
}

func (r *repository) GetSession(ctx context.Context, sessionID uuid.UUID) (model.Session, error) {
	stmt := `
        SELECT * FROM sessions
		WHERE id = @session_id
    `

	rows, err := r.DB(nil).Query(ctx, stmt, pgx.NamedArgs{
		"session_id": sessionID,
	})
	if err != nil {
		return model.Session{}, err
	}
	defer rows.Close()

	session, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Session])
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Session{}, fmt.Errorf("session: %w", db.ErrNotFound)
	}

	return session, err
}

func (r *repository) DeleteSession(ctx context.Context, tx pgx.Tx, sessionID uuid.UUID) error {
	stmt := `
        DELETE FROM sessions
		WHERE id = @id
    `

	_, err := r.DB(tx).Exec(ctx, stmt, pgx.NamedArgs{
		"id": sessionID,
	})
	return err
}
