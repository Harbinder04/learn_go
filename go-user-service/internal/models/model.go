package internal


type MyError struct {
	ReqId string `json:"request_id"`
	Error string `json:"error"`
}

type MyResposeType struct {
	ReqId string `json:"request_id"`
	Data interface{} `json:"data"`
}