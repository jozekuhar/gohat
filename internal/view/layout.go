package view

import (
	"tmpl/internal/shared/routes"

	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type layout struct{}

func newLayout() *layout {
	return &layout{}
}

func (v *layout) Guest(children ...g.Node) g.Node {
	return v.base(children...)
}

func (v *layout) Auth(children ...g.Node) g.Node {
	return v.base(
		v.navbar(),
		g.Map(children, func(node g.Node) g.Node { return node }),
	)
}

func (v *layout) base(children ...g.Node) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:       "",
		Description: "",
		Language:    "",
		Head: g.Group{
			g.Map(getCSSPaths("view"), func(path string) g.Node {
				return h.Link(h.Rel("stylesheet"), h.Href(path))
			}),
		},
		Body: g.Group{
			g.Map(children, func(node g.Node) g.Node { return node }),
			h.Script(h.Type("module"), h.Src(getAssetPath("view"))),
		},
		HTMLAttrs: g.Group{
			h.Class("scroll-smooth"),
		},
	})
}

func (v *layout) navbar() g.Node {
	return h.Header(
		h.Class("mx-auto max-w-7xl"),
		h.Nav(
			h.Class("flex justify-end items-center py-2 px-4"),
			h.Ul(
				h.Li(
					h.Button(
						hx.Post(routes.ILogout),
						hx.Swap("none"),
						h.Class(
							"py-1 px-4 text-sm text-white rounded-sm hover:bg-gray-800 bg-gray-950",
						),
						g.Text("Logout"),
					),
				),
			),
		),
	)
}
