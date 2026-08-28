package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/shared/config"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/web/cookie"
	"mimokocke/internal/web/view"

	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
	hx "maragu.dev/gomponents-htmx/http"
)

type authHandler struct {
	cfg         *config.Config
	logger      *slog.Logger
	authSrv     *auth.Service
	formDecoder *form.Decoder
	validator   *validator.Validate
	authView    *view.Auth
}

func NewAuthHandler(
	cfg *config.Config,
	logger *slog.Logger,
	authSrv *auth.Service,
	formDecoder *form.Decoder,
	validator *validator.Validate,
) *authHandler {
	return &authHandler{
		cfg:         cfg,
		logger:      logger,
		authSrv:     authSrv,
		formDecoder: formDecoder,
		validator:   validator,
		authView:    view.NewAuth(),
	}
}

func (h *authHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	render(w, h.authView.LoginPage())
}

type loginForm struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

func (h *authHandler) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parsing form", "err", err)
		return
	}

	var form loginForm
	err = h.formDecoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decoding form", "err", err)
		return
	}

	err = h.validator.Struct(form)
	if err != nil {
		h.logger.Error("validating form", "err", err)
		return
	}

	s, err := h.authSrv.LoginUserWithPassword(r.Context(), form.Email, form.Password)
	if err != nil {
		h.logger.Error("logging in user with password", "err", err)
		return
	}

	cookie.SetSession(w, s.ID.String(), s.ExpiresAt)
	hx.SetRedirect(w.Header(), routes.Index)
	w.WriteHeader(http.StatusSeeOther)
}

func (h *authHandler) GetRegister(w http.ResponseWriter, r *http.Request) {
	render(w, h.authView.RegisterPage())
}

type registerForm struct {
	Email           string
	Password        string
	ConfirmPassword string
}

func (h *authHandler) PostRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parsing form", "err", err)
		return
	}

	var form registerForm
	err = h.formDecoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decoding form", "err", err)
		return
	}

	if form.Password != form.ConfirmPassword {
		// TODO(jozekuhar): error happened
		return
	}

	s, err := h.authSrv.RegisterUserWithPassword(
		r.Context(),
		auth.RegisterUserWithPasswordParams{
			Email:    form.Email,
			Password: form.Password,
		},
	)
	if err != nil {
		// TODO(jozekuhar): error happened
		h.logger.Error("registering user", "err", err)
		return
	}

	fmt.Println("redirect?")

	cookie.SetSession(w, s.ID.String(), s.ExpiresAt)
	hx.SetRedirect(w.Header(), routes.Index)
}

func (h *authHandler) GetSignInWithGoogle(w http.ResponseWriter, r *http.Request) {
	redirectURL, state := h.authSrv.GenerateGoogleSignInURLAndState(r.Context())

	cookie.SetOAuthState(w, state)
	hx.SetRedirect(w.Header(), redirectURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *authHandler) GetSignInGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	cookieState, err := cookie.GetOAuthState(r)
	if err != nil {
		return
	}

	if state != cookieState {
		h.logger.Error("url and cookie state missmatch")
		return
	}

	s, err := h.authSrv.ProcessGoogleSignIn(r.Context(), code)
	if err != nil {
		h.logger.Error("processing google login", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie.SetSession(w, s.ID.String(), s.ExpiresAt)
	http.Redirect(w, r, routes.Index, http.StatusSeeOther)
}

func (h *authHandler) PostLogout(w http.ResponseWriter, r *http.Request) {
	sessionIDStr, err := cookie.GetSession(r)
	if err != nil {
		h.logger.Error("retrieving session cookie", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.authSrv.LogoutUser(r.Context(), sessionIDStr)
	if err != nil {
		h.logger.Error("logging out user", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie.ClearSession(w)
	hx.SetRefresh(w.Header())
	w.WriteHeader(http.StatusAccepted)
}
