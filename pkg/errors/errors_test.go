package errors

import (
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
			customError:     &CustomError{Type: BadRequest},
			expectedMessage: "invalid request",
			expectedStatus:  400,
		},
		{
			name:            "invalid credentials",
			customError:     &CustomError{Type: InvalidCredentials},
			expectedMessage: "invalid credentials",
			expectedStatus:  401,
		},
		{
			name:            "not found",
			customError:     &CustomError{Type: NotFound, EntityType: "user"},
			expectedMessage: "user not found",
			expectedStatus:  404,
		},
		{
			name:            "internal server error",
			customError:     &CustomError{Type: InternalServerError},
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
