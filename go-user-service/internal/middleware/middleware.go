package internal

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)


func CheckResTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		start := time.Now()
	
		next.ServeHTTP(w, r)

		duration := time.Since(start)
		if duration >  500*time.Millisecond {
			if log, ok := r.Context().Value("logger").(*slog.Logger); ok {
			log.Warn("Request takes > 500ms")
			}
		}
	})

}

func CheckTimeOut(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		timoutCtx, cancel := context.WithTimeout(r.Context(), 4 * time.Second)
		defer cancel()
		 
		next.ServeHTTP(w, r.WithContext(timoutCtx))
	})
}