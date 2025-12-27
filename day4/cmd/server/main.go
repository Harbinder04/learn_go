package main

import (
	internal "day4/Internal"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	st := internal.NewUserStore()
	userhandler := internal.NewUserHandler(st)

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)
	r.Route("/users", func(r chi.Router) {
		r.With(internal.OnlyPost).Post("/", userhandler.CreateUser)
		r.Get("/", userhandler.GetAllUsers)
		r.Get("/{id}", userhandler.GetUserbyId)
	})

    fmt.Println("Server is started at port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Unable to run a server %v", err)
	}
}