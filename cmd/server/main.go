package main

import (
	"log"
	"net/http"

	"hasanraj3100/ping-pong/internal/game"
	"hasanraj3100/ping-pong/internal/presentation/web"
	"hasanraj3100/ping-pong/internal/repository"
	"hasanraj3100/ping-pong/internal/user"
)

func main() {
	repo := repository.NewUserRepository()
	svc := user.NewUserService(repo)
	handler := user.NewHandler(svc)

	gameRepo := repository.NewGameRepository()
	gameSvc := game.NewGameService(gameRepo)
	gameHandler := game.NewHandler(gameSvc, svc.GetByUUID)

	mux := http.NewServeMux()
	mux.HandleFunc("/", web.Index)
	mux.HandleFunc("/welcome", web.NewWelcomeHandler(svc.GetByUUID))
	mux.Handle("/static/", web.StaticHandler())
	mux.HandleFunc("POST /users", handler.CreateUser)
	mux.HandleFunc("POST /games", gameHandler.CreateGame)
	mux.HandleFunc("POST /games/{id}/join", gameHandler.JoinGame)
	mux.HandleFunc("GET /games/{id}", web.NewGamePageHandler(svc.GetByUUID, gameSvc.GetGame, gameSvc.JoinGame))

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
