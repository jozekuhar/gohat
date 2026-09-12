package view

import (
	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/authz"

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

func (v *dashboard) OrganizationsPage(identity authz.Identity, orgs []db.Organization) g.Node {
	return v.layout.app(
		identity,
		h.Main(
			g.Text("Index"),
		),
	)
}

func (v *dashboard) DashboardPage(identity authz.Identity) g.Node {
	return v.layout.app(
		identity,
		h.Main(
			g.Text("Dashboard"),
		),
	)
}
