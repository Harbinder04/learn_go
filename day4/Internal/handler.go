package internal

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5/middleware"
)


func NewUserHandler(st *UserStore) *UserHandler{
 return &UserHandler{
	store: st,
 }
}

func CreateJsonError(w http.ResponseWriter, errorCode int, reqId string, message string){
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(errorCode)
	json.NewEncoder(w).Encode(MyError{ReqId: reqId,Error: message})
}

func CreateJsonResponse(w http.ResponseWriter, statusCode int, reqId string, user interface{}){
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(MyResposeType{ReqId: reqId,Data: user})
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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request){
	reqId := middleware.GetReqID(r.Context())
    var nu User
	err := json.NewDecoder(r.Body).Decode(&nu); if err != nil {
		CreateJsonError(w, http.StatusBadRequest, reqId, err.Error())
		return
	}

	Validationerr := CheckIsValid(nu.Name, nu.Email); if Validationerr != nil {
		CreateJsonError(w, http.StatusBadRequest, reqId, Validationerr.Error())
		return
	}

	prefix := "abcd"
	id := string(prefix[rand.IntN(3)]) + strconv.Itoa(rand.IntN(100))

	if err := h.store.Create(User{Id: id, Name: nu.Name, Email: nu.Email}); err != nil {
		CreateJsonError(w, http.StatusBadRequest, reqId, err.Error())
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id, "request_id": reqId})
}


func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request){
    reqId := middleware.GetReqID(r.Context())
   if r.Method != http.MethodGet {
		CreateJsonError(w, http.StatusBadRequest, reqId, "Method must be GET")
		return
   }

   allUsers := h.store.GetAllUser()

   CreateJsonResponse(w, http.StatusOK, reqId, allUsers)
}

func (h *UserHandler) GetUserbyId(w http.ResponseWriter, r *http.Request) {
	reqId := middleware.GetReqID(r.Context())
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodGet {
		CreateJsonError(w, http.StatusBadRequest, reqId, "Method must be GET")
		return
   }

   id := r.PathValue("id")

   u, err := h.store.GetByID(id); if err != nil {
		CreateJsonError(w, http.StatusBadRequest, reqId, err.Error())
		return
   }

   CreateJsonResponse(w, http.StatusOK, reqId, u)
}