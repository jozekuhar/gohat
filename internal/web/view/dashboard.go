package view

import (
	"mimokocke/internal/model"
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

func (v *Dashboard) OrganizationsPage(identity authz.Identity, orgs []model.Organization) g.Node {
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
