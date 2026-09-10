package tenant

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"uuid"

	"mimokocke/internal/model"
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
) ([]model.Organization, error) {
	activeStatus := model.MembershipStatusActive
	return s.tenantRepo.ListOrganizations(ctx, ListOrganizationsFilter{
		UserID:           &userID,
		MembershipStatus: &activeStatus,
	})
}

func (s *Service) RegisterOrganization(
	ctx context.Context,
	userID uuid.UUID,
	orgName string,
	orgSlug string,
	firstName string,
	lastName string,
) (db.Organization, error) {
	activeStatus := model.MembershipStatusActive
	activeMemberships, err := s.tenantRepo.ListOrganizations(ctx, ListOrganizationsFilter{
		UserID:           &userID,
		MembershipStatus: &activeStatus,
	})
	if err != nil {
		return db.Organization{}, fmt.Errorf(
			"listing organization memberships for user: %w",
			err,
		)
	}
	if len(activeMemberships) > 0 {
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

		_, err = s.tenantRepo.CreateMembership(ctx, tx, CreateMembershipParams{
			ID:             uuid.NewV7(),
			OrganizationID: organization.ID,
			UserID:         userID,
			FirstName:      firstName,
			LastName:       lastName,
			Role:           authz.RoleOwner,
			Permissions:    []authz.Permission{},
			Status:         model.MembershipStatusActive,
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
) (activeMembership, error) {
	return s.tenantRepo.GetActiveMemberhip(
		ctx,
		userID,
		orgSlug,
		model.MembershipStatusActive,
	)
}

type MembershipsData struct {
	Memberhips  []model.Membership
	Invitations []model.Invitation
}

// not good
func (s *Service) GetMembershipsData(
	ctx context.Context,
	identity authz.Identity,
) (MembershipsData, error) {
	if !identity.HasPermission(authz.PermMembershipRead) {
		return MembershipsData{}, authz.ErrPermissionDenied
	}

	var data MembershipsData
	var err error

	data.Memberhips, err = s.tenantRepo.ListMemberships(
		ctx,
		identity.OrgID,
		ListMembershipsFilter{},
	)
	if err != nil {
		return MembershipsData{}, err
	}

	data.Invitations, err = s.tenantRepo.ListInvitations(ctx, identity.OrgID)

	return data, err
}

// not good
func (s *Service) GetMembership(
	ctx context.Context,
	orgID, membershipID uuid.UUID,
) (model.Membership, error) {
	fmt.Println(orgID, membershipID)
	return s.tenantRepo.GetMembership(ctx, orgID, membershipID, GetMembershipFilter{})
}

// not good: kaj za vraga ta service updetja?
func (s *Service) UpdateMembership(ctx context.Context, m model.Membership) error {
	return s.tenantRepo.UpdateMembership(ctx, nil, m.OrganizationID, m.ID, UpdateMembershipParams{
		FirstName: &m.FirstName,
		LastName:  &m.LastName,
	})
}

// CancelMembership cancels a membership.
func (s *Service) CancelMembership(
	ctx context.Context,
	identity authz.Identity,
	membershipID uuid.UUID,
) error {
	membership, err := s.tenantRepo.GetMembership(
		ctx,
		identity.OrgID,
		membershipID,
		GetMembershipFilter{},
	)
	if err != nil {
		return err
	}

	isOwnMembership := membership.UserID == identity.ID
	hasDeletePermission := identity.HasPermission(authz.PermMembershipDelete)

	if !isOwnMembership && !hasDeletePermission {
		return authz.ErrPermissionDenied
	}

	return s.tenantRepo.UpdateMembershipCanceled(
		ctx,
		nil,
		identity.OrgID,
		membershipID,
		identity.ID,
	)
}

type InviteUserParams struct {
	Email       string
	FirstName   string
	LastName    string
	Role        authz.Role
	Permissions []authz.Permission
}

// InviteUser handles sending invitation to user to join organization.
func (s *Service) InviteUser(
	ctx context.Context,
	identity authz.Identity,
	params InviteUserParams,
) (model.Invitation, error) {
	if !identity.HasPermission(authz.PermMembershipCreate) {
		return model.Invitation{}, authz.ErrPermissionDenied
	}

	isMember, err := s.tenantRepo.CheckUserIsMember(ctx, identity.OrgID, params.Email)
	if err != nil {
		return model.Invitation{}, fmt.Errorf("checking membership: %w", err)
	}
	if isMember {
		return model.Invitation{}, ErrUserAlreadyMember
	}

	latestInvitation, err := s.tenantRepo.GetLatestInvitation(ctx, identity.OrgID, params.Email)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return model.Invitation{}, fmt.Errorf("getting latest invitation: %w", err)
	}
	if err == nil && latestInvitation.IsPending() {
		return model.Invitation{}, ErrInvitationAlreadyPending
	}

	rawToken, tokenHash, err := generateInvitationToken()
	if err != nil {
		return model.Invitation{}, err
	}

	invitation, err := s.tenantRepo.CreateInvitation(ctx, nil, CreateInvitationParams{
		ID:          uuid.NewV7(),
		OrgID:       identity.OrgID,
		InviterID:   identity.ID,
		Email:       params.Email,
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Role:        params.Role,
		Permissions: params.Permissions,
		TokenHash:   tokenHash,
		ExpiresAt:   s.clock.NowUTC().Add(24 * time.Hour),
	})
	if err != nil {
		return model.Invitation{}, fmt.Errorf("creating invitation: %w", err)
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
	if !identity.HasPermission(authz.PermMembershipDelete) {
		return authz.ErrPermissionDenied
	}

	return s.tenantRepo.UpdateInvitationCanceled(
		ctx,
		nil,
		invitationID,
		identity.OrgID,
		identity.ID,
	)
}

// not good
func (s *Service) GetInvitation(
	ctx context.Context,
	token string,
) (model.Invitation, error) {
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

	if !invitation.IsPending() {
		return ErrInvitationNotPending
	}

	if userEmail != invitation.Email {
		return fmt.Errorf("user email not same as invitation email")
	}

	return pgx.BeginFunc(ctx, s.tenantRepo.Pool, func(tx pgx.Tx) error {
		err = s.tenantRepo.UpdateInvitationAcceptedAt(ctx, tx, invitation.ID)
		if err != nil {
			return err
		}

		_, err = s.tenantRepo.CreateMembership(ctx, tx, CreateMembershipParams{
			ID:             uuid.NewV7(),
			OrganizationID: invitation.OrganizationID,
			UserID:         userID,
			FirstName:      invitation.FirstName,
			LastName:       invitation.LastName,
			Role:           invitation.Role,
			Permissions:    invitation.Permissions,
			Status:         model.MembershipStatusActive,
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

	if !invitation.IsPending() {
		return ErrInvitationNotPending
	}

	if userEmail != invitation.Email {
		return fmt.Errorf("user email not same as invitation email")
	}

	return s.tenantRepo.UpdateInvitationDeclinedAt(ctx, nil, invitation.ID)
}
