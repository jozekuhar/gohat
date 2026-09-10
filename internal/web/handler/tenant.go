package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/model"
	"mimokocke/internal/shared/authz"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/tenant"
	"mimokocke/internal/web/view"

	"github.com/go-playground/form"
	"github.com/goforj/godump"
	hx "maragu.dev/gomponents-htmx/http"
)

type tenantHandler struct {
	logger          *slog.Logger
	tenantSrv       *tenant.Service
	formDecoder     *form.Decoder
	layoutView      *view.Layout
	coreView        *view.Core
	dashboardView   *view.Dashboard
	membershipsView *view.Memberships
	toastView       *view.Toast
}

func NewTenantHandler(
	logger *slog.Logger,
	tenantSrv *tenant.Service,
	formDecoder *form.Decoder,
) *tenantHandler {
	return &tenantHandler{
		logger:          logger,
		tenantSrv:       tenantSrv,
		formDecoder:     formDecoder,
		layoutView:      view.NewLayout(),
		coreView:        view.NewCore(),
		dashboardView:   view.NewDashboard(),
		membershipsView: view.NewMemberships(),
		toastView:       view.NewToast(),
	}
}

func (h *tenantHandler) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	authCtx := auth.MustAuthFromContext(r.Context())

	organizations, err := h.tenantSrv.ListActiveOrganizations(r.Context(), authCtx.UserID)
	if err != nil {
		h.logger.Error("list organizations", "err", err)
		return
	}

	// TODO(jozekuhar): i don't need organizations here
	godump.Dump(organizations)

	render(w, h.dashboardView.OrganizationsPage(authz.Identity{}, nil), http.StatusOK)
}

func (h *tenantHandler) GetSidebarOrganizationsPartial(w http.ResponseWriter, r *http.Request) {
	authCtx := auth.MustAuthFromContext(r.Context())

	organizations, err := h.tenantSrv.ListActiveOrganizations(r.Context(), authCtx.UserID)
	if err != nil {
		h.logger.Error("list organizations", "err", err)
		return
	}

	render(w, h.layoutView.SidebarHeaderPopoverOrganizationsPartial(organizations), http.StatusOK)
}

func (h *tenantHandler) GetCreateOrganizationFormModal(w http.ResponseWriter, r *http.Request) {
	render(w, h.layoutView.OrganizationCreateFormModal(), http.StatusOK)
}

type createOrganizationForm struct {
	Name      string
	Slug      string
	FirstName string
	LastName  string
}

func (h *tenantHandler) PostCreateOrganization(w http.ResponseWriter, r *http.Request) {
	authCtx := auth.MustAuthFromContext(r.Context())

	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parsing create organization form", "err", err)
		render(w, h.toastView.Fragment("Something went wrong"), http.StatusInternalServerError)
		return
	}

	var form createOrganizationForm
	err = h.formDecoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decoding create organization form", "err", err)
		render(w, h.toastView.Fragment("Something went wrong"), http.StatusInternalServerError)
		return
	}

	organization, err := h.tenantSrv.RegisterOrganization(
		r.Context(),
		authCtx.UserID,
		form.Name,
		form.Slug,
		form.FirstName,
		form.LastName,
	)
	if errors.Is(err, tenant.ErrOrganizationLimitReached) {
		h.logger.Warn("user tries to create new organization when limit reached", "err", err)
		render(
			w,
			h.toastView.Fragment("You have reached maximum organization limit."),
			http.StatusBadRequest,
		)
		return
	}
	if err != nil {
		h.logger.Error("registering organization", "err", err)
		render(w, h.toastView.Fragment("Something went wrong"), http.StatusInternalServerError)
		return
	}

	hx.SetRedirect(w.Header(), fmt.Sprintf(routes.OrgDashboard, organization.Slug))
}

func (h *tenantHandler) GetRoot(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	http.Redirect(
		w,
		r,
		fmt.Sprintf(routes.OrgDashboard, identity.OrgSlug),
		http.StatusFound,
	)
}

func (h *tenantHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	render(w, h.dashboardView.DashboardPage(identity), http.StatusOK)
}

func (h *tenantHandler) GetMemberships(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	data, err := h.tenantSrv.GetMembershipsData(r.Context(), identity)
	if err != nil {
		h.logger.Error("getting memberships data", "err", err)
		return
	}

	godump.Dump(data)

	render(w, h.membershipsView.MembershipsPage(identity, data), http.StatusOK)
}

func (h *tenantHandler) GetUpdateMembershipFormModal(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	membershipID, err := pathValueUUID(r, routes.PathMembershipID)
	if err != nil {
		h.logger.Error("getting path value uuid", "err", err)
		render(w, h.toastView.Fragment("Something went wrong"), http.StatusBadRequest)
		return
	}

	// fetch membership and pass data to update form modal
	membership, err := h.tenantSrv.GetMembership(r.Context(), identity.OrgID, membershipID)
	if err != nil {
		h.logger.Error("getting membership", "err", err)
		render(w, h.toastView.Fragment("Something went wrong"), http.StatusBadRequest)
		return
	}

	render(w, h.membershipsView.MembershipUpdateFormModal(identity, membership), http.StatusOK)
}

func (h *tenantHandler) PatchUpdateMembership(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	membershipID, err := pathValueUUID(r, routes.PathMembershipID)
	if err != nil {
		h.logger.Error("getting path value uuid", "err", err)
		render(w, h.toastView.Fragment("Something went wrong"), http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		h.logger.Error("parsing form", "err", err)
		render(w, h.toastView.Fragment("Something went wrong"), http.StatusInternalServerError)
		return
	}

	var form struct {
		FirstName string
		LastName  string
	}

	err = h.formDecoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decoding form", "err", err)
		render(w, h.toastView.Fragment("Something went wrong"), http.StatusInternalServerError)
		return
	}

	err = h.tenantSrv.UpdateMembership(r.Context(), model.Membership{
		ID:             membershipID,
		OrganizationID: identity.OrgID,
		FirstName:      form.FirstName,
		LastName:       form.LastName,
	})
	if err != nil {
		h.logger.Error("updating membership", "err", err)
		render(w, h.toastView.Fragment("Something went wrong"), http.StatusInternalServerError)
		return
	}

	render(w, h.toastView.Fragment("Membership succesfully updated"), http.StatusOK)
}

func (h *tenantHandler) PostCancelMembership(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	membershipID, err := pathValueUUID(r, routes.PathMembershipID)
	if err != nil {
		h.logger.Error("getting path value uuid for membership cancel", "err", err)
		render(w, h.toastView.Fragment("Something went wrong"), http.StatusBadRequest)
		return
	}

	err = h.tenantSrv.CancelMembership(r.Context(), identity, membershipID)
	if err != nil {
		h.logger.Error("canceling membership", "err", err)
		render(w, h.toastView.Fragment("Something went wrong"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *tenantHandler) GetCreateInvitationFormModal(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())
	render(w, h.membershipsView.InvitationCreateFormModal(identity), http.StatusOK)
}

func (h *tenantHandler) PostCreateInvitation(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parsing invitation form", "err", err)
		return
	}

	var form struct {
		Email       string
		FirstName   string
		LastName    string
		Role        authz.Role
		Permissions []authz.Permission
	}

	err = h.formDecoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decoding invitation form", "err", err)
		return
	}

	invitation, err := h.tenantSrv.InviteUser(
		r.Context(),
		identity,
		tenant.InviteUserParams{
			Email:       form.Email,
			FirstName:   form.FirstName,
			LastName:    form.LastName,
			Role:        form.Role,
			Permissions: form.Permissions,
		},
	)
	if errors.Is(err, tenant.ErrUserAlreadyMember) {
		render(
			w,
			h.toastView.Fragment("User ("+form.Email+") is already member of organization."),
			http.StatusBadRequest,
		)
		return
	}
	if errors.Is(err, tenant.ErrInvitationAlreadyPending) {
		render(
			w,
			h.toastView.Fragment("User ("+form.Email+") has already pending invitation."),
			http.StatusBadRequest,
		)
		return
	}
	if err != nil {
		h.logger.Error("invite user to organization", "err", err)
		render(w, h.toastView.Fragment("Something went wrong"), http.StatusInternalServerError)
		return
	}

	godump.Dump(invitation)
	render(w, h.toastView.Fragment("Invited"), http.StatusAccepted)
}

func (h *tenantHandler) PostCancelInvitation(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	invitationID, err := pathValueUUID(r, routes.PathInvitationID)
	if err != nil {
		h.logger.Error("parsing invitation id", "err", err)
		return
	}

	err = h.tenantSrv.CancelInvite(r.Context(), identity, invitationID)
	if err != nil {
		h.logger.Error("canceling invitation", "err", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *tenantHandler) GetShowInvitation(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue(routes.PathInvitationToken)

	invitation, err := h.tenantSrv.GetInvitation(r.Context(), token)
	if errors.Is(err, tenant.ErrInvitationNotPending) {
		render(w, h.membershipsView.InvitationErrorPage(), http.StatusForbidden)
		return
	}
	if err != nil {
		h.logger.Error("getting invitation", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = auth.AuthFromContext(r.Context())
	if err != nil {
		render(w, h.membershipsView.InvitationRegisterPage(token, invitation), http.StatusOK)
		return
	}

	render(w, h.membershipsView.InvitationAcceptPage(token, invitation), http.StatusOK)
}

func (h *tenantHandler) PostAcceptInvitation(w http.ResponseWriter, r *http.Request) {
	authCtx := auth.MustAuthFromContext(r.Context())

	token := r.PathValue(routes.PathInvitationToken)

	err := h.tenantSrv.AcceptInvitation(r.Context(), authCtx.UserID, authCtx.UserEmail, token)
	if err != nil {
		h.logger.Error("accepting invitation", "err", err)
		return
	}

	w.WriteHeader(http.StatusSeeOther)
	hx.SetRedirect(w.Header(), routes.Root)
}

func (h *tenantHandler) PostDeclineInvitation(w http.ResponseWriter, r *http.Request) {
	authCtx := auth.MustAuthFromContext(r.Context())

	token := r.PathValue(routes.PathInvitationToken)

	err := h.tenantSrv.DeclineInvitation(r.Context(), authCtx.UserEmail, token)
	if err != nil {
		h.logger.Error("declining invitation", "err", err)
		return
	}

	w.WriteHeader(http.StatusSeeOther)
	hx.SetLocation(w.Header(), routes.Root)
}
