package realtime

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/coder/websocket"
)

type Hub struct {
	mu    sync.Mutex
	rooms map[string]map[*client]struct{}
}

type client struct {
	send chan []byte
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[*client]struct{})}
}

func (h *Hub) Broadcast(room string, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for c := range h.rooms[room] {
		select {
		case c.send <- data:
		default: // slow client, drop the message rather than block the broadcaster
		}
	}
}

func (h *Hub) Join(ctx context.Context, room string, conn *websocket.Conn, onMessage func(data []byte)) {
	c := &client{send: make(chan []byte, 8)}

	h.mu.Lock()
	if h.rooms[room] == nil {
		h.rooms[room] = make(map[*client]struct{})
	}
	h.rooms[room][c] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.rooms[room], c)
		h.mu.Unlock()
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-c.send:
				if conn.Write(ctx, websocket.MessageText, data) != nil {
					cancel()
					return
				}
			}
		}
	}()

	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}
		if onMessage != nil {
			onMessage(data)
		}
	}
}
