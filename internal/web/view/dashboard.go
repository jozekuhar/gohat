package view

import (
	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/identity"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

const (
	InputNameOrganizationName = "name"
	InputNameOrganizationSlug = "slug"
)

type dashboard struct {
	layout *Layout
}

func newDashboard(layout *Layout) *dashboard {
	return &dashboard{
		layout: layout,
	}
}

func (v *dashboard) OrganizationsPage(ident identity.Identity, orgs []db.Organization) g.Node {
	return v.layout.app(
		ident,
		h.Main(
			g.Text("Index"),
		),
	)
}

func (v *dashboard) DashboardPage(ident identity.Identity) g.Node {
	return v.layout.app(
		ident,
		h.Main(
			g.Text("Dashboard"),
		),
	)
}
