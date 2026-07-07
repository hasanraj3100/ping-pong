package user

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

type createUserRequest struct {
	Name string `json:"name"`
}

type createUserResponse struct {
	Username string `json:"username"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	created := h.svc.Register(req.Name)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createUserResponse{Username: created.Name})
}
