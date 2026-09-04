package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/shared/authz"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/tenant"
	"mimokocke/internal/web/view"

	"github.com/go-playground/form"
	"github.com/goforj/godump"
)

type tenantHandler struct {
	logger          *slog.Logger
	tenantSrv       *tenant.Service
	formDecoder     *form.Decoder
	layoutView      *view.Layout
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
		dashboardView:   view.NewDashboard(),
		membershipsView: view.NewMemberships(),
		toastView:       view.NewToast(),
	}
}

func (h *tenantHandler) GetOrganizations(w http.ResponseWriter, r *http.Request) {
	userID := auth.MustUserIDFomContext(r.Context())

	organizations, err := h.tenantSrv.ListOrganizations(r.Context(), userID)
	if err != nil {
		h.logger.Error("list organizations", "err", err)
		return
	}

	godump.Dump(organizations)
	fmt.Println("hello world")

	render(w, h.dashboardView.OrganizationsPage(authz.Identity{}, nil))
}

func (h *tenantHandler) GetSidebarOrganizationsPartial(w http.ResponseWriter, r *http.Request) {
	userID := auth.MustUserIDFomContext(r.Context())

	organizations, err := h.tenantSrv.ListOrganizations(r.Context(), userID)
	if err != nil {
		h.logger.Error("list organizations", "err", err)
		return
	}

	render(w, h.layoutView.SidebarHeaderPopoverOrganizationsPartial(organizations))
}

func (h *tenantHandler) GetOrganizationsCreateForm(w http.ResponseWriter, r *http.Request) {
	render(w, h.layoutView.OrganizationCreateFormModal())
}

type createOrganizationForm struct {
	Name      string
	Slug      string
	FirstName string
	LastName  string
}

func (h *tenantHandler) PostCreateOrganization(w http.ResponseWriter, r *http.Request) {
	userID := auth.MustUserIDFomContext(r.Context())

	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parsing create organization form", "err", err)
		return
	}

	var form createOrganizationForm
	err = h.formDecoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decoding create organization form", "err", err)
		return
	}

	organization, err := h.tenantSrv.RegisterOrganization(
		r.Context(),
		userID,
		form.Name,
		form.Slug,
		form.FirstName,
		form.LastName,
	)
	if errors.Is(err, tenant.ErrOrganizationLimitReached) {
		h.logger.Warn("user tries to create new organization when limit reached", "err", err)
		render(w, h.toastView.Fragment("You have reached maximum organization limit"))
		return
	}
	if err != nil {
		h.logger.Error("registering organization", "err", err)
		return
	}

	_ = organization
}

func (h *tenantHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	render(w, h.dashboardView.DashboardPage(identity))
}

func (h *tenantHandler) GetMemberships(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	data, err := h.tenantSrv.GetMembershipsData(r.Context(), identity)
	if err != nil {
		h.logger.Error("getting memberships data", "err", err)
		return
	}

	render(w, h.membershipsView.MembershipsPage(identity, data))
}

type createInvitationForm struct {
	Email       string
	FirstName   string
	LastName    string
	Role        authz.Role
	Permissions []authz.Permission
}

func (h *tenantHandler) PostCreateInvitation(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parsing invitation form", "err", err)
		return
	}

	var form createInvitationForm
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
	if err != nil {
		h.logger.Error("invite user to organization", "err", err)
		return
	}

	godump.Dump(invitation)
}

func (h *tenantHandler) DeleteRemoveInvitation(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	invitationID, err := pathValueUUID(r, routes.PathInvitationID)
	if err != nil {
		h.logger.Error("", "err", err)
		return
	}

	fmt.Println(invitationID)

	err = h.tenantSrv.CancelInvite(r.Context(), identity, invitationID)
	if err != nil {
		h.logger.Error("canceling invitation", "err", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *tenantHandler) GetShowInvitation(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFomContext(r.Context())
	if err != nil {
		// Redirect user to register page if not auth
		// user id not in context
		slog.Error("user id from context", "err", err)
		http.Redirect(w, r, routes.Register, http.StatusSeeOther)
		return
	}

	// Check if this is same user to accept invitation

	godump.Dump(userID)
}

func (h *tenantHandler) PostAcceptInvitation(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func (h *tenantHandler) DeleteDeclineInvitation(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
