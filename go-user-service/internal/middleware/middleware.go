package internal

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// All users share one bucket — one user spamming blocks everyone else
var userLimiter = rate.NewLimiter(rate.Every(time.Second/10), 5)

func CheckResTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		if duration > 500*time.Millisecond {
			if log, ok := r.Context().Value("logger").(*slog.Logger); ok {
				log.Warn("Request takes > 500ms")
			}
		}
	})

}

func CheckTimeOut(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timoutCtx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(timoutCtx))
	})
}

func RateLimitUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := userLimiter.Wait(r.Context()); err != nil {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
