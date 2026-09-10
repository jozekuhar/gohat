package handler

import (
	"net/http"

	"mimokocke/internal/web/view"
)

type Core struct {
	coreView *view.Core
}

func NewCore() *Core {
	return &Core{
		coreView: view.NewCore(),
	}
}

func (h *Core) GetNotFound(w http.ResponseWriter, r *http.Request) {
	render(w, h.coreView.NotFoundPage(), http.StatusNotFound)
}

func (h *Core) GetTermsOfService(w http.ResponseWriter, r *http.Request) {
	render(w, h.coreView.PrivacyPolicyPage(), http.StatusOK)
}

func (h *Core) GetPrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	render(w, h.coreView.TermsOfServicePage(), http.StatusOK)
}
