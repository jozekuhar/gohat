package auth

import (
	"log/slog"
	"net/http"

	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/web/cookie"
	"mimokocke/internal/web/response"
	"mimokocke/internal/web/view"

	"github.com/go-playground/form"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	logger    *slog.Logger
	decoder   *form.Decoder
	validator *validator.Validate
	srv       *Service
}

func NewHandler(
	logger *slog.Logger,
	decoder *form.Decoder,
	validator *validator.Validate,
	srv *Service,
) *Handler {
	return &Handler{
		logger:    logger,
		decoder:   decoder,
		validator: validator,
		srv:       srv,
	}
}

func (h *Handler) GetLogin(w http.ResponseWriter, r *http.Request) {
	response.RenderStatus(w, loginPage(), http.StatusOK)
}

func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parsing login form", "err", err)
		response.RenderStatus(
			w,
			view.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	var form struct {
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}

	err = h.decoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decoding login form", "err", err)
		response.RenderStatus(
			w,
			view.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	err = h.validator.Struct(form)
	if err != nil {
		h.logger.Error("validating login form", "err", err)
		response.RenderStatus(
			w,
			view.ToastFragment("Something went wrong"),
			http.StatusUnprocessableEntity,
		)
		return
	}

	s, err := h.srv.LoginUserWithPassword(r.Context(), form.Email, form.Password)
	if err != nil {
		h.logger.Error("logging in user with password", "err", err)
		response.RenderStatus(
			w,
			view.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	cookie.SetSession(w, s.ID.String(), s.ExpiresAt)
	response.HXRedirect(w, routes.AppRoot)
	response.Status(w, http.StatusSeeOther)
}

func (h *Handler) GetRegister(w http.ResponseWriter, r *http.Request) {
	response.RenderStatus(w, registerPage(), http.StatusOK)
}

func (h *Handler) PostRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parsing form", "err", err)
		response.RenderStatus(
			w,
			view.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	var form struct {
		Email           string
		Password        string
		ConfirmPassword string
	}

	err = h.decoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decoding form", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		response.RenderStatus(
			w,
			view.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	if form.Password != form.ConfirmPassword {
		response.RenderStatus(
			w,
			view.ToastFragment("Password missmatch"),
			http.StatusUnprocessableEntity,
		)
		return
	}

	s, err := h.srv.RegisterUserWithPassword(
		r.Context(),
		RegisterUserWithPasswordParams{
			Email:    form.Email,
			Password: form.Password,
		},
	)
	if err != nil {
		h.logger.Error("registering user", "err", err)
		response.RenderStatus(
			w,
			view.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	cookie.SetSession(w, s.ID.String(), s.ExpiresAt)
	response.HXRedirect(w, routes.AppRoot)
	response.Status(w, http.StatusSeeOther)
}

func (h *Handler) GetSignInWithGoogle(w http.ResponseWriter, r *http.Request) {
	redirectURL, state := h.srv.GenerateGoogleSignInURLAndState(r.Context())

	cookie.SetOAuthState(w, state)
	response.HXRedirect(w, redirectURL)
	response.Status(w, http.StatusTemporaryRedirect)
}

// For errors we don't use toasts because this is webhook page.
func (h *Handler) GetSignInGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	cookieState, err := cookie.GetOAuthState(r)
	if err != nil {
		h.logger.Error("getting oauth state cookie", "err", err)
		return
	}

	if state != cookieState {
		h.logger.Error("url and cookie state missmatch")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s, err := h.srv.ProcessGoogleSignIn(r.Context(), code)
	if err != nil {
		h.logger.Error("processing google login", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie.SetSession(w, s.ID.String(), s.ExpiresAt)
	http.Redirect(w, r, routes.AppRoot, http.StatusSeeOther)
}

func (h *Handler) PostLogout(w http.ResponseWriter, r *http.Request) {
	sessionIDStr, err := cookie.GetSession(r)
	if err != nil {
		h.logger.Error("retrieving session cookie", "err", err)
		response.RenderStatus(
			w,
			view.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	err = h.srv.LogoutUser(r.Context(), sessionIDStr)
	if err != nil {
		h.logger.Error("logging out user", "err", err)
		response.RenderStatus(
			w,
			view.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	cookie.ClearSession(w)
	response.HXRefresh(w)
	response.Status(w, http.StatusAccepted)
}

func (h *Handler) GetProfileSettings(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustIdentityFromContext(r.Context())

	response.RenderStatus(w, profileSettingsPage(ident), http.StatusOK)
}
