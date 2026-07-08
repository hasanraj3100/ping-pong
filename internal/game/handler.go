package game

import (
	"encoding/json"
	"net/http"

	"hasanraj3100/ping-pong/internal/domain"
)

type UserLookup func(uuid string) (domain.User, bool)

type Handler struct {
	svc    *GameService
	lookup UserLookup
}

func NewHandler(svc *GameService, lookup UserLookup) *Handler {
	return &Handler{svc: svc, lookup: lookup}
}

func (h *Handler) CreateGame(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("uuid")
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	creator, ok := h.lookup(cookie.Value)
	if !ok {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	created := h.svc.CreateGame(creator)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created)
}

func (h *Handler) JoinGame(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("uuid")
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	joiner, ok := h.lookup(cookie.Value)
	if !ok {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	joined, err := h.svc.JoinGame(r.PathValue("id"), joiner)
	if err != nil {
		switch err {
		case ErrGameNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case ErrGameFull, ErrAlreadyJoined:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(joined)
}
