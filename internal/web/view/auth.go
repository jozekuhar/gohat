package view

import (
	"mimokocke/internal/shared/routes"

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
			h.Class("flex flex-col h-screen items-center justify-center"),
			h.Form(
				hx.Post(routes.HXLogin),
				hx.Swap("none"),
				h.Class("flex flex-col border"),
				h.Input(
					h.Type("email"),
					h.Name("Email"),
					h.Placeholder("Email"),
				),
				h.Input(
					h.Type("password"),
					h.Name("Password"),
					h.Placeholder("Password"),
				),
				h.Button(
					g.Text("Login"),
				),
			),
			h.Button(
				hx.Get(routes.HXSignInGoogle),
				hx.Swap("none"),
				g.Text("Login with Google"),
			),
			h.A(
				h.Class("text-blue-500"),
				h.Href(routes.Register),
				g.Text("Register"),
			),
		),
	)
}

func (v *Auth) RegisterPage() g.Node {
	return v.layout.Guest(
		h.Main(
			h.Class("flex flex-col h-screen items-center justify-center"),
			h.Form(
				hx.Post(routes.HXRegister),
				hx.Swap("none"),
				h.Class("border flex flex-col"),
				h.Input(
					h.Type("email"),
					h.Name("Email"),
					h.Placeholder("Email"),
				),
				h.Input(
					h.Type("password"),
					h.Name("Password"),
					h.Placeholder("Password"),
				),
				h.Input(
					h.Type("password"),
					h.Name("ConfirmPassword"),
					h.Placeholder("Confirm password"),
				),
				h.Button(
					g.Text("Register"),
				),
			),
			h.Button(
				hx.Get(routes.HXSignInGoogle),
				hx.Swap("none"),
				g.Text("Login with Google"),
			),
			h.A(
				h.Class("text-blue-500"),
				h.Href(routes.Login),
				g.Text("Login instead"),
			),
		),
	)
}
