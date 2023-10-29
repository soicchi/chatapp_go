package errors

import (
	"log"
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
	Type  CustomErrorType
	Error error
}

var errorDetails = map[CustomErrorType]struct {
	Message string
	Status  int
}{
	BadRequest:          {Message: "invalid request", Status: http.StatusBadRequest},
	InvalidCredentials:  {Message: "invalid credentials", Status: http.StatusUnauthorized},
	NotFound:            {Message: "resource not found", Status: http.StatusNotFound},
	InternalServerError: {Message: "internal server error", Status: http.StatusInternalServerError},
}

func NewCustomError(t CustomErrorType, err error) *CustomError {
	return &CustomError{
		Type:  t,
		Error: err,
	}
}

func (e *CustomError) ErrorResponse(c echo.Context) error {
	log.Println(e.Error)
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
