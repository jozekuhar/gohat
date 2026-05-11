package view

import (
	"tmpl/internal/shared/routes"

	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

type Auth struct {
	layout *layout
}

func NewAuth() *Auth {
	return &Auth{
		layout: newLayout(),
	}
}

func (v *Auth) LoginPage() g.Node {
	return v.layout.Guest(
		h.Main(
			h.Class("flex justify-center items-center h-screen"),
			h.Button(
				hx.Get(routes.ILoginGoogle),
				hx.Swap("none"),
				h.Class("py-1 px-4 text-sm text-white rounded-sm hover:bg-gray-800 bg-gray-950"),
				g.Text("Login with Google"),
			),
		),
	)
}
