package errorweb

import (
	"net/http"

	"mimokocke/internal/web/response"
)

func GetNotFound(w http.ResponseWriter, r *http.Request) {
	response.RenderStatus(w, notFoundPage(), http.StatusNotFound)
}
