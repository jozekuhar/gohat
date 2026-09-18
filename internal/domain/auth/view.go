package auth

import (
	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/web/view"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

func loginPage() g.Node {
	return view.BlankLayout(
		h.Div(
			h.Class(
				"bg-muted container grid h-svh flex-col items-center justify-center lg:max-w-none lg:px-0",
			),
			h.Div(
				h.Class(
					"mx-auto flex w-full flex-col justify-center space-y-2 sm:w-120 lg:p-8",
				),
				h.Div(
					h.Class("mb-4 flex items-center justify-center"),
					h.Img(
						h.Alt("Artifact"),
						h.Loading("lazy"),
						h.Width("24"),
						h.Height("24"),
						g.Attr("decoding", "async"),
						h.Class("mr-2 dark:hidden"),
						h.Src(
							"https://artifact-nextjs-template.vercel.app/images/logo.svg?dpl=dpl_Gyjss8KqnbavQHpRArQ6yPL3hPJ1",
						),
						h.Style("color: transparent;"),
					),
					h.Img(
						h.Alt(""),
						h.Aria("hidden", "true"),
						h.Loading("lazy"),
						h.Width("24"),
						h.Height("24"),
						g.Attr("decoding", "async"),
						h.Class("mr-2 hidden dark:block"),
						h.Src(
							"https://artifact-nextjs-template.vercel.app/images/logo.svg?dpl=dpl_Gyjss8KqnbavQHpRArQ6yPL3hPJ1",
						),
						h.Style("color: transparent;"),
					),
					h.H1(
						h.Class("text-xl font-medium"),
						g.Text("Artifact Admin Kit"),
					),
				),
				h.Div(
					h.Class("bg-card text-card-foreground rounded-lg border p-6"),
					h.Div(
						h.Class("flex flex-col space-y-2 text-left"),
						h.H1(
							h.Class("text-2xl font-semibold tracking-tight"),
							g.Text("Login"),
						),
						h.P(
							h.Class("text-muted-foreground text-sm"),
							g.Text("Enter your email and password below"),
							h.Br(),
							g.Text("to log into your account"),
						),
					),
					h.Div(
						h.Class("grid gap-6"),
						h.Form(
							hx.Post(routes.HXLogin),
							hx.Swap("none"),
							h.Div(
								h.Class("grid gap-2"),
								h.Div(
									h.Class("space-y-1"),
									h.Label(
										h.Class(
											"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
										),
										h.For("_r_7_-form-item"),
										g.Text("Email"),
									),
									h.Input(
										h.Class(
											"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
										),
										h.Placeholder("name@example.com"),
										h.ID("_r_7_-form-item"),
										h.Aria("describedby", "_r_7_-form-item-description"),
										h.Aria("invalid", "false"),
										h.Value(""),
										h.Name("Email"),
										h.AutoComplete("username"),
									),
								),
								h.Div(
									h.Class("space-y-1"),
									h.Div(
										h.Class("flex items-center justify-between"),
										h.Label(
											h.Class(
												"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
											),
											h.For("_r_8_-form-item"),
											g.Text("Password"),
										),
										h.A(
											h.Class(
												"text-muted-foreground text-sm font-medium hover:opacity-75",
											),
											h.Href("/forgot-password"),
											g.Text("Forgot password?"),
										),
									),
									h.Div(
										x.Data(`{showPassword: false}`),
										x.Cloak(),
										h.Class("relative rounded-md"),
										h.Input(
											x.Bind("type", `showPassword ? "text" : "password"`),
											h.Class(
												"border-input placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-sm shadow-xs file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50",
											),
											h.Placeholder("********"),
											h.ID("_r_8_-form-item"),
											h.Aria("describedby", "_r_8_-form-item-description"),
											h.Aria("invalid", "false"),
											h.Type("password"),
											h.Value(""),
											h.Name("Password"),
											h.AutoComplete("current-password"),
										),
										h.Button(
											x.On("click", "showPassword = !showPassword"),
											h.Class(
												"inline-flex items-center justify-center gap-2 whitespace-nowrap text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground text-muted-foreground absolute top-1/2 right-1 h-6 w-6 -translate-y-1/2 rounded-md",
											),
											h.Type("button"),
											h.SVG(
												x.Show("!showPassword"),
												g.Attr("xmlns", "http://www.w3.org/2000/svg"),
												h.Width("18"),
												h.Height("18"),
												g.Attr("viewBox", "0 0 24 24"),
												g.Attr("fill", "none"),
												g.Attr("stroke", "currentColor"),
												g.Attr("stroke-width", "2"),
												g.Attr("stroke-linecap", "round"),
												g.Attr("stroke-linejoin", "round"),
												h.Class("tabler-icon tabler-icon-eye-off "),
												g.El(
													"path",
													g.Attr(
														"d",
														"M10.585 10.587a2 2 0 0 0 2.829 2.828",
													),
												),
												g.El(
													"path",
													g.Attr(
														"d",
														"M16.681 16.673a8.717 8.717 0 0 1 -4.681 1.327c-3.6 0 -6.6 -2 -9 -6c1.272 -2.12 2.712 -3.678 4.32 -4.674m2.86 -1.146a9.055 9.055 0 0 1 1.82 -.18c3.6 0 6.6 2 9 6c-.666 1.11 -1.379 2.067 -2.138 2.87",
													),
												),
												g.El("path",
													g.Attr("d", "M3 3l18 18"),
												),
											),
											h.SVG(
												x.Show("showPassword"),
												g.Attr("xmlns", "http://www.w3.org/2000/svg"),
												h.Width("18"),
												h.Height("18"),
												g.Attr("viewBox", "0 0 24 24"),
												g.Attr("fill", "none"),
												g.Attr("stroke", "currentColor"),
												g.Attr("stroke-width", "2"),
												g.Attr("stroke-linecap", "round"),
												g.Attr("stroke-linejoin", "round"),
												h.Class("tabler-icon tabler-icon-eye "),
												g.El(
													"path",
													g.Attr(
														"d",
														"M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0",
													),
												),
												g.El(
													"path",
													g.Attr(
														"d",
														"M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6",
													),
												),
											),
										),
									),
								),
								h.Button(
									h.Class(
										"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2 mt-2",
									),
									g.Text("Login"),
								),
								h.Div(
									h.Class("relative my-2"),
									h.Div(
										h.Class("absolute inset-0 flex items-center"),
										h.Span(
											h.Class("w-full border-t"),
										),
									),
									h.Div(
										h.Class("relative flex justify-center text-xs uppercase"),
										h.Span(
											h.Class("bg-background text-muted-foreground px-2"),
											g.Text("Or continue with"),
										),
									),
								),
								h.Div(
									h.Class("flex items-center gap-2"),
									h.Button(
										hx.Get(routes.HXSignInGoogle),
										hx.Swap("none"),
										c.Classes{
											"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 border border-input bg-background shadow-xs hover:bg-accent hover:text-accent-foreground h-9 px-4 py-2 w-full": true,
											"[&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0": true,
										},
										h.Type("button"),
										g.Raw(
											`<svg fill="currentColor" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px" viewBox="0 0 210 210" style="enable-background:new 0 0 210 210;" xml:space="preserve"><path d="M0,105C0,47.103,47.103,0,105,0c23.383,0,45.515,7.523,64.004,21.756l-24.4,31.696C133.172,44.652,119.477,40,105,40 c-35.841,0-65,29.159-65,65s29.159,65,65,65c28.867,0,53.398-18.913,61.852-45H105V85h105v20c0,57.897-47.103,105-105,105 S0,162.897,0,105z"></path></svg>`,
										),
										g.Text("Google"),
									),
									h.A(
										hx.Boost("true"),
										h.Href(routes.Register),
										c.Classes{
											"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 border border-input bg-background shadow-xs hover:bg-accent hover:text-accent-foreground h-9 px-4 py-2 w-full": true,
											"[&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0": true,
										},
										h.Type("button"),
										g.Raw(
											`<svg fill="currentColor" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" id="Capa_1" x="0px" y="0px" viewBox="0 0 512 512" style="enable-background:new 0 0 512 512;" xml:space="preserve"> <path d="M463.748,48.251c-64.336-64.336-169.013-64.335-233.349,0.001c-43.945,43.945-59.209,108.706-40.181,167.461    L4.396,401.536c-2.813,2.813-4.395,6.621-4.395,10.606V497c0,8.291,6.709,15,15,15h84.858c3.984,0,7.793-1.582,10.605-4.395    l21.211-21.226c3.237-3.237,4.819-7.778,4.292-12.334l-2.637-22.793l31.582-2.974c7.178-0.674,12.847-6.343,13.521-13.521    l2.974-31.582l22.793,2.651c4.233,0.571,8.496-0.85,11.704-3.691c3.193-2.856,5.024-6.929,5.024-11.206V363h27.422    c3.984,0,7.793-1.582,10.605-4.395l38.467-37.958c58.74,19.043,122.381,4.929,166.326-39.046    C528.084,217.266,528.084,112.587,463.748,48.251z M421.313,154.321c-17.549,17.549-46.084,17.549-63.633,0    s-17.549-46.084,0-63.633s46.084-17.549,63.633,0S438.861,136.772,421.313,154.321z"/></svg>`,
										),
										g.Text("Password"),
									),
								),
							),
						),
					),
					h.P(
						h.Class("text-muted-foreground mt-4 px-8 text-center text-sm"),
						g.Text("By clicking login, you agree to our "),
						h.A(
							hx.Boost("true"),
							h.Href(routes.TermsOfService),
							h.Class("hover:text-primary underline underline-offset-4"),
							g.Text("Terms of Service"),
						),
						g.Text(" and "),
						h.A(
							hx.Boost("true"),
							h.Href(routes.PrivacyPolicy),
							h.Class("hover:text-primary underline underline-offset-4"),
							g.Text("Privacy Policy"),
						),
						g.Text("."),
					),
				),
			),
		),
	)
}

func registerPage() g.Node {
	return view.BlankLayout(
		h.Div(
			h.Class(
				"bg-muted container grid h-svh flex-col items-center justify-center lg:max-w-none lg:px-0",
			),
			h.Div(
				h.Class(
					"mx-auto flex w-full flex-col justify-center space-y-2 sm:w-120 lg:p-8",
				),
				h.Div(
					h.Class("mb-4 flex items-center justify-center"),
					h.Img(
						h.Alt("Artifact"),
						h.Loading("lazy"),
						h.Width("24"),
						h.Height("24"),
						g.Attr("decoding", "async"),
						h.Class("mr-2 dark:hidden"),
						h.Src(
							"https://artifact-nextjs-template.vercel.app/images/logo.svg?dpl=dpl_Gyjss8KqnbavQHpRArQ6yPL3hPJ1",
						),
						h.Style("color: transparent;"),
					),
					h.Img(
						h.Alt(""),
						h.Aria("hidden", "true"),
						h.Loading("lazy"),
						h.Width("24"),
						h.Height("24"),
						g.Attr("decoding", "async"),
						h.Class("mr-2 hidden dark:block"),
						h.Src(
							"https://artifact-nextjs-template.vercel.app/images/logo.svg?dpl=dpl_Gyjss8KqnbavQHpRArQ6yPL3hPJ1",
						),
						h.Style("color: transparent;"),
					),
					h.H1(
						h.Class("text-xl font-medium"),
						g.Text("Artifact Admin Kit"),
					),
				),
				h.Div(
					h.Class("bg-card text-card-foreground rounded-lg border p-6"),
					h.Div(
						h.Class("mb-2 flex flex-col space-y-2 text-left"),
						h.H1(
							h.Class("text-lg font-semibold tracking-tight"),
							g.Text("Create an account"),
						),
						h.P(
							h.Class("text-muted-foreground text-sm"),
							g.Text("Enter your email and password to create an account."),
							h.Br(),
							g.Text("Already have an account? "),
							h.A(
								hx.Boost("true"),
								h.Class("hover:text-primary underline underline-offset-4"),
								h.Href(routes.Login),
								g.Text("Log In"),
							),
						),
					),
					h.Div(
						h.Class("grid gap-6"),
						h.Form(
							hx.Post(routes.HXRegister),
							hx.Swap("none"),
							h.Div(
								h.Class("grid gap-2"),
								h.Div(
									h.Class("space-y-1"),
									h.Label(
										h.Class(
											"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
										),
										h.For("_r_7_-form-item"),
										g.Text("Email"),
									),
									h.Input(
										h.Class(
											"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
										),
										h.Placeholder("name@example.com"),
										h.ID("_r_7_-form-item"),
										h.Aria("describedby", "_r_7_-form-item-description"),
										h.Aria("invalid", "false"),
										h.Value(""),
										h.Name("Email"),
										h.AutoComplete("username"),
									),
								),
								h.Div(
									h.Class("space-y-1"),
									h.Label(
										h.Class(
											"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
										),
										h.For("_r_8_-form-item"),
										g.Text("Password"),
									),
									h.Div(
										x.Data(`{showPassword: false}`),
										x.Cloak(),
										h.Class("relative rounded-md"),
										h.Input(
											x.Bind("type", `showPassword ? "text" : "password"`),
											h.Class(
												"border-input placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-sm shadow-xs file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50",
											),
											h.Placeholder("********"),
											h.ID("_r_8_-form-item"),
											h.Aria("describedby", "_r_8_-form-item-description"),
											h.Aria("invalid", "false"),
											h.Type("password"),
											h.Value(""),
											h.Name("Password"),
											h.AutoComplete("new-password"),
										),
										h.Button(
											x.On("click", "showPassword = !showPassword"),
											h.Class(
												"inline-flex items-center justify-center gap-2 whitespace-nowrap text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground text-muted-foreground absolute top-1/2 right-1 h-6 w-6 -translate-y-1/2 rounded-md",
											),
											h.Type("button"),
											h.SVG(
												x.Show("!showPassword"),
												g.Attr("xmlns", "http://www.w3.org/2000/svg"),
												h.Width("18"),
												h.Height("18"),
												g.Attr("viewBox", "0 0 24 24"),
												g.Attr("fill", "none"),
												g.Attr("stroke", "currentColor"),
												g.Attr("stroke-width", "2"),
												g.Attr("stroke-linecap", "round"),
												g.Attr("stroke-linejoin", "round"),
												h.Class("tabler-icon tabler-icon-eye-off "),
												g.El(
													"path",
													g.Attr(
														"d",
														"M10.585 10.587a2 2 0 0 0 2.829 2.828",
													),
												),
												g.El(
													"path",
													g.Attr(
														"d",
														"M16.681 16.673a8.717 8.717 0 0 1 -4.681 1.327c-3.6 0 -6.6 -2 -9 -6c1.272 -2.12 2.712 -3.678 4.32 -4.674m2.86 -1.146a9.055 9.055 0 0 1 1.82 -.18c3.6 0 6.6 2 9 6c-.666 1.11 -1.379 2.067 -2.138 2.87",
													),
												),
												g.El("path",
													g.Attr("d", "M3 3l18 18"),
												),
											),
											h.SVG(
												x.Show("showPassword"),
												g.Attr("xmlns", "http://www.w3.org/2000/svg"),
												h.Width("18"),
												h.Height("18"),
												g.Attr("viewBox", "0 0 24 24"),
												g.Attr("fill", "none"),
												g.Attr("stroke", "currentColor"),
												g.Attr("stroke-width", "2"),
												g.Attr("stroke-linecap", "round"),
												g.Attr("stroke-linejoin", "round"),
												h.Class("tabler-icon tabler-icon-eye "),
												g.El(
													"path",
													g.Attr(
														"d",
														"M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0",
													),
												),
												g.El(
													"path",
													g.Attr(
														"d",
														"M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6",
													),
												),
											),
										),
									),
								),
								h.Div(
									h.Class("space-y-1"),
									h.Label(
										h.Class(
											"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
										),
										h.For("_r_9_-form-item"),
										g.Text("Confirm Password"),
									),
									h.Div(
										x.Data(`{showPassword: false}`),
										x.Cloak(),
										h.Class("relative rounded-md"),
										h.Input(
											x.Bind("type", `showPassword ? "text" : "password"`),
											h.Class(
												"border-input placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-sm shadow-xs file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50",
											),
											h.Placeholder("********"),
											h.ID("_r_9_-form-item"),
											h.Aria("describedby", "_r_9_-form-item-description"),
											h.Aria("invalid", "false"),
											h.Type("password"),
											h.Value(""),
											h.Name("ConfirmPassword"),
											h.AutoComplete("new-password"),
										),
										h.Button(
											x.On("click", "showPassword = !showPassword"),
											h.Class(
												"inline-flex items-center justify-center gap-2 whitespace-nowrap text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground text-muted-foreground absolute top-1/2 right-1 h-6 w-6 -translate-y-1/2 rounded-md",
											),
											h.Type("button"),
											h.SVG(
												x.Show("!showPassword"),
												g.Attr("xmlns", "http://www.w3.org/2000/svg"),
												h.Width("18"),
												h.Height("18"),
												g.Attr("viewBox", "0 0 24 24"),
												g.Attr("fill", "none"),
												g.Attr("stroke", "currentColor"),
												g.Attr("stroke-width", "2"),
												g.Attr("stroke-linecap", "round"),
												g.Attr("stroke-linejoin", "round"),
												h.Class("tabler-icon tabler-icon-eye-off "),
												g.El(
													"path",
													g.Attr(
														"d",
														"M10.585 10.587a2 2 0 0 0 2.829 2.828",
													),
												),
												g.El(
													"path",
													g.Attr(
														"d",
														"M16.681 16.673a8.717 8.717 0 0 1 -4.681 1.327c-3.6 0 -6.6 -2 -9 -6c1.272 -2.12 2.712 -3.678 4.32 -4.674m2.86 -1.146a9.055 9.055 0 0 1 1.82 -.18c3.6 0 6.6 2 9 6c-.666 1.11 -1.379 2.067 -2.138 2.87",
													),
												),
												g.El("path",
													g.Attr("d", "M3 3l18 18"),
												),
											),
											h.SVG(
												x.Show("showPassword"),
												g.Attr("xmlns", "http://www.w3.org/2000/svg"),
												h.Width("18"),
												h.Height("18"),
												g.Attr("viewBox", "0 0 24 24"),
												g.Attr("fill", "none"),
												g.Attr("stroke", "currentColor"),
												g.Attr("stroke-width", "2"),
												g.Attr("stroke-linecap", "round"),
												g.Attr("stroke-linejoin", "round"),
												h.Class("tabler-icon tabler-icon-eye "),
												g.El(
													"path",
													g.Attr(
														"d",
														"M10 12a2 2 0 1 0 4 0a2 2 0 0 0 -4 0",
													),
												),
												g.El(
													"path",
													g.Attr(
														"d",
														"M21 12c-2.4 4 -5.4 6 -9 6c-3.6 0 -6.6 -2 -9 -6c2.4 -4 5.4 -6 9 -6c3.6 0 6.6 2 9 6",
													),
												),
											),
										),
									),
								),
								h.Button(
									h.Class(
										"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 h-9 px-4 py-2 mt-2",
									),
									g.Text("Create Account"),
								),
								h.Div(
									h.Class("relative my-2"),
									h.Div(
										h.Class("absolute inset-0 flex items-center"),
										h.Span(
											h.Class("w-full border-t"),
										),
									),
									h.Div(
										h.Class("relative flex justify-center text-xs uppercase"),
										h.Span(
											h.Class("bg-background text-muted-foreground px-2"),
											g.Text("Or continue with"),
										),
									),
								),
								h.Div(
									h.Class("flex items-center gap-2"),
									h.Button(
										hx.Get(routes.HXSignInGoogle),
										hx.Swap("none"),
										h.Class(
											"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 border border-input bg-background shadow-xs hover:bg-accent hover:text-accent-foreground h-9 px-4 py-2 w-full",
										),
										h.Type("button"),
										h.SVG(
											h.Class("size-6"),
											g.Attr("xmlns", "http://www.w3.org/2000/svg"),
											g.Attr("xlink", "http://www.w3.org/1999/xlink"),
											g.Attr("x", "0px"),
											g.Attr("y", "0px"),
											g.Attr("viewBox", "0 0 210 210"),
											h.Style("enable-background:new 0 0 210 210;"),
											g.Attr("space", "preserve"),
											g.El(
												"path",
												g.Attr(
													"d",
													"M0,105C0,47.103,47.103,0,105,0c23.383,0,45.515,7.523,64.004,21.756l-24.4,31.696C133.172,44.652,119.477,40,105,40\n\tc-35.841,0-65,29.159-65,65s29.159,65,65,65c28.867,0,53.398-18.913,61.852-45H105V85h105v20c0,57.897-47.103,105-105,105\n\tS0,162.897,0,105z",
												),
											),
										),
										g.Text("Google"),
									),
								),
							),
						),
					),
					h.P(
						h.Class("text-muted-foreground mt-4 px-8 text-center text-sm"),
						g.Text("By creating an account, you agree to our "),
						h.A(
							hx.Boost("true"),
							h.Href(routes.TermsOfService),
							h.Class("hover:text-primary underline underline-offset-4"),
							g.Text("Terms of Service"),
						),
						g.Text(" and "),
						h.A(
							hx.Boost("true"),
							h.Href(routes.PrivacyPolicy),
							h.Class("hover:text-primary underline underline-offset-4"),
							g.Text("Privacy Policy"),
						),
						g.Text("."),
					),
				),
			),
		),
	)
}

func profileSettingsPage(ident identity.IdentityCtx) g.Node {
	return view.AppLayout(
		ident,
		view.SettingsLayout(
			ident,
			view.TabProfile,
			profileSettingsContent(),
		),
	)
}

func profileSettingsContent() g.Node {
	return h.Div(
		g.Text("Profile settings"),
	)
}
