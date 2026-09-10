package handler

import (
	"log"
	"net/http"

	g "maragu.dev/gomponents"
)

func render(w http.ResponseWriter, n g.Node, statusCode int) {
	w.WriteHeader(statusCode)
	if err := n.Render(w); err != nil {
		log.Panicf("rendering node")
	}
}
