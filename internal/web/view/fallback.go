package view

import (
	g "maragu.dev/gomponents"
)

type Fallback struct {
	layout *layout
}

func NewFallback() *Fallback {
	return &Fallback{
		layout: newLayout(),
	}
}

func (v *Fallback) NotFoundPage() g.Node {
	return v.layout.Guest(g.Text("(404) Page not found"))
}

func (v *Fallback) MethodNotAllowedPage() g.Node {
	return v.layout.Guest(g.Text("(405) Method not allowed"))
}
