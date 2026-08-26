package main

import (
	"context"
	"log"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/channel"
	"mimokocke/internal/middleware"
	"mimokocke/internal/provider/postgres"
	"mimokocke/internal/shared/clock"
	"mimokocke/internal/shared/config"
	"mimokocke/internal/shared/logger"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/tenant"
	"mimokocke/internal/web/handler"
	"mimokocke/internal/web/view"

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
	_ = systemClock

	pool, err := postgres.LoadPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Panicf("loading postgres: %s", err)
	}

	resendClient := resend.NewClient(cfg.ResendAPIKey)

	err = view.LoadManifest()
	if err != nil {
		log.Panicf("loading manifest: %s", err)
	}

	r := chi.NewRouter()

	tenantRepo := tenant.NewRepository(pool)
	authRepo := auth.NewRepository(pool)
	channelRepo := channel.NewRepository(pool)

	tenantSrv := tenant.NewService(logger, resendClient, tenantRepo)
	authSrv := auth.NewService(cfg, authRepo)
	channelSrv := channel.NewService(channelRepo)

	staticHdl := handler.NewStatic()
	fallbackHdl := handler.NewFallback()
	tenantHdl := handler.NewTenantHandler(logger, tenantSrv)
	authHdl := handler.NewAuthHandler(cfg, logger, authSrv)
	channelHdl := handler.NewChanneHandler(logger, channelSrv)

	tenantMdw := middleware.NewTenantMiddleware(logger, tenantSrv)
	authMdw := middleware.NewAuthMiddleware(logger, authSrv)

	r.Group(func(r chi.Router) {
		r.Get(routes.Static, staticHdl.GetStatic)
		r.NotFound(fallbackHdl.GetNotFound)
		r.MethodNotAllowed(fallbackHdl.GetMethodNotAllowed)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMdw.RequireGuest)
		r.Get(routes.Login, authHdl.GetLogin)
		r.Get(routes.ILoginGoogle, authHdl.GetLoginWithGoogle)
		r.Get(routes.ILoginGoogleCallback, authHdl.GetLoginGoogleCallback)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMdw.RequireAuth)
		r.Post(routes.ILogout, authHdl.PostLogout)
		r.Get(routes.Index, tenantHdl.GetOrganizations)
		r.Post(routes.HXOrganizationCreate, tenantHdl.PostCreateOrganization)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMdw.RequireAuth)
		r.Use(tenantMdw.RequireMembership)
		r.Get(routes.DashboardPath, tenantHdl.GetDashboard)
		r.Get(routes.MembershipsPath, tenantHdl.GetMemberships)
		r.Get(routes.MembershipsCreatePath, tenantHdl.PostCreateInvitation)
		r.Get(routes.ChannelsPath, channelHdl.GetChannels)
	})

	if err := http.ListenAndServe(cfg.Port, r); err != nil {
		log.Panicf("running server on port: %s", cfg.Port)
	}
}
