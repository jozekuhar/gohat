package view

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type icon struct{}

func newIcon() *icon {
	return &icon{}
}

func (v *icon) xMark(class string) g.Node {
	return h.SVG(
		h.Class(class),
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
	)
}
