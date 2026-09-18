package channel

import (
	"fmt"

	"mimokocke/internal/provider/db/sqlc"
	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/web/view"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

const (
	idChannelList = "channel-list"
)

func idChannelItem(id string) string {
	return "channel-item-" + id
}

func channelsPage(ident identity.IdentityCtx, chns []sqlc.Channel) g.Node {
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
										g.Text("Channels"),
									),
									h.Span(
										h.Class(
											"bg-muted text-muted-foreground shrink-0 rounded-md px-1.5 py-0.5 text-xs dark:bg-zinc-900 dark:text-zinc-300",
										),
										g.Textf("%d", len(chns)),
									),
								),
								h.Button(
									hx.Get(
										fmt.Sprintf(
											routes.HXOrgChannelsCreate,
											ident.OrgSlug,
										),
									),
									hx.Swap("none"),
									h.Class(
										"inline-flex items-center justify-center whitespace-nowrap font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 rounded-md text-xs size-9 shrink-0 gap-1.5 p-0 sm:w-auto sm:px-3 lg:hidden",
									),
									h.Aria("label", "Invite member"),
									g.Raw(
										`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-plus size-3.5" aria-hidden="true"><path d="M5 12h14"></path><path d="M12 5v14"></path></svg>`,
									),
									h.Span(
										h.Class("hidden sm:inline"),
										g.Text("Add Channel"),
									),
								),
							),
							h.Div(
								hx.Get(
									fmt.Sprintf(
										routes.HXOrgChannelsCreate,
										ident.OrgSlug,
									),
								),
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
									g.Text("Add Channel"),
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
							h.ID(idChannelList),
							g.Map(chns, func(chn sqlc.Channel) g.Node {
								return channelItem(ident, chn)
							}),
						),
					),
				),
			),
		),
	)
}

func channelItem(ident identity.IdentityCtx, chn sqlc.Channel) g.Node {
	return h.Div(
		hx.Get(fmt.Sprintf(routes.HXOrgChannelsUpdate, ident.OrgSlug, chn.ID.String())),
		hx.Swap("none"),
		h.ID(idChannelItem(chn.ID.String())),
		h.Class(
			"hover:bg-muted/40 flex min-w-225 items-center border-b px-4 py-3 text-sm last:border-b-0 sm:px-6 lg:min-w-0 dark:border-white/8 dark:hover:bg-white/4",
		),
		h.Div(
			h.Class(
				"w-80 shrink-0 text-xs text-zinc-500 lg:w-[25%] lg:shrink dark:text-zinc-400",
			),
			g.Text(chn.Name),
		),
		h.Div(
			h.Class(
				"w-37.5 shrink-0 text-xs text-zinc-500 lg:w-[25%] lg:shrink dark:text-zinc-400",
			),
			func() g.Node {
				switch chn.Provider {
				case sqlc.ChannelProviderShopify:
					return g.Text("Shopify")
				case sqlc.ChannelProviderWooCommerce:
					return g.Text("WooCommerce")
				default:
					return g.Text("#N/A")
				}
			}(),
		),
		// h.Div(
		// 	h.Class(
		// 		"w-25 shrink-0 text-xs text-zinc-500 lg:w-[25%] lg:shrink dark:text-zinc-400",
		// 	),
		// ),
		h.Div(
			h.Class(
				"w-82.5 shrink-0 text-xs text-zinc-500 lg:w-[25%] lg:shrink dark:text-zinc-400",
			),
			func() g.Node {
				switch chn.Status {
				case sqlc.ChannelStatusActive:
					return g.Text("Active")
				case sqlc.ChannelStatusInactive:
					return g.Text("Inactive")
				default:
					return g.Text("#N/A")
				}
			}(),
		),
	)
}

func createChannelFormModal(ident identity.IdentityCtx) g.Node {
	formID := "channel-form"

	return view.ModalFragment(
		h.Div(
			h.Class("flex flex-col gap-2 text-center sm:text-left"),
			h.H2(
				h.Class("text-lg leading-none font-semibold"),
				g.Text("Add New Channel"),
			),
			h.P(
				h.Class("text-muted-foreground text-sm"),
				g.Text("Add new channel here. Click save when you're done."),
			),
		),
		h.Div(
			x.Data(`{ provider: null }`),
			x.Cloak(),
			h.Class("space-y-4"),
			view.Input(view.InputParams{
				Label:        "Name",
				Name:         "Name",
				Placeholder:  "GLS Germany",
				AutoFocus:    true,
				AutoComplete: "off",
				Form:         formID,
			}),
			view.Select(view.SelectParams{
				XModel:      "provider",
				Label:       "Provider",
				Placeholder: "Select provider",
				Options: []view.SelectOption{
					{Value: "woocommerce", Text: "WooCommerce"},
					{Value: "shopify", Text: "Shopify"},
					{Value: "manual", Text: "Manual"},
				},
			}),
			h.Template(
				x.If(`provider === 'woocommerce'`),
				h.Form(
					x.Init(`htmx.process($el)`),
					h.ID(formID),
					h.Class("space-y-4"),
					view.Input(view.InputParams{
						Label:        "Store URL",
						Name:         "StoreURL",
						AutoComplete: "off",
					}),
					view.Input(view.InputParams{
						Label:        "Consumer Key",
						Name:         "ConsumerKey",
						AutoComplete: "off",
					}),
					view.Input(view.InputParams{
						Label:        "Consumer Secret",
						Name:         "ConsumerSecret",
						AutoComplete: "off",
					}),
					h.Div(
						h.Class("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"),
						h.Button(
							hx.Post(
								fmt.Sprintf(routes.HXOrgChannelsWooCommerceTest, ident.OrgSlug),
							),
							hx.Swap("none"),
							h.Class(
								"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 border border-input bg-background shadow-xs hover:bg-accent hover:text-accent-foreground h-9 px-4 py-2",
							),
							h.Type("submit"),
							g.Text("Test connection"),
						),
						h.Button(
							hx.Post(
								fmt.Sprintf(routes.HXOrgChannelsWooCommerceCreate, ident.OrgSlug),
							),
							hx.Target("#"+idChannelList),
							hx.Swap("append"),
							h.Class(
								"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2",
							),
							h.Type("submit"),
							g.Text("Save changes"),
						),
					),
				),
			),
			h.Template(
				x.If(`provider === 'shopify'`),
				h.Form(
					x.Init(`htmx.process($el)`),
					h.ID(formID),
					h.Class("space-y-4"),
					view.Input(view.InputParams{
						Label:        "secret",
						Name:         "secret",
						AutoComplete: "off",
					}),
					h.Button(
						hx.Post(fmt.Sprintf(routes.HXOrgChannelsShopifyTest, ident.OrgSlug)),
						hx.Swap("none"),
						g.Text("Test"),
					),
					h.Div(
						hx.Post(fmt.Sprintf(routes.HXOrgChannelsShopifyCreate, ident.OrgSlug)),
						hx.Target("#"+idChannelList),
						hx.Swap("append"),
						h.Class("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"),
						h.Button(
							h.Class(
								"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2",
							),
							h.Type("submit"),
							g.Text("Save changes"),
						),
					),
				),
			),
			h.Template(
				x.If(`provider === 'manual'`),
				h.Form(
					x.Init(`htmx.process($el)`),
					h.ID(formID),
					h.Class("space-y-4"),
					h.Div(
						// hx.Post(fmt.Sprintf(routes.HXOrgChannelsShopifyCreate, ident.OrgSlug)),
						// hx.Target("#"+idChannelList),
						// hx.Swap("append"),
						h.Class("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"),
						h.Button(
							h.Class(
								"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2",
							),
							h.Type("submit"),
							g.Text("Save changes"),
						),
					),
				),
			),
		),
	)
}

func updateChannelFormModal(ident identity.IdentityCtx, chn ChannelDetails) g.Node {
	idChannelForm := "channel-form"

	return view.ModalFragment(
		h.Div(
			h.Class("flex flex-col gap-2 text-center sm:text-left"),
			h.H2(
				h.Class("text-lg leading-none font-semibold"),
				g.Text("Update Channel"),
			),
			h.P(
				h.Class("text-muted-foreground text-sm"),
				g.Text("Update channel here. Click save when you're done."),
			),
		),
		h.Div(
			x.Data(fmt.Sprintf(`{ provider: %q }`, chn.Channel.Provider.String())),
			x.Cloak(),
			h.Class("space-y-4"),
			view.Input(view.InputParams{
				Label:        "Name",
				Name:         "Name",
				Value:        chn.Channel.Name,
				AutoFocus:    true,
				AutoComplete: "off",
				Form:         idChannelForm,
			}),
			view.Select(view.SelectParams{
				XModel:      "provider",
				Label:       "Provider",
				Placeholder: "Select provider",
				Options: []view.SelectOption{
					{Value: "woocommerce", Text: "WooCommerce"},
					{Value: "shopify", Text: "Shopify"},
				},
			}),
			h.Template(
				x.If(`provider === 'woocommerce'`),
				h.Form(
					x.Init(`htmx.process($el)`),
					h.ID(idChannelForm),
					h.Class("space-y-4"),
					view.Input(view.InputParams{
						Label:        "Store URL",
						Name:         "StoreURL",
						Value:        chn.MaskedWooCommerceCredentials.StoreURL,
						AutoComplete: "off",
					}),
					view.Input(view.InputParams{
						Label:        "Consumer Key",
						Name:         "ConsumerKey",
						Value:        chn.MaskedWooCommerceCredentials.MaskedConsumerKey,
						AutoComplete: "off",
					}),
					view.Input(view.InputParams{
						Label:        "Consumer Secret",
						Name:         "ConsumerSecret",
						Value:        chn.MaskedWooCommerceCredentials.MaskedConsumerSecret,
						AutoComplete: "off",
					}),
					h.Div(
						h.Class("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"),
						h.Button(
							hx.Post(
								fmt.Sprintf(routes.HXOrgChannelsWooCommerceTest, ident.OrgSlug),
							),
							hx.Swap("none"),
							h.Class(
								"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 border border-input bg-background shadow-xs hover:bg-accent hover:text-accent-foreground h-9 px-4 py-2",
							),
							h.Type("submit"),
							g.Text("Test connection"),
						),
						h.Button(
							hx.Delete(
								fmt.Sprintf(
									routes.HXOrgChannelsDelete,
									ident.OrgSlug,
									chn.Channel.ID.String(),
								),
							),
							hx.Swap("delete"),
							hx.Target(fmt.Sprintf("#channel_%s", chn.Channel.ID.String())),
							h.Type("button"),
							h.Class(
								"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-destructive text-destructive-foreground shadow-xs hover:bg-destructive/90 h-9 px-4 py-2",
							),
							g.Text("Delete"),
						),
						h.Button(
							// hx.Post(
							// 	fmt.Sprintf(routes.HXOrgChannelsWooCommerceCreate, ident.OrgSlug),
							// ),
							hx.Swap("none"),
							h.Class(
								"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2",
							),
							h.Type("submit"),
							g.Text("Save changes"),
						),
					),
				),
			),
			h.Template(
				x.If(`provider === 'shopify'`),
				h.Form(
					x.Init(`htmx.process($el)`),
					h.ID(idChannelForm),
					h.Class("space-y-4"),
					view.Input(view.InputParams{
						Label:        "secret",
						Name:         "secret",
						AutoComplete: "off",
					}),
					h.Button(
						hx.Post(fmt.Sprintf(routes.HXOrgChannelsShopifyTest, ident.OrgSlug)),
						hx.Swap("none"),
						g.Text("Test"),
					),
					h.Div(
						// hx.Post(fmt.Sprintf(routes.HXOrgChannelsShopifyCreate, ident.OrgSlug)),
						hx.Swap("none"),
						h.Class("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"),
						h.Button(
							h.Class(
								"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium  focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2",
							),
							h.Type("submit"),
							g.Text("Save changes"),
						),
					),
				),
			),
		),
	)
}
