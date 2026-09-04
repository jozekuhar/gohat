package view

import (
	"fmt"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

const idModalPortal = "modal-portal"

type modal struct{}

func newModal() *modal {
	return &modal{}
}

func (v *modal) container() g.Node {
	return h.Div(
		h.ID(idModalPortal),
	)
}

func (v *modal) fragment(nodes ...g.Node) g.Node {
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

			h.Class("fixed inset-0 overflow-y-auto"),
			h.Role("dialog"),
			h.Aria("modal", "true"),
			h.Aria("labelledby", titleID),
			h.Div(
				x.Show("show"),
				x.Transition(".opacity"),
				h.Class("fixed inset-0 bg-black/80"),
				// z-50 je bil na template
			),

			h.Div(
				x.Show("show"),
				x.Transition(),
				h.Class("relative flex min-h-screen items-center justify-center p-4"),
				h.Div(
					x.On("click.outside", "destroy()"),
					g.Attr("x-trap.inert.noscroll", "show"),
					c.Classes{
						"relative bg-background grid w-full max-w-[calc(100%-2rem)] gap-4 rounded-lg border p-6 shadow-lg duration-200 sm:max-w-lg": true,
						"": true,
					},
					h.Button(
						x.On("click", "destroy()"),
						c.Classes{
							"ring-offset-background focus:ring-ring absolute top-4 right-4 rounded-xs opacity-70 transition-opacity hover:opacity-100 focus:ring-2 focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4": true,
							"bg-accent text-muted-foreground": false,
						},
						g.Raw(
							`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-x" aria-hidden="true"><path d="M18 6 6 18"></path><path d="m6 6 12 12"></path></svg>`,
						),
						h.Span(h.Class("sr-only"), g.Text("Close")),
					),
					g.Map(nodes, func(node g.Node) g.Node {
						return node
					}),
				),
			),
		),
	)
}

// func (v *modal) onClose() g.Node {
// 	return x.On("click", "destroy()")
// }
