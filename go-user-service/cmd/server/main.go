package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handlers "go-user-service/internal/handler"
	logger "go-user-service/internal/logger"
	customMiddleware "go-user-service/internal/middleware"
	store "go-user-service/internal/store"
	"go-user-service/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)
	
type Config struct {
	port string
	env string
	db string
}


func main() {
	godotenv.Load("../../.env.dev")
	cfg := &Config{
		port: os.Getenv("PORT"),
		env: os.Getenv("ENV"),
		db: os.Getenv("DB_URL"),
	}
	port := cfg.port
	env := cfg.env
	logger := logger.NewLogger(env)

	db := db.NewdbConnection(cfg.db)

	if cfg.db == "" {
		logger.Error("Database URL not provided")
		os.Exit(1)
	}
	
	r := chi.NewRouter()
	server := &http.Server{
		Addr : ":"+ port,
		Handler: r,
	}

	st := store.NewUserStore(db)
	handler := handlers.NewUserHandler(st, logger)
	
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	
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

	// start the server
	go func() {
       if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
           logger.Error("Server failed: " + err.Error())
       }
   }()

	shutDownCtx := make(chan struct{})

	// for handling gracefull shutdown.
	go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan
		logger.Info("Shutdown signal received")

        shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
        defer shutdownRelease()
		logger.Warn("Waiting for ongoing requests")

        if err := server.Shutdown(shutdownCtx); err != nil {
            logger.Error("Server shutdown error: " + err.Error())
        }
		logger.Info("Closing database connection")
   		db.Close()
        logger.Info("Server stopped gracefully")

		close(shutDownCtx)
    }()

	log.Println("Server started on 8080")
	
	<- shutDownCtx

}