package view

import (
	"encoding/json"
	"fmt"

	"gohat"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type viteEntry struct {
	File    string   `json:"file"`
	Name    string   `json:"name"`
	Src     string   `json:"src"`
	IsEntry bool     `json:"isEntry"`
	CSS     []string `json:"css"`
}

type viteManifest map[string]viteEntry

var vm = make(viteManifest)

func LoadManifest() error {
	data, err := gohat.Static.ReadFile("static/dist/manifest.json")
	if err != nil {
		return err
	}

	var tmp viteManifest
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	for _, entry := range tmp {
		vm[entry.Name] = entry
	}

	return nil
}

func getAssetPath(key string) string {
	return fmt.Sprintf("/static/dist/%s", vm[key].File)
}

func getCSSPaths(key string) []string {
	paths := []string{}
	for _, file := range vm[key].CSS {
		paths = append(paths, fmt.Sprintf("/static/dist/%s", file))
	}
	return paths
}

func favicons() g.Node {
	return g.Group{
		h.Link(
			h.Rel("icon"),
			h.Type("image/png"),
			h.Href("/static/favicon/favicon-96x96.png"),
			g.Attr("sizes", "96x96"),
		),
		h.Link(
			h.Rel("icon"),
			h.Type("image/svg+xml"),
			h.Href("/static/favicon/favicon.svg"),
		),
		h.Link(
			h.Rel("shortcut icon"),
			h.Href("/static/favicon/favicon.ico"),
		),
		h.Link(
			h.Rel("apple-touch-icon"),
			g.Attr("sizes", "180x180"),
			h.Href("/static/favicon/apple-touch-icon.png"),
		),
		h.Meta(
			h.Name("apple-mobile-web-app-title"),
			h.Content("Gohat"),
		),
		h.Link(
			h.Rel("manifest"),
			h.Href("/static/favicon/site.webmanifest"),
		),
	}
}
