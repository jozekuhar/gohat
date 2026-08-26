package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/tenant"
	"mimokocke/internal/web/view"
)

type tenantHandler struct {
	logger          *slog.Logger
	tenantSrv       *tenant.Service
	dashboardView   *view.Dashboard
	membershipsView *view.Memberships
	toastView       *view.Toast
}

func NewTenantHandler(logger *slog.Logger, tenantSrv *tenant.Service) *tenantHandler {
	return &tenantHandler{
		logger:          logger,
		tenantSrv:       tenantSrv,
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

func (h *tenantHandler) PostCreateOrganization(w http.ResponseWriter, r *http.Request) {
	userID := auth.MustUserIDFomContext(r.Context())

	name := r.PostFormValue(view.InputNameOrganizationName)
	if name == "" {
		render(w, h.toastView.FragmentText("Invalid input", "Organization name is required"))
		return
	}

	slug := r.PostFormValue(view.InputNameOrganizationSlug)
	if slug == "" {
		render(w, h.toastView.FragmentText("Invalid input", "Organization slug is required"))
		return
	}

	organization, err := h.tenantSrv.RegisterOrganization(r.Context(), userID, name, slug)
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

func (h *tenantHandler) PostCreateInvitation(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
