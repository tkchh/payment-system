package response

import (
	"github.com/go-chi/render"
	"net/http"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type Response struct {
	Status string `json:"status"`
	Code   int    `json:"code,omitempty"`
	Data   any    `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}

func Success(data any) Response {
	return Response{
		Status: StatusOK,
		Code:   http.StatusOK,
		Data:   data,
	}
}

func Error(r *http.Request, msg string, code int) Response {
	render.Status(r, code)
	return Response{
		Status: StatusError,
		Code:   code,
		Error:  msg,
	}
}
