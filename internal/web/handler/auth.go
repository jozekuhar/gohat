package handler

import (
	"log/slog"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/shared/config"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/shared/session"
	"mimokocke/internal/web/view"

	hx "maragu.dev/gomponents-htmx/http"
)

type authHandler struct {
	cfg      *config.Config
	logger   *slog.Logger
	authSrv  *auth.Service
	authView *view.Auth
}

func NewAuthHandler(cfg *config.Config, logger *slog.Logger, authSrv *auth.Service) *authHandler {
	return &authHandler{
		cfg:      cfg,
		logger:   logger,
		authSrv:  authSrv,
		authView: view.NewAuth(),
	}
}

func (h *authHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	render(w, h.authView.LoginPage())
}

func (h *authHandler) GetLoginWithGoogle(w http.ResponseWriter, r *http.Request) {
	redirectURL := h.authSrv.GenerateGoogleLoginURL(r.Context())
	hx.SetRedirect(w.Header(), redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *authHandler) GetLoginGoogleCallback(w http.ResponseWriter, r *http.Request) {
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

func (h *authHandler) PostLogout(w http.ResponseWriter, r *http.Request) {
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
