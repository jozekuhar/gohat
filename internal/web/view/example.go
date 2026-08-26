package view

import (
	"fmt"
	"time"

	"mimokocke/internal/shared/routes"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

const idSSEContent = "sse-content"

type Example struct {
	layout    *layout
	modal     *modal
	component *component
}

func NewExample() *Example {
	return &Example{
		layout:    newLayout(),
		modal:     newModal(),
		component: newComponent(),
	}
}

func (v *Example) ExamplePage(currentTime time.Time) g.Node {
	return v.layout.Auth(
		h.Main(
			h.Class("mx-auto max-w-7xl space-y-4 px-6"),
			h.Div(
				h.Span(g.Textf("Current date is: %s", currentTime.Format("02.01.2006"))),
			),
			h.Div(
				h.Class("flex gap-4"),
				v.component.primaryButton(
					"Open modal",
					hx.Get(routes.IExampleModal),
					hx.Swap("none"),
					h.Type("button"),
				),
				v.component.primaryButton(
					"Long request",
					hx.Get(routes.IExampleLongRequest),
					hx.Swap("none"),
					h.Type("button"),
				),
				v.component.primaryButton(
					"SSE Request",
					g.Attr("hx-sse:connect", routes.IExampleSSERequest),
					hx.Target(fmt.Sprintf("#%s", idSSEContent)),
					hx.Swap("append"),
					h.Type("button"),
				),
			),
			v.exampleSSEContent(),
		),
	)
}

func (v *Example) ExampleModalContent() g.Node {
	return v.modal.fragment(
		"This is modal title",
		h.Div(
			h.Class("mt-4 space-y-4"),
			h.P(
				h.Class("text-sm text-gray-950"),
				g.Text(
					"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
				),
			),
			h.Div(
				h.Class("flex justify-end gap-2"),
				v.component.outlineButton(
					"Cancel",
					v.modal.onClose(),
				),
				v.component.primaryButton(
					"Submit",
					x.On("click", "alert('submit')"),
				),
			),
		),
	)
}

func (v *Example) exampleSSEContent() g.Node {
	return h.Div(
		h.ID(idSSEContent),
	)
}
