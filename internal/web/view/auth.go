package view

import (
	"mimokocke/internal/shared/routes"

	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

type Auth struct {
	layout    *layout
	component *component
}

func NewAuth() *Auth {
	return &Auth{
		layout:    newLayout(),
		component: newComponent(),
	}
}

func (v *Auth) LoginPage() g.Node {
	return v.layout.Guest(
		h.Main(
			h.Class("flex h-screen items-center justify-center"),
			v.component.primaryButton(
				"Login with Google",
				hx.Get(routes.ILoginGoogle),
				hx.Swap("none"),
			),
		),
	)
}
