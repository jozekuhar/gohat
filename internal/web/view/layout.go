package view

import (
	"fmt"

	"mimokocke/internal/shared/authz"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/tenant"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

const (
	idSidebarPopoverItems = "sidebar_popover_items"
)

type Layout struct {
	toast *Toast
	modal *modal
}

func NewLayout() *Layout {
	return &Layout{
		toast: NewToast(),
		modal: newModal(),
	}
}

func (v *Layout) blank(children ...g.Node) g.Node {
	return v.base(
		g.Map(children, func(node g.Node) g.Node {
			return node
		}),
	)
}

func (v *Layout) app(identity authz.Identity, children ...g.Node) g.Node {
	return v.base(
		h.Div(
			x.Data(`{ sidebarOpen: false }`),
			x.Cloak(),
			h.Class("flex flex-1 flex-col"),
			h.Div(
				h.Class("flex min-h-svh w-full"),
				v.sidebar(identity),
				v.content(children...),
			),
		),
	)
}

func (v *Layout) base(children ...g.Node) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:       "",
		Description: "",
		Language:    "en",
		Head: g.Group{
			favicons(),
			g.Map(getCSSPaths("view"), func(path string) g.Node {
				return h.Link(h.Rel("stylesheet"), h.Href(path))
			}),
			// TODO(jozekuhar): fix fonts (all because there are lot of fonts in view.css)
			g.Raw(`
				<link rel="preconnect" href="https://fonts.googleapis.com">
				<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
				<link href="https://fonts.googleapis.com/css2?family=Geist:ital,wght@0,100..900;1,100..900&family=Plus+Jakarta+Sans:ital,wght@0,200..800;1,200..800&display=swap" rel="stylesheet">
			`),
		},
		Body: g.Group{
			h.Class("antialiased"),
			g.Map(children, func(node g.Node) g.Node {
				return node
			}),
			v.modal.container(),
			v.toast.container(),
			h.Script(h.Type("module"), h.Src(getAssetPath("view"))),
		},
		HTMLAttrs: g.Group{
			h.Class("font-sans"),
		},
	})
}

func (v *Layout) sidebar(identity authz.Identity) g.Node {
	return h.Div(
		h.Class("relative"),
		h.Div(
			h.Class("[--sidebar-width:16rem] text-sidebar-foreground hidden md:block"),
			h.Div(
				c.Classes{
					"relative h-svh w-(--sidebar-width) bg-transparent": true,
					"transition-[width] duration-200 ease-linear":       true,
				},
			),
			h.Div(
				c.Classes{
					"fixed inset-y-0 /z-10 hidden h-svh w-(--sidebar-width) md:flex left-0": true,
					"transition-[left,right,width] duration-200 ease-linear":                true,
				},
				h.Div(
					h.Class("bg-sidebar flex flex-col h-full w-full"),
					v.sidebarHeader(),
					g.If(identity.OrganizationSlug == "", v.sidebarContentEmpty()),
					g.If(identity.OrganizationSlug != "", v.sidebarContent(identity)),
					v.sidebarFooter(),
				),
			),
		),
	)
}

func (v *Layout) sidebarHeader() g.Node {
	return h.Div(
		x.Data(`{ menuOpen: false }`),
		x.Cloak(),
		c.Classes{
			"relative":                true,
			"flex flex-col gap-2 p-2": true,
		},
		h.Ul(
			h.Class("flex w-full min-w-0 flex-col gap-1"),
			h.Li(
				h.Class("relative"),
				h.Button(
					hx.Get(routes.HXSidebarOrganizations),
					hx.Swap("none"),
					x.On(
						"htmx:finally:request",
						`if (event.detail.ctx.status === 'swapped') menuOpen = true`,
					),
					x.Bind(
						"class",
						`menuOpen && "bg-sidebar-accent text-sidebar-accent-foreground"`,
					),
					c.Classes{
						"flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left outline-hidden hover:bg-sidebar-accent hover:text-sidebar-accent-foreground h-12 text-sm ring-sidebar-ring/50 focus-visible:ring-1": true,
						"active:bg-sidebar-accent active:text-sidebar-accent-foreground": true,
						"[&>svg]:size-4 [&>svg]:shrink-0":                                true,
					},
					h.Div(
						h.Class(
							"border-muted-foreground/25 flex aspect-square size-8 items-center justify-center rounded-lg border bg-transparent",
						),
						h.Img(
							h.Alt("Artifact"),
							h.Loading("lazy"),
							h.Width("18"),
							h.Height("18"),
							g.Attr("decoding", "async"),
							h.Data("nimg", "1"),
							h.Class("size-4 dark:hidden"),
							h.Style("color:transparent"),
							h.Src(
								"https://artifact-nextjs-template.vercel.app/images/logo.svg?dpl=dpl_Gyjss8KqnbavQHpRArQ6yPL3hPJ1",
							),
						),
					),
					h.Div(
						h.Class("grid flex-1 text-left text-xs leading-tight"),
						h.Span(
							h.Class("truncate font-semibold"),
							g.Text("Artifact Admin Kit"),
						),
						h.Span(
							h.Class("truncate text-xs"),
							g.Text("Nextjs + shadcn/ui"),
						),
					),
					g.Raw(
						`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-chevrons-up-down ml-auto" aria-hidden="true"><path d="m7 15 5 5 5-5"></path><path d="m7 9 5-5 5 5"></path></svg>`,
					),
				),
			),
		),
		v.sidebarHeaderPopover(),
	)
}

func (v *Layout) sidebarHeaderPopover() g.Node {
	return h.Div(
		x.Show("menuOpen"),
		x.On("click.outside", "menuOpen = false"),
		c.Classes{
			"absolute left-full top-0": true,
			"bg-popover text-popover-foreground z-50 overflow-hidden border p-1 shadow-md min-w-56 rounded-lg": true,
		},
		h.Div(
			h.Class("px-2 py-1.5 font-semibold text-muted-foreground text-xs"),
			g.Text("Organizations"),
		),
		// Content
		h.Div(
			h.ID(idSidebarPopoverItems),
		),
		h.Div(
			h.Class("bg-muted -mx-1 my-1 h-px"),
		),
		h.Button(
			hx.Get(routes.HXSidebarOrganizationsCreate),
			hx.Swap("none"),
			c.Classes{
				"w-full relative flex cursor-default items-center rounded-sm text-sm outline-hidden transition-colors select-none gap-2 p-2": true,
				"hocus:bg-accent hocus:text-accent-foreground": true,
				"[&>svg]:size-4 [&>svg]:shrink-0":              true,
			},
			h.Div(
				h.Class(
					"bg-background flex size-6 items-center justify-center rounded-md border",
				),
				g.Raw(
					`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-plus size-4" aria-hidden="true"><path d="M5 12h14"></path><path d="M12 5v14"></path></svg>`,
				),
			),
			h.Div(
				h.Class("text-muted-foreground font-medium"),
				g.Text("Add team"),
			),
		),
	)
}

func (v *Layout) SidebarHeaderPopoverOrganizationsPartial(
	organizations []tenant.Organization,
) g.Node {
	return g.El(
		"hx-partial",
		hx.Target(fmt.Sprintf("#%s", idSidebarPopoverItems)),
		hx.Swap("outerHTML"),
		h.Div(
			h.ID(idSidebarPopoverItems),
			g.Map(organizations, func(organization tenant.Organization) g.Node {
				return h.A(
					c.Classes{
						"relative flex cursor-default items-center rounded-sm text-sm outline-hidden transition-colors select-none gap-2 p-2 text-balance": true,
						"hocus:bg-accent hocus:text-accent-foreground": true,
						"[&>svg]:size-4 [&>svg]:shrink-0":              true,
					},
					h.Href(fmt.Sprintf(routes.OrgDashboard, organization.Slug)),
					h.Div(
						h.Class(
							"flex size-6 items-center justify-center rounded-sm border",
						),
						h.Img(
							h.Alt("Artifact"),
							h.Loading("lazy"),
							h.Width("18"),
							h.Height("18"),
							g.Attr("decoding", "async"),
							h.Data("nimg", "1"),
							h.Class("size-4 shrink-0 dark:hidden"),
							h.Src(organization.Slug),
							h.Style("color: transparent;"),
						),
					),
					g.Text(organization.Name),
					h.Span(
						h.Class("ml-auto text-xs tracking-widest opacity-60"),
						g.Text("⌘1"),
					),
				)
			}),
		),
	)
}

func (v *Layout) OrganizationCreateFormModal() g.Node {
	return v.modal.fragment(
		h.Div(
			h.Class("flex flex-col gap-2 text-center sm:text-left"),
			h.H2(
				h.Class("text-lg leading-none font-semibold"),
				g.Text("Create New Organization"),
			),
			h.P(
				h.Class("text-muted-foreground text-sm"),
				g.Text("Create new organization here. Click save when you're done."),
			),
		),
		h.Form(
			hx.Post(routes.HXSidebarOrganizationsCreate),
			hx.Swap("none"),
			h.Class("space-y-4"),
			h.Div(
				h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
				h.Label(
					h.Class(
						"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
					),
					g.Text("Name"),
				),
				h.Input(
					g.Attr("autofocus"),
					h.Class(
						"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden md:text-sm col-span-4",
					),
					h.Placeholder("Example Company"),
					h.AutoComplete("off"),
					h.Value(""),
					h.Name("Name"),
				),
			),
			h.Div(
				h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
				h.Label(
					h.Class(
						"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
					),
					g.Text("Slug"),
				),
				h.Input(
					h.Class(
						"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm col-span-4",
					),
					h.Placeholder("example-company"),
					h.AutoComplete("off"),
					h.Value(""),
					h.Name("Slug"),
				),
			),
			h.Div(
				h.Class("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"),
				h.Button(
					h.Class(
						"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2",
					),
					g.Text("Save"),
				),
			),
		),
	)
}

func (v *Layout) sidebarContentEmpty() g.Node {
	return h.Div(
		h.Class("flex min-h-0 flex-1 flex-col gap-2 overflow-auto"),
	)
}

func (v *Layout) sidebarContent(identity authz.Identity) g.Node {
	return h.Div(
		h.Class("flex min-h-0 flex-1 flex-col gap-2 overflow-auto"),
		// Group
		h.Div(
			h.Class("relative flex w-full min-w-0 flex-col p-2"),
			h.Div(
				c.Classes{
					"text-sidebar-foreground/70 ring-sidebar-ring flex h-8 shrink-0 items-center rounded-md px-2 text-xs font-medium outline-hidden ease-linear focus-visible:ring-2": true,
					"transition-[margin,opacity] duration-200": true,
					"[&>svg]:size-4 [&>svg]:shrink-0":          true,
				},
				g.Text("Order Management"),
			),
			h.Ul(
				h.Class("flex w-full min-w-0 flex-col gap-1"),
				h.Li(
					h.Class("relative"),
					h.A(
						c.Classes{
							"flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left outline-hidden ring-sidebar-ring h-8 text-sm": true,
							"transition-[width,height,padding] ": true,
							"focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground": true,
							"[&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0 ":                                                                                    true,
						},
						h.Href(fmt.Sprintf(routes.OrgDashboard, identity.OrganizationSlug)),
						g.Raw(
							`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tabler-icon tabler-icon-layout-dashboard "><path d="M5 4h4a1 1 0 0 1 1 1v6a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1v-6a1 1 0 0 1 1 -1"></path><path d="M5 16h4a1 1 0 0 1 1 1v2a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1v-2a1 1 0 0 1 1 -1"></path><path d="M15 12h4a1 1 0 0 1 1 1v6a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1v-6a1 1 0 0 1 1 -1"></path><path d="M15 4h4a1 1 0 0 1 1 1v2a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1v-2a1 1 0 0 1 1 -1"></path></svg>`,
						),
						h.Span(
							g.Text("Dashboard"),
						),
					),
				),
				h.Li(
					h.Class("relative"),
					h.A(
						c.Classes{
							"flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left outline-hidden ring-sidebar-ring h-8 text-sm": true,
							"transition-[width,height,padding]": true,
							"focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground": true,
							"[&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0":                                                                                     true,
						},
						h.Href(fmt.Sprintf(routes.OrgMemberships, identity.OrganizationSlug)),
						g.Raw(
							`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tabler-icon tabler-icon-users "><path d="M5 7a4 4 0 1 0 8 0a4 4 0 1 0 -8 0"></path><path d="M3 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path><path d="M21 21v-2a4 4 0 0 0 -3 -3.85"></path></svg>`,
						),
						h.Span(
							g.Text("Users"),
						),
					),
				),
			),
		),
	)
}

func (v *Layout) sidebarFooter() g.Node {
	type menuItem struct {
		href string
		svg  string
		text string
	}

	return h.Div(
		x.Data(`{ menuOpen: false }`),
		x.Cloak(),
		c.Classes{
			"relative":                true,
			"flex flex-col gap-2 p-2": true,
		},
		// Button
		h.Ul(
			h.Class("flex w-full min-w-0 flex-col gap-1"),
			h.Li(
				h.Class("relative"),
				h.Button(
					x.On("click", "menuOpen = true"),
					x.Bind(
						"class",
						`menuOpen && "bg-sidebar-accent text-sidebar-accent-foreground"`,
					),
					c.Classes{
						"flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left outline-hidden ring-sidebar-ring h-12 text-sm": true,
						"focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground ":                                  true,
						"hover:bg-sidebar-accent hover:text-sidebar-accent-foreground":                                                          true,
						"[&>svg]:size-4 [&>svg]:shrink-0": true,
					},
					h.Span(
						h.Class("relative flex shrink-0 overflow-hidden h-8 w-8 rounded-lg"),
						h.Img(
							h.Class("aspect-square h-full w-full"),
							h.Alt(""),
							h.Src(
								"https://artifact-nextjs-template.vercel.app/avatars/cruip-avatar.svg",
							),
						),
					),
					h.Div(
						h.Class("grid flex-1 text-left text-sm leading-tight"),
						h.Span(
							h.Class("truncate font-semibold"),
							g.Text("Cruip"),
						),
						h.Span(
							h.Class("truncate text-xs"),
							g.Text("cruip@example.com"),
						),
					),
					g.Raw(
						`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-chevrons-up-down ml-auto size-4" aria-hidden="true"><path d="m7 15 5 5 5-5"></path><path d="m7 9 5-5 5 5"></path></svg>`,
					),
				),
			),
		),
		// Popover
		h.Div(
			x.Show("menuOpen"),
			x.On("click.outside", "menuOpen = false"),
			c.Classes{
				"absolute left-full bottom-0": true,
				"bg-popover text-popover-foreground z-50 overflow-hidden border p-1 shadow-md min-w-56 rounded-lg": true,
			},
			h.Div(
				h.Class("text-sm p-0 font-normal"),
				h.Div(
					h.Class("flex items-center gap-2 px-1 py-1.5 text-left text-sm"),
					h.Span(
						h.Class("relative flex shrink-0 overflow-hidden h-8 w-8 rounded-lg"),
						h.Img(
							h.Class("aspect-square h-full w-full"),
							h.Alt("Cruip"),
							h.Src(
								"https://artifact-nextjs-template.vercel.app/avatars/cruip-avatar.svg",
							),
						),
					),
					h.Div(
						h.Class("grid flex-1 text-left text-sm leading-tight"),
						h.Span(
							h.Class("truncate font-semibold"),
							g.Text("Cruip"),
						),
						h.Span(
							h.Class("truncate text-xs"),
							g.Text("cruip@example.com"),
						),
					),
				),
			),
			h.Div(
				h.Class("bg-muted -mx-1 my-1 h-px"),
			),
			h.Div(
				g.Map([]menuItem{
					{
						href: "#",
						text: "Profile",
						svg:  `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-badge-check" aria-hidden="true"><path d="M3.85 8.62a4 4 0 0 1 4.78-4.77 4 4 0 0 1 6.74 0 4 4 0 0 1 4.78 4.78 4 4 0 0 1 0 6.74 4 4 0 0 1-4.77 4.78 4 4 0 0 1-6.75 0 4 4 0 0 1-4.78-4.77 4 4 0 0 1 0-6.76Z"></path><path d="m9 12 2 2 4-4"></path></svg>`,
					},
					{
						href: "#",
						text: "Billing",
						svg:  `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-credit-card" aria-hidden="true"><rect width="20" height="14" x="2" y="5" rx="2"></rect><line x1="2" x2="22" y1="10" y2="10"></line></svg>`,
					},
					{
						href: "#",
						text: "Notifications",
						svg:  `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-bell" aria-hidden="true"><path d="M10.268 21a2 2 0 0 0 3.464 0"></path><path d="M3.262 15.326A1 1 0 0 0 4 17h16a1 1 0 0 0 .74-1.673C19.41 13.956 18 12.499 18 8A6 6 0 0 0 6 8c0 4.499-1.411 5.956-2.738 7.326"></path></svg>`,
					},
				}, func(item menuItem) g.Node {
					return h.A(
						c.Classes{
							"relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden transition-colors select-none ": true,
							"hocus:bg-accent hocus:text-accent-foreground": true,
							"[&>svg]:size-4 [&>svg]:shrink-0":              true,
						},
						h.Href(item.href),
						g.Raw(item.svg),
						g.Text(item.text),
					)
				}),
			),
			h.Div(
				h.Class("bg-muted -mx-1 my-1 h-px"),
			),
			h.Div(
				hx.Post(routes.HXLogout),
				hx.Swap("none"),
				c.Classes{
					"relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-hidden transition-colors select-none": true,
					"hocus:bg-accent hocus:text-accent-foreground": true,
					"[&>svg]:size-4 [&>svg]:shrink-0":              true,
				},
				g.Raw(
					`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-log-out" aria-hidden="true"><path d="m16 17 5-5-5-5"></path><path d="M21 12H9"></path><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path></svg>`,
				),
				g.Text("Logout"),
			),
		),
	)
}

func (v *Layout) content(children ...g.Node) g.Node {
	return h.Div(
		h.Class("flex h-full w-full min-w-0 flex-col"),
		v.contentHeader(),
		g.Map(children, func(node g.Node) g.Node {
			return node
		}),
	)
}

func (v *Layout) contentHeader() g.Node {
	return h.Header(
		h.Class(
			"bg-background grid w-full min-w-0 grid-cols-[auto_minmax(0,1fr)_auto] items-center gap-2 border-b px-4 py-4 sm:gap-3 sm:px-6",
		),
		//
		h.Button(
			h.Class(
				"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground size-8 shrink-0",
			),
			g.Raw(
				`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-panel-left" aria-hidden="true"><rect width="18" height="18" x="3" y="3" rx="2"></rect><path d="M9 3v18"></path></svg>`,
			),
			h.Span(
				h.Class("sr-only"),
				g.Text("Toggle Sidebar"),
			),
		),
		h.Div(
			h.Class("min-w-0"),
			h.H1(
				h.Class("truncate text-base font-medium"),
				// g.Text("Ecommerce App"),
			),
		),
		h.Button(
			x.On("click", "console.log('theme should be dark')"),
			h.Class(
				"inline-flex items-center justify-center gap-2 whitespace-nowrap text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground size-9 scale-100 rounded-lg",
			),
			h.Type("button"),
			g.Raw(
				`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tabler-icon tabler-icon-sun size-[1.2rem] scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90"><path d="M8 12a4 4 0 1 0 8 0a4 4 0 1 0 -8 0"></path><path d="M3 12h1m8 -9v1m8 8h1m-9 8v1m-6.4 -15.4l.7 .7m12.1 -.7l-.7 .7m0 11.4l.7 .7m-12.1 -.7l-.7 .7"></path></svg>`,
			),
			g.Raw(
				`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tabler-icon tabler-icon-moon absolute size-[1.2rem] scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0"><path d="M12 3c.132 0 .263 0 .393 0a7.5 7.5 0 0 0 7.92 12.446a9 9 0 1 1 -8.313 -12.454l0 .008"></path></svg>`,
			),
			h.Span(
				h.Class("sr-only"),
				g.Text("Toggle theme"),
			),
		),
	)
}
