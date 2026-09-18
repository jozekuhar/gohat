package orderweb

import (
	"fmt"

	"mimokocke/internal/web/view"

	"mimokocke/internal/provider/db/sqlc"
	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/routes"

	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

const idOrderList = "order-list"

func ordersPage(ident identity.IdentityCtx, orders []sqlc.ListOrdersWithDetailsRow) g.Node {
	return view.AppLayout(
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
										g.Text("Orders"),
									),
									h.Span(
										h.Class(
											"bg-muted text-muted-foreground shrink-0 rounded-md px-1.5 py-0.5 text-xs dark:bg-zinc-900 dark:text-zinc-300",
										),
										g.Textf("%d", len(orders)),
									),
								),
								h.Button(
									hx.Get(fmt.Sprintf(routes.HXOrgOrderCreate, ident.OrgSlug)),
									hx.Swap("none"),
									h.Class(
										"inline-flex items-center justify-center whitespace-nowrap font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 rounded-md text-xs size-9 shrink-0 gap-1.5 p-0 sm:w-auto sm:px-3 lg:hidden",
									),
									g.Raw(
										`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-plus size-3.5" aria-hidden="true"><path d="M5 12h14"></path><path d="M12 5v14"></path></svg>`,
									),
									h.Span(
										h.Class("hidden sm:inline"),
										g.Text("Add Order"),
									),
								),
							),
							h.Div(
								hx.Get(fmt.Sprintf(routes.HXOrgOrderCreate, ident.OrgSlug)),
								hx.Swap("none"),
								h.Class(
									"grid min-w-0 grid-cols-2 items-center gap-2 lg:flex lg:flex-wrap lg:justify-end",
								),
								h.Button(
									h.Class(
										"items-center justify-center whitespace-nowrap font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 rounded-md px-3 text-xs hidden h-8 gap-1.5 lg:inline-flex",
									),
									g.Raw(
										`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-plus size-3.5" aria-hidden="true"><path d="M5 12h14"></path><path d="M12 5v14"></path></svg>`,
									),
									g.Text("Add Order"),
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
									h.Class("w-80 shrink-0 lg:w-[25%] lg:shrink"),
									g.Text("Name"),
								),
								h.Div(
									h.Class("w-37.5 shrink-0 lg:w-[25%] lg:shrink"),
									g.Text("Provider"),
								),
								// h.Div(
								// 	h.Class("w-25 shrink-0 lg:w-[25%] lg:shrink"),
								// ),
								h.Div(
									h.Class("w-82.5 shrink-0 lg:w-[25%] lg:shrink"),
									g.Text("Status"),
								),
							),
						),
						h.Div(
							h.ID(idOrderList),
							g.Map(orders, func(order sqlc.ListOrdersWithDetailsRow) g.Node {
								return orderItem(ident, order.Order, order.Channel)
							}),
						),
					),
				),
			),
		),
	)
}

func orderItem(ident identity.IdentityCtx, order sqlc.Order, channel sqlc.Channel) g.Node {
	return h.Div(
		c.Classes{
			"hover:bg-muted/40 flex min-w-225 items-center border-b px-4 py-3 text-sm last:border-b-0 sm:px-6 lg:min-w-0 dark:border-white/8 dark:hover:bg-white/4": true,
			"a": true,
		},
		h.Div(
			h.Class(
				"w-80 shrink-0 text-xs text-zinc-500 lg:w-[25%] lg:shrink dark:text-zinc-400",
			),
			g.Text(order.ChannelID.String()),
		),
		h.Div(
			h.Class(
				"w-37.5 shrink-0 text-xs text-zinc-500 lg:w-[25%] lg:shrink dark:text-zinc-400",
			),
		),
		h.Div(
			h.Class(
				"w-25 shrink-0 text-xs text-zinc-500 lg:w-[25%] lg:shrink dark:text-zinc-400",
			),
		),
		h.Div(
			h.Class(
				"w-82.5 shrink-0 text-xs text-zinc-500 lg:w-[25%] lg:shrink dark:text-zinc-400",
			),
		),
	)
}

func createOrderFormModal() g.Node {
	return view.ModalFragment(
		h.Div(
			h.Class("flex flex-col gap-2 text-center sm:text-left"),
			h.H2(
				h.Class("text-lg leading-none font-semibold"),
				g.Text("Add Order"),
			),
			h.P(
				h.Class("text-muted-foreground text-sm"),
				g.Text("Add order here. Click save when you're done."),
			),
		),
		h.Form(
			h.Class("space-y-4"),
			view.Input(view.InputParams{
				Label:        "First Name",
				Name:         "FirstName",
				AutoFocus:    true,
				AutoComplete: "off",
			}),
			view.Input(view.InputParams{
				Label:        "Last Name",
				Name:         "LastName",
				AutoComplete: "off",
			}),
			view.Select(view.SelectParams{
				Label:       "Channel",
				Placeholder: "Select channel",
				Options: []view.SelectOption{
					{Value: "woocommerce", Text: "WooCommerce"},
					{Value: "shopify", Text: "Shopify"},
				},
			}),
			view.Select(view.SelectParams{
				Label:       "Courier",
				Placeholder: "Select courier",
				Options: []view.SelectOption{
					{Value: "woocommerce", Text: "WooCommerce"},
					{Value: "shopify", Text: "Shopify"},
				},
			}),
		),
	)
}
