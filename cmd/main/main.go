package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"tmpl/internal/handler"
	"tmpl/internal/middleware"
	"tmpl/internal/provider/postgres"
	"tmpl/internal/repository"
	"tmpl/internal/service"
	"tmpl/internal/shared/config"
	"tmpl/internal/shared/routes"
	"tmpl/internal/view"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Panicf("loading config: %s", err)
	}

	logger := slog.Default()

	pool, err := postgres.LoadPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Panicf("loading postgres: %s", err)
	}

	err = view.LoadManifest()
	if err != nil {
		log.Panicf("loading manifest: %s", err)
	}

	err = routes.Load(cfg.GoogleRedirectURL)
	if err != nil {
		log.Panicf("loading routes: %s", err)
	}

	r := chi.NewRouter()

	userRepo := repository.NewUser(pool)
	authRepo := repository.NewAuth(pool)

	userSrv := service.NewUser(userRepo)
	authSrv := service.NewAuth(cfg, authRepo, userSrv)

	coreHdl := handler.NewCore()
	authHdl := handler.NewAuth(cfg, logger, authSrv)
	counterHdl := handler.NewCounter()

	authMdw := middleware.NewAuth(logger, authSrv)

	r.Group(func(r chi.Router) {
		r.Get(routes.Static, coreHdl.GetStatic)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMdw.RequireGuest)
		r.Get(routes.Login, authHdl.GetLogin)
		r.Get(routes.ILoginGoogle, authHdl.GetLoginWithGoogle)
		r.Get(routes.ILoginGoogleCallback, authHdl.GetLoginGoogleCallback)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMdw.RequireAuth)
		r.Get(routes.Index, counterHdl.GetCounter)
		r.Post(routes.ILogout, authHdl.PostLogout)
	})

	if err := http.ListenAndServe(cfg.Port, r); err != nil {
		log.Panicf("running server on port: %s", cfg.Port)
	}
}
