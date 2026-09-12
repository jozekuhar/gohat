package view

import (
	"fmt"

	"mimokocke/internal/shared/authz"
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

func (v *channel) ChannelsPage(identity authz.Identity) g.Node {
	return v.layout.app(
		identity,
		h.Main(
			h.H1(g.Text("Channels")),

			h.Button(
				hx.Get(fmt.Sprintf(routes.HXOrgChannelsCreate, identity.OrgSlug)),
				hx.Swap("none"),
				g.Text("New channel"),
			),
		),
	)
}

func (v *channel) CreateChannelFormModal(identity authz.Identity) g.Node {
	return v.modal.fragment(
		h.Div(
			h.Class(
				"bg-background data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 fixed top-[50%] left-[50%] z-50 grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 rounded-lg border p-6 shadow-lg duration-200 sm:max-w-lg",
			),
			h.Div(
				h.Class("flex flex-col gap-2 text-center sm:text-left"),
				h.H2(
					h.Class("text-lg leading-none font-semibold"),
					g.Text("Add New User"),
				),
				h.P(
					h.Class("text-muted-foreground text-sm"),
					g.Text("Add new channel here. Click save when you're done."),
				),
			),
			h.Form(
				// TODO(jozekuhar): think how can i inject data to component
				x.Data(`{ 
					provider: null,
					options: [
						{ value: "woocommerce", text: "WooCommerce" },
						{ value: "shopify", text: "Shopify" }
					]
				}`),
				h.Class("space-y-4"),
				ui.Input(
					"Name",
					h.Name("Name"),
					h.Placeholder("GLS Germany"),
					h.AutoFocus(),
					h.AutoComplete("off"),
				),
				// START
				ui.Select(ui.SelectParams{
					Label:       "Provider",
					Placeholder: "Select provider",
					XModel:      "provider",
					Options: []ui.SelectOption{
						{Value: "woocommerce", Text: "WooCommerce"},
						{Value: "shopify", Text: "Shopify"},
					},
				}),
				// END
				h.Template(
					x.If("provider === 'woocommerce'"),
					h.Div(
						h.Class("space-y-4"),
						ui.Input(
							"WooCommerce URL",
							h.Name("URL"),
							h.Placeholder(""),
							h.AutoComplete("off"),
						),
						ui.Input(
							"Consumer Key",
							h.Name("ConsumerKey"),
							h.Placeholder(""),
							h.AutoComplete("off"),
						),
						ui.Input(
							"Consumer Secret",
							h.Name("ConsumerSecret"),
							h.Placeholder(""),
							h.AutoComplete("off"),
						),
					),
				),
				h.Template(
					x.If("provider === 'shopify'"),
					h.Div(
						ui.Input(
							"todo",
							h.Placeholder("todo"),
						),
					),
				),
				h.Div(
					h.Data("slot", "dialog-footer"),
					h.Class("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"),
					h.Button(
						h.Class(
							"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2",
						),
						h.Type("submit"),
						h.FormAttr("user-form"),
						g.Text("Save changes"),
					),
				),
				h.Button(
					h.Type("button"),
					h.Class(
						"ring-offset-background focus:ring-ring data-[state=open]:bg-accent data-[state=open]:text-muted-foreground absolute top-4 right-4 rounded-xs opacity-70 transition-opacity hover:opacity-100 focus:ring-2 focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
					),
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
						h.Class("lucide lucide-x"),
						h.Aria("hidden", "true"),
						g.El("path",
							g.Attr("d", "M18 6 6 18"),
						),
						g.El("path",
							g.Attr("d", "m6 6 12 12"),
						),
					),
					h.Span(
						h.Class("sr-only"),
						g.Text("Close"),
					),
				),
			),
		),
	)
}
