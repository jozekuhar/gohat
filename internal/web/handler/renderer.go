package handler

import (
	"log"
	"net/http"

	g "maragu.dev/gomponents"
)

func render(w http.ResponseWriter, n g.Node) {
	if err := n.Render(w); err != nil {
		log.Panicf("rendering node")
	}
}
