package staticweb

import (
	"net/http"

	"mimokocke/internal/web/response"

	"mimokocke"
)

func GetStatic(w http.ResponseWriter, r *http.Request) {
	http.FileServerFS(mimokocke.Static).ServeHTTP(w, r)
}

func GetTermsOfService(w http.ResponseWriter, r *http.Request) {
	response.RenderStatus(w, privacyPolicyPage(), http.StatusOK)
}

func GetPrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	response.RenderStatus(w, termsOfServicePage(), http.StatusOK)
}
