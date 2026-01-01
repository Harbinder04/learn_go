package Internal

import (
	"encoding/json"
	"errors"
	models "go-user-service/internal/models"
	internal "go-user-service/internal/store"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
)

type UserHandler struct {
	store  *internal.UserStore
	logger *slog.Logger
}

func NewUserHandler(st *internal.UserStore, lg *slog.Logger) *UserHandler {
	return &UserHandler{
		store:  st,
		logger: lg,
	}
}

func CreateJsonError(w http.ResponseWriter, errorCode int, reqId string, logger *slog.Logger, message string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(errorCode)
	logger.Error(message, "Req_id:", reqId)
	json.NewEncoder(w).Encode(models.MyError{ReqId: reqId, Error: message})
}

func CreateJsonResponse(w http.ResponseWriter, statusCode int, reqId string, user interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(models.MyResposeType{ReqId: reqId, Data: user})
}

func CheckIsValid(name string, email string) error {
	// Todo: remove this delay later
	// time.Sleep(5 * time.Second)
	if len(name) < 2 {
		return errors.New("Name must be greater than 2 characters")
	}
	if email == "" {
		return errors.New("Email must not be empty")
	}

	return nil
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	reqId := middleware.GetReqID(r.Context())

	var nu internal.User
	err := json.NewDecoder(r.Body).Decode(&nu)
	if err != nil {
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, err.Error())
		return
	}

	Validationerr := CheckIsValid(nu.Name, nu.Email)
	if Validationerr != nil {
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, Validationerr.Error())
		return
	}

	prefix := "abcd"
	id := string(prefix[rand.IntN(3)]) + strconv.Itoa(rand.IntN(100))

	resId, err := h.store.Create(internal.User{Id: id, Name: nu.Name, Email: nu.Email})
	if err != nil {
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, err.Error())
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Todo; How to convert int64 to string i.e. resId
	json.NewEncoder(w).Encode(map[string]string{"id": resId, "request_id": reqId})
	h.logger.Info("user created", "id", id)

}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	reqId := middleware.GetReqID(r.Context())
	if r.Method != http.MethodGet {
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, "Method must be GET")
		return
	}

	allUsers, err := h.store.GetAllUser()
	if err != nil {
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, err.Error())
	}

	CreateJsonResponse(w, http.StatusOK, reqId, allUsers)
}

func (h *UserHandler) GetUserbyId(w http.ResponseWriter, r *http.Request) {
	reqId := middleware.GetReqID(r.Context())
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodGet {
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, "Method must be GET")
		return
	}

	id := r.PathValue("id")

	u, err := h.store.GetByID(id)
	if err != nil {
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, err.Error())
		return
	}

	CreateJsonResponse(w, http.StatusOK, reqId, u)
}
