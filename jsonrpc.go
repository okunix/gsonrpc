package gsonrpc

import (
	"encoding/json"
)

type Request struct {
	Version string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params,omitempty"`
	ID      any    `json:"id,omitempty"`
}

type Response struct {
	Version string    `json:"jsonrpc"`
	Result  any       `json:"result,omitempty"`
	Error   *RPCError `json:"error,omitempty"`
	ID      any       `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (err RPCError) Error() string {
	errBytes, _ := json.Marshal(err)
	return string(errBytes)
}

func NewError(code int, message string, data any) RPCError {
	return RPCError{Code: code, Message: message, Data: data}
}

var (
	ErrParseError     = NewError(-32700, "Parse Error", nil)
	ErrInvalidRequest = NewError(-32600, "Invalid Request", nil)
	ErrMethodNotFound = NewError(-32601, "Method Not Found", nil)
	ErrInvalidParams  = NewError(-32602, "Invalid Params", nil)
	ErrInternalError  = NewError(-32603, "Internal Error", nil)
)

type JsonRPC struct {
	*RPC
}

func NewJsonRPC(rpc *RPC) *JsonRPC {
	return &JsonRPC{RPC: rpc}
}

func (j *JsonRPC) ProcessRequest(request []byte) []byte {
	responseObject := &Response{
		Version: "2.0",
	}

	var requestObject Request
	if err := json.Unmarshal(request, &requestObject); err != nil {
		responseObject.Error = &ErrParseError
		responseObject.ID = nil
		resp, _ := json.Marshal(responseObject)
		return resp
	}
	if requestObject.ID == nil {
		responseObject = nil
	} else {
		responseObject.ID = requestObject.ID
	}

	var out []any
	var err error
	out, err = j.RPC.Call(requestObject.Method, requestObject.Params...)

	if err != nil {
		responseObject.Error = new(err.(RPCError))
		resp, _ := json.Marshal(responseObject)
		return resp
	}

	if responseObject == nil {
		return nil
	}

	responseObject.Result = out
	if len(out) == 1 {
		responseObject.Result = out[0]
	}

	resp, _ := json.Marshal(responseObject)
	return resp
}
