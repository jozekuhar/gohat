package view

import (
	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type Counter struct {
	layout *layout
}

func NewCounter() *Counter {
	return &Counter{}
}

func (v *Counter) CounterPage() g.Node {
	return v.layout.Auth(
		h.Main(
			x.Data("CounterComponent"),
			x.Cloak(),
			h.Class("flex flex-col justify-center items-center"),
			h.Div(
				x.Text("count"),
				h.Class("text-7xl"),
			),
			h.Div(
				h.Class("flex gap-4 mt-4"),
				h.Button(
					x.On("click", "inc()"),
					h.Class(
						"py-1 px-4 text-sm text-white rounded-sm hover:bg-gray-800 bg-gray-950",
					),
					g.Text("Inc"),
				),
				h.Button(
					x.On("click", "dec()"),
					h.Class(
						"py-1 px-4 text-sm text-white rounded-sm hover:bg-gray-800 bg-gray-950",
					),
					g.Text("Dec"),
				),
			),
		),
	)
}
