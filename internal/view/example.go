package view

import (
	"fmt"
	"time"

	"gohat/internal/shared/routes"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

const idSSEContent = "sse-content"

type Example struct {
	layout    *layout
	modal     *modal
	toast     *toast
	component *component
}

func NewExample() *Example {
	return &Example{
		layout:    newLayout(),
		modal:     newModal(),
		toast:     newToast(),
		component: newComponent(),
	}
}

func (v *Example) ExamplePage(currentTime time.Time) g.Node {
	fmt.Println(currentTime)
	return v.layout.Auth(
		h.Main(
			h.Class("px-6 mx-auto space-y-4 max-w-7xl"),
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
					"Create toast",
					hx.Get(routes.IExampleToast),
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
				h.Class("flex gap-2 justify-end"),
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

func (v *Example) ExampleToastContent() g.Node {
	return v.toast.fragment(
		"Task succesfully shipped",
		h.Div(
			h.Class("text-sm"),
			g.Text(
				"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			),
		),
	)
}

func (v *Example) exampleSSEContent() g.Node {
	return h.Div(
		h.ID(idSSEContent),
	)
}
