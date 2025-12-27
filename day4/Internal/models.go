package internal

import (
	"sync"
)

type UserStore struct {
	mu sync.RWMutex
	users map[string]User
}

type User struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type UserHandler struct {
	store *UserStore
}

type MyError struct {
	ReqId string `json:"request_id"`
	Error string `json:"error"`
}

type MyResposeType struct {
	ReqId string `json:"request_id"`
	Data interface{} `json:"data"`
}