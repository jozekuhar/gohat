package view

import (
	"fmt"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

const idModalPortal = "modal-portal"

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
		h.ID("modal-portal"),
	)
}

func (v *modal) fragment(title string, content g.Node) g.Node {
	titleID := "modal-title"

	return g.El(
		"hx-partial",
		hx.Swap("innerHTML"),
		hx.Target(fmt.Sprintf("#%s", idModalPortal)),
		h.Div(
			x.Cloak(),
			x.Data("modal"),
			x.Show("show"),
			x.On("keydown.escape.window.prevent.stop", "destroy()"),
			h.Style("display: none"),

			h.Class("overflow-y-auto fixed inset-0"),
			h.Role("dialog"),
			h.Aria("modal", "true"),
			h.Aria("labelledby", titleID),
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
					x.On("click.outside", "destroy()"),
					g.Attr("x-trap.inert.noscroll", "show"),
					h.Class(
						"relative p-6 max-w-xl bg-white rounded-lg border border-gray-100 shadow-sm min-w-96",
					),
					h.Div(
						h.Class("flex gap-8 justify-between items-center"),
						h.H2(
							h.ID(titleID),
							h.Class("text-base font-medium"),
							g.Text(title),
						),
						h.Button(
							x.On("click", "destroy()"),
							h.Class(
								"relative rounded-full focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 focus-visible:ring-offset-2 before:content-[''] before:absolute before:-inset-2",
							),
							h.Div(
								h.Class("flex justify-center items-center size-10"),
								v.icon.xMark("size-3.5 fill-gray-400"),
							),
						),
					),
					content,
				),
			),
		),
	)
}

func (v *modal) onClose() g.Node {
	return x.On("click", "destroy()")
}
