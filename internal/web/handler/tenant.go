package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/tenant"
	"mimokocke/internal/web/view"

	"github.com/go-playground/form"
	"github.com/goforj/godump"
)

type tenantHandler struct {
	logger          *slog.Logger
	tenantSrv       *tenant.Service
	formDecoder     *form.Decoder
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

	render(w, h.dashboardView.OrganizationsPage(organizations))
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
		h.logger.Error("decodeing create organization form", "err", err)
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
		render(
			w,
			h.toastView.FragmentText(
				"Maximum limit reached",
				"You have reached maximum organization limit",
			),
		)
		return
	}
	if err != nil {
		h.logger.Error("registering organization", "err", err)
		return
	}

	_ = organization
}

func (h *tenantHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	userID := auth.MustUserIDFomContext(r.Context())
	membership := tenant.MustIdentityFromContext(r.Context())

	render(w, h.dashboardView.DashboardPage(userID, membership))
}

func (h *tenantHandler) GetMemberships(w http.ResponseWriter, r *http.Request) {
	userID := auth.MustUserIDFomContext(r.Context())
	_ = userID
	// still need to log it?

	identity := tenant.MustIdentityFromContext(r.Context())

	memberships, err := h.tenantSrv.GetMemberships(r.Context(), identity)
	if err != nil {
		return
	}

	render(w, h.membershipsView.MembershipsPage(identity, memberships))
}

type createInvitationForm struct {
	Email       string
	FirstName   string
	LastName    string
	Role        string
	Permissions []string
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

	godump.Dump(identity)
	godump.Dump(form)

	// invitation, err := h.tenantSrv.InviteUser(r.Context(), identity, email)
	// if err != nil {
	// 	h.logger.Error("invite user to organization", "err", err)
	// 	return
	// }
	//
	// godump.Dump(invitation)
	//
	// _, _ = w.Write([]byte("Hello world"))
}

func (h *tenantHandler) GetShowInvitation(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func (h *tenantHandler) PostAcceptInvitation(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func (h *tenantHandler) DeleteDeclineInvitation(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
