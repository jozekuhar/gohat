package handler

import (
	"net/http"

	"mimokocke"
)

func GetStatic(w http.ResponseWriter, r *http.Request) {
	http.FileServerFS(mimokocke.Static).ServeHTTP(w, r)
}
