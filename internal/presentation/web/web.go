package web

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"hasanraj3100/ping-pong/internal/domain"
)

//go:embed views static
var files embed.FS

var welcomeTmpl = template.Must(template.ParseFS(files, "views/welcome.html"))

type UserLookup func(uuid string) (domain.User, bool)

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

func StaticHandler() http.Handler {
	static, err := fs.Sub(files, "static")
	if err != nil {
		panic(err)
	}

	return http.StripPrefix("/static/", http.FileServerFS(static))
}
