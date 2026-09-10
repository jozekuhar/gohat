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
			"fixed right-0 bottom-0 z-200 flex max-w-sm flex-col items-end gap-4 pr-4 pb-4 sm:justify-start",
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
			h.Class("flex gap-4 bg-background w-fit rounded-lg border p-6 shadow-lg"),
			h.P(
				h.Class("w-fit text-pretty md:text-sm"),
				g.Text(text),
			),
			h.Button(
				x.On("click", "destroy()"),
				h.Class(
					"h-fit ring-offset-background focus:ring-ring rounded-xs opacity-70 transition-opacity hover:opacity-100 focus:ring-2 focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
				),
				g.Raw(
					`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-x" aria-hidden="true"><path d="M18 6 6 18"></path><path d="m6 6 12 12"></path></svg>`,
				),
			),
		),
	)
}
