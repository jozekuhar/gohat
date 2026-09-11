package view

import (
	"fmt"
	"time"
	"uuid"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/authz"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/tenant"

	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

type Memberships struct {
	layout *Layout
}

func NewMemberships() *Memberships {
	return &Memberships{
		layout: NewLayout(),
	}
}

func (v *Memberships) MembershipsPage(
	identity authz.Identity,
	data tenant.MembershipsData,
) g.Node {
	return v.layout.app(
		identity,
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
											identity.OrgSlug,
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
										identity.OrgSlug,
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
								"min-w-[900px] border-b px-4 py-2 text-sm text-zinc-500 sm:px-6 lg:min-w-0 dark:border-white/8 dark:bg-[#141414]/95 dark:text-zinc-400",
							),
							h.Div(
								h.Class("flex items-center"),
								h.Div(
									h.Class("w-[320px] shrink-0 lg:w-[36%] lg:shrink"),
									g.Text("Member"),
								),
								h.Div(
									h.Class("w-[150px] shrink-0 lg:w-[16%] lg:shrink"),
									g.Text("Availability"),
								),
								h.Div(
									h.Class("w-[100px] shrink-0 lg:w-[14%] lg:shrink"),
									g.Text("Role"),
								),
								// h.Div(
								// 	h.Class("w-[220px] shrink-0 lg:w-[20%] lg:shrink"),
								// 	g.Text("Teams"),
								// ),
								h.Div(
									h.Class("w-[330px] shrink-0 lg:w-[14%] lg:shrink"),
									g.Text("Joined"),
								),
							),
						),
						h.Div(
							g.Map(data.Invitations, func(invite db.Invitation) g.Node {
								return v.membershipItem(
									identity.OrgSlug,
									uuid.UUID{},
									invite.FirstName,
									invite.LastName,
									invite.Email,
									"Invited",
									invite.Role.String(),
									invite.CreatedAt,
								)
							}),
							g.Map(data.Memberhips, func(m db.ListMembershipsRow) g.Node {
								return v.membershipItem(
									identity.OrgSlug,
									m.Membership.ID,
									m.Membership.FirstName,
									m.Membership.LastName,
									m.User.Email,
									m.Membership.Status.String(),
									m.Membership.Role.String(),
									m.Membership.CreatedAt,
								)
							}),
						),
					),
				),
			),
		),
	)
}

func (v *Memberships) membershipItem(
	orgSlug string,
	id uuid.UUID,
	firstName, lastName, email, status, role string,
	createdAt time.Time,
) g.Node {
	return h.Div(
		hx.Get(fmt.Sprintf(routes.HXOrgMembershipsUpdate, orgSlug, id)),
		hx.Swap("none"),
		h.Class(
			"hover:bg-muted/40 flex min-w-[900px] items-center border-b px-4 py-3 text-sm transition-colors last:border-b-0 sm:px-6 lg:min-w-0 dark:border-white/8 dark:hover:bg-white/4",
		),
		h.Div(
			h.Class("w-[320px] shrink-0 lg:w-[36%] lg:shrink"),
			h.Div(
				h.Class("flex items-center gap-2"),
				h.Div(
					h.Class("relative"),
					h.Span(
						h.Class(
							"relative flex overflow-hidden rounded-full size-8 shrink-0",
						),
						g.Raw(
							`<svg viewBox="-42 0 512 512.001" xmlns="http://www.w3.org/2000/svg" id="fi_1077063"><path d="m210.351562 246.632812c33.882813 0 63.21875-12.152343 87.195313-36.128906 23.96875-23.972656 36.125-53.304687 36.125-87.191406 0-33.875-12.152344-63.210938-36.128906-87.191406-23.976563-23.96875-53.3125-36.121094-87.191407-36.121094-33.886718 0-63.21875 12.152344-87.191406 36.125s-36.128906 53.308594-36.128906 87.1875c0 33.886719 12.15625 63.222656 36.128906 87.195312 23.980469 23.96875 53.316406 36.125 87.191406 36.125zm-65.972656-189.292968c18.394532-18.394532 39.972656-27.335938 65.972656-27.335938 25.996094 0 47.578126 8.941406 65.976563 27.335938 18.394531 18.398437 27.339844 39.980468 27.339844 65.972656 0 26-8.945313 47.578125-27.339844 65.976562-18.398437 18.398438-39.980469 27.339844-65.976563 27.339844-25.992187 0-47.570312-8.945312-65.972656-27.339844-18.398437-18.394531-27.34375-39.976562-27.34375-65.976562 0-25.992188 8.945313-47.574219 27.34375-65.972656zm0 0"></path><path d="m426.128906 393.703125c-.691406-9.976563-2.089844-20.859375-4.148437-32.351563-2.078125-11.578124-4.753907-22.523437-7.957031-32.527343-3.3125-10.339844-7.808594-20.550781-13.375-30.335938-5.769532-10.15625-12.550782-19-20.160157-26.277343-7.957031-7.613282-17.699219-13.734376-28.964843-18.199219-11.226563-4.441407-23.667969-6.691407-36.976563-6.691407-5.226563 0-10.28125 2.144532-20.042969 8.5-6.007812 3.917969-13.035156 8.449219-20.878906 13.460938-6.707031 4.273438-15.792969 8.277344-27.015625 11.902344-10.949219 3.542968-22.066406 5.339844-33.042969 5.339844-10.96875 0-22.085937-1.796876-33.042968-5.339844-11.210938-3.621094-20.300782-7.625-26.996094-11.898438-7.769532-4.964844-14.800782-9.496094-20.898438-13.46875-9.753906-6.355468-14.808594-8.5-20.035156-8.5-13.3125 0-25.75 2.253906-36.972656 6.699219-11.257813 4.457031-21.003906 10.578125-28.96875 18.199219-7.609375 7.28125-14.390625 16.121094-20.15625 26.273437-5.558594 9.785157-10.058594 19.992188-13.371094 30.339844-3.199219 10.003906-5.875 20.945313-7.953125 32.523437-2.0625 11.476563-3.457031 22.363282-4.148437 32.363282-.679688 9.777344-1.023438 19.953125-1.023438 30.234375 0 26.726562 8.496094 48.363281 25.25 64.320312 16.546875 15.746094 38.4375 23.730469 65.066406 23.730469h246.53125c26.621094 0 48.511719-7.984375 65.0625-23.730469 16.757813-15.945312 25.253906-37.589843 25.253906-64.324219-.003906-10.316406-.351562-20.492187-1.035156-30.242187zm-44.90625 72.828125c-10.933594 10.40625-25.449218 15.464844-44.378906 15.464844h-246.527344c-18.933594 0-33.449218-5.058594-44.378906-15.460938-10.722656-10.207031-15.933594-24.140625-15.933594-42.585937 0-9.59375.316406-19.066407.949219-28.160157.617187-8.921874 1.878906-18.722656 3.75-29.136718 1.847656-10.285156 4.199219-19.9375 6.996094-28.675782 2.683593-8.378906 6.34375-16.675781 10.882812-24.667968 4.332031-7.617188 9.316407-14.152344 14.816407-19.417969 5.144531-4.925781 11.628906-8.957031 19.269531-11.980469 7.066406-2.796875 15.007812-4.328125 23.628906-4.558594 1.050781.558594 2.921875 1.625 5.953125 3.601563 6.167969 4.019531 13.277344 8.605469 21.136719 13.625 8.859375 5.648437 20.273437 10.75 33.910156 15.152344 13.941406 4.507812 28.160156 6.796875 42.273437 6.796875 14.113282 0 28.335938-2.289063 42.269532-6.792969 13.648437-4.410156 25.058594-9.507813 33.929687-15.164063 8.042969-5.140624 14.953125-9.59375 21.121094-13.617187 3.03125-1.972656 4.902344-3.042969 5.953125-3.601563 8.625.230469 16.566406 1.761719 23.636719 4.558594 7.636719 3.023438 14.121093 7.058594 19.265625 11.980469 5.5 5.261719 10.484375 11.796875 14.816406 19.421875 4.542969 7.988281 8.207031 16.289062 10.886719 24.660156 2.800781 8.75 5.15625 18.398438 7 28.675782 1.867187 10.433593 3.132812 20.238281 3.75 29.144531v.007812c.636719 9.058594.957031 18.527344.960937 28.148438-.003906 18.449219-5.214844 32.378906-15.9375 42.582031zm0 0"></path></svg>`,
						),
						// h.Img(
						// 	h.Class("aspect-square h-full w-full"),
						// 	h.Alt("Ava Reed"),
						// 	h.Src("/avatars/avatar-3.png"),
						// ),
					),
				),
				h.Div(
					h.Class("min-w-0"),
					h.P(
						h.Class(
							"truncate font-medium dark:text-zinc-100",
						),
						g.Text(fmt.Sprintf("%s %s", firstName, lastName)),
					),
					h.P(
						h.Class(
							"truncate text-xs text-zinc-500 dark:text-zinc-400",
						),
						g.Text(email),
					),
				),
			),
		),
		h.Div(
			h.Class("w-[150px] shrink-0 lg:w-[16%] lg:shrink"),
			h.Div(
				h.Class("flex items-center gap-2"),
				h.SVG(
					g.Attr("xmlns", "http://www.w3.org/2000/svg"),
					h.Width("24"),
					h.Height("24"),
					g.Attr("viewBox", "0 0 24 24"),
					g.Attr("fill", "none"),
					g.Attr("stroke", "currentColor"),
					g.Attr("stroke-width", "2"),
					g.Attr("stroke-linecap", "round"),
					g.Attr("stroke-linejoin", "round"),
					h.Class(
						"lucide lucide-circle-check size-4 shrink-0 text-green-500",
					),
					h.Aria("hidden", "true"),
					g.El("circle",
						g.Attr("cx", "12"),
						g.Attr("cy", "12"),
						g.Attr("r", "10"),
					),
					g.El("path",
						g.Attr("d", "m9 12 2 2 4-4"),
					),
				),
				h.Span(
					h.Class(
						"truncate text-xs font-medium dark:text-zinc-100 capitalize",
					),
					g.Text(status),
				),
			),
		),
		h.Div(
			h.Class(
				"w-25 shrink-0 text-xs text-zinc-500 lg:w-[14%] lg:shrink dark:text-zinc-400 capitalize",
			),
			g.Text(role),
		),
		h.Div(
			h.Class(
				"w-[330px] shrink-0 text-xs text-zinc-500 lg:w-[14%] lg:shrink dark:text-zinc-400",
			),
			g.Text(createdAt.Format("2006-01-02")),
		),
	)
}

func (v *Memberships) MembershipUpdateFormModal(
	identity authz.Identity,
	membership db.GetMembershipRow,
) g.Node {
	return v.layout.modal.fragment(
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
			hx.Patch(
				fmt.Sprintf(
					routes.HXOrgMembershipsUpdate,
					identity.OrgSlug,
					membership.Membership.ID,
				),
			),
			hx.Swap("none"),
			h.Class("space-y-4"),
			h.Div(
				h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
				h.Label(
					h.Class(
						"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
					),
					g.Text("First Name"),
				),
				h.Input(
					h.Class(
						"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm col-span-4",
					),
					h.AutoFocus(),
					h.Placeholder("John"),
					h.AutoComplete("off"),
					h.Value(membership.Membership.FirstName),
					h.Name("FirstName"),
				),
			),
			h.Div(
				h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
				h.Label(
					h.Class(
						"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
					),
					h.For("_r_4t_-form-item"),
					g.Text("Last Name"),
				),
				h.Input(
					h.Class(
						"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm col-span-4",
					),
					h.Placeholder("Doe"),
					h.AutoComplete("off"),
					h.Value(membership.Membership.LastName),
					h.Name("LastName"),
				),
			),
			h.Div(
				h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
				h.Label(
					h.Class(
						"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
					),
					g.Text("Email"),
				),
				h.Input(
					h.Class(
						"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm col-span-4",
					),
					h.Value(membership.User.Email),
					h.Name("Email"),
					h.Disabled(),
				),
			),
			h.Div(
				h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
				h.Label(
					h.Class(
						"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
					),
					g.Text("Phone Number"),
				),
				h.Input(
					h.Class(
						"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm col-span-4",
					),
					h.Placeholder("+123456789"),
					h.Value(""),
					h.Name("PhoneNumber"),
				),
			),
			// TODO(jozekuhar): state?
			h.Div(
				h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
				h.Label(
					h.Class(
						"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
					),
					g.Text("Role"),
				),
				h.Button(
					h.Type("button"),
					h.Aria("autocomplete", "none"),
					h.Class(
						"border-input ring-offset-background placeholder:text-muted-foreground focus:ring-ring flex h-9 w-full items-center justify-between rounded-md border bg-transparent px-3 py-2 text-sm whitespace-nowrap shadow-xs focus:ring-1 focus:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1 col-span-4",
					),
					h.Span(
						h.Style("pointer-events: none;"),
						g.Text("Select a role"),
					),
					g.Raw(
						`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-chevron-down h-4 w-4 opacity-50" aria-hidden="true"><path d="m6 9 6 6 6-6"></path></svg>`,
					),
				),
				h.Select(
					h.Style(
						"position: absolute; border: 0px; width: 1px; height: 1px; padding: 0px; margin: -1px; overflow: hidden; clip: rect(0px, 0px, 0px, 0px); white-space: nowrap; overflow-wrap: normal;",
					),
					h.Name("Role"),
					h.Option(
						h.Value("owner"),
						g.Text("Owner"),
					),
					h.Option(
						h.Value("admin"),
						g.Text("Admin"),
					),
					h.Option(
						h.Value("member"),
						g.Text("Member"),
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

func (v *Memberships) InvitationCreateFormModal(identity authz.Identity) g.Node {
	return v.layout.modal.fragment(
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
			hx.Post(fmt.Sprintf(routes.HXOrgInvitationsCreate, identity.OrgSlug)),
			hx.Swap("none"),
			h.Class("space-y-4"),
			h.Div(
				h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
				h.Label(
					h.Class(
						"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
					),
					g.Text("First Name"),
				),
				h.Input(
					h.Class(
						"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm col-span-4",
					),
					h.AutoFocus(),
					h.Placeholder("John"),
					h.AutoComplete("off"),
					h.Value(""),
					h.Name("FirstName"),
				),
			),
			h.Div(
				h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
				h.Label(
					h.Class(
						"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
					),
					h.For("_r_4t_-form-item"),
					g.Text("Last Name"),
				),
				h.Input(
					h.Class(
						"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm col-span-4",
					),
					h.Placeholder("Doe"),
					h.AutoComplete("off"),
					h.Value(""),
					h.Name("LastName"),
				),
			),
			h.Div(
				h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
				h.Label(
					h.Class(
						"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
					),
					g.Text("Email"),
				),
				h.Input(
					h.Class(
						"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm col-span-4",
					),
					h.Placeholder("john.doe@gmail.com"),
					h.Value(""),
					h.Name("Email"),
				),
			),
			h.Div(
				h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
				h.Label(
					h.Class(
						"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
					),
					g.Text("Phone Number"),
				),
				h.Input(
					h.Class(
						"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm col-span-4",
					),
					h.Placeholder("+123456789"),
					h.Value(""),
					h.Name("PhoneNumber"),
				),
			),
			// TODO(jozekuhar): state?
			h.Div(
				h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
				h.Label(
					h.Class(
						"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
					),
					g.Text("Role"),
				),
				h.Button(
					h.Type("button"),
					h.Aria("autocomplete", "none"),
					h.Class(
						"border-input ring-offset-background placeholder:text-muted-foreground focus:ring-ring flex h-9 w-full items-center justify-between rounded-md border bg-transparent px-3 py-2 text-sm whitespace-nowrap shadow-xs focus:ring-1 focus:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1 col-span-4",
					),
					h.Span(
						h.Style("pointer-events: none;"),
						g.Text("Select a role"),
					),
					g.Raw(
						`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-chevron-down h-4 w-4 opacity-50" aria-hidden="true"><path d="m6 9 6 6 6-6"></path></svg>`,
					),
				),
				h.Select(
					h.Style(
						"position: absolute; border: 0px; width: 1px; height: 1px; padding: 0px; margin: -1px; overflow: hidden; clip: rect(0px, 0px, 0px, 0px); white-space: nowrap; overflow-wrap: normal;",
					),
					h.Name("Role"),
					h.Option(
						h.Value("owner"),
						g.Text("Owner"),
					),
					h.Option(
						h.Value("admin"),
						g.Text("Admin"),
					),
					h.Option(
						h.Value("member"),
						g.Text("Member"),
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

func (v *Memberships) InvitationRegisterPage(
	token string,
	invitation db.GetInvitationByTokenHashRow,
) g.Node {
	return v.layout.blank(
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
						invitation.Invitation.FirstName,
						invitation.Invitation.LastName,
						invitation.Organization.Name,
						invitation.Organization.Slug,
						"inviter email?",
						// invitation.Inviter.Email,
						invitation.Invitation.Email,
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

func (v *Memberships) InvitationAcceptPage(
	token string,
	invitation db.GetInvitationByTokenHashRow,
) g.Node {
	return v.layout.blank(
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
						invitation.Invitation.FirstName,
						invitation.Invitation.LastName,
						invitation.Organization.Name,
						invitation.Organization.Slug,
						"inviter email",
						invitation.Invitation.Email,
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

func (v *Memberships) InvitationErrorPage() g.Node {
	return v.layout.blank(
		g.Text("Invitaion is not available anymore"),
	)
}
