package apperror

import "net/http"

type AppError interface {
	error
	StatusCode() int
	Message() string
}

type appErrorImpl struct {
	statusCode int
	message    string
}

func (e *appErrorImpl) Error() string {
	return e.message
}

func (e *appErrorImpl) StatusCode() int {
	return e.statusCode
}

func (e *appErrorImpl) Message() string {
	return e.message
}

// --- Helper Functions ---

// Conflict (409)
func NewConflict(message string) AppError {
	return &appErrorImpl{
		statusCode: http.StatusConflict,
		message:    message,
	}
}

// Bad Request (400)
func NewBadRequest(message string) AppError {
	return &appErrorImpl{
		statusCode: http.StatusBadRequest,
		message:    message,
	}
}

// Not Found (404)
func NewNotFound(message string) AppError {
	return &appErrorImpl{
		statusCode: http.StatusNotFound,
		message:    message,
	}
}

// Unauthorized (401)
func NewUnauthorized(message string) AppError {
	return &appErrorImpl{
		statusCode: http.StatusUnauthorized,
		message:    message,
	}
}

// Forbidden (403)
func NewForbidden(message string) AppError {
	return &appErrorImpl{
		statusCode: http.StatusForbidden,
		message:    message,
	}
}

// InternalServerError (500)
func NewInternalServerError(message string) AppError {
	return &appErrorImpl{
		statusCode: http.StatusInternalServerError,
		message:    message,
	}
}

// MethodNotAllowed (405)
func NewMethodNotAllowed(message string) AppError {
	return &appErrorImpl{
		statusCode: http.StatusMethodNotAllowed,
		message:    message,
	}
}
