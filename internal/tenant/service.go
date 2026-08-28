package tenant

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"uuid"

	"mimokocke/internal/shared/authz"
	"mimokocke/internal/shared/clock"
	"mimokocke/internal/shared/config"
	"mimokocke/internal/shared/routes"

	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/resend/resend-go/v3"
)

var ErrOrganizationLimitReached = errors.New("user has reached maximum organizasion limit")

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

func (s *Service) ListOrganizations(
	ctx context.Context,
	userID uuid.UUID,
) ([]Organization, error) {
	organizations, err := s.tenantRepo.ListOrganizationsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return organizations, nil
}

func (s *Service) RegisterOrganization(
	ctx context.Context,
	userID uuid.UUID,
	orgName string,
	orgSlug string,
	firstName string,
	lastName string,
) (Organization, error) {
	// user can be only part of single organization
	memberships, err := s.tenantRepo.ListActiveMemberships(ctx, userID)
	if err != nil {
		return Organization{}, fmt.Errorf("listing organization memberships for user: %w", err)
	}
	if len(memberships) > 0 {
		return Organization{}, ErrOrganizationLimitReached
	}

	tx, err := s.tenantRepo.pool.Begin(ctx)
	if err != nil {
		return Organization{}, fmt.Errorf("transaction begin in organization registration: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			s.logger.Error("rollback error", "err", err)
		}
	}()

	orgID := uuid.NewV7()
	orgSlug = slug.Make(orgSlug)
	organization, err := s.tenantRepo.CreateOrganization(
		ctx,
		tx,
		orgID,
		orgName,
		orgSlug,
	)
	if err != nil {
		return Organization{}, fmt.Errorf("creating organization: %w", err)
	}

	membershipID := uuid.NewV7()
	membershipPermissions := []string{}

	membership, err := s.tenantRepo.CreateMembership(
		ctx,
		tx,
		membershipID,
		organization.ID,
		userID,
		firstName,
		lastName,
		string(authz.RoleOwner),
		membershipPermissions,
		membershipStatusActive,
	)
	if err != nil {
		return Organization{}, fmt.Errorf("creating organization membership: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return Organization{}, fmt.Errorf("commit organization registration: %w", err)
	}

	_ = membership

	return organization, nil
}

func (s *Service) GetActiveMembership(
	ctx context.Context,
	userID uuid.UUID,
	orgSlug string,
) (activeMembership, error) {
	am, err := s.tenantRepo.GetActiveMemberhip(
		ctx,
		userID,
		orgSlug,
		membershipStatusActive,
	)
	if err != nil {
		return activeMembership{}, err
	}

	return am, nil
}

type MembershipsData struct {
	Memberhips  []Membership
	Invitations []Invitation
}

func (s *Service) GetMembershipsData(
	ctx context.Context,
	identity authz.Identity,
) (MembershipsData, error) {
	// TODO(jozekuhar): permissions to view members

	memberships, err := s.tenantRepo.ListMemberships(ctx, identity.OrganizationID)
	if err != nil {
		return MembershipsData{}, err
	}

	invitations, err := s.tenantRepo.ListInvitations(ctx, identity.OrganizationID)
	if err != nil {
		return MembershipsData{}, err
	}

	return MembershipsData{
		Memberhips:  memberships,
		Invitations: invitations,
	}, nil
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
) (Invitation, error) {
	if !identity.HasPermission(authz.PermMembershipCreate) {
		return Invitation{}, fmt.Errorf("you don't have permission to do this thing")
	}

	tx, err := s.tenantRepo.pool.Begin(ctx)
	if err != nil {
		return Invitation{}, fmt.Errorf("transaction begin: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			s.logger.Error("rollback error", "err", err)
		}
	}()

	token, tokenHash, err := generateInviteToken()
	if err != nil {
		return Invitation{}, err
	}

	invitation, err := s.tenantRepo.CreateInvitation(ctx, tx, Invitation{
		ID:             uuid.NewV7(),
		OrganizationID: identity.OrganizationID,
		Email:          params.Email,
		FirstName:      params.FirstName,
		LastName:       params.LastName,
		Role:           params.Role,
		Permissions:    params.Permissions,
		TokenHash:      tokenHash,
		ExpiresAt:      s.clock.NowUTC().Add(24 * time.Hour),
	})
	if err != nil {
		return Invitation{}, fmt.Errorf("creating invitation: %w", err)
	}

	invitationURL := fmt.Sprintf(
		"http://%s%s?token=%s",
		s.cfg.Host,
		fmt.Sprintf(routes.InvitationsJoin, identity.OrganizationSlug),
		token,
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
	if err != nil {
		return Invitation{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return Invitation{}, err
	}

	return invitation, nil
}

func (s *Service) CancelInvite(
	ctx context.Context,
	identity authz.Identity,
	invitationID uuid.UUID,
) error {
	// TODO(jozekuhar): check permissions

	err := s.tenantRepo.DeleteInvitation(ctx, invitationID, identity.OrganizationID)
	return err
}

// ConsumeInvitation makes user accept invitation to join organization.
func (s *Service) ConsumeInvitation() {
	panic("unimplemented")
}
