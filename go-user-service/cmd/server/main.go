package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	handlers "go-user-service/internal/handler"
	logger "go-user-service/internal/logger"
	customMiddleware "go-user-service/internal/middleware"
	store "go-user-service/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)
	

func main() {
	godotenv.Load("../../.env.dev")

	port := os.Getenv("PORT")
	env := os.Getenv("ENV")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	
	st := store.NewUserStore()
	logger := logger.NewLogger(env)
	handler := handlers.NewUserHandler(st, logger)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "logger", logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	
	r.With(customMiddleware.CheckResTime).Route("/users", func(r chi.Router) {
		r.Post("/", handler.CreateUser)
		r.Get("/", handler.GetAllUsers)
		r.Get("/{id}", handler.GetUserbyId)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"msg":"Pong"})
	})

	log.Println("Server started on 8080")
	err := http.ListenAndServe(":"+port, r); if err != nil {
		log.Fatalf("Server not started %v", err)
	}
}