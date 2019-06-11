package vk

import (
	"encoding/json"
	"sync"
)

// API - главный объект
type API struct {
	AccessToken    string
	retryCount     int
	httpRetryCount int
	ErrorToSkip    []string
	sync.Mutex
	ExecuteErrors []ExecuteErrors
	ExecuteCode   string
}

// ExecuteErrors - объект ошибок execute
type ExecuteErrors struct {
	Method    string `json:"method"`
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// Response - объект ответа VK
type Response struct {
	Response      json.RawMessage `json:"response"`
	Error         ResponseError   `json:"error"`
	ExecuteErrors []ExecuteErrors `json:"execute_errors"`
}

// ResponseError - объект ошибки выболнения запроса
type ResponseError struct {
	ErrorCode     int                 `json:"error_code"`
	ErrorMsg      string              `json:"error_msg"`
	RequestParams []map[string]string `json:"request_params"`
}
