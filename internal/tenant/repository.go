package tenant

import (
	"context"
	"errors"
	"strings"
	"time"
	"uuid"

	"mimokocke/internal/model"
	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/authz"

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

type CreateOrganizationParams struct {
	ID   uuid.UUID
	Name string
	Slug string
}

func (r *repository) CreateOrganization(
	ctx context.Context,
	tx pgx.Tx,
	params db.CreateOrganizationParams,
) (db.Organization, error) {
	return r.Queries.CreateOrganization(ctx, params)
	// stmt := `
	//        INSERT INTO organizations (id, name, slug)
	// 	VALUES (@id, @name, @slug)
	// 	RETURNING *
	//    `
	//
	// rows, err := r.DB(tx).Query(ctx, stmt, pgx.NamedArgs{
	// 	"id":   params.ID,
	// 	"name": params.Name,
	// 	"slug": params.Slug,
	// })
	// if err != nil {
	// 	// TODO(jozekuhar): org with slug already exists
	// 	return model.Organization{}, err
	// }
	// defer rows.Close()
	//
	// return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Organization])
}

type ListOrganizationsFilter struct {
	UserID           *uuid.UUID
	MembershipStatus *model.MembershipStatus
}

func (r *repository) ListOrganizations(
	ctx context.Context,
	filter ListOrganizationsFilter,
) ([]model.Organization, error) {
	stmt := `
		SELECT o.id, o.name, o.slug 
		FROM organizations AS o
		JOIN memberships AS m ON o.id = m.organization_id
		WHERE 1 = 1
    `

	args := pgx.NamedArgs{}

	if filter.UserID != nil {
		stmt += " AND m.user_id = @user_id"
		args["user_id"] = filter.UserID
	}

	if filter.MembershipStatus != nil {
		stmt += " AND m.status = @membership_status"
		args["membership_status"] = filter.MembershipStatus
	}

	rows, err := r.DB(nil).Query(ctx, stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Organization, error) {
		o := model.Organization{}
		err := row.Scan(&o.ID, &o.Name, &o.Slug)
		return o, err
	})
}

type CreateMembershipParams struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	FirstName      string
	LastName       string
	Role           authz.Role
	Permissions    []authz.Permission
	Status         model.MembershipStatus
}

func (r *repository) CreateMembership(
	ctx context.Context,
	tx pgx.Tx,
	params CreateMembershipParams,
) (model.Membership, error) {
	stmt := `
        INSERT INTO memberships (id, organization_id, user_id, first_name, last_name, role, permissions, status)
		VALUES (@id, @organization_id, @user_id, @first_name, @last_name, @role, @permissions, @status)
		RETURNING *
    `

	rows, err := r.DB(tx).Query(ctx, stmt, pgx.NamedArgs{
		"id":              params.ID,
		"organization_id": params.OrganizationID,
		"user_id":         params.UserID,
		"first_name":      params.FirstName,
		"last_name":       params.LastName,
		"role":            params.Role,
		"permissions":     params.Permissions,
		"status":          params.Status,
	})
	if err != nil {
		return model.Membership{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Membership])
}

type ListMembershipsFilter struct {
	UserID *uuid.UUID
	Status *model.MembershipStatus
}

func (r *repository) ListMemberships(
	ctx context.Context,
	orgID uuid.UUID,
	filter ListMembershipsFilter,
) ([]model.Membership, error) {
	stmt := `
		SELECT m.id, m.first_name, m.last_name, m.role, m.status, m.created_at, u.email
		FROM memberships AS m
		JOIN users AS u ON m.user_id = u.id
		WHERE m.organization_id = @organization_id
    `

	args := pgx.NamedArgs{
		"organization_id": orgID,
	}

	if filter.UserID != nil {
		stmt += " AND m.user_id = @user_id"
		args["user_id"] = filter.UserID
	}
	if filter.Status != nil {
		stmt += " AND m.status = @status"
		args["status"] = filter.Status
	}

	rows, err := r.DB(nil).Query(ctx, stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Membership, error) {
		m := model.Membership{
			User: &model.User{},
		}
		err := row.Scan(
			&m.ID,
			&m.FirstName,
			&m.LastName,
			&m.Role,
			&m.Status,
			&m.CreatedAt,
			&m.User.Email,
		)
		return m, err
	})
}

type GetMembershipFilter struct {
	Status *model.MembershipStatus
}

func (r *repository) GetMembership(
	ctx context.Context,
	orgID, userID uuid.UUID,
	filter GetMembershipFilter,
) (model.Membership, error) {
	stmt := `
        SELECT m.id, m.organization_id, m.first_name, m.last_name, m.role, m.permissions, u.email
		FROM memberships AS m
		JOIN users AS u ON m.user_id = u.id
		WHERE m.organization_id = @organization_id
		  AND m.user_id = @user_id
    `

	args := pgx.NamedArgs{
		"organization_id": orgID,
		"user_id":         userID,
	}

	if filter.Status != nil {
		stmt += " AND m.status = @membership_status"
		args["membership_status"] = filter.Status
	}

	rows, err := r.DB(nil).Query(ctx, stmt, args)
	if err != nil {
		return model.Membership{}, err
	}

	m, err := pgx.CollectExactlyOneRow(
		rows,
		func(row pgx.CollectableRow) (model.Membership, error) {
			m := model.Membership{
				User: &model.User{},
			}
			err := row.Scan(
				&m.ID,
				&m.OrganizationID,
				&m.FirstName,
				&m.LastName,
				&m.Role,
				&m.Permissions,
				&m.User.Email,
			)
			return m, err
		},
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Membership{}, db.ErrNotFound
	}
	return m, err
}

// THIS NOT OK START
type activeMembership struct {
	FirstName   string
	LastName    string
	OrgID       uuid.UUID
	OrgName     string
	Role        authz.Role
	Permissions []authz.Permission
}

func (r *repository) GetActiveMemberhip(
	ctx context.Context,
	userID uuid.UUID,
	orgSlug string,
	status model.MembershipStatus,
) (activeMembership, error) {
	stmt := `
		SELECT o.id, o.name, om.first_name, om.last_name, om.role, om.permissions
		FROM organizations AS o
		JOIN memberships AS om ON o.id = om.organization_id
		WHERE o.slug = @slug
		  AND om.user_id = @user_id
		  AND om.status = @status
    `

	var am activeMembership
	err := r.DB(nil).QueryRow(ctx, stmt, pgx.NamedArgs{
		"slug":    orgSlug,
		"user_id": userID,
		"status":  status,
	}).Scan(&am.OrgID, &am.OrgName, &am.FirstName, &am.LastName, &am.Role, &am.Permissions)
	if errors.Is(err, pgx.ErrNoRows) {
		return activeMembership{}, db.ErrNotFound
	}
	return am, err
}

// THIS NOT OK END

func (r *repository) CheckUserIsMember(
	ctx context.Context,
	orgID uuid.UUID,
	userEmail string,
) (bool, error) {
	stmt := `
		SELECT EXISTS (
			SELECT 1 
			FROM memberships AS m
			JOIN users AS u ON m.user_id = u.id 
			WHERE m.organization_id = @organization_id
			  AND u.email = @email
		)
    `

	var exists bool
	err := r.DB(nil).QueryRow(ctx, stmt, pgx.NamedArgs{
		"organization_id": orgID,
		"email":           userEmail,
	}).Scan(&exists)

	return exists, err
}

type UpdateMembershipParams struct {
	FirstName *string
	LastName  *string
}

func (r *repository) UpdateMembership(
	ctx context.Context,
	tx pgx.Tx,
	orgId, membershipID uuid.UUID,
	params UpdateMembershipParams,
) error {
	stmt := `
		UPDATE memberships
		SET {set_clause}
		WHERE organization_id = @organization_id
		  AND id = @membership_id
	`

	args := pgx.NamedArgs{
		"organization_id": orgId,
		"membership_id":   membershipID,
	}

	setClauses := []string{}

	if params.FirstName != nil {
		setClauses = append(setClauses, "first_name = @first_name")
		args["first_name"] = params.FirstName
	}

	if params.LastName != nil {
		setClauses = append(setClauses, "last_name = @last_name")
		args["last_name"] = params.LastName
	}

	if len(setClauses) == 0 {
		return nil
	}

	stmt = strings.Replace(stmt, "{set_clause}", strings.Join(setClauses, ", "), 1)

	cmdTag, err := r.DB(tx).Exec(ctx, stmt, args)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return db.ErrNotFound
	}

	return nil
}

func (r *repository) UpdateMembershipCanceled(
	ctx context.Context,
	tx pgx.Tx,
	orgID uuid.UUID,
	membershipID uuid.UUID,
	canceledByID uuid.UUID,
) error {
	stmt := `
		UPDATE memberships
		SET status = @status,
		    canceled_at = NOW(),
		    canceled_by_id = @canceled_by_id
		WHERE organization_id = @organization_id
		  AND id = @membership_id
		  AND status = @active_status
    `

	res, err := r.DB(tx).Exec(ctx, stmt, pgx.NamedArgs{
		"organization_id": orgID,
		"membership_id":   membershipID,
		"canceled_by_id":  canceledByID,
		"status":          model.MembershipStatusInactive,
		"active_status":   model.MembershipStatusActive,
	})
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return db.ErrNotFound
	}

	return nil
}

type CreateInvitationParams struct {
	ID          uuid.UUID
	OrgID       uuid.UUID
	InviterID   uuid.UUID
	Email       string
	FirstName   string
	LastName    string
	Role        authz.Role
	Permissions []authz.Permission
	TokenHash   string
	ExpiresAt   time.Time
}

func (r *repository) CreateInvitation(
	ctx context.Context,
	tx pgx.Tx,
	params CreateInvitationParams,
) (model.Invitation, error) {
	stmt := `
		INSERT INTO invitations (id, organization_id, inviter_id, email, first_name, last_name, role, permissions, token_hash, expires_at)
		VALUES (@id, @organization_id, @inviter_id, @email, @first_name, @last_name, @role, @permissions, @token_hash, @expires_at)
		RETURNING *
	`

	rows, err := r.DB(tx).Query(ctx, stmt, pgx.NamedArgs{
		"id":              params.ID,
		"organization_id": params.OrgID,
		"inviter_id":      params.InviterID,
		"email":           params.Email,
		"first_name":      params.FirstName,
		"last_name":       params.LastName,
		"role":            params.Role,
		"permissions":     params.Permissions,
		"token_hash":      params.TokenHash,
		"expires_at":      params.ExpiresAt,
	})
	if err != nil {
		return model.Invitation{}, err
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Invitation])
}

func (r *repository) ListInvitations(
	ctx context.Context,
	orgID uuid.UUID,
) ([]model.Invitation, error) {
	stmt := `
		SELECT id, organization_id, email, first_name, last_name, expires_at
		FROM invitations
		WHERE organization_id = @organization_id
    `

	args := pgx.NamedArgs{
		"organization_id": orgID,
	}

	rows, err := r.DB(nil).Query(ctx, stmt, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Invitation, error) {
		i := model.Invitation{}
		err := row.Scan(&i.ID, &i.OrganizationID, &i.Email, &i.FirstName, &i.LastName, &i.ExpiresAt)
		return i, err
	})
}

func (r *repository) GetInvitationByTokenHash(
	ctx context.Context,
	tokenHash string,
) (model.Invitation, error) {
	stmt := `
		SELECT i.id, i.email, o.name, o.slug, u.email, m.first_name, m.last_name
		FROM invitations AS i
		JOIN organizations AS o ON i.organization_id = o.id
		JOIN users AS u ON i.inviter_id = u.id
		JOIN memberships AS m ON i.inviter_id = m.user_id AND i.organization_id = m.organization_id
		WHERE i.token_hash = @token_hash
    `

	rows, err := r.DB(nil).Query(ctx, stmt, pgx.NamedArgs{
		"token_hash": tokenHash,
	})
	if err != nil {
		return model.Invitation{}, err
	}

	i, err := pgx.CollectExactlyOneRow(
		rows,
		func(row pgx.CollectableRow) (model.Invitation, error) {
			i := model.Invitation{
				Organization: &model.Organization{},
				Inviter: &model.User{
					Membership: &model.Membership{},
				},
			}
			err := row.Scan(
				&i.ID, &i.Email,
				&i.Organization.Name, &i.Organization.Slug,
				&i.Inviter.Email,
				&i.Inviter.Membership.FirstName, &i.Inviter.Membership.LastName,
			)
			return i, err
		},
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Invitation{}, db.ErrNotFound
	}
	return i, err
}

func (r *repository) GetLatestInvitation(
	ctx context.Context,
	orgID uuid.UUID,
	userEmail string,
) (model.Invitation, error) {
	stmt := `
		SELECT id, organization_id, email, accepted_at, declined_at, canceled_at, expires_at
		FROM invitations
		WHERE organization_id = @organization_id
		  AND email = @email
		ORDER BY created_at DESC
		LIMIT 1
    `

	rows, err := r.DB(nil).Query(ctx, stmt, pgx.NamedArgs{
		"organization_id": orgID,
		"email":           userEmail,
	})
	if err != nil {
		return model.Invitation{}, err
	}

	i, err := pgx.CollectExactlyOneRow(
		rows,
		func(row pgx.CollectableRow) (model.Invitation, error) {
			i := model.Invitation{}
			err := row.Scan(
				&i.ID,
				&i.OrganizationID,
				&i.Email,
				&i.AcceptedAt,
				&i.DeclinedAt,
				&i.CanceledAt,
				&i.ExpiresAt,
			)
			return i, err
		},
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Invitation{}, db.ErrNotFound
	}
	return i, err
}

func (r *repository) UpdateInvitationAcceptedAt(
	ctx context.Context,
	tx pgx.Tx,
	invitationID uuid.UUID,
) error {
	stmt := `
		UPDATE invitations
		SET accepted_at = NOW()
		WHERE id = @invitation_id
		  AND accepted_at IS NULL
		  AND declined_at IS NULL
		  AND canceled_at IS NULL
	`

	cmdTag, err := r.DB(tx).Exec(ctx, stmt, pgx.NamedArgs{
		"invitation_id": invitationID,
	})
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return db.ErrNotFound
	}

	return nil
}

func (r *repository) UpdateInvitationDeclinedAt(
	ctx context.Context,
	tx pgx.Tx,
	invitationID uuid.UUID,
) error {
	stmt := `
		UPDATE invitations
		SET declined_at = NOW()
		WHERE id = @invitation_id
		  AND accepted_at IS NULL
		  AND declined_at IS NULL
		  AND canceled_at IS NULL
	`

	cmdTag, err := r.DB(tx).Exec(ctx, stmt, pgx.NamedArgs{
		"invitation_id": invitationID,
	})
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return db.ErrNotFound
	}

	return nil
}

func (r *repository) UpdateInvitationCanceled(
	ctx context.Context,
	tx pgx.Tx,
	orgID uuid.UUID,
	invitationID uuid.UUID,
	canceledByID uuid.UUID,
) error {
	stmt := `
		UPDATE invitations
		SET canceled_at = NOW(),
		    canceled_by_id = @canceled_by_id
		WHERE organization_id = @organization_id
		  AND id = @invitation_id
		  AND accepted_at IS NULL
		  AND declined_at IS NULL
		  AND canceled_at IS NULL
    `

	res, err := r.DB(tx).Exec(ctx, stmt, pgx.NamedArgs{
		"organization_id": orgID,
		"invitation_id":   invitationID,
		"canceled_by_id":  canceledByID,
	})
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return db.ErrNotFound
	}

	return nil
}
