package view

import (
	"fmt"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/web/view/ui"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

type channel struct {
	layout *Layout
	modal  *modal
}

func newChannel(layout *Layout, modal *modal) *channel {
	return &channel{
		layout: layout,
		modal:  modal,
	}
}

func (v *channel) ChannelsPage(ident identity.Identity, channels []db.Channel) g.Node {
	return v.layout.app(
		ident,
		h.Main(
			h.H1(g.Text("Channels")),

			h.Button(
				hx.Get(fmt.Sprintf(routes.HXOrgChannelsCreate, ident.OrgSlug)),
				hx.Swap("none"),
				g.Text("New channel"),
			),

			h.Ul(
				g.Map(channels, func(channel db.Channel) g.Node {
					return h.Li(
						h.Span(g.Text(channel.Name)),
					)
				}),
			),
		),
	)
}

func (v *channel) CreateChannelFormModal(ident identity.Identity) g.Node {
	return v.modal.fragment(
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
		h.Form(
			hx.Post(fmt.Sprintf(routes.HXOrgChannelsCreate, ident.OrgSlug)),
			hx.Swap("none"),
			x.Data(`{ provider: null }`),
			h.Class("space-y-4"),
			ui.Input(ui.InputParams{
				Label:        "Name",
				Name:         "Name",
				Placeholder:  "GLS Germany",
				AutoFocus:    true,
				AutoComplete: "off",
			}),
			ui.Select(ui.SelectParams{
				XModel:      "provider",
				Label:       "Provider",
				Name:        "Provider",
				Placeholder: "Select provider",
				Options: []ui.SelectOption{
					{Value: "woocommerce", Text: "WooCommerce"},
					{Value: "shopify", Text: "Shopify"},
				},
			}),
			h.Template(
				x.If("provider === 'woocommerce'"),
				h.Div(
					h.Class("space-y-4"),
					ui.Input(ui.InputParams{
						Label:        "WooCommerce URL",
						Name:         "URL",
						AutoComplete: "off",
					}),
					ui.Input(ui.InputParams{
						Label:        "Consumer Key",
						Name:         "ConsumerKey",
						AutoComplete: "off",
					}),
					ui.Input(ui.InputParams{
						Label:        "Consumer Secret",
						Name:         "ConsumerSecret",
						AutoComplete: "off",
					}),
				),
			),
			h.Template(
				x.If("provider === 'shopify'"),
				h.Div(
					ui.Input(ui.InputParams{
						Label:        "secret",
						Name:         "secret",
						AutoComplete: "off",
					}),
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
