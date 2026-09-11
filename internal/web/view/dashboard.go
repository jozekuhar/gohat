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

type Dashboard struct {
	layout *Layout
}

func NewDashboard() *Dashboard {
	return &Dashboard{
		layout: NewLayout(),
	}
}

func (v *Dashboard) OrganizationsPage(identity authz.Identity, orgs []db.Organization) g.Node {
	return v.layout.app(
		identity,
		h.Main(
			g.Text("Index"),
		),
	)
}

func (v *Dashboard) DashboardPage(identity authz.Identity) g.Node {
	return v.layout.app(
		identity,
		h.Main(
			g.Text("Dashboard"),
		),
	)
}
