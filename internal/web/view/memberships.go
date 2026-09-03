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
	identity authz.Identity,
	data tenant.MembershipsData,
) g.Node {
	return v.layout.App(
		h.Main(
			h.Class("max-w-7xl mx-auto"),
			h.H1(g.Text("Members")),
			h.Div(
				g.Text("Members list"),
			),
			h.H2(g.Text("Memberhips")),
			h.Ul(
				g.Map(data.Memberhips, func(membership tenant.Membership) g.Node {
					return h.Li(
						h.Span(
							g.Text(membership.ID.String()),
						),
						h.Span(
							g.Text(string(membership.Role)),
						),
					)
				}),
			),
			h.H2(g.Text("Invitations")),
			h.Ul(
				g.Map(data.Invitations, func(invitation tenant.Invitation) g.Node {
					return h.Li(
						h.Class("flex gap-4"),
						h.Span(
							g.Text(invitation.ID.String()),
						),
						h.Span(
							g.Text(string(invitation.Role)),
						),
						h.Button(
							hx.Delete(
								fmt.Sprintf(
									routes.InvitationsDetail,
									identity.OrganizationSlug,
									invitation.ID.String(),
								),
							),
							hx.Swap("none"),
							g.Text("Delete"),
						),
					)
				}),
			),
			h.Form(
				hx.Post(fmt.Sprintf(routes.Invitations, identity.OrganizationSlug)),
				hx.Swap("none"),
				h.Class("border flex flex-col"),
				h.Input(
					h.Class("border"),
					h.Placeholder("First name"),
					h.Name("FirstName"),
				),
				h.Input(
					h.Class("border"),
					h.Placeholder("Last name"),
					h.Name("LastName"),
				),
				h.Input(
					h.Class("border"),
					h.Placeholder("Email"),
					h.Name("Email"),
				),
				h.Select(
					h.Name("Role"),
					h.Option(
						h.Value(string(authz.RoleAdmin)),
						g.Text("Owner"),
					),
					h.Option(
						h.Value(string(authz.RoleOwner)),
						g.Text("Admin"),
					),
					h.Option(
						h.Value(string(authz.RoleMember)),
						h.Selected(),
						g.Text("Member"),
					),
				),
				h.FieldSet(
					h.Class("flex flex-col border"),
					h.Legend(g.Text("Memberships")),
					h.Label(
						g.Text("Memberships (Read)"),
						h.Input(
							h.Type("checkbox"),
							h.Name("Permissions"),
							h.Value(string(authz.PermMembershipRead)),
						),
					),
					h.Label(
						g.Text("Memberships (Create)"),
						h.Input(
							h.Type("checkbox"),
							h.Name("Permissions"),
							h.Value(string(authz.PermMembershipCreate)),
						),
					),
					h.Label(
						g.Text("Memberships (Update)"),
						h.Input(
							h.Type("checkbox"),
							h.Name("Permissions"),
							h.Value(string(authz.PermMembershipUpdate)),
						),
					),
					h.Label(
						g.Text("Memberships (Delete)"),
						h.Input(
							h.Type("checkbox"),
							h.Name("Permissions"),
							h.Value(string(authz.PermMembershipDelete)),
						),
					),
				),
				// Permissions
				h.Button(
					h.Class("bg-black text-white"),
					g.Text("Invite membership"),
				),
			),
		),
	)
}
