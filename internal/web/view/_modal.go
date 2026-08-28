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
		h.ID(idModalPortal),
	)
}

func (v *modal) fragment(title string, content g.Node) g.Node {
	titleID := "modal-title"

	_ = v.fragment
	_ = v.onClose

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

			h.Class("fixed inset-0 overflow-y-auto"),
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
				h.Class("relative flex min-h-screen items-center justify-center p-4"),
				h.Div(
					x.On("click.outside", "destroy()"),
					g.Attr("x-trap.inert.noscroll", "show"),
					h.Class(
						"relative max-w-xl min-w-96 rounded-lg border border-gray-100 bg-white p-6 shadow-sm",
					),
					h.Div(
						h.Class("flex items-center justify-between gap-8"),
						h.H2(
							h.ID(titleID),
							h.Class("text-base font-medium"),
							g.Text(title),
						),
						h.Button(
							x.On("click", "destroy()"),
							h.Class(
								"relative rounded-full before:absolute before:-inset-2 before:content-[''] focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 focus-visible:ring-offset-2",
							),
							h.Div(
								h.Class("flex size-10 items-center justify-center"),
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
