package tenant

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
}

func (r *repository) ListActiveOrganizations(
	ctx context.Context,
	userID uuid.UUID,
) ([]db.Organization, error) {
	return r.Queries.ListActiveOrganizations(ctx, userID)
}

func (r *repository) CreateMembership(
	ctx context.Context,
	tx pgx.Tx,
	params db.CreateMembershipParams,
) (db.Membership, error) {
	return r.Queries.CreateMembership(ctx, params)
}

func (r *repository) ListMemberships(
	ctx context.Context,
	orgID uuid.UUID,
) ([]db.ListMembershipsRow, error) {
	return r.Queries.ListMemberships(ctx, orgID)
}

func (r *repository) GetMembership(
	ctx context.Context,
	params db.GetMembershipParams,
) (db.GetMembershipRow, error) {
	return r.Queries.GetMembership(ctx, params)
}

func (r *repository) GetActiveMembershipByUserID(
	ctx context.Context,
	params db.GetActiveMembershipByUserIDParams,
) (db.GetActiveMembershipByUserIDRow, error) {
	return r.Queries.GetActiveMembershipByUserID(ctx, params)
}

func (r *repository) CheckMembershipByEmail(
	ctx context.Context,
	params db.CheckMembershipByEmailParams,
) (bool, error) {
	return r.Queries.CheckMembershipByEmail(ctx, params)
}

func (r *repository) UpdateMembership(
	ctx context.Context,
	tx pgx.Tx,
	params db.UpdateMembershipParams,
) (db.Membership, error) {
	return r.Queries.UpdateMembership(ctx, params)
}

func (r *repository) UpdateMembershipCanceled(
	ctx context.Context,
	tx pgx.Tx,
	params db.UpdateMembershipCanceledParams,
) error {
	return r.Queries.UpdateMembershipCanceled(ctx, params)
}

func (r *repository) CreateInvitation(
	ctx context.Context,
	tx pgx.Tx,
	params db.CreateInvitationParams,
) (db.Invitation, error) {
	return r.Queries.CreateInvitation(ctx, params)
}

func (r *repository) ListInvitations(
	ctx context.Context,
	orgID uuid.UUID,
) ([]db.Invitation, error) {
	return r.Queries.ListInvitations(ctx, orgID)
}

func (r *repository) GetInvitationByTokenHash(
	ctx context.Context,
	tokenHash string,
) (db.GetInvitationByTokenHashRow, error) {
	return r.Queries.GetInvitationByTokenHash(ctx, tokenHash)
}

func (r *repository) GetLatestInvitationByEmail(
	ctx context.Context,
	params db.GetLatestInvitationByEmailParams,
) (db.Invitation, error) {
	return r.Queries.GetLatestInvitationByEmail(ctx, params)
}

func (r *repository) UpdateInvitationAcceptedAt(
	ctx context.Context,
	tx pgx.Tx,
	invitationID uuid.UUID,
) error {
	return r.Queries.UpdateInvitationAcceptedAt(ctx, invitationID)
}

func (r *repository) UpdateInvitationDeclinedAt(
	ctx context.Context,
	tx pgx.Tx,
	invitationID uuid.UUID,
) error {
	return r.Queries.UpdateInvitationDeclinedAt(ctx, invitationID)
}

func (r *repository) UpdateInvitationCanceled(
	ctx context.Context,
	tx pgx.Tx,
	params db.UpdateInvitationCanceledParams,
) error {
	return r.Queries.UpdateInvitationCanceled(ctx, params)
}
