package handler

import (
	"net/http"

	"gohat"
)

type Static struct{}

func NewStatic() *Static {
	return &Static{}
}

func (v *Static) GetStatic(w http.ResponseWriter, r *http.Request) {
	http.FileServerFS(gohat.Static).ServeHTTP(w, r)
}
