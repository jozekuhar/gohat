package tenant

import (
	"context"
	"errors"
	"fmt"
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
) (Organization, error) {
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
		return Organization{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Organization])
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
	role string,
	permissions []string,
	status membershipStatus,
) (Membership, error) {
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
		return Membership{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Membership])
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
		"status":  membershipStatusActive,
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

type activeMembership struct {
	OrganizationID uuid.UUID
	Role           authz.Role
	Permissions    []authz.Permission
}

func (r *repository) GetActiveMemberhip(
	ctx context.Context,
	userID uuid.UUID,
	orgSlug string,
	status membershipStatus,
) (activeMembership, error) {
	stmt := `
		SELECT o.id, om.role, om.permissions
		FROM organizations AS o
		LEFT JOIN memberships AS om ON o.id = om.organization_id
		WHERE o.slug = @slug
		  AND om.user_id = @user_id
		  AND om.status = @status
    `

	var am activeMembership
	err := r.pool.QueryRow(ctx, stmt, pgx.NamedArgs{
		"slug":    orgSlug,
		"user_id": userID,
		"status":  status,
	}).Scan(&am.OrganizationID, &am.Role, &am.Permissions)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return activeMembership{}, ErrMembershipNotFound
		}
		return activeMembership{}, err
	}

	return am, nil
}

func (r *repository) CreateInvitation(
	ctx context.Context,
	tx pgx.Tx,
	i Invitation,
) (Invitation, error) {
	stmt := `
		INSERT INTO invitations (id, organization_id, email, first_name, last_name, role, permissions, token_hash, expires_at)
		VALUES (@id, @organization_id, @email, @first_name, @last_name, @role, @permissions, @token_hash, @expires_at)
		RETURNING *
	`

	query := r.pool.Query
	if tx != nil {
		query = tx.Query
	}

	rows, err := query(ctx, stmt, pgx.NamedArgs{
		"id":              i.ID,
		"organization_id": i.OrganizationID,
		"email":           i.Email,
		"first_name":      i.FirstName,
		"last_name":       i.LastName,
		"role":            i.Role,
		"permissions":     i.Permissions,
		"token_hash":      i.TokenHash,
		"expires_at":      i.ExpiresAt,
	})
	if err != nil {
		fmt.Println("error is here")
		return Invitation{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Invitation])
}

func (r *repository) ListInvitations(ctx context.Context, orgID uuid.UUID) ([]Invitation, error) {
	stmt := `
		SELECT * FROM invitations
		WHERE organization_id = @organization_id
    `

	rows, err := r.pool.Query(ctx, stmt, pgx.NamedArgs{
		"organization_id": orgID,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Invitation])
}

func (r *repository) DeleteInvitation(ctx context.Context, id, orgID uuid.UUID) error {
	stmt := `
		DELETE FROM invitations
		WHERE id = @id 
		  AND organization_id = @organization_id
    `

	_, err := r.pool.Exec(ctx, stmt, pgx.NamedArgs{
		"id":              id,
		"organization_id": orgID,
	})
	return err
}
