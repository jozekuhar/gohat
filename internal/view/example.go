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
				hx.Get(routes.IExampleModal),
				hx.Swap("none"),
				h.Type("button"),
				g.Text("Open modal"),
			),
			v.component.primaryButton(
				hx.Get(routes.IExampleToast),
				hx.Swap("none"),
				h.Type("button"),
				g.Text("Create toast"),
			),
		),
	)
}

func (v *Example) ExampleModalContent() g.Node {
	return v.modal.content(
		"This is title",
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
					v.modal.xOnClickClose(),
					g.Text("Cancel"),
				),
				v.component.primaryButton(
					x.On("click", "alert('submit')"),
					g.Text("Submit"),
				),
			),
		),
	)
}

func (v *Example) ExampleToastContent() g.Node {
	return v.toast.content(
		"Task succesfully shipped",
		h.Div(
			h.Class("text-sm"),
			g.Text(
				"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
			),
		),
	)
}
