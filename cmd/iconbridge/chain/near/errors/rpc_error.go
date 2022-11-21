package errors

import (
	"github.com/icon-project/icon-bridge/common/jsonrpc"
)

type ErrorCause struct {
	Name string      `json:"name,omitempty"`
	Info interface{} `json:"info,omitempty"`
	Code uint8
}

type RpcError struct {
	Name    string            `json:"name,omitempty"`
	Cause   ErrorCause        `json:"cause,omitempty"`
	Code    jsonrpc.ErrorCode `json:"code"`
	Message string            `json:"message"`
	Data    interface{}       `json:"data,omitempty"`
}

func (e *RpcError) Error() string {
	return e.Cause.Name
}
