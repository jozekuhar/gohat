package middleware

import (
	"log/slog"
	"net/http"

	"gohat/internal/service"
	"gohat/internal/shared/ctxutil"
	"gohat/internal/shared/routes"
	"gohat/internal/shared/session"
)

type Auth struct {
	logger  *slog.Logger
	authSrv *service.Auth
}

func NewAuth(logger *slog.Logger, authSrv *service.Auth) *Auth {
	return &Auth{
		logger:  logger,
		authSrv: authSrv,
	}
}

func (m *Auth) RequireGuest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionIDStr, err := session.GetCookie(r)
		if err != nil {
			handler.ServeHTTP(w, r)
			return
		}

		_, err = m.authSrv.VerifySession(r.Context(), sessionIDStr)
		if err != nil {
			m.logger.Warn("verifying session", "err", err)
			session.ClearCookie(w)
			handler.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, routes.Index, http.StatusSeeOther)
	})
}

func (m *Auth) RequireAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionIDStr, err := session.GetCookie(r)
		if err != nil {
			http.Redirect(w, r, routes.Login, http.StatusSeeOther)
			return
		}

		s, err := m.authSrv.VerifySession(r.Context(), sessionIDStr)
		if err != nil {
			m.logger.Warn("verifying session", "err", err)
			session.ClearCookie(w)
			http.Redirect(w, r, routes.Login, http.StatusSeeOther)
			return
		}

		ctx := ctxutil.WithUserID(r.Context(), s.UserID)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
