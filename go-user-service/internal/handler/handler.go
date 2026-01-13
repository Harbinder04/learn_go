package Internal

import (
	"context"
	"encoding/json"
	"errors"
	models "go-user-service/internal/models"
	internal "go-user-service/internal/store"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type UserHandler struct {
	store  internal.UserRepository
	logger *slog.Logger
}

func NewUserHandler(st internal.UserRepository, lg *slog.Logger) *UserHandler {
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
	ctx := r.Context()

	// fmt.Print(r.Context())

	reqId := middleware.GetReqID(ctx)

	var newUsr internal.User
	err := json.NewDecoder(r.Body).Decode(&newUsr)
	if err != nil {
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, err.Error())
		return
	}

	Validationerr := CheckIsValid(newUsr.Name, newUsr.Email)
	if Validationerr != nil {
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, Validationerr.Error())
		return
	}

	dur, exists, err := h.store.UserExists(ctx, newUsr.Email)
	if dur > 3 * time.Second {
		h.logger.Info("DB query takes more than 300ms")
	}
	if err != nil {
		if ctx.Err() == context.Canceled {
			h.logger.Info("Request cancelled during user creation", "request_id", reqId)
			return
		}
		h.logger.Error(err.Error())
		CreateJsonError(w, http.StatusInternalServerError, reqId, h.logger, "Failed to check user existence")
		return
	}
	if exists {
		CreateJsonError(w, http.StatusConflict, reqId, h.logger, "User with this email already exists")
		return
	}

	prefix := "abcd"
	id := string(prefix[rand.IntN(3)]) + strconv.Itoa(rand.IntN(100))

	//⚠️todo: Remove this delay
	// time.Sleep(5 * time.Second)

	dur, resId, err := h.store.Create(ctx, internal.User{Id: id, Name: newUsr.Name, Email: newUsr.Email})
	if dur > 3*time.Millisecond {
		h.logger.Info("DB query takes more than 300ms")
	}
	if err != nil {

		// ⚠️Working but not idomatic way like the driver can return error in some another form or by wrapping it
		// if errors.Is(err, context.Canceled) {
		// 	h.logger.Info("Request cancelled during user creation", "request_id", reqId)
		// 	return
		// }

		// ✅directly checking context for cancellation
		if ctx.Err() == context.Canceled {
			h.logger.Info("Request cancelled during user creation", "request_id", reqId)
			return
		}

		if ctx.Err() == context.DeadlineExceeded {
			h.logger.Info("Request Timeout", "request_id", reqId)
			CreateJsonError(w, http.StatusGatewayTimeout, reqId, h.logger, "Request Timeout")
			return
		}

		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, err.Error())
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{"id": resId, "request_id": reqId})
	h.logger.Info("user created", "id", id)

}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqId := middleware.GetReqID(r.Context())
	if r.Method != http.MethodGet {
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, "Method must be GET")
		return
	}

	dur, allUsers, err := h.store.GetAllUser(ctx)
	if dur > 3*time.Millisecond {
		h.logger.Info("DB query takes more than 300ms")
	}

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			h.logger.Info("Request Timeout", "request_id", reqId)
			CreateJsonError(w, http.StatusGatewayTimeout, reqId, h.logger, "Request Timeout")
			return
		}
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, err.Error())
	}

	if ctx.Err() == context.DeadlineExceeded {
		h.logger.Info("Request Timeout", "request_id", reqId)
		CreateJsonError(w, http.StatusGatewayTimeout, reqId, h.logger, "Request Timeout")
		return
	}
	CreateJsonResponse(w, http.StatusOK, reqId, allUsers)
}

func (h *UserHandler) GetUserbyId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqId := middleware.GetReqID(r.Context())
	w.Header().Set("Content-type", "application/json")
	if r.Method != http.MethodGet {
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, "Method must be GET")
		return
	}

	id := r.PathValue("id")

	dur, u, err := h.store.GetByID(ctx, id)
	if dur > 3*time.Millisecond {
		h.logger.Info("DB query takes more than 300ms")
	}
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			h.logger.Info("Request Timeout", "request_id", reqId)
			CreateJsonError(w, http.StatusGatewayTimeout, reqId, h.logger, "Request Timeout")
			return
		}
		CreateJsonError(w, http.StatusBadRequest, reqId, h.logger, err.Error())
		return
	}

	CreateJsonResponse(w, http.StatusOK, reqId, u)
}
