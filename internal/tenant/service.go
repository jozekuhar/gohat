package tenant

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"uuid"

	"mimokocke/internal/shared/authz"

	"github.com/goforj/godump"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
	"github.com/resend/resend-go/v3"
)

var ErrOrganizationLimitReached = errors.New("user has reached maximum organizasion limit")

type Service struct {
	logger       *slog.Logger
	resendClient *resend.Client
	tenantRepo   *repository
}

func NewService(logger *slog.Logger, resendClient *resend.Client, tenantRepo *repository) *Service {
	return &Service{
		logger:       logger,
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
) (*Organization, error) {
	// user can be only part of single organization
	memberships, err := s.tenantRepo.ListActiveMemberships(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("listing organization memberships for user: %w", err)
	}
	if len(memberships) > 0 {
		return nil, ErrOrganizationLimitReached
	}

	tx, err := s.tenantRepo.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("transaction begin in organization registration: %w", err)
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
		return nil, fmt.Errorf("creating organization: %w", err)
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
		MembershipRoleOwner,
		membershipPermissions,
		MembershipStatusActive,
	)
	if err != nil {
		return nil, fmt.Errorf("creating organization membership: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit organization registration: %w", err)
	}

	_ = membership

	return organization, nil
}

func (s *Service) VerifyMembership(
	ctx context.Context,
	userID uuid.UUID,
	orgSlug string,
) (*authz.Identity, error) {
	identity, err := s.tenantRepo.GetIdentity(
		ctx,
		userID,
		orgSlug,
		MembershipStatusActive,
	)
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func (s *Service) GetMemberships(
	ctx context.Context,
	identity *authz.Identity,
) ([]Membership, error) {
	// TODO(jozekuhar): permissions to view members

	memberships, err := s.tenantRepo.ListMemberships(ctx, identity.OrganizationID)
	if err != nil {
		return nil, err
	}
	return memberships, nil
}

// InviteUser handles sending invitation to user to join organization.
func (s *Service) InviteUser(
	ctx context.Context,
	identity *authz.Identity,
	email string,
) (*Invitation, error) {
	if !identity.HasPermission(authz.MembershipCreate) {
		return nil, fmt.Errorf("you don't have permission to do this thing")
	}

	tx, err := s.tenantRepo.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("transaction begin: %w", err)
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			s.logger.Error("rollback error", "err", err)
		}
	}()

	token, tokenHash, err := generateInviteToken()
	if err != nil {
		return nil, err
	}

	invitationID := uuid.NewV7()
	invitation, err := s.tenantRepo.CreateInvitation(ctx, tx, invitationID, tokenHash)
	if err != nil {
		return nil, err
	}

	godump.Dump(invitation)

	params := &resend.SendEmailRequest{
		From:    "invites@gajogroup.com",
		To:      []string{email},
		Subject: "Invitation",
		ReplyTo: "support@gajogroup.com",
		Text:    "Hi, you have been invited to our organization. Token: " + token,
		Tags: []resend.Tag{
			{Name: "type", Value: "invite"},
		},
	}
	resp, err := s.resendClient.Emails.Send(params)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	godump.Dump(resp)

	// membershipID := uuid.NewV7()
	// role := MembershipRoleMember
	// permissions := []string{}
	// membership, err := s.tenantRepo.CreateMembership(
	// 	ctx,
	// 	nil,
	// 	membershipID,
	// 	identity.OrganizationID,
	// 	uuid.UUID{},
	// 	role,
	// 	permissions,
	// 	MembershipStatusInvited,
	// )
	// if err != nil {
	// 	return nil, err
	// }
	// return membership, nil
	return nil, nil
}

// ConsumeInvitation makes user accept invitation to join organization.
func (s *Service) ConsumeInvitation() {}
