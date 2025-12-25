package main

import (
	"fmt"
	"time"
	"errors"
	"log"
	"net/http"
	"http_day2/cmd/server"
	"os"
    "os/signal"
    "syscall"
	"context"
)

func main() {
	serverConfig := &http.Server{
        Addr: ":8080",
    }

    fmt.Println("Starting server...")
	// task1
	http.HandleFunc("/health", server.HealthCheck)
	
	//task 2
	// http.HandleFunc("/user", server.CreateUser)

	//task 4
	myHandler := http.HandlerFunc(server.CreateUser)
	
	wrappedHandler := server.LoggingMiddleware(myHandler)

	http.Handle("/user", wrappedHandler)
	
	fmt.Println("Server is running on http://localhost:8080")

    shutdownChan := make(chan bool, 1)

    go func() {
        if err := serverConfig.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
            log.Fatalf("HTTP server error: %v", err)
        }

        log.Println("Stopped serving new connections.")
        shutdownChan <- true
    }()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownRelease()

    if err := serverConfig.Shutdown(shutdownCtx); err != nil {
        log.Fatalf("HTTP shutdown error: %v", err)
    }

    <-shutdownChan
    log.Println("Graceful shutdown complete.")
	// ------------------------------------
	// fmt.Println("Starting server...")
	// // task1
	// http.HandleFunc("/health", server.HealthCheck)
	
	// //task 2
	// // http.HandleFunc("/user", server.CreateUser)

	// //task 4
	// myHandler := http.HandlerFunc(server.CreateUser)
	
	// wrappedHandler := server.LoggingMiddleware(myHandler)

	// http.Handle("/user", wrappedHandler)
	
	// fmt.Println("Server is running on http://localhost:8080")

	// if err := http.ListenAndServe(":8080", nil); !errors.Is(err, http.ErrServerClosed) {
	// 	log.Fatal("Unable to start localhost: ",err)
	// }

	
}
