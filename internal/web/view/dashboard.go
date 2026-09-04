package view

import (
	"mimokocke/internal/shared/authz"
	"mimokocke/internal/tenant"

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

func (v *Dashboard) OrganizationsPage(identity authz.Identity, orgs []tenant.Organization) g.Node {
	return v.layout.app(
		identity,
		h.Main(
			g.Text("Index page"),
		),
	)
}

func (v *Dashboard) DashboardPage(identity authz.Identity) g.Node {
	return v.layout.app(
		identity,
		h.Main(
			g.Text("Organizations Detail Page"),
		),
	)
}
