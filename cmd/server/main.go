package main

import (
	"log"
	"net/http"

	"hasanraj3100/ping-pong/internal/repository"
	"hasanraj3100/ping-pong/internal/user"
)

func main() {
	repo := repository.NewUserRepository()
	svc := user.NewUserService(repo)
	handler := user.NewHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", handler.CreateUser)

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
