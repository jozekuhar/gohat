package handler

import (
	"net/http"

	"mimokocke/internal/web/components"
	"mimokocke/internal/web/response"
)

func GetNotFound(w http.ResponseWriter, r *http.Request) {
	response.RenderStatus(w, components.NotFoundPage(), http.StatusNotFound)
}

func GetTermsOfService(w http.ResponseWriter, r *http.Request) {
	response.RenderStatus(w, components.PrivacyPolicyPage(), http.StatusOK)
}

func GetPrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	response.RenderStatus(w, components.TermsOfServicePage(), http.StatusOK)
}
