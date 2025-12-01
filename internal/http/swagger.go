package http

import (
	"net/http"
)

func (h *Handler) Swagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./docs/swagger.yaml")
}
