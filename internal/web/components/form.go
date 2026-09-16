package components

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	x "github.com/glsubri/gomponents-alpine"
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type InputParams struct {
	Label        string
	Name         string
	Value        string
	Placeholder  string
	Disabled     bool
	AutoFocus    bool
	AutoComplete string
	Form         string
}

func Input(params InputParams) g.Node {
	id := generateID()

	return h.Div(
		h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
		h.Label(
			h.For(id),
			h.Class(
				"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
			),
			g.Text(params.Label),
		),
		h.Input(
			h.ID(id),
			h.Class(
				"border-input file:text-foreground placeholder:text-muted-foreground focus-visible:ring-ring flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium focus-visible:ring-1 focus-visible:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 md:text-sm col-span-4",
			),
			g.If(params.Name != "", h.Name(params.Name)),
			g.If(params.Value != "", h.Value(params.Value)),
			g.If(params.Placeholder != "", h.Placeholder(params.Placeholder)),
			g.If(params.Disabled, h.Disabled()),
			g.If(params.AutoFocus, h.AutoFocus()),
			g.If(params.AutoComplete != "", h.AutoComplete(params.AutoComplete)),
			g.If(params.Form != "", h.FormAttr(params.Form)),
		),
	)
}

type SelectParams struct {
	Label       string
	Name        string
	Value       string
	Placeholder string
	XModel      string
	Options     []SelectOption
}

type SelectOption struct {
	Value string `json:"value"`
	Text  string `json:"text"`
}

func Select(params SelectParams) g.Node {
	// id := generateID()
	options, _ := json.Marshal(params.Options)

	return h.Div(
		x.Data(fmt.Sprintf(`{ value: null, options: %s }`, options)),
		x.Cloak(),
		x.Modelable("value"),
		x.Model(params.XModel),
		g.Attr("x-listbox"),
		h.Class("grid grid-cols-6 items-center space-y-0 gap-x-4 gap-y-1"),
		h.Input(
			h.Class("hidden"),
			x.Bind("value", "value"),
			h.Name(params.Name),
			// g.If(params.Value != "", h.Value(params.Value)),
		),
		h.Label(
			// h.For(id),
			h.Class(
				"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 col-span-2 text-right",
			),
			g.Text(params.Label),
		),
		h.Div(
			h.Class("relative col-span-4"),
			h.Button(
				g.Attr("x-listbox:button"),
				// h.ID(id),
				h.Type("button"),
				h.Class(
					"border-input ring-offset-background placeholder:text-muted-foreground focus:ring-ring flex h-9 w-full items-center justify-between rounded-md border bg-transparent px-3 py-2 text-sm whitespace-nowrap shadow-xs focus:ring-1 focus:outline-hidden disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1",
				),
				h.Span(
					x.Text(
						fmt.Sprintf(
							`options.find(opt => opt.value === value)?.text || %q`,
							params.Placeholder,
						),
					),
					h.Style("pointer-events: none;"),
				),
				g.Raw(
					`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-chevron-down h-4 w-4 opacity-50" aria-hidden="true"><path d="m6 9 6 6 6-6"></path></svg>`,
				),
			),

			h.Ul(
				g.Attr("x-listbox:options"),
				x.Cloak(),
				h.Class(
					"absolute mt-1 bg-popover text-popover-foreground z-50 max-h-96 min-w-[8rem] overflow-hidden rounded-md border shadow-md p-1 w-full",
					// data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2
					// h-[var(--radix-select-trigger-height)] min-w-[var(--radix-select-trigger-width)]
				),
				h.Div(
					h.Class(
						"p-1 h-[var(--radix-select-trigger-height)] w-full min-w-[var(--radix-select-trigger-width)]",
					),
					h.Template(
						x.For("option in options"),
						x.Bind("key", "option.value"),
						h.Li(
							g.Attr("x-listbox:option"),
							x.Bind("value", "option.value"),
							c.Classes{
								"relative flex w-full cursor-default items-center rounded-sm py-1.5 pr-8 pl-2 text-sm outline-hidden select-none data-disabled:pointer-events-none data-disabled:opacity-50": true,
								"hocus:bg-accent hocus:text-accent-foreground": true,
							},
							h.Span(
								x.Text("option.text"),
							),
							h.Span(
								x.Show("option.value == value"),
								h.Class(
									"absolute right-2 flex h-3.5 w-3.5 items-center justify-center",
								),
								g.Raw(
									`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-check h-4 w-4" aria-hidden="true"><path d="M20 6 9 17l-5-5"></path></svg>`,
								),
							),
						),
					),
				),
			),
		),
	)
}

type CheckboxParams struct {
	Label        string
	Name         string
	Value        string
	AutoComplete string
}

func Checkbox(params CheckboxParams) g.Node {
	return h.Label(
		h.Class("flex items-center gap-2 px-2 hover:bg-gray-100 dark:hover:bg-gray-800"),
		h.Input(
			h.Type("checkbox"),
			h.Class("peer sr-only"),
			g.If(params.Name != "", h.Name(params.Name)),
			g.If(params.Value != "", h.Value(params.Value)),
			g.If(params.AutoComplete != "", h.AutoComplete(params.AutoComplete)),
		),
		h.Button(
			h.Type("button"),
			h.Role("checkbox"),
			c.Classes{
				"border-primary focus-visible:ring-ring h-4 w-4 shrink-0 rounded-sm border shadow-sm focus-visible:ring-1 focus-visible:outline-hidden ": true,
				"peer-checked:bg-primary peer-checked:text-primary-foreground":                                                                           true,
				"disabled:cursor-not-allowed disabled:opacity-50":                                                                                        true,
			},
		),
		h.Span(
			c.Classes{
				"text-sm font-medium leading-none flex h-full flex-1 cursor-pointer items-center justify-between py-2": true,
				"peer-disabled:cursor-not-allowed peer-disabled:opacity-70":                                            true,
			},
			h.P(
				h.Class("text-xs"),
				g.Text(params.Label),
			),
		),
	)
}

func generateID() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return "id_" + hex.EncodeToString(bytes)
}
