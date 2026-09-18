package view

import (
	"fmt"

	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/routes"

	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type tabKey string

const (
	TabGeneral tabKey = "general"
	TabProfile tabKey = "profile"
)

func SettingsLayout(ident identity.IdentityCtx, tab tabKey, children ...g.Node) g.Node {
	type menuItem struct {
		tab  tabKey
		href string
		svg  string
		text string
	}

	return h.Div(
		h.Class("flex flex-1 flex-col gap-4 overflow-hidden p-4"),
		h.Div(
			h.Class("space-y-0.5"),
			h.H1(
				h.Class("text-xl font-bold tracking-tight md:text-2xl"),
				g.Text("Settings"),
			),
			h.P(
				h.Class("text-muted-foreground"),
				g.Text("Update account preferences and manage integrations."),
			),
		),
		h.Div(
			h.Class(
				"flex flex-1 flex-col space-y-8 overflow-auto md:space-y-2 md:overflow-hidden lg:flex-row lg:space-y-0 lg:space-x-12",
			),
			h.Aside(
				h.Class("lg:sticky lg:w-1/5"),
				h.Div(
					h.Class("p-1 md:hidden"),
					h.Button(
						h.Type("button"),
						h.Role("combobox"),
						h.Aria("expanded", "false"),
						h.Aria("autocomplete", "none"),
						h.Dir("ltr"),
						h.Class(
							"border-input ring-offset-background placeholder:text-muted-foreground focus:ring-ring flex w-full items-center justify-between rounded-md border bg-transparent px-3 py-2 text-sm whitespace-nowrap shadow-xs focus:ring-1 focus:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1 h-10 sm:w-48",
						),
						h.Span(
							h.Style("pointer-events: none;"),
							h.Div(
								h.Class("flex gap-x-4 px-2 py-0.5"),
								h.Span(
									h.Class("scale-125 [&_svg]:size-4.5"),
									g.Raw(
										`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tabler-icon tabler-icon-tool "><path d="M7 10h3v-3l-3.5 -3.5a6 6 0 0 1 8 8l6 6a2 2 0 0 1 -3 3l-6 -6a6 6 0 0 1 -8 -8l3.5 3.5"></path></svg>`,
									),
								),
								h.Span(
									h.Class("text-md"),
									g.Text("General"),
								),
							),
						),
						g.Raw(
							`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-chevron-down h-4 w-4 opacity-50" aria-hidden="true"><path d="m6 9 6 6 6-6"></path></svg>`,
						),
					),
				),
				h.Div(
					h.Dir("ltr"),
					h.Class(
						"relative overflow-hidden bg-background hidden w-full min-w-48 px-1 py-2 md:block",
					),
					h.Div(
						h.Class("h-full w-full rounded-[inherit] overflow-x-auto!"),
						h.Style("overflow: scroll hidden;"),
						h.Div(
							h.Style("min-width: 100%; display: table;"),
							h.Nav(
								h.Class(
									"flex space-x-2 py-1 lg:flex-col lg:space-y-1 lg:space-x-0",
								),
								g.Map([]menuItem{
									{
										text: "General",
										tab:  "general",
										href: fmt.Sprintf(routes.SettingsGeneral, ident.OrgSlug),
										svg:  `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tabler-icon tabler-icon-tool "><path d="M7 10h3v-3l-3.5 -3.5a6 6 0 0 1 8 8l6 6a2 2 0 0 1 -3 3l-6 -6a6 6 0 0 1 -8 -8l3.5 3.5"></path></svg>`,
									},
									{
										text: "Profile",
										tab:  "profile",
										href: fmt.Sprintf(routes.SettingsProfile, ident.OrgSlug),
										svg:  `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tabler-icon tabler-icon-user "><path d="M8 7a4 4 0 1 0 8 0a4 4 0 0 0 -8 0"></path><path d="M6 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2"></path></svg>`,
									},
								}, func(item menuItem) g.Node {
									return h.A(
										hx.Boost("true"),
										c.Classes{
											"inline-flex items-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 hover:text-accent-foreground h-9 px-4 py-2  justify-start": true,
											"bg-muted hover:bg-muted":              tab == item.tab,
											"hover:bg-transparent hover:underline": tab != item.tab,
										},
										h.Href(item.href),
										h.Span(
											h.Class("mr-2 [&_svg]:size-4.5"),
											g.Raw(item.svg),
										),
										g.Text(item.text),
									)
								}),
							),
						),
					),
					h.Div(
						h.Class(
							"flex touch-none select-none h-2.5 flex-col border-t border-t-transparent p-px",
						),
						h.Div(
							h.Class("bg-border relative flex-1 rounded-full"),
						),
					),
				),
			),
			h.Div(
				h.Class("flex w-full overflow-y-scroll p-1 pr-4 md:overflow-y-hidden"),
				g.Group(children),
			),
		),
	)
}
