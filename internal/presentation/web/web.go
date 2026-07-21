package web

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"

	"hasanraj3100/ping-pong/internal/domain"
)

//go:embed views static
var files embed.FS

var (
	welcomeTmpl = template.Must(template.ParseFS(files, "views/welcome.html"))
	gameTmpl    = template.Must(template.ParseFS(files, "views/game.html"))
)

type (
	UserLookup func(uuid string) (domain.User, bool)
	GameLookup func(id string) (domain.Game, bool)
	GameJoiner func(id string, player domain.User) (domain.Game, error)
)

type gamePageData struct {
	Game     domain.Game
	ShareURL string
	Role     string
}

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, files, "views/index.html")
}

func NewWelcomeHandler(lookup UserLookup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("uuid")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		found, ok := lookup(cookie.Value)
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		welcomeTmpl.Execute(w, found)
	}
}

func NewGamePageHandler(lookupUser UserLookup, lookupGame GameLookup, joinGame GameJoiner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := r.PathValue("id")
		loginRedirect := "/?redirect=" + url.QueryEscape("/games/"+gameID)

		cookie, err := r.Cookie("uuid")
		if err != nil {
			http.Redirect(w, r, loginRedirect, http.StatusSeeOther)
			return
		}

		user, ok := lookupUser(cookie.Value)
		if !ok {
			http.Redirect(w, r, loginRedirect, http.StatusSeeOther)
			return
		}

		g, ok := lookupGame(gameID)
		if !ok {
			http.NotFound(w, r)
			return
		}

		isPlayer := g.Player1.UUID == user.UUID || g.Player2.UUID == user.UUID
		if !isPlayer && g.Player2.UUID == "" {
			if joined, err := joinGame(gameID, user); err == nil {
				g = joined
			}
		}

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		role := "spectator"

		switch user.UUID {
		case g.Player1.UUID:
			role = "player1"
		case g.Player2.UUID:
			role = "player2"
		}

		data := gamePageData{
			Game:     g,
			ShareURL: scheme + "://" + r.Host + "/games/" + g.ID,
			Role:     role,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		gameTmpl.Execute(w, data)
	}
}

func StaticHandler() http.Handler {
	static, err := fs.Sub(files, "static")
	if err != nil {
		panic(err)
	}

	return http.StripPrefix("/static/", http.FileServerFS(static))
}
