package tenant

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"uuid"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/clock"
	"mimokocke/internal/shared/config"
	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/permissions"
	"mimokocke/internal/shared/routes"

	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/resend/resend-go/v3"
)

var (
	ErrOrganizationLimitReached = errors.New("user has reached maximum organization limit")
	ErrInvitationNotPending     = errors.New("invitation is not pending")
	ErrInvitationAlreadyPending = errors.New("invitation is already pending")
	ErrUserAlreadyMember        = errors.New("user already member of organization")
)

type Service struct {
	cfg          *config.Config
	logger       *slog.Logger
	clock        clock.Clock
	resendClient *resend.Client
	tenantRepo   *repository
}

func NewService(
	cfg *config.Config,
	logger *slog.Logger,
	clock clock.Clock,
	resendClient *resend.Client,
	tenantRepo *repository,
) *Service {
	return &Service{
		cfg:          cfg,
		logger:       logger,
		clock:        clock,
		resendClient: resendClient,
		tenantRepo:   tenantRepo,
	}
}

// not good
func (s *Service) ListActiveOrganizations(
	ctx context.Context,
	userID uuid.UUID,
) ([]db.Organization, error) {
	return s.tenantRepo.ListActiveOrganizations(ctx, userID)
}

type RegisterOrganizationParams struct {
	UserID    uuid.UUID
	OrgName   string
	OrgSlug   string
	FirstName string
	LastName  string
}

func (s *Service) RegisterOrganization(
	ctx context.Context,
	params RegisterOrganizationParams,
) (db.Organization, error) {
	activeOrgs, err := s.tenantRepo.ListActiveOrganizations(ctx, params.UserID)
	if err != nil {
		return db.Organization{}, fmt.Errorf(
			"listing organization memberships for user: %w",
			err,
		)
	}
	if len(activeOrgs) > 0 {
		return db.Organization{}, ErrOrganizationLimitReached
	}

	var organization db.Organization
	err = pgx.BeginFunc(ctx, s.tenantRepo.Pool(), func(tx pgx.Tx) error {
		organization, err = s.tenantRepo.CreateOrganization(ctx, tx, db.CreateOrganizationParams{
			ID:   uuid.NewV7(),
			Name: params.OrgName,
			Slug: slug.Make(params.OrgSlug),
		})
		if err != nil {
			return fmt.Errorf("creating organization: %w", err)
		}

		_, err = s.tenantRepo.CreateMembership(ctx, tx, db.CreateMembershipParams{
			ID:             uuid.NewV7(),
			OrganizationID: organization.ID,
			UserID:         params.UserID,
			FirstName:      params.FirstName,
			LastName:       params.LastName,
			Role:           permissions.RoleOwner,
			Permissions:    []permissions.MembershipPermission{},
			Status:         db.MemberStatusActive,
		},
		)
		if err != nil {
			return fmt.Errorf("creating organization membership: %w", err)
		}

		return nil
	})

	return organization, err
}

func (s *Service) GetActiveMembership(
	ctx context.Context,
	userID uuid.UUID,
	orgSlug string,
) (db.GetActiveMembershipByUserIDRow, error) {
	return s.tenantRepo.GetActiveMembershipByUserID(ctx, db.GetActiveMembershipByUserIDParams{
		Slug:   orgSlug,
		UserID: userID,
	})
}

type MembershipsOverview struct {
	Memberhips  []db.ListMembershipsRow
	Invitations []db.Invitation
}

func (s *Service) GetMembershipsOverview(
	ctx context.Context,
	ident identity.IdentityCtx,
) (MembershipsOverview, error) {
	if !ident.HasPermission(permissions.MembershipRead) {
		return MembershipsOverview{}, permissions.ErrDenied
	}

	var data MembershipsOverview
	var err error

	data.Memberhips, err = s.tenantRepo.ListMemberships(
		ctx,
		ident.OrgID,
	)
	if err != nil {
		return MembershipsOverview{}, err
	}

	data.Invitations, err = s.tenantRepo.ListInvitations(ctx, ident.OrgID)

	return data, err
}

func (s *Service) GetMembershipDetails(
	ctx context.Context,
	ident identity.IdentityCtx,
	membershipID uuid.UUID,
) (db.GetMembershipRow, error) {
	if !ident.HasPermission(permissions.MembershipRead) {
		return db.GetMembershipRow{}, permissions.ErrDenied
	}

	return s.tenantRepo.GetMembership(ctx, db.GetMembershipParams{
		OrganizationID: ident.OrgID,
		MembershipID:   membershipID,
	})
}

// not good: kaj za vraga ta service updetja?
func (s *Service) UpdateMembership(ctx context.Context, m db.Membership) (db.Membership, error) {
	return s.tenantRepo.UpdateMembership(ctx, nil, db.UpdateMembershipParams{
		FirstName:      &m.FirstName,
		LastName:       &m.LastName,
		OrganizationID: m.OrganizationID,
		MembershipID:   m.ID,
	})
}

// CancelMembership cancels a membership.
func (s *Service) CancelMembership(
	ctx context.Context,
	ident identity.IdentityCtx,
	membershipID uuid.UUID,
) error {
	membership, err := s.tenantRepo.GetMembership(ctx, db.GetMembershipParams{
		OrganizationID: ident.OrgID,
		MembershipID:   membershipID,
	})
	if err != nil {
		return err
	}

	isOwnMembership := membership.Membership.UserID == ident.ID
	hasDeletePermission := ident.HasPermission(permissions.MembershipDelete)

	if !isOwnMembership && !hasDeletePermission {
		return permissions.ErrDenied
	}

	return s.tenantRepo.UpdateMembershipCanceled(
		ctx,
		nil,
		db.UpdateMembershipCanceledParams{
			CanceledByID:   &ident.ID,
			OrganizationID: ident.OrgID,
			MembershipID:   membershipID,
		},
	)
}

type InviteUserParams struct {
	Email       string
	FirstName   string
	LastName    string
	Role        permissions.MembershipRole
	Permissions []permissions.MembershipPermission
}

// InviteUser handles sending invitation to user to join organization.
func (s *Service) InviteUser(
	ctx context.Context,
	ident identity.IdentityCtx,
	params InviteUserParams,
) (db.Invitation, error) {
	if !ident.HasPermission(permissions.MembershipCreate) {
		return db.Invitation{}, permissions.ErrDenied
	}

	isMember, err := s.tenantRepo.CheckMembershipByEmail(ctx, db.CheckMembershipByEmailParams{
		OrganizationID: ident.OrgID,
		Email:          params.Email,
	})
	if err != nil {
		return db.Invitation{}, fmt.Errorf("checking membership: %w", err)
	}
	if isMember {
		return db.Invitation{}, ErrUserAlreadyMember
	}

	latestInvitation, err := s.tenantRepo.GetLatestInvitationByEmail(
		ctx,
		db.GetLatestInvitationByEmailParams{
			OrganizationID: ident.OrgID,
			Email:          params.Email,
		},
	)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return db.Invitation{}, fmt.Errorf("getting latest invitation: %w", err)
	}
	if err == nil && latestInvitation.IsPending() {
		return db.Invitation{}, ErrInvitationAlreadyPending
	}

	rawToken, tokenHash, err := generateInvitationToken()
	if err != nil {
		return db.Invitation{}, err
	}

	invitation, err := s.tenantRepo.CreateInvitation(ctx, nil, db.CreateInvitationParams{
		ID:             uuid.NewV7(),
		OrganizationID: ident.OrgID,
		InviterID:      ident.ID,
		Email:          params.Email,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		Role:           params.Role,
		Permissions:    params.Permissions,
		TokenHash:      tokenHash,
		ExpiresAt:      s.clock.NowUTC().Add(24 * time.Hour),
	})
	if err != nil {
		return db.Invitation{}, fmt.Errorf("creating invitation: %w", err)
	}

	invitationURL := fmt.Sprintf(
		"http://%s%s",
		s.cfg.Host,
		fmt.Sprintf(routes.InvitationsJoin, rawToken),
	)

	emailParams := &resend.SendEmailRequest{
		From:    "invites@gajogroup.com",
		To:      []string{invitation.Email},
		Subject: "Invitation",
		ReplyTo: "support@gajogroup.com",
		Text:    "Hi, you have been invited to our organization. URL: " + invitationURL,
		Tags: []resend.Tag{
			{Name: "type", Value: "invite"},
		},
	}
	_, err = s.resendClient.Emails.Send(emailParams)
	return invitation, err
}

// CancelInvite cancels invite.
func (s *Service) CancelInvite(
	ctx context.Context,
	ident identity.IdentityCtx,
	invitationID uuid.UUID,
) error {
	if !ident.HasPermission(permissions.MembershipDelete) {
		return permissions.ErrDenied
	}

	return s.tenantRepo.UpdateInvitationCanceled(ctx, nil, db.UpdateInvitationCanceledParams{
		CanceledByID:   &ident.ID,
		OrganizationID: ident.OrgID,
		InvitationID:   invitationID,
	})
}

// GetInvitation retrieves invite.
func (s *Service) GetInvitation(
	ctx context.Context,
	token string,
) (db.GetInvitationByTokenHashRow, error) {
	tokenHash := hashToken(token)
	return s.tenantRepo.GetInvitationByTokenHash(ctx, tokenHash)
}

// AcceptInvitation makes user accept invitation to join organization.
func (s *Service) AcceptInvitation(
	ctx context.Context,
	userID uuid.UUID,
	userEmail, token string,
) error {
	tokenHash := hashToken(token)

	invitation, err := s.tenantRepo.GetInvitationByTokenHash(ctx, tokenHash)
	if err != nil {
		return err
	}

	if !invitation.Invitation.IsPending() {
		return ErrInvitationNotPending
	}

	if userEmail != invitation.Invitation.Email {
		return fmt.Errorf("user email not same as invitation email")
	}

	return pgx.BeginFunc(ctx, s.tenantRepo.Pool(), func(tx pgx.Tx) error {
		err = s.tenantRepo.UpdateInvitationAcceptedAt(ctx, tx, invitation.Invitation.ID)
		if err != nil {
			return err
		}

		_, err = s.tenantRepo.CreateMembership(ctx, tx, db.CreateMembershipParams{
			ID:             uuid.NewV7(),
			OrganizationID: invitation.Organization.ID,
			UserID:         userID,
			FirstName:      invitation.Invitation.FirstName,
			LastName:       invitation.Invitation.LastName,
			Role:           invitation.Invitation.Role,
			Permissions:    invitation.Invitation.Permissions,
			Status:         db.MemberStatusActive,
		})
		return err
	})
}

// DeclineInvitation declines invitation to join organization.
func (s *Service) DeclineInvitation(ctx context.Context, userEmail, token string) error {
	tokenHash := hashToken(token)

	invitation, err := s.tenantRepo.GetInvitationByTokenHash(ctx, tokenHash)
	if err != nil {
		return err
	}

	if !invitation.Invitation.IsPending() {
		return ErrInvitationNotPending
	}

	if userEmail != invitation.Invitation.Email {
		return fmt.Errorf("user email not same as invitation email")
	}

	return s.tenantRepo.UpdateInvitationDeclinedAt(ctx, nil, invitation.Invitation.ID)
}
