package middleware

import (
	"log/slog"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/tenant"
)

type tenantMiddleware struct {
	logger    *slog.Logger
	tenantSrv *tenant.Service
}

func NewTenantMiddleware(logger *slog.Logger, tenantSrv *tenant.Service) *tenantMiddleware {
	return &tenantMiddleware{
		logger:    logger,
		tenantSrv: tenantSrv,
	}
}

func (m *tenantMiddleware) RequireMembership(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := auth.MustUserIDFomContext(r.Context())
		orgSlug := r.PathValue(routes.PathOrganizationSlug)

		i, err := m.tenantSrv.VerifyMembership(r.Context(), userID, orgSlug)
		if err != nil {
			m.logger.Error("verifying membership", "err", err)
			return
		}

		ctx := tenant.WithIdentity(r.Context(), i)
		newR := r.WithContext(ctx)

		handler.ServeHTTP(w, newR)
	})
}
