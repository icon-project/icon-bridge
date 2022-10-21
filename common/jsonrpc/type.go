package jsonrpc

import (
	"encoding/json"
	"fmt"
)

const Version = "2.0"

type Request struct {
	Version string          `json:"jsonrpc" validate:"required,version"`
	Method  string          `json:"method" validate:"required"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id"`
}

type Response struct {
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	ID      interface{} `json:"id,omitempty"`
}

type ErrorCode int

const (
	ErrorCodeJsonParse      ErrorCode = -32700
	ErrorCodeInvalidRequest ErrorCode = -32600
	ErrorCodeMethodNotFound ErrorCode = -32601
	ErrorCodeInvalidParams  ErrorCode = -32602
	ErrorCodeInternal       ErrorCode = -32603
	ErrorCodeServer         ErrorCode = -32000
)

type Error struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("jsonrpc: code: %d, message: %s, data: %+v", e.Code, e.Message, e.Data)
}

type Response2 struct {
	Version string          `json:"jsonrpc"`
	Result  interface{}     `json:"result,omitempty"`
	Error   json.RawMessage `json:"error,omitempty"`
	ID      interface{}     `json:"id,omitempty"`
}
