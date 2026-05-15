package handler

import (
	"net/http"

	"gohat/internal/shared/web"
	"gohat/internal/view"

	hx "maragu.dev/gomponents-htmx/http"
)

type Example struct {
	exampleView *view.Example
}

func NewExample() *Example {
	return &Example{
		exampleView: view.NewExample(),
	}
}

func (h *Example) GetExample(w http.ResponseWriter, r *http.Request) {
	web.Render(w, h.exampleView.ExamplePage())
}

func (h *Example) GetExampleModal(w http.ResponseWriter, r *http.Request) {
	hx.SetTrigger(w.Header(), view.EventOpenModal)
	web.Render(w, h.exampleView.ExampleModalContent())
}

func (h *Example) GetExampleToast(w http.ResponseWriter, r *http.Request) {
	web.Render(w, h.exampleView.ExampleToastContent())
}
