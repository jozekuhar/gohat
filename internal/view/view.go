package view

import (
	"encoding/json"
	"fmt"

	gotmpl "tmpl"
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
	data, err := gotmpl.Static.ReadFile("static/dist/manifest.json")
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
	return fmt.Sprintf("/static/%s", vm[key].File)
}

func getCSSPaths(key string) []string {
	paths := []string{}
	for _, file := range vm[key].CSS {
		paths = append(paths, fmt.Sprintf("/static/%s", file))
	}
	return paths
}
