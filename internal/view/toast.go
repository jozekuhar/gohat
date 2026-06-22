package view

import (
	"fmt"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

const idToastPortal = "toast-portal"

type toast struct {
	icon *icon
}

func newToast() *toast {
	return &toast{
		icon: newIcon(),
	}
}

func (v *toast) container() g.Node {
	return h.Div(
		h.ID(idToastPortal),
		h.Class("fixed right-0 bottom-0 pr-4 pb-4 space-y-4 max-w-sm sm:justify-start"),
		h.Role("status"),
		h.Aria("live", "polite"),
	)
}

func (v *toast) fragment(title string, content g.Node) g.Node {
	return g.El(
		"hx-partial",
		hx.Target(fmt.Sprintf("#%s", idToastPortal)),
		hx.Swap("append"),
		h.Div(
			x.Cloak(),
			x.Data("toast"),
			x.Show("show"),
			x.Transition(".duration.500ms"),
			h.Class("p-6 space-y-2 bg-white rounded-lg border border-gray-100 shadow-sm"),
			h.Div(
				h.Class("flex gap-8 justify-between items-center"),
				h.H1(
					h.Class("text-sm font-medium"),
					g.Text(title),
				),
				h.Button(
					x.On("click", "destroy()"),
					h.Class(
						"relative rounded-full focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 focus-visible:ring-offset-2 before:absolute before:content-[''] before:-inset-4.5",
					),
					h.Div(
						h.Class("flex justify-center items-center size-4"),
						v.icon.xMark("size-3 fill-gray-400"),
					),
				),
			),
			content,
		),
	)
}
