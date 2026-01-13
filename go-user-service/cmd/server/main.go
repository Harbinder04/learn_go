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

	"go-user-service/config"
	dbConfig "go-user-service/internal/db"
	handlers "go-user-service/internal/handler"
	logger "go-user-service/internal/logger"
	customMiddleware "go-user-service/internal/middleware"
	store "go-user-service/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func main() {
	cfg := config.NewConfig()

	port := cfg.ServerConfig.Port
	env := cfg.ServerConfig.Env

	logger := logger.NewLogger(env)

	db, err := dbConfig.NewdbConnection(logger, cfg.DatabaseConfig.GetConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
		ReadTimeout: cfg.ServerConfig.ReadTimeout,
		WriteTimeout: cfg.ServerConfig.WriteTimeout,
		IdleTimeout: cfg.ServerConfig.IdleTimeout,
		ReadHeaderTimeout: cfg.ServerConfig.ReadHeaderTimeout,
	}

	st := store.NewSQLUserStore(db)
	handler := handlers.NewUserHandler(st, logger)

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// attaches logger to the request
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "logger", logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.With(customMiddleware.CheckResTime, customMiddleware.CheckTimeOut).Route("/users", func(r chi.Router) {
		r.Post("/", handler.CreateUser)
		r.Get("/", handler.GetAllUsers)
		r.Get("/{id}", handler.GetUserbyId)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"msg": "Pong"})
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
	<-shutDownCtx

}
