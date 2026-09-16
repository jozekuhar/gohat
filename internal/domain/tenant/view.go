package tenant

import (
	"fmt"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/permissions"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/web/components"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

func appPage(orgs []db.Organization) g.Node {
	return components.BlankLayout(
		h.Div(
			h.Class(
				"bg-muted flex h-svh flex-col items-center justify-center lg:max-w-none lg:px-0",
			),
			h.Div(
				h.Class("flex flex-row gap-4"),
				g.Map(orgs, func(org db.Organization) g.Node {
					return h.A(
						hx.Boost("true"),
						h.Class("bg-card flex flex-col gap-4 rounded-2xl border p-5"),
						h.Href(fmt.Sprintf(routes.OrgDashboard, org.Slug)),
						h.Span(
							g.Text(org.Name),
						),
						h.Span(
							g.Textf("@%s", org.Slug),
						),
					)
				}),
			),
		),
	)
}

func dashboardPage(ident identity.IdentityCtx) g.Node {
	return components.AppLayout(
		ident,
		h.Main(
			g.Text("Dashboard"),
		),
	)
}

func membershipsPage(ident identity.IdentityCtx, data MembershipsOverview) g.Node {
	return components.AppLayout(
		ident,
		h.Div(
			h.Class("flex h-full flex-col overflow-hidden"),
			h.Div(
				h.Class("flex min-h-0 flex-1 flex-col overflow-hidden"),
				h.Main(
					h.Class(
						"bg-background flex min-h-0 flex-1 flex-col overflow-hidden dark:bg-[radial-gradient(circle_at_top,rgba(255,255,255,0.035),transparent_28%)]",
					),
					h.Div(
						h.Class("border-b dark:border-white/8 dark:bg-[#111111]/95"),
						h.Div(
							h.Class(
								"flex flex-col gap-2 px-4 py-3 sm:px-6 lg:min-h-14 lg:flex-row lg:items-center lg:justify-between lg:gap-4 lg:py-0",
							),
							h.Div(
								h.Class("flex min-w-0 items-center justify-between gap-3"),
								h.Div(
									h.Class(
										"flex min-w-0 items-center gap-2.5 text-sm font-medium dark:text-zinc-100",
									),
									h.Span(
										h.Class("truncate"),
										g.Text("Members"),
									),
									h.Span(
										h.Class(
											"bg-muted text-muted-foreground shrink-0 rounded-md px-1.5 py-0.5 text-xs dark:bg-zinc-900 dark:text-zinc-300",
										),
										g.Textf("%d", len(data.Memberhips)),
									),
								),
								h.Button(
									hx.Get(
										fmt.Sprintf(
											routes.HXOrgInvitationsCreate,
											ident.OrgSlug,
										),
									),
									hx.Swap("none"),
									h.Class(
										"inline-flex items-center justify-center whitespace-nowrap font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 rounded-md text-xs size-9 shrink-0 gap-1.5 p-0 sm:w-auto sm:px-3 lg:hidden",
									),
									h.Aria("label", "Invite member"),
									g.Raw(
										`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-plus size-3.5" aria-hidden="true"><path d="M5 12h14"></path><path d="M12 5v14"></path></svg>`,
									),
									h.Span(
										h.Class("hidden sm:inline"),
										g.Text("Invite Member"),
									),
								),
							),
							h.Div(
								hx.Get(
									fmt.Sprintf(
										routes.HXOrgInvitationsCreate,
										ident.OrgSlug,
									),
								),
								hx.Swap("none"),
								h.Class(
									"grid min-w-0 grid-cols-2 items-center gap-2 lg:flex lg:flex-wrap lg:justify-end",
								),
								h.Button(
									h.Class(
										"items-center justify-center whitespace-nowrap font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 rounded-md px-3 text-xs hidden h-8 gap-1.5 lg:inline-flex",
									),
									g.Raw(
										`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-plus size-3.5" aria-hidden="true"><path d="M5 12h14"></path><path d="M12 5v14"></path></svg>`,
									),
									g.Text("Invite Member"),
								),
							),
						),
					),
					h.Div(
						h.Class(
							"flex-1 overflow-auto dark:bg-[linear-gradient(180deg,rgba(255,255,255,0.01),transparent_12%)]",
						),
						h.Div(
							h.Class(
								"min-w-225 border-b px-4 py-2 text-sm text-zinc-500 sm:px-6 lg:min-w-0 dark:border-white/8 dark:bg-[#141414]/95 dark:text-zinc-400",
							),
							h.Div(
								h.Class("flex items-center"),
								h.Div(
									h.Class("w-45 shrink-0 lg:w-[20%] lg:shrink"),
									g.Text("Member"),
								),
								h.Div(
									h.Class("w-45 shrink-0 lg:w-[20%] lg:shrink"),
									g.Text("Availability"),
								),
								h.Div(
									h.Class("w-45 shrink-0 lg:w-[20%] lg:shrink"),
									g.Text("Role"),
								),
								h.Div(
									h.Class("w-45 shrink-0 lg:w-[20%] lg:shrink"),
									g.Text("Teams"),
								),
								h.Div(
									h.Class("w-45 shrink-0 lg:w-[20%] lg:shrink"),
									g.Text("Joined"),
								),
							),
						),
						h.Div(
							h.ID("invitations_list"),
							g.Map(data.Invitations, func(invite db.Invitation) g.Node {
								return invitationItem(invite)
							}),
						),
					),
					h.Div(
						h.Class(
							"flex-1 overflow-auto dark:bg-[linear-gradient(180deg,rgba(255,255,255,0.01),transparent_12%)]",
						),
						h.Div(
							h.Class(
								"min-w-225 border-b px-4 py-2 text-sm text-zinc-500 sm:px-6 lg:min-w-0 dark:border-white/8 dark:bg-[#141414]/95 dark:text-zinc-400",
							),
							h.Div(
								h.Class("flex items-center"),
								h.Div(
									h.Class("w-45 shrink-0 lg:w-[20%] lg:shrink"),
									g.Text("Member"),
								),
								h.Div(
									h.Class("w-45 shrink-0 lg:w-[20%] lg:shrink"),
									g.Text("Availability"),
								),
								h.Div(
									h.Class("w-45 shrink-0 lg:w-[20%] lg:shrink"),
									g.Text("Role"),
								),
								h.Div(
									h.Class("w-45 shrink-0 lg:w-[20%] lg:shrink"),
									g.Text("Teams"),
								),
								h.Div(
									h.Class("w-45 shrink-0 lg:w-[20%] lg:shrink"),
									g.Text("Joined"),
								),
							),
						),
						h.Div(
							h.ID("memberships_list"),
							g.Map(data.Memberhips, func(member db.ListMembershipsRow) g.Node {
								return membershipItem(
									ident.OrgSlug,
									member.Membership,
									member.User,
								)
							}),
						),
					),
				),
			),
		),
	)
}

func invitationItem(invite db.Invitation) g.Node {
	return h.Div(
		h.Class(
			"hover:bg-muted/40 flex min-w-225 items-center border-b px-4 py-3 text-sm last:border-b-0 sm:px-6 lg:min-w-0 dark:border-white/8 dark:hover:bg-white/4",
		),
		h.Div(
			h.Class("w-45 shrink-0 lg:w-[20%] lg:shrink"),
			h.Div(
				h.Class("flex items-center gap-2"),
				h.Div(
					h.Class("min-w-0"),
					h.P(
						h.Class(
							"truncate font-medium dark:text-zinc-100",
						),
						g.Text(fmt.Sprintf("%s %s", invite.FirstName, invite.LastName)),
					),
					h.P(
						h.Class(
							"truncate text-xs text-zinc-500 dark:text-zinc-400",
						),
						g.Text(invite.Email),
					),
				),
			),
		),
		h.Div(
			h.Class(
				"w-45 shrink-0 text-xs text-zinc-500 lg:w-[20%] lg:shrink dark:text-zinc-400 capitalize",
			),
			g.Text(invite.FirstName),
		),
		h.Div(
			h.Class(
				"w-45 shrink-0 text-xs text-zinc-500 lg:w-[20%] lg:shrink dark:text-zinc-400 capitalize",
			),
			g.Text(invite.LastName),
		),
		h.Div(
			h.Class(
				"w-45 shrink-0 text-xs text-zinc-500 lg:w-[20%] lg:shrink dark:text-zinc-400 capitalize",
			),
			g.Text("todo"),
		),
		h.Div(
			h.Class(
				"w-45 shrink-0 text-xs text-zinc-500 lg:w-[20%] lg:shrink dark:text-zinc-400",
			),
			g.Text(invite.CreatedAt.Format("2006-01-02")),
		),
	)
}

func membershipItem(orgSlug string, member db.Membership, user db.User) g.Node {
	return h.Div(
		hx.Get(fmt.Sprintf(routes.HXOrgMembershipsUpdate, orgSlug, member.ID)),
		hx.Swap("none"),
		h.ID(fmt.Sprintf("item-%s", member.ID.String())),
		h.Class(
			"hover:bg-muted/40 flex min-w-225 items-center border-b px-4 py-3 text-sm last:border-b-0 sm:px-6 lg:min-w-0 dark:border-white/8 dark:hover:bg-white/4",
		),
		h.Div(
			h.Class("w-45 shrink-0 lg:w-[20%] lg:shrink"),
			h.Div(
				h.Class("flex items-center gap-2"),
				h.Div(
					h.Class("min-w-0"),
					h.P(
						h.Class(
							"truncate font-medium dark:text-zinc-100",
						),
						g.Text(
							fmt.Sprintf(
								"%s %s",
								member.FirstName,
								member.LastName,
							),
						),
					),
					h.P(
						h.Class(
							"truncate text-xs text-zinc-500 dark:text-zinc-400",
						),
						g.Text(user.Email),
					),
				),
			),
		),
		h.Div(
			h.Class(
				"w-45 shrink-0 text-xs text-zinc-500 lg:w-[20%] lg:shrink dark:text-zinc-400 capitalize",
			),
			g.Text(member.Status.String()),
		),
		h.Div(
			h.Class(
				"w-45 shrink-0 text-xs text-zinc-500 lg:w-[20%] lg:shrink dark:text-zinc-400 capitalize",
			),
			g.Text(member.Role.String()),
		),
		h.Div(
			h.Class(
				"w-45 shrink-0 text-xs text-zinc-500 lg:w-[20%] lg:shrink dark:text-zinc-400 capitalize",
			),
			g.Text("todo"),
		),
		h.Div(
			h.Class(
				"w-45 shrink-0 text-xs text-zinc-500 lg:w-[20%] lg:shrink dark:text-zinc-400",
			),
			g.Text(member.CreatedAt.Format("2006-01-02")),
		),
	)
}

func membershipUpdateFormModal(ident identity.IdentityCtx, member db.GetMembershipRow) g.Node {
	return components.ModalFragment(
		h.Div(
			h.Class("flex flex-col gap-2 text-center sm:text-left"),
			h.H2(
				h.Class("text-lg leading-none font-semibold"),
				g.Text("Update user"),
			),
			h.P(
				h.Class("text-muted-foreground text-sm"),
				g.Text("Update user here. Click save when you're done."),
			),
		),
		h.Form(

			x.Data(fmt.Sprintf(`{ role: %q }`, member.Membership.Role.String())),
			x.Cloak(),
			hx.Patch(
				fmt.Sprintf(
					routes.HXOrgMembershipsUpdate,
					ident.OrgSlug,
					member.Membership.ID,
				),
			),
			hx.Swap("outerHTML"),
			hx.Target(fmt.Sprintf("#item-%s", member.Membership.ID)),
			h.Class("space-y-4"),
			components.Input(components.InputParams{
				Label:        "First Name",
				Name:         "FirstName",
				Value:        member.Membership.FirstName,
				Placeholder:  "John",
				AutoFocus:    true,
				AutoComplete: "off",
			}),
			components.Input(components.InputParams{
				Label:        "Last Name",
				Name:         "LastName",
				Value:        member.Membership.LastName,
				Placeholder:  "Doe",
				AutoComplete: "off",
			}),
			components.Input(components.InputParams{
				Label:        "Email",
				Name:         "Email",
				Value:        member.User.Email,
				Placeholder:  "john.doe@gmail.com",
				AutoComplete: "off",
				Disabled:     true,
			}),
			components.Select(components.SelectParams{
				XModel:      "role",
				Label:       "Role",
				Name:        "role",
				Value:       member.Membership.Role.String(),
				Placeholder: "Select a role",
				Options: []components.SelectOption{
					{
						Value: "owner",
						Text:  "Owner",
					},
					{
						Value: "admin",
						Text:  "Admin",
					},
					{
						Value: "member",
						Text:  "Member",
					},
				},
			}),
			h.Template(
				x.If("role === 'member'"),
				h.Div(
					h.Class("grid grid-cols-2 gap-4"),
					h.Div(
						h.Class("border-muted flex flex-col overflow-hidden rounded-md border p-2"),
						components.Checkbox(components.CheckboxParams{
							Label:        "Membership Read",
							Name:         "Permissions",
							Value:        permissions.MembershipRead.String(),
							AutoComplete: "off",
						}),
						components.Checkbox(components.CheckboxParams{
							Label:        "Membership Create",
							Name:         "Permissions",
							Value:        permissions.MembershipCreate.String(),
							AutoComplete: "off",
						}),
						components.Checkbox(components.CheckboxParams{
							Label:        "Membership Update",
							Name:         "Permissions",
							Value:        permissions.MembershipUpdate.String(),
							AutoComplete: "off",
						}),
						components.Checkbox(components.CheckboxParams{
							Label:        "Membership Delete",
							Name:         "Permissions",
							Value:        permissions.MembershipDelete.String(),
							AutoComplete: "off",
						}),
					),
					h.Div(
						h.Class("border-muted flex flex-col overflow-hidden rounded-md border p-2"),
						components.Checkbox(components.CheckboxParams{
							Label:        "Read Channels",
							Name:         "Permissions",
							Value:        permissions.ChannelRead.String(),
							AutoComplete: "off",
						}),
						components.Checkbox(components.CheckboxParams{
							Label:        "Create Channels",
							Name:         "Permissions",
							Value:        permissions.ChannelCreate.String(),
							AutoComplete: "off",
						}),
						components.Checkbox(components.CheckboxParams{
							Label:        "Update Channels",
							Name:         "Permissions",
							Value:        permissions.ChannelUpdate.String(),
							AutoComplete: "off",
						}),
						components.Checkbox(components.CheckboxParams{
							Label:        "Delete Channels",
							Name:         "Permissions",
							Value:        permissions.ChannelDelete.String(),
							AutoComplete: "off",
						}),
					),
				),
			),
			h.Div(
				h.Class("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"),
				h.Button(
					h.Class(
						"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2",
					),
					h.Type("submit"),
					g.Text("Save changes"),
				),
			),
		),
	)
}

func invitationCreateFormModal(ident identity.IdentityCtx) g.Node {
	return components.ModalFragment(
		h.Div(
			h.Class("flex flex-col gap-2 text-center sm:text-left"),
			h.H2(
				h.Class("text-lg leading-none font-semibold"),
				g.Text("Invite New User"),
			),
			h.P(
				h.Class("text-muted-foreground text-sm"),
				g.Text("Invite new user here. Click save when you're done."),
			),
		),
		h.Form(
			x.Data(`{ role: null }`),
			x.Cloak(),
			hx.Post(fmt.Sprintf(routes.HXOrgInvitationsCreate, ident.OrgSlug)),
			hx.Swap("prepend"),
			hx.Target("#invitations_list"),
			h.Class("space-y-4"),
			components.Input(components.InputParams{
				Label:        "First Name",
				Name:         "FirstName",
				Placeholder:  "John",
				AutoFocus:    true,
				AutoComplete: "off",
			}),
			components.Input(components.InputParams{
				Label:        "Last Name",
				Name:         "LastName",
				Placeholder:  "Doe",
				AutoComplete: "off",
			}),
			components.Input(components.InputParams{
				Label:        "Email",
				Name:         "Email",
				Placeholder:  "john.doe@gmail.com",
				AutoComplete: "off",
			}),
			components.Select(components.SelectParams{
				Name:        "Role",
				Label:       "Role",
				Placeholder: "Select a role",
				XModel:      "role",
				Options: []components.SelectOption{
					{
						Value: "owner",
						Text:  "Owner",
					},
					{
						Value: "admin",
						Text:  "Admin",
					},
					{
						Value: "member",
						Text:  "Member",
					},
				},
			}),
			h.Template(
				x.If("role === 'member'"),
				h.Div(
					h.Class("grid grid-cols-2 gap-4"),
					h.Div(
						h.Class("border-muted flex flex-col overflow-hidden rounded-md border p-2"),
						components.Checkbox(components.CheckboxParams{
							Label:        "Membership Read",
							Name:         "Permissions",
							Value:        permissions.MembershipRead.String(),
							AutoComplete: "off",
						}),
						components.Checkbox(components.CheckboxParams{
							Label:        "Membership Create",
							Name:         "Permissions",
							Value:        permissions.MembershipCreate.String(),
							AutoComplete: "off",
						}),
						components.Checkbox(components.CheckboxParams{
							Label:        "Membership Update",
							Name:         "Permissions",
							Value:        permissions.MembershipUpdate.String(),
							AutoComplete: "off",
						}),
						components.Checkbox(components.CheckboxParams{
							Label:        "Membership Delete",
							Name:         "Permissions",
							Value:        permissions.MembershipDelete.String(),
							AutoComplete: "off",
						}),
					),
					h.Div(
						h.Class("border-muted flex flex-col overflow-hidden rounded-md border p-2"),
						components.Checkbox(components.CheckboxParams{
							Label:        "Read Channels",
							Name:         "Permissions",
							Value:        permissions.ChannelRead.String(),
							AutoComplete: "off",
						}),
						components.Checkbox(components.CheckboxParams{
							Label:        "Create Channels",
							Name:         "Permissions",
							Value:        permissions.ChannelCreate.String(),
							AutoComplete: "off",
						}),
						components.Checkbox(components.CheckboxParams{
							Label:        "Update Channels",
							Name:         "Permissions",
							Value:        permissions.ChannelUpdate.String(),
							AutoComplete: "off",
						}),
						components.Checkbox(components.CheckboxParams{
							Label:        "Delete Channels",
							Name:         "Permissions",
							Value:        permissions.ChannelDelete.String(),
							AutoComplete: "off",
						}),
					),
				),
			),
			h.Div(
				h.Class("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"),
				h.Button(
					h.Class(
						"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2",
					),
					h.Type("submit"),
					g.Text("Save changes"),
				),
			),
		),
	)
}

func invitationRegisterPage(token string, invite db.GetInvitationByTokenHashRow) g.Node {
	return components.BlankLayout(
		h.Div(
			h.Class("h-svh"),
			h.Div(
				h.Class("m-auto flex h-full w-full flex-col items-center justify-center gap-2"),
				h.H1(
					h.Class("text-[7rem] leading-tight font-bold"),
					g.Text("401"),
				),
				h.Span(
					h.Class("font-medium"),
					g.Text("Unauthorized Access"),
				),
				h.P(
					h.Class("text-muted-foreground text-center"),
					g.Textf(
						"Hi %s %s, you have been invided to <br> organization %s (@%s) by %s with email %s",
						invite.Invitation.FirstName,
						invite.Invitation.LastName,
						invite.Organization.Name,
						invite.Organization.Slug,
						"inviter email?",
						// invitation.Inviter.Email,
						invite.Invitation.Email,
					),
					g.Text("Please log in with the appropriate credentials"),
					h.Br(),
					g.Text("to access this resource."),
				),
				h.Div(
					h.Class("mt-6 flex gap-4"),
					h.Button(
						h.Class(
							"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2",
						),
						g.Text("Create account"),
					),
					h.Button(
						hx.Post(fmt.Sprintf(routes.HXInvitationsAccept, token)),
						hx.Swap("none"),
						h.Class(
							"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 border border-input bg-background shadow-xs hover:bg-accent hover:text-accent-foreground h-9 px-4 py-2",
						),
						g.Text("Already have an account"),
					),
				),
			),
		),
	)
}

func invitationAcceptPage(token string, invite db.GetInvitationByTokenHashRow) g.Node {
	return components.BlankLayout(
		h.Div(
			h.Class("h-svh"),
			h.Div(
				h.Class("m-auto flex h-full w-full flex-col items-center justify-center gap-2"),
				h.H1(
					h.Class("text-[7rem] leading-tight font-bold"),
					g.Text("401"),
				),
				h.Span(
					h.Class("font-medium"),
					g.Text("Unauthorized Access"),
				),
				h.P(
					h.Class("text-muted-foreground text-center"),
					g.Textf(
						"Hi %s %s, you have been invided to <br> organization %s (@%s) by %s with email %s",
						invite.Invitation.FirstName,
						invite.Invitation.LastName,
						invite.Organization.Name,
						invite.Organization.Slug,
						"inviter email",
						invite.Invitation.Email,
					),
					g.Text("Please log in with the appropriate credentials"),
					h.Br(),
					g.Text("to access this resource."),
				),
				h.Div(
					h.Class("mt-6 flex gap-4"),
					h.Button(
						hx.Post(fmt.Sprintf(routes.HXInvitationsDecline, token)),
						hx.Swap("none"),
						h.Class(
							"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2",
						),
						g.Text("Accept invitation"),
					),
					h.Button(
						hx.Post(fmt.Sprintf(routes.HXInvitationsAccept, token)),
						hx.Swap("none"),
						h.Class(
							"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 border border-input bg-background shadow-xs hover:bg-accent hover:text-accent-foreground h-9 px-4 py-2",
						),
						g.Text("Decline invitation"),
					),
				),
			),
		),
	)
}

func invitationErrorPage() g.Node {
	return components.BlankLayout(
		g.Text("Invitaion is not available anymore"),
	)
}
