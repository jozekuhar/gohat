package view

import (
	"fmt"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/authz"
	"mimokocke/internal/shared/routes"

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
			x.Cloak(),
			x.Data(`{ sidebarOpen: false }`),
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
			x.Cloak(),
			x.Data(`{ darkMode: localStorage.getItem('theme') === 'dark' }`),
			x.Init(
				`$watch('darkMode', value => localStorage.setItem('theme', value ? 'dark' : 'light'))`,
			),
			x.Bind("class", `darkMode && "dark"`),
			h.Class("font-sans"),
		},
	})
}

func (v *Layout) sidebar(identity authz.Identity) g.Node {
	return h.Div(
		h.Class("relative"),
		h.Div(
			c.Classes{
				"[--sidebar-width:16rem] text-sidebar-foreground md:block": true,
			},
			h.Div(
				x.Show("sidebarOpen"),
				x.On("click", "sidebarOpen = false"),
				x.Transition(`:enter="transition-opacity ease-out duration-300"`),
				x.Transition(`:enter-start="opacity-0"`),
				x.Transition(`:enter-end="opacity-100"`),
				x.Transition(`:leave="transition ease-in duration-300"`),
				x.Transition(`:leave-start="opacity-100"`),
				x.Transition(`:leave-end="opacity-0"`),
				c.Classes{
					"fixed inset-0 z-50 bg-black/80 md:hidden!": true,
				},
			),
			h.Div(
				x.Show("sidebarOpen"),
				c.Classes{
					"relative h-svh w-(--sidebar-width) hidden bg-transparent md:block!": true,
				},
			),
			h.Div(
				x.Show("sidebarOpen"),
				x.Transition(`:enter="transition-[left] ease-in-out duration-300"`),
				x.Transition(`:enter-start="-left-full"`),
				x.Transition(`:enter-end="left-0"`),
				x.Transition(`:leave="transition-[left] ease-in-out duration-300"`),
				x.Transition(`:leave-start="left-0"`),
				x.Transition(`:leave-end="-left-full"`),
				c.Classes{
					"fixed inset-y-0 z-50 h-svh w-(--sidebar-width) md:flex! /left-0": true,
				},
				h.Div(
					h.Class("bg-sidebar flex flex-col h-full w-full"),
					v.sidebarHeader(identity),
					g.If(identity.OrgSlug == "", v.sidebarContentEmpty()),
					g.If(identity.OrgSlug != "", v.sidebarContent(identity)),
					v.sidebarFooter(identity),
				),
			),
		),
	)
}

func (v *Layout) sidebarHeader(identity authz.Identity) g.Node {
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
							g.Text(identity.OrgName),
						),
						h.Span(
							h.Class("truncate text-xs"),
							g.Text(identity.OrgSlug),
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
			"absolute top-full inset-x-2 md:left-full md:ml-2 md:top-0 md:mt-2":                                true,
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
	organizations []db.Organization,
) g.Node {
	return g.El(
		"hx-partial",
		hx.Target(fmt.Sprintf("#%s", idSidebarPopoverItems)),
		hx.Swap("outerHTML"),
		h.Div(
			h.ID(idSidebarPopoverItems),
			g.Map(organizations, func(organization db.Organization) g.Node {
				return h.A(
					hx.Boost("true"),
					c.Classes{
						"relative flex cursor-default items-center rounded-sm text-sm outline-hidden transition-colors select-none gap-2 p-2 text-balance": true,
						"hocus:bg-accent hocus:text-accent-foreground": true,
						"[&>svg]:size-4 [&>svg]:shrink-0":              true,
					},
					h.Href(fmt.Sprintf(routes.OrgDashboard, organization.Slug)),
					h.Div(
						c.Classes{
							"flex size-6 items-center justify-center rounded-sm border":   true,
							"[&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0": true,
						},
						g.Raw(
							`<svg id="fi_2891415" enable-background="new 0 0 24 24" height="512" viewBox="0 0 24 24" width="512" xmlns="http://www.w3.org/2000/svg"><g id="XMLID_1_"><path d="m14 9.09 8.81 1.75c.71.15 1.19.75 1.19 1.46v10.2c0 .83-.67 1.5-1.5 1.5h-9c.28 0 .5-.22.5-.5v-.5h8.5c.27 0 .5-.22.5-.5v-10.2c0-.23-.16-.44-.39-.49l-8.61-1.7z"></path><path d="m19.5 14c.28 0 .5.22.5.5s-.22.5-.5.5h-2c-.28 0-.5-.22-.5-.5s.22-.5.5-.5z"></path><path d="m19.5 17c.28 0 .5.22.5.5s-.22.5-.5.5h-2c-.28 0-.5-.22-.5-.5s.22-.5.5-.5z"></path><path d="m19.5 20c.28 0 .5.22.5.5s-.22.5-.5.5h-2c-.28 0-.5-.22-.5-.5s.22-.5.5-.5z"></path><path d="m14 23.5c0 .28-.22.5-.5.5s-.5-.22-.5-.5v-.5-13.5c0-.15.07-.29.18-.39.12-.09.27-.13.42-.1l.4.08v1.02 12.89z"></path><path d="m13 23v.5c0 .28.22.5.5.5h-4c.28 0 .5-.22.5-.5v-.5z"></path><path d="m10.5 5c.28 0 .5.22.5.5s-.22.5-.5.5h-2c-.28 0-.5-.22-.5-.5s.22-.5.5-.5z"></path><path d="m11 8.5c0 .28-.22.5-.5.5h-2c-.28 0-.5-.22-.5-.5s.22-.5.5-.5h2c.28 0 .5.22.5.5z"></path><path d="m10.5 11c.28 0 .5.22.5.5s-.22.5-.5.5h-2c-.28 0-.5-.22-.5-.5s.22-.5.5-.5z"></path><path d="m10.5 14c.28 0 .5.22.5.5s-.22.5-.5.5h-2c-.28 0-.5-.22-.5-.5s.22-.5.5-.5z"></path><path d="m6 14.5c0 .28-.22.5-.5.5h-2c-.28 0-.5-.22-.5-.5s.22-.5.5-.5h2c.28 0 .5.22.5.5z"></path><path d="m5.5 5c.28 0 .5.22.5.5s-.22.5-.5.5h-2c-.28 0-.5-.22-.5-.5s.22-.5.5-.5z"></path><path d="m5.5 8c.28 0 .5.22.5.5s-.22.5-.5.5h-2c-.28 0-.5-.22-.5-.5s.22-.5.5-.5z"></path><path d="m5.5 11c.28 0 .5.22.5.5s-.22.5-.5.5h-2c-.28 0-.5-.22-.5-.5s.22-.5.5-.5z"></path><path d="m9 18.5c0-.28-.23-.5-.5-.5h-3c-.28 0-.5.22-.5.5v4.5h-1v-4.5c0-.83.67-1.5 1.5-1.5h3c.83 0 1.5.67 1.5 1.5v4.5h-1z"></path><path d="m5 23h4 1v.5c0 .28-.22.5-.5.5h-5c-.28 0-.5-.22-.5-.5v-.5z"></path><path d="m1.75.2 10.99 1.67c.73.12 1.26.74 1.26 1.48v5.74l-.4-.08c-.15-.03-.3.01-.42.1-.11.1-.18.24-.18.39v-6.15c0-.25-.18-.46-.42-.5l-10.99-1.66c-.03-.01-.06-.01-.09-.01-.12 0-.23.04-.32.12-.12.1-.18.23-.18.38v20.82c0 .28.23.5.5.5h2.5v.5c0 .28.22.5.5.5h-3c-.83 0-1.5-.67-1.5-1.5v-20.82c0-.44.19-.86.53-1.14.34-.29.78-.41 1.22-.34z"></path></g></svg>`,
						),
						// h.Img(
						// 	h.Alt("Artifact"),
						// 	h.Loading("lazy"),
						// 	h.Width("18"),
						// 	h.Height("18"),
						// 	g.Attr("decoding", "async"),
						// 	h.Data("nimg", "1"),
						// 	h.Class("size-4 shrink-0 dark:hidden"),
						// 	h.Src(""),
						// 	h.Style("color: transparent;"),
						// ),
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
					h.Class(
						"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden md:text-sm col-span-4",
					),
					h.AutoFocus(),
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
						hx.Boost("true"),
						c.Classes{
							"flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left outline-hidden ring-sidebar-ring h-8 text-sm": true,
							"transition-[width,height,padding] ": true,
							"focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground": true,
							"[&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0 ":                                                                                    true,
						},
						h.Href(fmt.Sprintf(routes.OrgDashboard, identity.OrgSlug)),
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
						hx.Boost("true"),
						c.Classes{
							"flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left outline-hidden ring-sidebar-ring h-8 text-sm": true,
							"transition-[width,height,padding]": true,
							"focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground": true,
							"[&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0":                                                                                     true,
						},
						h.Href(fmt.Sprintf(routes.OrgMemberships, identity.OrgSlug)),
						g.Raw(
							`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tabler-icon tabler-icon-users "><path d="M5 7a4 4 0 1 0 8 0a4 4 0 1 0 -8 0"></path><path d="M3 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2"></path><path d="M16 3.13a4 4 0 0 1 0 7.75"></path><path d="M21 21v-2a4 4 0 0 0 -3 -3.85"></path></svg>`,
						),
						h.Span(
							g.Text("Users"),
						),
					),
				),
				h.Li(
					h.Class("relative"),
					h.A(
						hx.Boost("true"),
						c.Classes{
							"flex w-full items-center gap-2 overflow-hidden rounded-md p-2 text-left outline-hidden ring-sidebar-ring h-8 text-sm": true,
							"transition-[width,height,padding] ": true,
							"focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground": true,
							"[&>span:last-child]:truncate [&>svg]:size-4 [&>svg]:shrink-0 ":                                                                                    true,
						},
						h.Href(fmt.Sprintf(routes.OrgChannels, identity.OrgSlug)),
						g.Raw(
							`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tabler-icon tabler-icon-apps "><path d="M4 5a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v4a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1l0 -4"></path><path d="M4 15a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v4a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1l0 -4"></path><path d="M14 15a1 1 0 0 1 1 -1h4a1 1 0 0 1 1 1v4a1 1 0 0 1 -1 1h-4a1 1 0 0 1 -1 -1l0 -4"></path><path d="M14 7l6 0"></path><path d="M17 4l0 6"></path></svg>`,
						),
						h.Span(
							g.Text("Channels"),
						),
					),
				),
			),
		),
	)
}

func (v *Layout) sidebarFooter(identity authz.Identity) g.Node {
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
							g.Textf("%s %s", identity.FirstName, identity.LastName),
						),
						h.Span(
							h.Class("truncate text-xs"),
							g.Text(identity.Email),
						),
					),
					g.Raw(
						`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-chevrons-up-down ml-auto size-4" aria-hidden="true"><path d="m7 15 5 5 5-5"></path><path d="m7 9 5-5 5 5"></path></svg>`,
					),
				),
			),
		),
		v.sidebarFooterPopover(identity),
	)
}

func (v *Layout) sidebarFooterPopover(identity authz.Identity) g.Node {
	type menuItem struct {
		href string
		svg  string
		text string
	}

	return h.Div(
		x.Show("menuOpen"),
		x.On("click.outside", "menuOpen = false"),
		c.Classes{
			"absolute bottom-full inset-x-2 md:left-full md:ml-2 md:bottom-0 md:mb-2":                          true,
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
						g.Textf("%s %s", identity.FirstName, identity.LastName),
					),
					h.Span(
						h.Class("truncate text-xs"),
						g.Text(identity.Email),
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
					hx.Boost("true"),
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
		h.Button(
			x.On("click", `sidebarOpen = !sidebarOpen; console.log(sidebarOpen)`),
			c.Classes{
				"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground size-8 shrink-0": true,
				"md:invisible": true,
			},
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
				g.Text(""),
			),
		),
		h.Button(
			x.On("click", "darkMode = !darkMode; console.log(darkMode)"),
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
