package tenant

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/permissions"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/web/components"
	"mimokocke/internal/web/request"
	"mimokocke/internal/web/response"

	hx "maragu.dev/gomponents-htmx/http"

	"github.com/go-playground/form"
)

type Handler struct {
	logger  *slog.Logger
	decoder *form.Decoder
	srv     *Service
}

func NewHandler(
	logger *slog.Logger,
	decoder *form.Decoder,
	srv *Service,
) *Handler {
	return &Handler{
		logger:  logger,
		decoder: decoder,
		srv:     srv,
	}
}

func (h *Handler) GetAppRoot(w http.ResponseWriter, r *http.Request) {
	authCtx := identity.MustAuthFromContext(r.Context())

	orgs, err := h.srv.ListActiveOrganizations(r.Context(), authCtx.UserID)
	if err != nil {
		h.logger.Error("list organizations", "err", err)
		return
	}

	response.RenderStatus(
		w,
		appPage(orgs),
		http.StatusOK,
	)
}

func (h *Handler) GetSidebarOrganizationsPartial(w http.ResponseWriter, r *http.Request) {
	authCtx := identity.MustAuthFromContext(r.Context())

	organizations, err := h.srv.ListActiveOrganizations(r.Context(), authCtx.UserID)
	if err != nil {
		h.logger.Error("list organizations", "err", err)
		return
	}

	response.RenderStatus(
		w,
		components.SidebarHeaderPopoverOrganizationsPartial(organizations),
		http.StatusOK,
	)
}

func (h *Handler) GetCreateOrganizationFormModal(w http.ResponseWriter, r *http.Request) {
	response.RenderStatus(w, components.OrganizationCreateFormModal(), http.StatusOK)
}

func (h *Handler) PostCreateOrganization(w http.ResponseWriter, r *http.Request) {
	authCtx := identity.MustAuthFromContext(r.Context())

	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parsing create organization form", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	var form struct {
		Name      string
		Slug      string
		FirstName string
		LastName  string
	}

	err = h.decoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decoding create organization form", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	organization, err := h.srv.RegisterOrganization(
		r.Context(),
		RegisterOrganizationParams{
			UserID:    authCtx.UserID,
			OrgName:   form.Name,
			OrgSlug:   form.Slug,
			FirstName: form.FirstName,
			LastName:  form.LastName,
		},
	)
	if errors.Is(err, ErrOrganizationLimitReached) {
		h.logger.Warn("user tries to create new organization when limit reached", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("You have reached maximum organization limit."),
			http.StatusBadRequest,
		)
		return
	}
	if err != nil {
		h.logger.Error("registering organization", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	hx.SetRedirect(w.Header(), fmt.Sprintf(routes.OrgDashboard, organization.Slug))
}

func (h *Handler) GetOrgRoot(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustIdentityFromContext(r.Context())

	http.Redirect(
		w,
		r,
		fmt.Sprintf(routes.OrgDashboard, ident.OrgSlug),
		http.StatusFound,
	)
}

func (h *Handler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustIdentityFromContext(r.Context())

	response.RenderStatus(w, dashboardPage(ident), http.StatusOK)
}

func (h *Handler) GetMemberships(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustIdentityFromContext(r.Context())

	data, err := h.srv.GetMembershipsOverview(r.Context(), ident)
	if err != nil {
		h.logger.Error("getting memberships data", "err", err)
		return
	}

	response.RenderStatus(w, membershipsPage(ident, data), http.StatusOK)
}

func (h *Handler) GetUpdateMembershipFormModal(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustIdentityFromContext(r.Context())

	memberID, err := request.PathValueUUID(r, routes.PathMembershipID)
	if err != nil {
		h.logger.Error("getting path value uuid", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusBadRequest,
		)
		return
	}

	member, err := h.srv.GetMembershipDetails(r.Context(), ident, memberID)
	if err != nil {
		h.logger.Error("getting membership", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusBadRequest,
		)
		return
	}

	response.RenderStatus(w, membershipUpdateFormModal(ident, member), http.StatusOK)
}

func (h *Handler) PatchUpdateMembership(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustIdentityFromContext(r.Context())

	memberID, err := request.PathValueUUID(r, routes.PathMembershipID)
	if err != nil {
		h.logger.Error("getting path value uuid", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusBadRequest,
		)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.logger.Error("parsing form", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	var form struct {
		FirstName string
		LastName  string
	}

	err = h.decoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decoding form", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	updatedMember, err := h.srv.UpdateMembership(r.Context(), db.Membership{
		ID:             memberID,
		OrganizationID: ident.OrgID,
		FirstName:      form.FirstName,
		LastName:       form.LastName,
	})
	if err != nil {
		h.logger.Error("update membership", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	member, err := h.srv.GetMembershipDetails(r.Context(), ident, updatedMember.ID)
	if err != nil {
		h.logger.Error("get membership", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	hx.SetTrigger(w.Header(), components.EventModalClose)
	response.RenderStatus(
		w,
		membershipItem(ident.OrgSlug, member.Membership, member.User),
		http.StatusOK,
	)
	response.RenderStatus(
		w,
		components.ToastFragment("Membership succesfully updated"),
		http.StatusOK,
	)
}

func (h *Handler) PostCancelMembership(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustIdentityFromContext(r.Context())

	membershipID, err := request.PathValueUUID(r, routes.PathMembershipID)
	if err != nil {
		h.logger.Error("getting path value uuid for membership cancel", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusBadRequest,
		)
		return
	}

	err = h.srv.CancelMembership(r.Context(), ident, membershipID)
	if err != nil {
		h.logger.Error("canceling membership", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCreateInvitationFormModal(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustIdentityFromContext(r.Context())
	response.RenderStatus(w, invitationCreateFormModal(ident), http.StatusOK)
}

func (h *Handler) PostCreateInvitation(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustIdentityFromContext(r.Context())

	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parsing invitation form", "err", err)
		return
	}

	var form struct {
		Email       string
		FirstName   string
		LastName    string
		Role        permissions.MembershipRole
		Permissions []permissions.MembershipPermission
	}

	err = h.decoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decoding invitation form", "err", err)
		return
	}

	invite, err := h.srv.InviteUser(
		r.Context(),
		ident,
		InviteUserParams{
			Email:       form.Email,
			FirstName:   form.FirstName,
			LastName:    form.LastName,
			Role:        form.Role,
			Permissions: form.Permissions,
		},
	)
	if errors.Is(err, ErrUserAlreadyMember) {
		response.RenderStatus(
			w,
			components.ToastFragment("User ("+form.Email+") is already member of organization."),
			http.StatusBadRequest,
		)
		return
	}
	if errors.Is(err, ErrInvitationAlreadyPending) {
		response.RenderStatus(
			w,
			components.ToastFragment("User ("+form.Email+") has already pending invitation."),
			http.StatusBadRequest,
		)
		return
	}
	if err != nil {
		h.logger.Error("invite user to organization", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	response.HXTrigger(w, components.EventModalClose)
	response.Render(w, components.ToastFragment("Invited"))
	response.RenderStatus(w, invitationItem(invite), http.StatusCreated)
}

func (h *Handler) PostCancelInvitation(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustIdentityFromContext(r.Context())

	inviteID, err := request.PathValueUUID(r, routes.PathInvitationID)
	if err != nil {
		h.logger.Error("parsing invitation id", "err", err)
		return
	}

	err = h.srv.CancelInvite(r.Context(), ident, inviteID)
	if err != nil {
		h.logger.Error("canceling invitation", "err", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetShowInvitation(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue(routes.PathInvitationToken)

	invitation, err := h.srv.GetInvitation(r.Context(), token)
	if errors.Is(err, ErrInvitationNotPending) {
		response.RenderStatus(w, invitationErrorPage(), http.StatusForbidden)
		return
	}
	if err != nil {
		h.logger.Error("getting invitation", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = identity.AuthFromContext(r.Context())
	if err != nil {
		response.RenderStatus(w, invitationRegisterPage(token, invitation), http.StatusOK)
		return
	}

	response.RenderStatus(w, invitationAcceptPage(token, invitation), http.StatusOK)
}

func (h *Handler) PostAcceptInvitation(w http.ResponseWriter, r *http.Request) {
	authCtx := identity.MustAuthFromContext(r.Context())

	token := r.PathValue(routes.PathInvitationToken)

	err := h.srv.AcceptInvitation(r.Context(), authCtx.UserID, authCtx.UserEmail, token)
	if err != nil {
		h.logger.Error("accepting invitation", "err", err)
		return
	}

	w.WriteHeader(http.StatusSeeOther)
	hx.SetRedirect(w.Header(), routes.AppRoot)
}

func (h *Handler) PostDeclineInvitation(w http.ResponseWriter, r *http.Request) {
	authCtx := identity.MustAuthFromContext(r.Context())

	token := r.PathValue(routes.PathInvitationToken)

	err := h.srv.DeclineInvitation(r.Context(), authCtx.UserEmail, token)
	if err != nil {
		h.logger.Error("declining invitation", "err", err)
		return
	}

	w.WriteHeader(http.StatusSeeOther)
	hx.SetLocation(w.Header(), routes.AppRoot)
}
