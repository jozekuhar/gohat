package view

import (
	"fmt"

	"gohat/internal/shared/routes"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

type Counter struct {
	layout    *layout
	component *component
}

func NewCounter() *Counter {
	return &Counter{
		layout:    newLayout(),
		component: newComponent(),
	}
}

func (v *Counter) CounterPage() g.Node {
	return v.layout.Auth(
		h.Main(
			x.Data("counter"),
			x.Cloak(),
			h.Class("flex flex-col justify-center items-center"),
			v.component.primaryButton(
				hx.Get(routes.IModal),
				hx.Target(fmt.Sprintf("#%s", IDModalContent)),
				hx.Swap("innerHTML"),
				h.Type("button"),
				g.Text("Open modal"),
			),
			v.component.outlineButton(g.Text("hello")),
			h.Div(
				h.Class("mt-4 h-screen"),
				h.Div(
					x.Text("count"),
					h.Class("text-7xl rounded-sm border border-gray-950"),
				),
				h.Div(
					h.Class("flex gap-4 mt-4"),
					v.component.primaryButton(
						x.On("click", "inc()"),
						g.Text("Increment"),
					),
					v.component.primaryButton(
						x.On("click", "dec()"),
						g.Text("Decrement"),
					),
				),
			),
		),
	)
}

func (v *Counter) CounterModalContent() g.Node {
	return h.Div(
		h.Class("flex flex-col gap-6"),
		h.P(
			g.Text("Counter modal content"),
		),

		h.Div(
			h.Class("flex gap-2 justify-end"),
			v.component.outlineButton(
				x.On("click", fmt.Sprintf("$dispatch('%s')", EventCloseModal)),
				g.Text("Cancel"),
			),
			v.component.primaryButton(g.Text("Submit")),
		),
	)
}
