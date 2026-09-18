package main

import (
	"context"
	"log"
	"net/http"

	"mimokocke/internal/domain/auth"
	"mimokocke/internal/domain/channel"
	"mimokocke/internal/domain/courier"
	"mimokocke/internal/domain/order"
	"mimokocke/internal/domain/tenant"
	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/clock"
	"mimokocke/internal/shared/config"
	"mimokocke/internal/shared/logger"
	"mimokocke/internal/shared/routes"
	orderuc "mimokocke/internal/usecase/order"
	errorweb "mimokocke/internal/web/handler/error"
	orderweb "mimokocke/internal/web/handler/order"
	staticweb "mimokocke/internal/web/handler/static"
	"mimokocke/internal/web/middleware"
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

	pool, err := db.ConnectAndMigrate(ctx, cfg.DatabaseURL, logger, cfg.Debug)
	if err != nil {
		log.Panicf("connect and migrate db")
	}

	resendClient := resend.NewClient(cfg.ResendAPIKey)

	baseRepo := db.NewBaseRepo(pool)

	tenantRepo := tenant.NewRepository(baseRepo)
	authRepo := auth.NewRepository(baseRepo)
	channelRepo := channel.NewRepository(baseRepo)
	orderRepo := order.NewRepository(baseRepo)
	courierRepo := courier.NewRepository(baseRepo)

	tenantSrv := tenant.NewService(cfg, logger, systemClock, resendClient, tenantRepo)
	authSrv := auth.NewService(cfg, logger, systemClock, authRepo)
	channelSrv := channel.NewService(cfg, channelRepo)
	orderSrv := order.NewService(orderRepo)
	courierSrv := courier.NewService(courierRepo)

	orderUCRepo := orderuc.NewRepository(baseRepo)
	listOrdersUC := orderuc.NewListOrdersUseCase(orderUCRepo)
	createOrderUC := orderuc.NewCreateOrderUseCase(orderSrv)
	_ = createOrderUC

	decoder := form.NewDecoder()
	validator := validator.New(validator.WithRequiredStructEnabled())

	err = view.LoadManifest()
	if err != nil {
		log.Panicf("load manifest: %s", err)
	}

	authHdl := auth.NewHandler(logger, decoder, validator, authSrv)
	tenantHdl := tenant.NewHandler(logger, decoder, tenantSrv)
	channelHdl := channel.NewHandler(logger, decoder, channelSrv)

	orderHdl := orderweb.NewHandler(listOrdersUC, channelSrv, courierSrv, logger)

	tenantMdw := middleware.NewTenantMiddleware(logger, tenantSrv)
	authMdw := middleware.NewAuthMiddleware(logger, authSrv)

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get(routes.Static, staticweb.GetStatic)
		r.NotFound(errorweb.GetNotFound)
		r.MethodNotAllowed(errorweb.GetNotFound)
	})

	r.Group(func(r chi.Router) {
		r.Use(chimiddleware.Logger)

		// Core
		r.Get(routes.TermsOfService, staticweb.GetTermsOfService)
		r.Get(routes.PrivacyPolicy, staticweb.GetPrivacyPolicy)

		// Auth
		r.Get(routes.CallbackSignInGoogle, authHdl.GetSignInGoogleCallback)
	})

	r.Group(func(r chi.Router) {
		r.Use(chimiddleware.Logger)
		r.Use(authMdw.OptionalAuth)

		// Invitations
		r.Get(routes.InvitationsJoinPath, tenantHdl.GetShowInvitation)
	})

	r.Group(func(r chi.Router) {
		r.Use(chimiddleware.Logger)
		r.Use(authMdw.RequireGuest)

		r.Get(routes.Login, authHdl.GetLogin)
		r.Post(routes.HXLogin, authHdl.PostLogin)
		r.Get(routes.Register, authHdl.GetRegister)
		r.Post(routes.HXRegister, authHdl.PostRegister)
		r.Get(routes.HXSignInGoogle, authHdl.GetSignInWithGoogle)
	})

	r.Group(func(r chi.Router) {
		r.Use(chimiddleware.Logger)
		r.Use(authMdw.RequireAuth)

		// App
		r.Post(routes.HXLogout, authHdl.PostLogout)
		r.Get(routes.AppRoot, tenantHdl.GetAppRoot)

		// Sidebar
		r.Get(routes.HXSidebarOrganizations, tenantHdl.GetSidebarOrganizationsPartial)
		r.Get(routes.HXSidebarOrganizationsCreate, tenantHdl.GetCreateOrganizationFormModal)
		r.Post(routes.HXSidebarOrganizationsCreate, tenantHdl.PostCreateOrganization)

		// Invitations
		r.Post(routes.HXInvitationsAcceptPath, tenantHdl.PostAcceptInvitation)
		r.Post(routes.HXInvitationsDeclinePath, tenantHdl.PostDeclineInvitation)
	})

	r.Group(func(r chi.Router) {
		r.Use(chimiddleware.Logger)
		r.Use(authMdw.RequireAuth)
		r.Use(tenantMdw.RequireIdentity)

		// App
		r.Get(routes.OrgRootPath, tenantHdl.GetOrgRoot)
		r.Get(routes.OrgDashboardPath, tenantHdl.GetDashboard)

		// Auth
		r.Get(routes.SettingsProfilePath, authHdl.GetProfileSettings)
		r.Get(routes.SettingsGeneralPath, tenantHdl.GetGeneralSettings)

		// Memberships
		r.Get(routes.OrgMembershipsPath, tenantHdl.GetMemberships)
		r.Get(routes.HXOrgMembershipsUpdatePath, tenantHdl.GetUpdateMembershipFormModal)
		r.Patch(routes.HXOrgMembershipsUpdatePath, tenantHdl.PatchUpdateMembership)
		r.Get(routes.HXOrgInvitationsCreatePath, tenantHdl.GetCreateInvitationFormModal)
		r.Post(routes.HXOrgInvitationsCreatePath, tenantHdl.PostCreateInvitation)

		// Channels
		r.Get(routes.OrgChannelsPath, channelHdl.GetChannels)
		r.Get(routes.HXOrgChannelsCreatePath, channelHdl.GetCreateChannelFormModal)
		r.Post(routes.HXOrgChannelsWooCommerceCreatePath, channelHdl.PostCreateWooCommerceChannel)
		r.Post(routes.HXOrgChannelsWooCommerceTestPath, channelHdl.PostTestWooCommerceChannel)
		r.Get(routes.HXOrgChannelsUpdatePath, channelHdl.GetUpdateChannelModalForm)
		r.Post(routes.HXOrgChannelsDeactivatePath, channelHdl.PostDeactivateChannel)
		r.Delete(routes.HXOrgChannelsDeletePath, channelHdl.DeleteRemoveChannel)

		// Orders
		r.Get(routes.OrgOrdersPath, orderHdl.GetOrders)
		r.Get(routes.HXOrgOrderCreatePath, orderHdl.GetCreateOrderFormModal)
		r.Post(routes.HXOrgOrderCreatePath, orderHdl.PostCreateOrder)

		// Couriers
		// r.Get(routes.OrgCouriersPath, courierHdl.GetCouriers)
	})

	if err := http.ListenAndServe(cfg.Port, r); err != nil {
		log.Panicf("running server on port: %s", cfg.Port)
	}
}
