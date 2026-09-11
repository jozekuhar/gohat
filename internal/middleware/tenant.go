package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/provider/db"
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
		authCtx := auth.MustAuthFromContext(r.Context())
		orgSlug := r.PathValue(routes.PathOrganizationSlug)

		// TODO(jozekuhar): change to identity, so i just go for membership
		am, err := m.tenantSrv.GetActiveMembership(r.Context(), authCtx.UserID, orgSlug)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				m.coreHdl.GetNotFound(w, r)
				return
			}
			m.logger.Error("verifying membership", "err", err)
			return
		}

		identity := authz.Identity{
			ID:          authCtx.UserID,
			Email:       authCtx.UserEmail,
			FirstName:   am.Membership.FirstName,
			LastName:    am.Membership.LastName,
			OrgID:       am.Organization.ID,
			OrgName:     am.Organization.Name,
			OrgSlug:     orgSlug,
			Role:        am.Membership.Role,
			Permissions: am.Membership.Permissions,
		}

		ctx := tenant.WithIdentity(r.Context(), identity)
		newR := r.WithContext(ctx)

		handler.ServeHTTP(w, newR)
	})
}
