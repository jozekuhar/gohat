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
										g.Text("16"),
									),
								),
								h.Button(
									hx.Get(
										fmt.Sprintf(
											routes.HXOrgMembershipsCreate,
											identity.OrganizationSlug,
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
								h.Div(
									h.Class("w-[220px] shrink-0 lg:w-[20%] lg:shrink"),
									g.Text("Teams"),
								),
								h.Div(
									h.Class("w-[110px] shrink-0 lg:w-[14%] lg:shrink"),
									g.Text("Joined"),
								),
							),
						),
						h.Div(
							g.Map(data.Memberhips, func(member tenant.Membership) g.Node {
								return h.Div(
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
													h.Img(
														h.Class("aspect-square h-full w-full"),
														h.Alt("Ava Reed"),
														h.Src("/avatars/avatar-3.png"),
													),
												),
												h.Span(
													h.Class(
														"border-background absolute -right-0.5 -bottom-0.5 size-2.5 rounded-full border-2 dark:border-[#141414]",
													),
													h.Style("background-color:#16a34a"),
													h.Span(
														h.Class("sr-only"),
														g.Text("online"),
													),
												),
											),
											h.Div(
												h.Class("min-w-0"),
												h.P(
													h.Class(
														"truncate font-medium dark:text-zinc-100",
													),
													g.Text("Ava Reed"),
												),
												h.P(
													h.Class(
														"truncate text-xs text-zinc-500 dark:text-zinc-400",
													),
													g.Text("todo: email"),
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
												g.Text(member.Status.String()),
											),
										),
									),
									h.Div(
										h.Class(
											"w-[100px] shrink-0 text-xs text-zinc-500 lg:w-[14%] lg:shrink dark:text-zinc-400 capitalize",
										),
										g.Text(member.Role.String()),
									),
									h.Div(
										h.Class(
											"w-[110px] shrink-0 text-xs text-zinc-500 lg:w-[14%] lg:shrink dark:text-zinc-400",
										),
										g.Text("todo: joined"),
									),
								)
							}),
						),
					),
				),
			),
		),
	)
}
