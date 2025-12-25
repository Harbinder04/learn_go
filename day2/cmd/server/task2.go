// Problem 2: Request & Response Structs

// Create a /users endpoint (POST).

// Request JSON

// {
//   "name": "Harbinder",
//   "email": "test@example.com"
// }


// Response JSON

// {
//   "id": "generated-id",
//   "name": "Harbinder",
//   "email": "test@example.com"
// }


// Tasks

// Decode JSON request body

// Validate input using Day 1 logic

// Encode JSON response


// 🔹 Problem 3: HTTP Status Codes (Important)

// Modify /users:

// 201 Created → success

// 400 Bad Request → validation error

// 405 Method Not Allowed → wrong HTTP method


package server

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
	"net/http"
)

type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
}

type RepUser struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func CheckIsValid(name string, email string) error {
	if len(name) < 2 {
		return errors.New("name must be of 2 characters")
	}
	if email == "" {
		return errors.New("email is required")
	}
	return nil
}

func generateId() int {
	return rand.IntN(1000)
}

func CreateUser(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPost {
		http.Error(w, "Mehtod not allowed", http.StatusMethodNotAllowed)
	}

	if r.Header.Get("Content-Type") != "application/json" {
	     http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
	     return
	 }

	 var newUser User

	 // How to get data from the json request body in go? 
	 err := json.NewDecoder(r.Body).Decode(&newUser); if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err1 := CheckIsValid(newUser.Name, newUser.Email); err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}


	id := generateId()

    createdUser := RepUser{
		Id: id,
		Name: newUser.Name,
		Email: newUser.Email,
	}

	requestID := r.Context().Value(requestIDKey).(string)
	w.Header().Set("X-Request-ID", requestID)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}