package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"mimokocke/internal/auth"
	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/identity"
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
		authCtx := auth.MustFromContext(r.Context())
		orgSlug := r.PathValue(routes.PathOrganizationSlug)

		membership, err := m.tenantSrv.GetActiveMembership(r.Context(), authCtx.UserID, orgSlug)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				m.coreHdl.GetNotFound(w, r)
				return
			}
			m.logger.Error("verifying membership", "err", err)
			return
		}

		ident := identity.Identity{
			ID:          authCtx.UserID,
			Email:       membership.User.Email,
			FirstName:   membership.Membership.FirstName,
			LastName:    membership.Membership.LastName,
			OrgID:       membership.Organization.ID,
			OrgName:     membership.Organization.Name,
			OrgSlug:     orgSlug,
			Role:        membership.Membership.Role,
			Permissions: membership.Membership.Permissions,
		}

		ctx := identity.WithContext(r.Context(), ident)
		newR := r.WithContext(ctx)

		handler.ServeHTTP(w, newR)
	})
}
