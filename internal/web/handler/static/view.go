package staticweb

import (
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/web/view"

	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	h "maragu.dev/gomponents/html"
)

func privacyPolicyPage() g.Node {
	return view.BlankLayout(
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

func termsOfServicePage() g.Node {
	return view.BlankLayout(
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
