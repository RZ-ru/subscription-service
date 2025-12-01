package http

import "net/http"

func (h *Handler) RegisterRoutes() {
	http.HandleFunc("/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodGet:
			// Если есть GET /subscriptions?param=value → List
			if r.URL.Path == "/subscriptions" {
				h.List(w, r)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// CRUD по id
	http.HandleFunc("/subscriptions/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.Read(w, r)
		case http.MethodPut:
			h.Update(w, r)
		case http.MethodDelete:
			h.Delete(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// SUM
	http.HandleFunc("/subscriptions/sum", h.Sum)
	// Swagger
	http.HandleFunc("/swagger", h.Swagger)

}
