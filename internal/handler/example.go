package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"gohat/internal/shared/clock"
	"gohat/internal/shared/web"
	"gohat/internal/view"
)

type Example struct {
	logger      *slog.Logger
	clock       clock.Clock
	exampleView *view.Example
}

func NewExample(logger *slog.Logger, clock clock.Clock) *Example {
	return &Example{
		logger:      logger,
		clock:       clock,
		exampleView: view.NewExample(),
	}
}

func (h *Example) GetExample(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("example page visited")
	currentTime := h.clock.NowUTC()
	web.Render(w, h.exampleView.ExamplePage(currentTime))
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

func (h *Example) GetExampleSSERequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	rc := http.NewResponseController(w)
	i := 0

	for {
		_, err := fmt.Fprintf(w, "data: streaming response: %d, \n\n", i)
		if err != nil {
			return
		}

		err = rc.Flush()
		if err != nil {
			return
		}

		i++
		time.Sleep(200 * time.Millisecond)
	}
}
