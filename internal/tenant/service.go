package tenant

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"uuid"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/authz"
	"mimokocke/internal/shared/clock"
	"mimokocke/internal/shared/config"
	"mimokocke/internal/shared/routes"

	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/resend/resend-go/v3"
)

var (
	ErrOrganizationLimitReached = errors.New("user has reached maximum organizasion limit")
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

func (s *Service) RegisterOrganization(
	ctx context.Context,
	userID uuid.UUID,
	orgName string,
	orgSlug string,
	firstName string,
	lastName string,
) (db.Organization, error) {
	activeOrgs, err := s.tenantRepo.ListActiveOrganizations(ctx, userID)
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
	err = pgx.BeginFunc(ctx, s.tenantRepo.Pool, func(tx pgx.Tx) error {
		organization, err = s.tenantRepo.CreateOrganization(ctx, tx, db.CreateOrganizationParams{
			ID:   uuid.NewV7(),
			Name: orgName,
			Slug: slug.Make(orgSlug),
		})
		if err != nil {
			return fmt.Errorf("creating organization: %w", err)
		}

		_, err = s.tenantRepo.CreateMembership(ctx, tx, db.CreateMembershipParams{
			ID:             uuid.NewV7(),
			OrganizationID: organization.ID,
			UserID:         userID,
			FirstName:      firstName,
			LastName:       lastName,
			Role:           db.RoleOwner,
			Permissions:    []db.MembershipPermission{},
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

// not good
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

type MembershipsData struct {
	Memberhips  []db.ListMembershipsRow
	Invitations []db.Invitation
}

// not good
func (s *Service) GetMembershipsData(
	ctx context.Context,
	identity authz.Identity,
) (MembershipsData, error) {
	if !identity.HasPermission(db.PermMembershipRead) {
		return MembershipsData{}, authz.ErrPermissionDenied
	}

	var data MembershipsData
	var err error

	data.Memberhips, err = s.tenantRepo.ListMemberships(
		ctx,
		identity.OrgID,
	)
	if err != nil {
		return MembershipsData{}, err
	}

	data.Invitations, err = s.tenantRepo.ListInvitations(ctx, identity.OrgID)

	return data, err
}

// not good
func (s *Service) GetMembershipDetails(
	ctx context.Context,
	identity authz.Identity,
	membershipID uuid.UUID,
) (db.GetMembershipRow, error) {
	// todo: permissions
	if identity.HasPermission(db.PermMembershipRead) {
		return db.GetMembershipRow{}, authz.ErrPermissionDenied
	}

	return s.tenantRepo.GetMembership(ctx, db.GetMembershipParams{
		OrganizationID: identity.OrgID,
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
	identity authz.Identity,
	membershipID uuid.UUID,
) error {
	membership, err := s.tenantRepo.GetMembership(ctx, db.GetMembershipParams{
		OrganizationID: identity.OrgID,
		MembershipID:   membershipID,
	})
	if err != nil {
		return err
	}

	isOwnMembership := membership.Membership.UserID == identity.ID
	hasDeletePermission := identity.HasPermission(db.PermMembershipDelete)

	if !isOwnMembership && !hasDeletePermission {
		return authz.ErrPermissionDenied
	}

	return s.tenantRepo.UpdateMembershipCanceled(
		ctx,
		nil,
		db.UpdateMembershipCanceledParams{
			CanceledByID:   &identity.ID,
			OrganizationID: identity.OrgID,
			MembershipID:   membershipID,
		},
	)
}

type InviteUserParams struct {
	Email       string
	FirstName   string
	LastName    string
	Role        db.MembershipRole
	Permissions []db.MembershipPermission
}

// InviteUser handles sending invitation to user to join organization.
func (s *Service) InviteUser(
	ctx context.Context,
	identity authz.Identity,
	params InviteUserParams,
) (db.Invitation, error) {
	if !identity.HasPermission(db.PermMembershipCreate) {
		return db.Invitation{}, authz.ErrPermissionDenied
	}

	isMember, err := s.tenantRepo.CheckMembershipByEmail(ctx, db.CheckMembershipByEmailParams{
		OrganizationID: identity.OrgID,
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
			OrganizationID: identity.OrgID,
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
		OrganizationID: identity.OrgID,
		InviterID:      identity.ID,
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

// CancelInvite maeks invitation cancelled.
func (s *Service) CancelInvite(
	ctx context.Context,
	identity authz.Identity,
	invitationID uuid.UUID,
) error {
	if !identity.HasPermission(db.PermMembershipDelete) {
		return authz.ErrPermissionDenied
	}

	return s.tenantRepo.UpdateInvitationCanceled(ctx, nil, db.UpdateInvitationCanceledParams{
		CanceledByID:   &identity.ID,
		OrganizationID: identity.OrgID,
		InvitationID:   invitationID,
	})
}

// not good
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

	// if !invitation.IsPending() {
	// 	return ErrInvitationNotPending
	// }

	if userEmail != invitation.Invitation.Email {
		return fmt.Errorf("user email not same as invitation email")
	}

	return pgx.BeginFunc(ctx, s.tenantRepo.Pool, func(tx pgx.Tx) error {
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
