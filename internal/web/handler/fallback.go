package handler

import (
	"net/http"

	"mimokocke/internal/web/view"
)

type Fallback struct {
	fallbackView *view.Fallback
}

func NewFallback() *Fallback {
	return &Fallback{
		fallbackView: view.NewFallback(),
	}
}

func (h *Fallback) GetNotFound(w http.ResponseWriter, r *http.Request) {
	render(w, h.fallbackView.NotFoundPage())
}

func (h *Fallback) GetMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	render(w, h.fallbackView.MethodNotAllowedPage())
}
