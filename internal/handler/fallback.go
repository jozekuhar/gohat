package handler

import (
	"net/http"

	"gohat/internal/shared/web"
	"gohat/internal/view"
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
	web.Render(w, h.fallbackView.NotFoundPage())
}

func (h *Fallback) GetMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	web.Render(w, h.fallbackView.MethodNotAllowedPage())
}
