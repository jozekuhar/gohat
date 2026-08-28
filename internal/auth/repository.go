package auth

import (
	"context"
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

func (r *repository) CreateUser(ctx context.Context, tx pgx.Tx, u User) (User, error) {
	stmt := `
        INSERT INTO users (id, email)
		VALUES (@id, @email)
		RETURNING *
    `

	queryFn := r.pool.Query
	if tx != nil {
		queryFn = tx.Query
	}

	rows, err := queryFn(ctx, stmt, pgx.NamedArgs{
		"id":    u.ID,
		"email": u.Email,
	})
	if err != nil {
		return User{}, err
	}
	defer rows.Close()

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		// TODO(jozekuhar): check error and return ErrUserAlreadyExiss
		return User{}, err
	}

	return user, nil
}

func (r *repository) CreateAuthentication(
	ctx context.Context,
	tx pgx.Tx,
	a Authentication,
) (Authentication, error) {
	stmt := `
        INSERT INTO authentications (id, user_id, provider, provider_id, password_hash)
		VALUES (@id, @user_id, @provider, @provider_id, @password_hash)
		RETURNING *
    `

	queryFn := r.pool.Query
	if tx != nil {
		queryFn = tx.Query
	}

	rows, err := queryFn(ctx, stmt, pgx.NamedArgs{
		"id":            a.ID,
		"user_id":       a.UserID,
		"provider":      a.Provider,
		"provider_id":   a.ProviderID,
		"password_hash": a.PasswordHash,
	})
	if err != nil {
		return Authentication{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Authentication])
}

func (r *repository) GetAuthenticationByEmail(
	ctx context.Context,
	email string,
	provider AuthProvider,
) (Authentication, error) {
	stmt := `
		SELECT a.*
		FROM authentications a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE u.email = @email
		  AND a.provider = @provider
    `

	rows, err := r.pool.Query(ctx, stmt, pgx.NamedArgs{
		"email":    email,
		"provider": provider,
	})
	if err != nil {
		return Authentication{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Authentication])
}

func (r *repository) GetAuthenticationByProvider(
	ctx context.Context,
	provider AuthProvider,
	providerID string,
) (Authentication, error) {
	stmt := `
		SELECT *
		FROM authentications
		WHERE provider = @provider
		  AND provider_id = @provider_id
    `

	rows, err := r.pool.Query(ctx, stmt, pgx.NamedArgs{
		"provider":    provider,
		"provider_id": providerID,
	})
	if err != nil {
		return Authentication{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Authentication])
}

func (r *repository) CreateSession(ctx context.Context, tx pgx.Tx, s Session) (Session, error) {
	stmt := `
        INSERT INTO sessions (id, user_id, expires_at)
		VALUES (@id, @user_id, @expires_at)
		RETURNING *
    `

	queryFn := r.pool.Query
	if tx != nil {
		queryFn = tx.Query
	}

	rows, err := queryFn(ctx, stmt, pgx.NamedArgs{
		"id":         s.ID,
		"user_id":    s.UserID,
		"expires_at": s.ExpiresAt,
	})
	if err != nil {
		return Session{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Session])
}

func (r *repository) GetSession(ctx context.Context, sessionID uuid.UUID) (Session, error) {
	stmt := `
        SELECT * FROM sessions
		WHERE id = @session_id
    `

	rows, err := r.pool.Query(ctx, stmt, pgx.NamedArgs{
		"session_id": sessionID,
	})
	if err != nil {
		return Session{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Session])
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
