package view

import (
	"mimokocke/internal/shared/identity"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

const (
	idSidebarPopoverItems = "sidebar_popover_items"
)

func BlankLayout(children ...g.Node) g.Node {
	return baseLayout(
		g.Map(children, func(node g.Node) g.Node {
			return node
		}),
	)
}

func AppLayout(ident identity.IdentityCtx, children ...g.Node) g.Node {
	return baseLayout(
		h.Div(
			x.Data(`{ sidebarOpen: false }`),
			x.Cloak(),
			h.Class("flex flex-1 flex-col"),
			h.Div(
				h.Class("flex min-h-svh w-full"),
				sidebar(ident),
				content(children...),
			),
		),
	)
}

func baseLayout(children ...g.Node) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:       "",
		Description: "",
		Language:    "en",
		Head: g.Group{
			favicons(),
			g.Map(getCSSPaths("view"), func(path string) g.Node {
				return h.Link(h.Rel("stylesheet"), h.Href(path))
			}),
		},
		Body: g.Group{
			h.Class("antialiased"),
			g.Map(children, func(node g.Node) g.Node {
				return node
			}),
			modalContainer(),
			toastContainer(),
			h.Script(h.Type("module"), h.Src(getAssetPath("view"))),
		},
		HTMLAttrs: g.Group{
			x.Data(`{ darkMode: localStorage.getItem('theme') === 'dark' }`),
			x.Cloak(),
			x.Init(
				`$watch('darkMode', value => localStorage.setItem('theme', value ? 'dark' : 'light'))`,
			),
			x.Bind("class", `darkMode && "dark"`),
			h.Class("font-sans"),
		},
	})
}

func content(children ...g.Node) g.Node {
	return h.Div(
		h.Class("flex h-full w-full min-w-0 flex-col"),
		contentHeader(),
		g.Map(children, func(node g.Node) g.Node {
			return node
		}),
	)
}

func contentHeader() g.Node {
	return h.Header(
		h.Class(
			"bg-background grid w-full min-w-0 grid-cols-[auto_minmax(0,1fr)_auto] items-center gap-2 border-b px-4 py-4 sm:gap-3 sm:px-6",
		),
		h.Button(
			x.On("click", `sidebarOpen = !sidebarOpen; console.log(sidebarOpen)`),
			c.Classes{
				"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground size-8 shrink-0": true,
				"md:invisible": true,
			},
			g.Raw(
				`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-panel-left" aria-hidden="true"><rect width="18" height="18" x="3" y="3" rx="2"></rect><path d="M9 3v18"></path></svg>`,
			),
			h.Span(
				h.Class("sr-only"),
				g.Text("Toggle Sidebar"),
			),
		),
		h.Div(
			h.Class("min-w-0"),
			h.H1(
				h.Class("truncate text-base font-medium"),
				g.Text(""),
			),
		),
		h.Button(
			x.On("click", "darkMode = !darkMode; console.log(darkMode)"),
			h.Class(
				"inline-flex items-center justify-center gap-2 whitespace-nowrap text-sm font-medium transition-colors focus-visible:outline-hidden focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 hover:bg-accent hover:text-accent-foreground size-9 scale-100 rounded-lg",
			),
			h.Type("button"),
			g.Raw(
				`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tabler-icon tabler-icon-sun size-[1.2rem] scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90"><path d="M8 12a4 4 0 1 0 8 0a4 4 0 1 0 -8 0"></path><path d="M3 12h1m8 -9v1m8 8h1m-9 8v1m-6.4 -15.4l.7 .7m12.1 -.7l-.7 .7m0 11.4l.7 .7m-12.1 -.7l-.7 .7"></path></svg>`,
			),
			g.Raw(
				`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="tabler-icon tabler-icon-moon absolute size-[1.2rem] scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0"><path d="M12 3c.132 0 .263 0 .393 0a7.5 7.5 0 0 0 7.92 12.446a9 9 0 1 1 -8.313 -12.454l0 .008"></path></svg>`,
			),
			h.Span(
				h.Class("sr-only"),
				g.Text("Toggle theme"),
			),
		),
	)
}
