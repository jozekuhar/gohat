package view

import (
	"mimokocke/internal/shared/authz"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type Channel struct {
	layout *Layout
}

func NewChannel() *Channel {
	return &Channel{
		layout: NewLayout(),
	}
}

func (v *Channel) ChannelsPage(identity authz.Identity) g.Node {
	return v.layout.app(
		identity,
		h.Main(
			h.H1(g.Text("Channels")),
		),
	)
}
