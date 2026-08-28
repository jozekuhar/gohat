package view

import (
	"fmt"

	"mimokocke/internal/shared/authz"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/tenant"

	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

const (
	InputNameOrganizationName = "name"
	InputNameOrganizationSlug = "slug"
)

type Dashboard struct {
	layout *layout
}

func NewDashboard() *Dashboard {
	return &Dashboard{
		layout: newLayout(),
	}
}

func (v *Dashboard) OrganizationsPage(orgs []tenant.Organization) g.Node {
	return v.layout.Auth(
		h.Main(
			h.Class("mx-auto max-w-7xl"),
			h.Div(
				h.H1(g.Text("Organizations")),
				h.Ul(
					g.Map(orgs, func(org tenant.Organization) g.Node {
						return h.Li(
							h.Class("flex gap-4"),
							h.Span(g.Text(org.ID.String())),
							h.Span(g.Text(org.Name)),
							h.A(
								h.Class("text-blue-500"),
								h.Href(fmt.Sprintf(routes.Dashboard, org.Slug)),
								g.Text("select"),
							),
						)
					}),
				),
			),
			h.Div(
				h.Class("border"),
				h.Form(
					hx.Post(routes.HXOrganizationCreate),
					hx.Swap("none"),
					h.Input(
						h.Name("Name"),
						h.Placeholder("Organization name"),
					),
					h.Input(
						h.Name("Slug"),
						h.Placeholder("Organization slug"),
					),
					h.Input(
						h.Name("FirstName"),
						h.Placeholder("First name"),
					),
					h.Input(
						h.Name("LastName"),
						h.Placeholder("Last name"),
					),
					h.Button(g.Text("Create")),
				),
			),
		),
	)
}

func (v *Dashboard) DashboardPage(identity authz.Identity) g.Node {
	return v.layout.Auth(
		h.Main(
			h.Class("mx-auto max-w-7xl"),
			h.Div(
				g.Text("USER_ID"),
				g.Text(identity.UserID.String()),
			),
			h.Div(
				g.Text("ORG_ID"),
				g.Text(identity.OrganizationID.String()),
			),
			h.Div(
				g.Text("ORG_ROLE"),
				g.Text(identity.Role.String()),
			),
			h.Div(
				h.Class("flex flex-col"),
				h.A(
					h.Class("text-blue-500"),
					h.Href(fmt.Sprintf(routes.Memberships, identity.OrganizationSlug)),
					g.Text("Members"),
				),
				h.A(
					g.Text("Orders"),
				),
				h.A(
					h.Class("text-blue-500"),
					h.Href(fmt.Sprintf(routes.Channels, identity.OrganizationSlug)),
					g.Text("Channels"),
				),
				h.A(
					g.Text("Couriers"),
				),
			),
		),
	)
}
