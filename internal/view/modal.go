package view

import (
	"fmt"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

// TODO consts needs to be done
const (
	IDModalContent  = "modal-content"
	EventOpenModal  = "open-modal"
	EventCloseModal = "close-modal"
)

type modal struct{}

func newModal() *modal { return &modal{} }

func (v *modal) container() g.Node {
	return h.Div(
		x.Data("modal"),
		x.Show("isOpen"),
		x.On(fmt.Sprintf("%s.window", EventOpenModal), "open($event)"),
		x.On(fmt.Sprintf("%s.window", EventCloseModal), "close()"),
		x.On("keydown.escape.window.prevent.stop", "close()"),
		h.Style("display: none"),
		h.Class("overflow-y-auto fixed inset-0"),
		h.Div(
			x.Show("isOpen"),
			x.Transition(".opacity"),
			h.Class("fixed inset-0 bg-gray-950/15"),
		),
		h.Div(
			x.Show("isOpen"),
			x.Transition(),
			h.Class(
				"flex relative justify-center items-center p-4 min-h-screen",
			),
			h.Div(
				x.On("click.outside", "close()"),
				g.Attr("x-trap.inert.noscroll", "isOpen"),
				h.Class(
					"relative p-6 pb-8 space-y-3 max-w-xl bg-white rounded-sm shadow-lg min-w-96",
				),
				h.Div(
					h.Class("flex justify-between items-center"),
					h.H1(
						x.Text("title"),
						h.Class("text-lg"),
					),
					h.Button(
						h.Type("button"),
						x.On("click", "close()"),
						h.Class(
							"flex justify-center items-center rounded-full hover:bg-gray-100 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-600 focus-visible:ring-offset-2 size-10",
						),
						h.SVG(
							h.Class("size-4 fill-gray-950"),
							g.Attr("viewBox", "0 0 512 512"),
							g.El(
								"g",
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
				h.Div(
					h.ID(IDModalContent),
				),
			),
		),
	)
}
