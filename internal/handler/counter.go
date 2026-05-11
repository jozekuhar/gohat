package handler

import (
	"net/http"

	"tmpl/internal/shared/web"
	"tmpl/internal/view"
)

type Counter struct {
	counterView *view.Counter
}

func NewCounter() *Counter {
	return &Counter{
		counterView: view.NewCounter(),
	}
}

func (h *Counter) GetCounter(w http.ResponseWriter, r *http.Request) {
	web.Render(w, h.counterView.CounterPage())
}
