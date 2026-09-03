package view

import (
	"mimokocke/internal/channel"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type Channel struct {
	layout *layout
}

func NewChannel() *Channel {
	return &Channel{
		layout: newLayout(),
	}
}

func (v *Channel) ChannelsPage(channels []channel.Channel) g.Node {
	return v.layout.App(
		h.Main(
			h.Class("max-w-7xl mx-auto"),
			h.H1(g.Text("Channels")),
			h.Div(
				g.Text("list channels"),
			),
			h.Form(
				h.Class("border flex flex-col"),
				h.Input(
					h.Placeholder("Name"),
				),
				h.Select(
					h.Option(g.Text("WooCommerce")),
					h.Option(g.Text("Shopify")),
				),
				// Based on selcetion we render credentials
				h.Button(
					g.Text("Save channel"),
				),
			),
		),
	)
}
