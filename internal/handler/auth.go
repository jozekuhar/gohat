package handler

import (
	"log/slog"
	"net/http"

	"tmpl/internal/service"
	"tmpl/internal/shared/config"
	"tmpl/internal/shared/routes"
	"tmpl/internal/shared/session"
	"tmpl/internal/shared/web"
	"tmpl/internal/view"

	hx "maragu.dev/gomponents-htmx/http"
)

type Auth struct {
	cfg      *config.Config
	logger   *slog.Logger
	authSrv  *service.Auth
	authView *view.Auth
}

func NewAuth(cfg *config.Config, logger *slog.Logger, authSrv *service.Auth) *Auth {
	return &Auth{
		cfg:      cfg,
		logger:   logger,
		authSrv:  authSrv,
		authView: view.NewAuth(),
	}
}

func (h *Auth) GetLogin(w http.ResponseWriter, r *http.Request) {
	web.Render(w, h.authView.LoginPage())
}

func (h *Auth) GetLoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	redirectURL := h.authSrv.GenerateGoogleLoginURL(r.Context())
	hx.SetRedirect(w.Header(), redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Auth) GetLoginGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	s, err := h.authSrv.ProcessGoogleLogin(r.Context(), code, state)
	if err != nil {
		h.logger.Error("processing google login", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session.SetCookie(w, s.ID.String(), s.ExpiresAt)
	http.Redirect(w, r, routes.Index, http.StatusSeeOther)
}

func (h *Auth) PostLogout(w http.ResponseWriter, r *http.Request) {
	sessionIDStr, err := session.GetCookie(r)
	if err != nil {
		h.logger.Error("retrieving session cookie", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.authSrv.Logout(r.Context(), sessionIDStr)
	if err != nil {
		h.logger.Error("logging out user", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.ClearCookie(w)
	hx.SetRefresh(w.Header())
	w.WriteHeader(http.StatusAccepted)
}
