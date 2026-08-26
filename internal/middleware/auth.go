package middleware

import (
	"log/slog"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/shared/session"
)

type authMiddleware struct {
	logger  *slog.Logger
	authSrv *auth.Service
}

func NewAuthMiddleware(logger *slog.Logger, authSrv *auth.Service) *authMiddleware {
	return &authMiddleware{
		logger:  logger,
		authSrv: authSrv,
	}
}

func (m *authMiddleware) RequireGuest(handler http.Handler) http.Handler {
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

func (m *authMiddleware) RequireAuth(handler http.Handler) http.Handler {
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

		ctx := auth.WithUserID(r.Context(), s.UserID)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
