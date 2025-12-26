package handlers

import (
	"day3/internal/store"
	"encoding/json"
	"errors"
	"math/rand/v2"
	"net/http"
	"strconv"
)

// var userStore *store.UserStore
//task 5

// func InitStore(st *store.UserStore) {
//     userStore = st
// }

type UserHandler struct {
	store *store.UserStore
}

func NewUserHandler(st *store.UserStore) *UserHandler {
	return &UserHandler{
		store: st,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    err := json.NewEncoder(w).Encode(data); if err != nil {
		return err
	}
	return nil
}

//task 4
func writeJSONError(w http.ResponseWriter, code int, message string){
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
	})
}

func CheckIsValid(name string, email string) error {
	if len(name) < 2 {
		return errors.New("Name must be greater than 2 characters")
	}
	if email == "" {
		return errors.New("Email must not be empty")
	}

	return  nil
}


func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		writeJSONError(w, http.StatusBadRequest, "Content-Type must be application/json")
	    return
	}

	var newUser store.User

	err := json.NewDecoder(r.Body).Decode(&newUser); if err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	userErr := CheckIsValid(newUser.Name, newUser.Email); if userErr != nil {
		writeJSONError(w, http.StatusBadRequest, userErr.Error())
		return
	}

	prefix := "abcd"
	id := string(prefix[rand.IntN(3)]) + strconv.Itoa(rand.IntN(100))

	if err := h.store.Create(store.User{ID: id, Name: newUser.Name, Email: newUser.Email}); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{id: id}); err != nil {
        writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-type", "application/json")
   if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Method must be GET",
		})
		return
   }

   allUsers := h.store.List()

   w.WriteHeader(http.StatusOK)

   if err := json.NewEncoder(w).Encode(allUsers); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Internal Server Error",
		})
		return
    }
}

func (h *UserHandler) GetUserbyId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Method must be GET",
		})
		return
   }

   id := r.PathValue("id")

   u, err := h.store.GetByID(id); if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: err.Error(),
		})
		return
   }

   w.WriteHeader(http.StatusOK)

   if err := json.NewEncoder(w).Encode(u); err != nil {
	w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: err.Error(),
		})
	return
   }
}