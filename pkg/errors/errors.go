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

var errorDetails = map[CustomErrorType]struct {
	Message string
	Status  int
}{
	BadRequest:          {Message: "invalid request", Status: http.StatusBadRequest},
	InvalidCredentials:  {Message: "invalid credentials", Status: http.StatusUnauthorized},
	NotFound:            {Message: "resource not found", Status: http.StatusNotFound},
	InternalServerError: {Message: "internal server error", Status: http.StatusInternalServerError},
}

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
	detail, ok := errorDetails[e.Type]
	if !ok {
		return "unknown error", http.StatusInternalServerError
	}

	return detail.Message, detail.Status
}

func responseData(message string, status int) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
		"status":  status,
	}
}
