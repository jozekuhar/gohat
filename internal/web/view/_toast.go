package view

import (
	"fmt"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

const idToastPortal = "toast-portal"

type Toast struct {
	icon *icon
}

func NewToast() *Toast {
	return &Toast{
		icon: newIcon(),
	}
}

func (v *Toast) container() g.Node {
	return h.Div(
		h.ID(idToastPortal),
		h.Class(
			"fixed right-0 bottom-0 flex max-w-sm flex-col items-end gap-4 pr-4 pb-4 sm:justify-start",
		),
		h.Role("status"),
		h.Aria("live", "polite"),
	)
}

func (v *Toast) FragmentText(title string, text string) g.Node {
	return v.Fragment(title, g.Text(text))
}

func (v *Toast) Fragment(title string, content g.Node) g.Node {
	return g.El(
		"hx-partial",
		hx.Target(fmt.Sprintf("#%s", idToastPortal)),
		hx.Swap("append"),
		h.Div(
			x.Cloak(),
			x.Data("toast"),
			x.Show("show"),
			x.Transition(".duration.500ms"),
			h.Class("w-fit space-y-2 rounded-lg border border-gray-100 bg-white p-6 shadow-sm"),
			h.Div(
				h.Class("flex items-center justify-between gap-8"),
				h.H1(
					h.Class("text-sm font-medium"),
					g.Text(title),
				),
				h.Button(
					x.On("click", "destroy()"),
					h.Class(
						"relative rounded-full before:absolute before:-inset-4.5 before:content-[''] focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 focus-visible:ring-offset-2",
					),
					h.Div(
						h.Class("flex size-4 items-center justify-center"),
						v.icon.xMark("size-3 fill-gray-400"),
					),
				),
			),
			content,
		),
	)
}
