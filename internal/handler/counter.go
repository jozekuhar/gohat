package handler

import (
	"net/http"

	"gohat/internal/shared/web"
	"gohat/internal/view"

	hx "maragu.dev/gomponents-htmx/http"
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

func (h *Counter) GetCounterModal(w http.ResponseWriter, r *http.Request) {
	jsonto := `{"open-modal": {"title": "Counter"}}`

	hx.SetTrigger(w.Header(), jsonto)
	web.Render(w, h.counterView.CounterModalContent())
}
