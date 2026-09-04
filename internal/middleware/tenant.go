package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/shared/authz"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/tenant"
	"mimokocke/internal/web/handler"
)

type tenantMiddleware struct {
	logger    *slog.Logger
	tenantSrv *tenant.Service
	coreHdl   *handler.Core
}

func NewTenantMiddleware(
	logger *slog.Logger,
	tenantSrv *tenant.Service,
	coreHdl *handler.Core,
) *tenantMiddleware {
	return &tenantMiddleware{
		logger:    logger,
		tenantSrv: tenantSrv,
		coreHdl:   coreHdl,
	}
}

func (m *tenantMiddleware) RequireIdentity(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := auth.MustUserIDFomContext(r.Context())
		orgSlug := r.PathValue(routes.PathOrganizationSlug)

		// TODO(jozekuhar): change to identity, so i just go for membership
		am, err := m.tenantSrv.GetActiveMembership(r.Context(), userID, orgSlug)
		if err != nil {
			if errors.Is(err, tenant.ErrMembershipNotFound) {
				w.WriteHeader(http.StatusNotFound)
				m.coreHdl.GetNotFound(w, r)
				return
			}
			m.logger.Error("verifying membership", "err", err)
			return
		}

		identity := authz.Identity{
			UserID:           userID,
			OrganizationID:   am.OrganizationID,
			OrganizationSlug: orgSlug,
			Role:             am.Role,
			Permissions:      am.Permissions,
		}

		identity.UserID = userID

		ctx := tenant.WithIdentity(r.Context(), identity)
		newR := r.WithContext(ctx)

		handler.ServeHTTP(w, newR)
	})
}
