// 🔹 Problem 4: Middleware – Logging

// Create a logging middleware that logs:

// [METHOD] /path - duration

// Example:

// POST /users - 2ms

// Rules

// Middleware must wrap handlers

// Measure time taken per request

package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

/*The provided key must be comparable and should not be of type string or any other built-in type to avoid collisions between packages using context. Users of WithValue should define their own types for keys. To avoid allocating when assigning to an interface{}, context keys often have concrete type struct{}. Alternatively, exported context key variables' static type should be a pointer or interface.
*/

var reqCounter atomic.Uint64

type contextKey string

const requestIDKey contextKey = "requestID"

func generateRequestID() string {
	return fmt.Sprintf("req-%d", reqCounter.Add(1))
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // --- Code that runs BEFORE the wrapped handler ---
        timeTaken := time.Now()
        // task 5  
        requestId := generateRequestID()

        w.Header().Set("X-Request-ID", requestId)

        ctx := context.WithValue(r.Context(), requestIDKey, requestId)
        // Call the next handler in the chain
        next.ServeHTTP(w, r.WithContext(ctx))

        // --- Code that runs AFTER the wrapped handler ---
        duration := time.Since(timeTaken)

        log.Printf("[%s] %+10v   %+10v   %+10v", requestId , r.Method, r.URL.Path, duration.String())
    })
}
