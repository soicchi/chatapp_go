package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name            string
		customError     *CustomError
		expectedMessage string
		expectedStatus  int
	}{
		{
			name:            "bad request",
			customError:     &CustomError{Type: BadRequest, Error: errors.New("errors")},
			expectedMessage: "invalid request",
			expectedStatus:  400,
		},
		{
			name:            "invalid credentials",
			customError:     &CustomError{Type: InvalidCredentials, Error: errors.New("errors")},
			expectedMessage: "invalid credentials",
			expectedStatus:  401,
		},
		{
			name:            "internal server error",
			customError:     &CustomError{Type: InternalServerError, Error: errors.New("errors")},
			expectedMessage: "internal server error",
			expectedStatus:  500,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			message, status := test.customError.getErrorDetails()
			assert.Equal(t, test.expectedMessage, message)
			assert.Equal(t, test.expectedStatus, status)
		})
	}
}
