package view

import (
	"fmt"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

const idToastPortal = "toast-portal"

type Toast struct{}

func NewToast() *Toast {
	return &Toast{}
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

func (v *Toast) Fragment(text string) g.Node {
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
				h.Span(g.Text(text)),
				h.Button(
					x.On("click", "destroy()"),
					h.Class(
						"relative rounded-full before:absolute before:-inset-4.5 before:content-[''] focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 focus-visible:ring-offset-2",
					),
					h.Div(
						h.Class("flex size-3 items-center justify-center"),
						h.SVG(
							g.Attr("viewBox", "0 0 512 512"),
							g.El("g",
								g.El(
									"path",
									g.Attr(
										"d",
										"m25 512a25 25 0 0 1 -17.68-42.68l462-462a25 25 0 0 1 35.36 35.36l-462 462a24.93 24.93 0 0 1 -17.68 7.32z",
									),
								),
								g.El(
									"path",
									g.Attr(
										"d",
										"m487 512a24.93 24.93 0 0 1 -17.68-7.32l-462-462a25 25 0 0 1 35.36-35.36l462 462a25 25 0 0 1 -17.68 42.68z",
									),
								),
							),
						),
					),
				),
			),
		),
	)
}
