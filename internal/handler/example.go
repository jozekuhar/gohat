package handler

import (
	"log/slog"
	"net/http"
	"time"

	"gohat/internal/shared/web"
	"gohat/internal/view"
)

type Example struct {
	logger      *slog.Logger
	exampleView *view.Example
}

func NewExample(logger *slog.Logger) *Example {
	return &Example{
		logger:      logger,
		exampleView: view.NewExample(),
	}
}

func (h *Example) GetExample(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("example page visited")
	web.Render(w, h.exampleView.ExamplePage())
}

func (h *Example) GetExampleModal(w http.ResponseWriter, r *http.Request) {
	web.Render(w, h.exampleView.ExampleModalContent())
}

func (h *Example) GetExampleToast(w http.ResponseWriter, r *http.Request) {
	web.Render(w, h.exampleView.ExampleToastContent())
}

func (h *Example) GetExampleLongRequest(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 5)
}
