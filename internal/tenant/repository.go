package tenant

import (
	"context"
	"errors"
	"uuid"

	"mimokocke/internal/shared/authz"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotExist = errors.New("user not exist")

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) CreateOrganization(
	ctx context.Context,
	tx pgx.Tx,
	id uuid.UUID,
	name string,
	slug string,
) (*Organization, error) {
	stmt := `
        INSERT INTO organizations (id, name, slug)
		VALUES (@id, @name, @slug)
		RETURNING *
    `

	query := r.pool.Query
	if tx != nil {
		query = tx.Query
	}

	rows, err := query(ctx, stmt, pgx.NamedArgs{
		"id":   id,
		"name": name,
		"slug": slug,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Organization])
}

func (r *repository) ListOrganizationsForUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]Organization, error) {
	stmt := `
		SELECT * 
		FROM organizations
    `

	rows, err := r.pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Organization])
}

func (r *repository) CreateMembership(
	ctx context.Context,
	tx pgx.Tx,
	id, organizationID, userID uuid.UUID,
	firstName, lastName string,
	role membershipRole,
	permissions []string,
	status membershipStatus,
) (*Membership, error) {
	stmt := `
        INSERT INTO memberships (id, organization_id, user_id, first_name, last_name, role, permissions, status)
		VALUES (@id, @organization_id, @user_id, @first_name, @last_name, @role, @permissions, @status)
		RETURNING *
    `

	query := r.pool.Query
	if tx != nil {
		query = tx.Query
	}

	rows, err := query(ctx, stmt, pgx.NamedArgs{
		"id":              id,
		"organization_id": organizationID,
		"user_id":         userID,
		"first_name":      firstName,
		"last_name":       lastName,
		"role":            role,
		"permissions":     permissions,
		"status":          status,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Membership])
}

func (r *repository) ListActiveMemberships(
	ctx context.Context,
	userID uuid.UUID,
) ([]Membership, error) {
	stmt := `
		SELECT * 
		FROM memberships
		WHERE user_id = @user_id
		  AND status = @status
    `

	rows, err := r.pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
		"status":  MembershipStatusActive,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Membership])
}

func (r *repository) ListMemberships(ctx context.Context, orgID uuid.UUID) ([]Membership, error) {
	stmt := `
		SELECT * 
		FROM memberships 
		WHERE organization_id = @organization_id
	`

	rows, err := r.pool.Query(ctx, stmt, pgx.NamedArgs{
		"organization_id": orgID,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Membership])
}

func (r *repository) GetIdentity(
	ctx context.Context,
	userID uuid.UUID,
	orgSlug string,
	status membershipStatus,
) (*authz.Identity, error) {
	stmt := `
		SELECT o.id, o.slug, om.role, om.permissions
		FROM organizations AS o
		LEFT JOIN memberships AS om ON o.id = om.organization_id
		WHERE o.slug = @slug
		  AND om.user_id = @user_id
		  AND om.status = @status
    `

	i := &authz.Identity{}
	err := r.pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"slug":    orgSlug,
		"user_id": userID,
		"status":  status,
	}).Scan(&i.OrganizationID, &i.OrganizationSlug, &i.Role, &i.Permissions)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMembershipNotFound
		}
		return nil, err
	}

	return i, nil
}

func (r *repository) CreateInvitation(
	ctx context.Context,
	tx pgx.Tx,
	id uuid.UUID,
	tokenHash string,
) (*Invitation, error) {
	stmt := `
		INSERT INTO invitations (id, token_hash)
		VALUES (@id, @token_hash)
		RETURNING *
	`

	query := r.pool.Query
	if tx != nil {
		query = tx.Query
	}

	rows, err := query(ctx, stmt, pgx.NamedArgs{
		"id":         id,
		"token_hash": tokenHash,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Invitation])
}
