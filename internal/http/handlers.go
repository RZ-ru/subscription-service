package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"subs-service/internal/service"
)

type Handler struct {
	service *service.SubscriptionService
}

func NewHandler(s *service.SubscriptionService) *Handler {
	return &Handler{service: s}
}

// ===================== UTILS =====================

// парсим id из /subscriptions/123
func extractID(path string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return 0, http.ErrMissingFile
	}
	return strconv.ParseInt(parts[len(parts)-1], 10, 64)
}

// ===================== HANDLERS =====================

// CREATE
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ServiceName string  `json:"service_name"`
		Price       int     `json:"price"`
		UserID      string  `json:"user_id"`
		Start       string  `json:"start_date"`
		End         *string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.service.Logger().Error("invalid json", "error", err)
		http.Error(w, "invalid json", 400)
		return
	}

	id, err := h.service.Create(r.Context(), req.ServiceName, req.Price, req.UserID, req.Start, req.End)
	if err != nil {
		h.service.Logger().Error("create failed", "error", err)
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"id": id})
}

// READ
func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	sub, err := h.service.ReadByID(r.Context(), id)
	if err != nil {
		h.service.Logger().Error("read failed", "error", err)
		http.Error(w, err.Error(), 404)
		return
	}

	json.NewEncoder(w).Encode(sub)
}

// UPDATE
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	var req struct {
		ServiceName string  `json:"service_name"`
		Price       int     `json:"price"`
		UserID      string  `json:"user_id"`
		Start       string  `json:"start_date"`
		End         *string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.service.Logger().Error("invalid json", "error", err)
		http.Error(w, "invalid json", 400)
		return
	}

	if err := h.service.Update(r.Context(), id, req.ServiceName, req.Price, req.UserID, req.Start, req.End); err != nil {
		h.service.Logger().Error("update failed", "error", err)
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(204)
}

// DELETE
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.service.Logger().Error("delete failed", "error", err)
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(204)
}

// LIST
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	userID := query.Get("user_id")
	serviceName := query.Get("service_name")

	var userPtr *string
	if userID != "" {
		userPtr = &userID
	}

	var servicePtr *string
	if serviceName != "" {
		servicePtr = &serviceName
	}

	list, err := h.service.List(
		r.Context(),
		userPtr,
		servicePtr,
		0, 0,
	)

	if err != nil {
		h.service.Logger().Error("list failed", "error", err)
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(list)
}

// SUM
func (h *Handler) Sum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	from := query.Get("from")
	to := query.Get("to")
	if from == "" || to == "" {
		http.Error(w, "from and to are required", 400)
		return
	}

	userID := query.Get("user_id")
	serviceName := query.Get("service_name")

	var userPtr *string
	if userID != "" {
		userPtr = &userID
	}

	var servicePtr *string
	if serviceName != "" {
		servicePtr = &serviceName
	}

	sum, err := h.service.SumByPeriod(r.Context(), from, to, userPtr, servicePtr)
	if err != nil {
		h.service.Logger().Error("sum failed", "error", err)
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{"total": sum})
}
