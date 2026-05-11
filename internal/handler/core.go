package handler

import (
	"net/http"

	gotmpl "tmpl"
)

type Core struct{}

func NewCore() *Core {
	return &Core{}
}

func (v *Core) GetStatic(w http.ResponseWriter, r *http.Request) {
	http.FileServerFS(gotmpl.Static).ServeHTTP(w, r)
}
