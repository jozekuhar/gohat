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
	return r.GetQueries(tx).CreateOrganization(ctx, params)
}

func (r *repository) ListActiveOrganizations(
	ctx context.Context,
	userID uuid.UUID,
) ([]db.Organization, error) {
	return r.GetQueries(nil).ListActiveOrganizations(ctx, userID)
}

func (r *repository) CreateMembership(
	ctx context.Context,
	tx pgx.Tx,
	params db.CreateMembershipParams,
) (db.Membership, error) {
	return r.GetQueries(tx).CreateMembership(ctx, params)
}

func (r *repository) ListMemberships(
	ctx context.Context,
	orgID uuid.UUID,
) ([]db.ListMembershipsRow, error) {
	return r.GetQueries(nil).ListMemberships(ctx, orgID)
}

func (r *repository) GetMembership(
	ctx context.Context,
	params db.GetMembershipParams,
) (db.GetMembershipRow, error) {
	return r.GetQueries(nil).GetMembership(ctx, params)
}

func (r *repository) GetActiveMembershipByUserID(
	ctx context.Context,
	params db.GetActiveMembershipByUserIDParams,
) (db.GetActiveMembershipByUserIDRow, error) {
	return r.GetQueries(nil).GetActiveMembershipByUserID(ctx, params)
}

func (r *repository) CheckMembershipByEmail(
	ctx context.Context,
	params db.CheckMembershipByEmailParams,
) (bool, error) {
	return r.GetQueries(nil).CheckMembershipByEmail(ctx, params)
}

func (r *repository) UpdateMembership(
	ctx context.Context,
	tx pgx.Tx,
	params db.UpdateMembershipParams,
) (db.Membership, error) {
	return r.GetQueries(tx).UpdateMembership(ctx, params)
}

func (r *repository) UpdateMembershipCanceled(
	ctx context.Context,
	tx pgx.Tx,
	params db.UpdateMembershipCanceledParams,
) error {
	return r.GetQueries(tx).UpdateMembershipCanceled(ctx, params)
}

func (r *repository) CreateInvitation(
	ctx context.Context,
	tx pgx.Tx,
	params db.CreateInvitationParams,
) (db.Invitation, error) {
	return r.GetQueries(tx).CreateInvitation(ctx, params)
}

func (r *repository) ListInvitations(
	ctx context.Context,
	orgID uuid.UUID,
) ([]db.Invitation, error) {
	return r.GetQueries(nil).ListInvitations(ctx, orgID)
}

func (r *repository) GetInvitationByTokenHash(
	ctx context.Context,
	tokenHash string,
) (db.GetInvitationByTokenHashRow, error) {
	return r.GetQueries(nil).GetInvitationByTokenHash(ctx, tokenHash)
}

func (r *repository) GetLatestInvitationByEmail(
	ctx context.Context,
	params db.GetLatestInvitationByEmailParams,
) (db.Invitation, error) {
	return r.GetQueries(nil).GetLatestInvitationByEmail(ctx, params)
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
	params db.UpdateInvitationCanceledParams,
) error {
	return r.GetQueries(tx).UpdateInvitationCanceled(ctx, params)
}
