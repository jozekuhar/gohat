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

type Memberships struct {
	layout *layout
}

func NewMemberships() *Memberships {
	return &Memberships{
		layout: newLayout(),
	}
}

func (v *Memberships) MembershipsPage(
	identity *authz.Identity,
	memberships []tenant.Membership,
) g.Node {
	return v.layout.Auth(
		h.Main(
			h.Class("max-w-7xl mx-auto"),
			h.H1(g.Text("Members")),
			h.Div(
				g.Text("Members list"),
			),
			h.Ul(
				g.Map(memberships, func(m tenant.Membership) g.Node {
					return h.Li(
						h.Span(
							g.Text(m.ID.String()),
						),
						h.Span(
							g.Text(string(m.Role)),
						),
					)
				}),
			),
			h.Form(
				hx.Post(fmt.Sprintf(routes.MembershipsCreate, identity.OrganizationSlug)),
				h.Class("border"),
				h.Input(h.Placeholder("owner")),
				h.Select(
					h.Option(g.Text("Owner")),
					h.Option(g.Text("Admin")),
					h.Option(g.Text("Member")),
				),
				h.Input(
					h.Placeholder("Email"),
					h.Name("email"),
				),
				// Permissions
				h.Button(
					g.Text("Invite membership"),
				),
			),
		),
	)
}
