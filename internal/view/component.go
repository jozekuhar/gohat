package view

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type component struct{}

func newComponent() *component {
	return &component{}
}

func (v *component) primaryButton(children ...g.Node) g.Node {
	return h.Button(
		c.Classes{
			"py-1 px-4 text-sm text-white rounded-sm hover:bg-gray-800 bg-gray-950 active:bg-gray-950":                             true,
			"focus:outline-none focus-visible:z-10 focus-visible:ring-2 focus-visible:ring-indigo-500 focus-visible:ring-offset-2": true,
		},
		g.Map(children, func(child g.Node) g.Node { return child }),
	)
}

func (v *component) outlineButton(children ...g.Node) g.Node {
	return h.Button(
		c.Classes{
			"py-1 px-4 text-sm text-gray-950 rounded-sm hover:bg-gray-50 border border-gray-950 active:bg-transparent":             true,
			"focus:outline-none focus-visible:z-10 focus-visible:ring-2 focus-visible:ring-indigo-500 focus-visible:ring-offset-2": true,
		},
		g.Map(children, func(child g.Node) g.Node { return child }),
	)
}
