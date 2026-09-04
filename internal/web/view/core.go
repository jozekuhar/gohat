package view

import (
	"mimokocke/internal/shared/routes"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

type Core struct {
	layout *Layout
}

func NewCore() *Core {
	return &Core{
		layout: NewLayout(),
	}
}

func (v *Core) NotFoundPage() g.Node {
	return v.layout.blank(
		h.Div(
			h.Class("h-svh"),
			h.Div(
				h.Class("m-auto flex h-full w-full flex-col items-center justify-center gap-2"),
				h.H1(
					h.Class("text-[7rem] leading-tight font-bold"),
					g.Text("404"),
				),
				h.Span(
					h.Class("font-medium"),
					g.Text("Oops! Page Not Found!"),
				),
				h.P(
					h.Class("text-muted-foreground text-center"),
					g.Text("It seems like the page you're looking for"),
					h.Br(),
					g.Text("does not exist or might have been removed."),
				),
				h.Div(
					h.Class("mt-6 flex gap-4"),
					h.Button(
						// TODO(jozekuhar): back as anchor (htmx?) boost must work
						x.Data(""),
						x.On("click", "history.back()"),
						h.Class(
							"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 border border-input bg-background shadow-xs hover:bg-accent hover:text-accent-foreground h-9 px-4 py-2",
						),
						g.Text("Go Back"),
					),
					h.A(
						hx.Boost("true"),
						h.Class(
							"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2",
						),
						h.Href(routes.Root),
						g.Text("Back to Home"),
					),
				),
			),
		),
	)
}

func (v *Core) PrivacyPolicyPage() g.Node {
	return v.layout.blank(
		h.Div(
			h.Class("mx-auto max-w-prose space-y-4 p-8"),
			h.H1(
				h.Class("text-2xl font-semibold"),
				g.Text("Privacy Policy"),
			),
			h.P(
				h.Class("text-muted-foreground text-sm"),
				g.Text(
					"This is placeholder copy for the admin template. Replace with your real privacy policy before production.",
				),
			),
			h.P(
				h.A(
					hx.Boost("true"),
					h.Class("text-primary underline underline-offset-4"),
					h.Href(routes.Login),
					g.Text("Back to login"),
				),
			),
		),
	)
}

func (v *Core) TermsOfServicePage() g.Node {
	return v.layout.blank(
		h.Div(
			h.Class("mx-auto max-w-prose space-y-4 p-8"),
			h.H1(
				h.Class("text-2xl font-semibold"),
				g.Text("Terms of Service"),
			),
			h.P(
				h.Class("text-muted-foreground text-sm"),
				g.Text(
					"This is placeholder copy for the admin template. Replace with your real terms before production.",
				),
			),
			h.P(
				h.A(
					hx.Boost("true"),
					h.Class("text-primary underline underline-offset-4"),
					h.Href(routes.Login),
					g.Text("Back to login"),
				),
			),
		),
	)
}
