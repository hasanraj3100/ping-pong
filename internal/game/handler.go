package game

import (
	"context"
	"encoding/json"
	"net/http"

	"hasanraj3100/ping-pong/internal/domain"

	"github.com/coder/websocket"
)

type UserLookup func(uuid string) (domain.User, bool)

type Hub interface {
	Broadcaster
	Join(ctx context.Context, room string, conn *websocket.Conn, onMessage func(data []byte))
}

type Handler struct {
	svc    *GameService
	lookup UserLookup
	hub    Hub
}

func NewHandler(svc *GameService, lookup UserLookup, hub Hub) *Handler {
	return &Handler{svc: svc, lookup: lookup, hub: hub}
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

func (h *Handler) WS(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("uuid")
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	user, ok := h.lookup(cookie.Value)
	if !ok {
		http.Error(w, "not logged in", http.StatusUnauthorized)
		return
	}

	gameID := r.PathValue("id")
	g, ok := h.svc.GetGame(gameID)
	if !ok {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer conn.CloseNow()

	ctx := r.Context()

	data, err := json.Marshal(g)
	if err != nil {
		return
	}
	if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
		return
	}

	onMessage := func(data []byte) {
		var msg ClientMsg
		if err := json.Unmarshal(data, &msg); err != nil {
			return
		}
		h.svc.HandleClientMsg(msg, user, gameID)
	}

	h.hub.Join(ctx, gameID, conn, onMessage)
}
