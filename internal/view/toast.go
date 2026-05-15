package view

import (
	"fmt"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

const idToasts = "toasts"

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
		h.ID(idToasts),
		h.Class("fixed right-0 bottom-0 pr-4 pb-4 space-y-4 max-w-sm sm:justify-start"),
		h.Role("status"),
		h.Aria("live", "polite"),
	)
}

func (v *toast) content(title string, content g.Node) g.Node {
	return g.El(
		"hx-partial",
		hx.Target(fmt.Sprintf("#%s", idToasts)),
		hx.Swap("append"),
		h.Div(
			x.Data("toast"),
			x.Show("show"),
			x.Transition(".duration.500ms"),
			h.Class("p-4 space-y-2 bg-white rounded-lg border border-gray-200 shadow-sm"),
			h.Div(
				h.Class("flex gap-8 justify-between items-center font-medium"),
				h.H1(
					h.Class("text-sm"),
					g.Text(title),
				),
				h.Button(
					x.On("click", "remove()"),
					v.icon.xMark("size-2.5"),
				),
			),
			content,
		),
	)
}
