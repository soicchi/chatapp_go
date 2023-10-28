package errors

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CustomErrorType int

const (
	BadRequest CustomErrorType = iota
	InvalidCredentials
	NotFound
	InternalServerError
)

type CustomError struct {
	Type CustomErrorType
}

func NewCustomError(t CustomErrorType) *CustomError {
	return &CustomError{Type: t}
}

func (e *CustomError) ErrorResponse(c echo.Context) error {
	message, status := e.getErrorDetails()
	return c.JSON(status, responseData(message, status))
}

func (e *CustomError) getErrorDetails() (string, int) {
	switch e.Type {
	case BadRequest:
		return "invalid request", http.StatusBadRequest
	case InvalidCredentials:
		return "invalid credentials", http.StatusUnauthorized
	case NotFound:
		return "not found", http.StatusNotFound
	default:
		return "internal server error", http.StatusInternalServerError
	}
}

func responseData(message string, status int) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
		"status":  status,
	}
}
