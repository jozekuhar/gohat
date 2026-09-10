package main

import (
	"context"
	"log"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/channel"
	"mimokocke/internal/middleware"
	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/clock"
	"mimokocke/internal/shared/config"
	"mimokocke/internal/shared/logger"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/tenant"
	"mimokocke/internal/web/handler"
	"mimokocke/internal/web/view"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	_ "github.com/goforj/godump"
	"github.com/resend/resend-go/v3"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panicf("loading config: %s", err)
	}

	logger := logger.Init(cfg.Debug)

	systemClock := clock.NewSystemClock()

	pool, err := db.LoadPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Panicf("loading postgres: %s", err)
	}

	resendClient := resend.NewClient(cfg.ResendAPIKey)

	err = view.LoadManifest()
	if err != nil {
		log.Panicf("loading manifest: %s", err)
	}

	baseRepo := db.NewBaseRepo(pool)

	tenantRepo := tenant.NewRepository(baseRepo)
	authRepo := auth.NewRepository(baseRepo)
	channelRepo := channel.NewRepository(baseRepo)

	tenantSrv := tenant.NewService(cfg, logger, systemClock, resendClient, tenantRepo)
	authSrv := auth.NewService(cfg, logger, systemClock, authRepo)
	channelSrv := channel.NewService(channelRepo)

	formDecoder := form.NewDecoder()
	validator := validator.New(validator.WithRequiredStructEnabled())

	staticHdl := handler.NewStatic()
	coreHdl := handler.NewCore()
	tenantHdl := handler.NewTenantHandler(logger, tenantSrv, formDecoder)
	authHdl := handler.NewAuthHandler(cfg, logger, authSrv, formDecoder, validator)
	channelHdl := handler.NewChanneHandler(logger, channelSrv)

	tenantMdw := middleware.NewTenantMiddleware(logger, tenantSrv, coreHdl)
	authMdw := middleware.NewAuthMiddleware(logger, authSrv)

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get(routes.Static, staticHdl.GetStatic)
		r.NotFound(coreHdl.GetNotFound)
		r.MethodNotAllowed(coreHdl.GetNotFound)
	})

	r.Group(func(r chi.Router) {
		r.Use(chimiddleware.Logger)

		r.Group(func(r chi.Router) {
			// Core
			r.Get(routes.TermsOfService, coreHdl.GetTermsOfService)
			r.Get(routes.PrivacyPolicy, coreHdl.GetPrivacyPolicy)

			// Auth
			r.Get(routes.CallbackSignInGoogle, authHdl.GetSignInGoogleCallback)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMdw.OptionalAuth)

			// Invitations
			r.Get(routes.InvitationsJoinPath, tenantHdl.GetShowInvitation)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMdw.RequireGuest)

			r.Get(routes.Login, authHdl.GetLogin)
			r.Post(routes.HXLogin, authHdl.PostLogin)
			r.Get(routes.Register, authHdl.GetRegister)
			r.Post(routes.HXRegister, authHdl.PostRegister)
			r.Get(routes.HXSignInGoogle, authHdl.GetSignInWithGoogle)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMdw.RequireAuth)

			r.Post(routes.HXLogout, authHdl.PostLogout)
			r.Get(routes.Root, tenantHdl.GetOrganizations)

			// Sidebar
			r.Get(routes.HXSidebarOrganizations, tenantHdl.GetSidebarOrganizationsPartial)
			r.Get(routes.HXSidebarOrganizationsCreate, tenantHdl.GetCreateOrganizationFormModal)
			r.Post(routes.HXSidebarOrganizationsCreate, tenantHdl.PostCreateOrganization)

			// Invitations
			r.Post(routes.HXInvitationsAcceptPath, tenantHdl.PostAcceptInvitation)
			r.Post(routes.HXInvitationsDeclinePath, tenantHdl.PostDeclineInvitation)

			// Organization
			r.Group(func(r chi.Router) {
				r.Use(tenantMdw.RequireIdentity)
				r.Get(routes.OrgRootPath, tenantHdl.GetRoot)
				r.Get(routes.OrgDashboardPath, tenantHdl.GetDashboard)

				// Memberships
				r.Get(routes.OrgMembershipsPath, tenantHdl.GetMemberships)
				r.Get(routes.HXOrgMembershipsUpdatePath, tenantHdl.GetUpdateMembershipFormModal)
				r.Patch(routes.HXOrgMembershipsUpdatePath, tenantHdl.PatchUpdateMembership)
				r.Get(routes.HXOrgInvitationsCreatePath, tenantHdl.GetCreateInvitationFormModal)
				r.Post(routes.HXOrgInvitationsCreatePath, tenantHdl.PostCreateInvitation)

				// Channels
				r.Get(routes.OrgChannelsPath, channelHdl.GetChannels)
			})
		})
	})

	if err := http.ListenAndServe(cfg.Port, r); err != nil {
		log.Panicf("running server on port: %s", cfg.Port)
	}
}
