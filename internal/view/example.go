package view

import (
	"gohat/internal/shared/routes"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

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

func (v *Example) ExamplePage() g.Node {
	return v.layout.Auth(
		h.Main(
			h.Class("flex gap-4 px-6 mx-auto max-w-7xl"),
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
