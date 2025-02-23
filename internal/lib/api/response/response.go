// Package response предоставляет унифицированный формат HTTP-ответов.
package response

import (
	"github.com/go-chi/render"
	"net/http"
)

// Статусы для JSON-ответов
const (
	StatusOK    = "OK"    // Успешный статус операции
	StatusError = "Error" // Ошибочный статус операции
)

// Response - базовая структура для всех HTTP-ответов
// Содержит статус выполнения, код ответа и данные/ошибку
type Response struct {
	Status string `json:"status"`          // Статус операции: OK/Error
	Code   int    `json:"code,omitempty"`  // HTTP-статус код (только для ошибок)
	Data   any    `json:"data,omitempty"`  // Тело успешного ответа
	Error  string `json:"error,omitempty"` // Сообщение об ошибке
}

// Success создает успешный JSON-ответ
// Принимает данные для возврата клиенту
func Success(data any) Response {
	return Response{
		Status: StatusOK,
		Code:   http.StatusOK,
		Data:   data,
	}
}

// Error создает JSON-ответ с ошибкой
// Автоматически устанавливает HTTP-статус через render
func Error(r *http.Request, msg string, code int) Response {
	render.Status(r, code)
	return Response{
		Status: StatusError,
		Code:   code,
		Error:  msg,
	}
}
