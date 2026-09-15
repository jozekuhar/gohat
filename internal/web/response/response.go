package response

import (
	"log"
	"net/http"

	g "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx/http"
)

func Status(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func Render(w http.ResponseWriter, n g.Node) {
	if err := n.Render(w); err != nil {
		log.Panicf("render node")
	}
}

func RenderStatus(w http.ResponseWriter, n g.Node, statusCode int) {
	w.WriteHeader(statusCode)
	if err := n.Render(w); err != nil {
		log.Panicf("render node")
	}
}

func HXRefresh(w http.ResponseWriter) {
	hx.SetRefresh(w.Header())
}

func HXTrigger(w http.ResponseWriter, s string) {
	hx.SetTrigger(w.Header(), s)
}

func HXRedirect(w http.ResponseWriter, s string) {
	hx.SetRedirect(w.Header(), s)
}
