// task 1 ; create a simple health check
package server

import (
	"encoding/json"
	"net/http"
)

type Resposnse struct {
	Status string `json:"status"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := Resposnse {
		Status: "ok",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}