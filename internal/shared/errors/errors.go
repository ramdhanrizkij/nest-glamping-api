package errors

import "net/http"

type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func BadRequest(msg string) *AppError    { return NewAppError(http.StatusBadRequest, msg) }
func Unauthorized(msg string) *AppError  { return NewAppError(http.StatusUnauthorized, msg) }
func Forbidden(msg string) *AppError     { return NewAppError(http.StatusForbidden, msg) }
func NotFound(msg string) *AppError      { return NewAppError(http.StatusNotFound, msg) }
func Conflict(msg string) *AppError      { return NewAppError(http.StatusConflict, msg) }
func Internal(msg string) *AppError      { return NewAppError(http.StatusInternalServerError, msg) }
