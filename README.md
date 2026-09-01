# Ping Pong

A real-time, two-player Pong game written in Go. The server owns the ball physics, scoring, and game state; the browser client renders the game with [Kaplay](https://kaplayjs.com/) and syncs over WebSockets.

## Features

- Server-authoritative ball physics and scoring — clients cannot cheat the simulation
- Shareable game links: create a match, send the URL to a friend, they auto-join as Player 2
- Ready-up flow before each round starts
- Live paddle and ball sync over WebSockets, with client-side ball extrapolation for smooth motion between server ticks
- Ball speeds up on every paddle hit
- Win condition at 10 points with a game-over broadcast to both players
- No database — everything lives in-memory for the life of the process

## Tech Stack

- **Backend:** Go (standard library `net/http`, `html/template`), [`coder/websocket`](https://github.com/coder/websocket) for WebSocket connections, [`google/uuid`](https://github.com/google/uuid) for IDs
- **Frontend:** vanilla JS + [Kaplay](https://kaplayjs.com/) for canvas rendering, no build step

## Project Layout

```
cmd/server/            entry point, wires everything together and starts the HTTP server
internal/domain/       core types (User, Game, GameData)
internal/user/         user registration/lookup (service, handler, DTOs)
internal/game/         game service, HTTP + WebSocket handlers, ball loop, scoring
internal/realtime/     WebSocket hub — manages per-room connections and broadcasts
internal/repository/   in-memory repositories for users and games
internal/presentation/web/
  views/               HTML templates (welcome, game board)
  static/              client JS/CSS + the Kaplay engine bundle
```

## Getting Started

### Prerequisites

- Go 1.25+, **or** Docker + Docker Compose

### Run the server

**Locally with Go:**

```bash
go run ./cmd/server
```

**With Docker Compose:**

```bash
docker compose up --build
```

The server listens on `:8080` either way.

### Play a game

1. Open `http://localhost:8080` and enter a name to get a session cookie.
2. From `/welcome`, create a game — you'll be redirected to `/games/{id}` as Player 1.
3. Share the game URL with a second player (or open it in another browser/incognito window). They'll auto-join as Player 2.
4. Both players click **Ready**. Once both are ready, the ball starts moving.
5. Use the **Up/Down arrow keys** to move your paddle.
6. First to 10 points wins.

## How It Works

- **Auth:** lightweight — registering a name issues a UUID stored in a `uuid` cookie. There's no password or session expiry.
- **Game creation/joining:** `POST /games` creates a game owned by the caller; `POST /games/{id}/join` (or simply visiting the game page as a second user) assigns Player 2.
- **Realtime sync:** each game is a "room" in the WebSocket hub. Clients connect via `GET /games/{id}/ws`, send `move`/`ready` messages, and receive broadcasts (`move`, `ready`, `ball`, `score`, `game_over`) as JSON.
- **Physics loop:** once both players ready up, the server runs a fixed-tick (50ms) goroutine per game that advances the ball, resolves wall/paddle collisions, updates scores, and broadcasts state — the client only renders what it's told.

## API Reference

| Method | Path                  | Description                                  |
|--------|-----------------------|-----------------------------------------------|
| GET    | `/`                    | Landing/name entry page                      |
| GET    | `/welcome`             | Post-login landing page                      |
| POST   | `/users`               | Register a user, returns `{ uuid, username }` |
| POST   | `/games`               | Create a new game (requires `uuid` cookie)    |
| POST   | `/games/{id}/join`     | Join an existing game as Player 2             |
| GET    | `/games/{id}`          | Game board page (auto-joins if a slot is open)|
| GET    | `/games/{id}/ws`       | WebSocket connection for live game state      |


