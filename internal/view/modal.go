package view

import (
	"fmt"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

const (
	idModal         = "modal"
	EventOpenModal  = "open-modal"
	EventCloseModal = "close-modal"
)

type modal struct {
	icon *icon
}

func newModal() *modal {
	return &modal{
		icon: newIcon(),
	}
}

func (v *modal) container() g.Node {
	return h.Div(
		x.Data("modal"),
		x.Show("show"),
		x.On(fmt.Sprintf("%s.window", EventOpenModal), "open($event)"),
		x.On(fmt.Sprintf("%s.window", EventCloseModal), "close()"),
		x.On("keydown.escape.window.prevent.stop", "close()"),
		h.Style("display: none"),
		h.Class("overflow-y-auto fixed inset-0"),
		h.Div(
			x.Show("show"),
			x.Transition(".opacity"),
			h.Class("fixed inset-0 bg-gray-950/15"),
		),
		h.Div(
			x.Show("show"),
			x.Transition(),
			h.Class("flex relative justify-center items-center p-4 min-h-screen"),
			h.Div(
				h.ID(idModal),
				h.Class(
					"relative p-6 max-w-xl bg-white rounded-lg border border-gray-200 shadow-sm min-w-96",
				),
			),
		),
	)
}

func (v *modal) content(title string, content g.Node) g.Node {
	return g.El(
		"hx-partial",
		hx.Swap("innerHTML"),
		hx.Target(fmt.Sprintf("#%s", idModal)),
		g.Attr("x-trap.inert.noscroll", "show"),
		h.Div(
			x.On("click.outside", "close()"),
			h.Div(
				h.Class("flex gap-8 justify-between items-center"),
				h.H1(
					h.Class("text-sm font-medium"),
					g.Text(title),
				),
				h.Button(
					x.On("click", "close()"),
					v.icon.xMark("size-3 fill-gray-950"),
				),
			),
			content,
		),
	)
}

func (v *modal) xOnClickClose() g.Node {
	return x.On("click", "close()")
}
