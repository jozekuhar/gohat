package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"mimokocke/internal/domain/tenant"
	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/routes"
	errorweb "mimokocke/internal/web/handler/error"
)

type tenantMiddleware struct {
	logger    *slog.Logger
	tenantSrv *tenant.Service
}

func NewTenantMiddleware(
	logger *slog.Logger,
	tenantSrv *tenant.Service,
) *tenantMiddleware {
	return &tenantMiddleware{
		logger:    logger,
		tenantSrv: tenantSrv,
	}
}

func (m *tenantMiddleware) RequireIdentity(hdl http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCtx := identity.MustAuthFromContext(r.Context())
		orgSlug := r.PathValue(routes.PathOrganizationSlug)

		membership, err := m.tenantSrv.GetActiveMembership(r.Context(), authCtx.UserID, orgSlug)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				errorweb.GetNotFound(w, r)
				return
			}
			m.logger.Error("verifying membership", "err", err)
			return
		}

		ident := identity.IdentityCtx{
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

		ctx := identity.WithIdentity(r.Context(), ident)
		newR := r.WithContext(ctx)

		hdl.ServeHTTP(w, newR)
	})
}
