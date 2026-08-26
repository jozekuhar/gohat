package handler

import (
	"net/http"

	"mimokocke"
)

type Static struct{}

func NewStatic() *Static {
	return &Static{}
}

func (v *Static) GetStatic(w http.ResponseWriter, r *http.Request) {
	http.FileServerFS(mimokocke.Static).ServeHTTP(w, r)
}
