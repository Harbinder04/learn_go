package internal

import (
	"encoding/json"
	"net/http"
)


func OnlyPost(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Method not allowed", 
			})
			return 
		}
		next.ServeHTTP(w, r)
	})
}

