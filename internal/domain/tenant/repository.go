package tenant

import (
	"context"
	"errors"
	"fmt"
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

type CreateOrganizationParams struct {
	ID   uuid.UUID
	Name string
	Slug string
}

func (r *repository) CreateOrganization(
	ctx context.Context,
	tx pgx.Tx,
	params sqlc.CreateOrganizationParams,
) (sqlc.Organization, error) {
	return r.GetQueries(tx).CreateOrganization(ctx, params)
}

func (r *repository) ListActiveOrganizations(
	ctx context.Context,
	userID uuid.UUID,
) ([]sqlc.Organization, error) {
	return r.GetQueries(nil).ListActiveOrganizations(ctx, userID)
}

func (r *repository) CreateMembership(
	ctx context.Context,
	tx pgx.Tx,
	params sqlc.CreateMembershipParams,
) (sqlc.Membership, error) {
	return r.GetQueries(tx).CreateMembership(ctx, params)
}

func (r *repository) ListMemberships(
	ctx context.Context,
	orgID uuid.UUID,
) ([]sqlc.ListMembershipsRow, error) {
	return r.GetQueries(nil).ListMemberships(ctx, orgID)
}

func (r *repository) GetMembership(
	ctx context.Context,
	params sqlc.GetMembershipParams,
) (sqlc.GetMembershipRow, error) {
	return r.GetQueries(nil).GetMembership(ctx, params)
}

func (r *repository) GetActiveMembershipByUserID(
	ctx context.Context,
	params sqlc.GetActiveMembershipByUserIDParams,
) (sqlc.GetActiveMembershipByUserIDRow, error) {
	return r.GetQueries(nil).GetActiveMembershipByUserID(ctx, params)
}

func (r *repository) CheckMembershipByEmail(
	ctx context.Context,
	params sqlc.CheckMembershipByEmailParams,
) (bool, error) {
	return r.GetQueries(nil).CheckMembershipByEmail(ctx, params)
}

func (r *repository) UpdateMembership(
	ctx context.Context,
	tx pgx.Tx,
	params sqlc.UpdateMembershipParams,
) (sqlc.Membership, error) {
	return r.GetQueries(tx).UpdateMembership(ctx, params)
}

func (r *repository) UpdateMembershipCanceled(
	ctx context.Context,
	tx pgx.Tx,
	params sqlc.UpdateMembershipCanceledParams,
) error {
	return r.GetQueries(tx).UpdateMembershipCanceled(ctx, params)
}

func (r *repository) CreateInvitation(
	ctx context.Context,
	tx pgx.Tx,
	params sqlc.CreateInvitationParams,
) (sqlc.Invitation, error) {
	return r.GetQueries(tx).CreateInvitation(ctx, params)
}

func (r *repository) ListInvitations(
	ctx context.Context,
	orgID uuid.UUID,
) ([]sqlc.Invitation, error) {
	return r.GetQueries(nil).ListInvitations(ctx, orgID)
}

func (r *repository) GetInvitationByTokenHash(
	ctx context.Context,
	tokenHash string,
) (sqlc.GetInvitationByTokenHashRow, error) {
	return r.GetQueries(nil).GetInvitationByTokenHash(ctx, tokenHash)
}

func (r *repository) GetLatestInvitationByEmail(
	ctx context.Context,
	params sqlc.GetLatestInvitationByEmailParams,
) (sqlc.Invitation, error) {
	i, err := r.GetQueries(nil).GetLatestInvitationByEmail(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.Invitation{}, db.ErrNotFound
		}
		return sqlc.Invitation{}, fmt.Errorf("GetLatestInvitationByEmail: %w", err)
	}
	return i, nil
}

func (r *repository) UpdateInvitationAcceptedAt(
	ctx context.Context,
	tx pgx.Tx,
	invitationID uuid.UUID,
) error {
	return r.GetQueries(tx).UpdateInvitationAcceptedAt(ctx, invitationID)
}

func (r *repository) UpdateInvitationDeclinedAt(
	ctx context.Context,
	tx pgx.Tx,
	invitationID uuid.UUID,
) error {
	return r.GetQueries(tx).UpdateInvitationDeclinedAt(ctx, invitationID)
}

func (r *repository) UpdateInvitationCanceled(
	ctx context.Context,
	tx pgx.Tx,
	params sqlc.UpdateInvitationCanceledParams,
) error {
	return r.GetQueries(tx).UpdateInvitationCanceled(ctx, params)
}

func (r *repository) DeleteInvitation(
	ctx context.Context,
	tx pgx.Tx,
	params sqlc.DeleteInvitationParams,
) error {
	return r.GetQueries(tx).DeleteInvitation(ctx, params)
}
