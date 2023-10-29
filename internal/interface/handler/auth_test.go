package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"chatapp/internal/usecase"
	"chatapp/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) CreateUser(input *usecase.CreateUserInput) (*usecase.UserResponse, *errors.CustomError) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.CustomError)
	}

	if args.Get(1) == nil {
		return args.Get(0).(*usecase.UserResponse), nil
	}

	return args.Get(0).(*usecase.UserResponse), args.Get(1).(*errors.CustomError)
}

func (m *MockAuthUseCase) AuthenticateUser(input *usecase.AuthenticateUserInput) (*usecase.UserResponse, *errors.CustomError) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.CustomError)
	}

	if args.Get(1) == nil {
		return args.Get(0).(*usecase.UserResponse), nil
	}

	return args.Get(0).(*usecase.UserResponse), args.Get(1).(*errors.CustomError)
}

func TestSignUp(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedStatus int
		mockReturn     []interface{}
	}{
		{
			name:           "success",
			input:          `{"name": "test", "email": "test@test.com", "password": "password"}`,
			expectedStatus: http.StatusOK,
			mockReturn: []interface{}{
				&usecase.UserResponse{
					ID:   1,
					Name: "test",
				},
				nil,
			},
		},
		{
			name:           "invalid request for binding error",
			input:          `{"name": "test", "email": "test@test.com, "password": }`,
			expectedStatus: http.StatusBadRequest,
			mockReturn:     []interface{}{nil, nil},
		},
		{
			name:           "invalid request",
			input:          `{"name": "test", "email": "", "password": "password"}`,
			expectedStatus: http.StatusBadRequest,
			mockReturn: []interface{}{
				nil,
				&errors.CustomError{Type: errors.BadRequest},
			},
		},
		{
			name:           "internal server error",
			input:          `{"name": "test", "email": "test@test.com", "password": "password"}`,
			expectedStatus: http.StatusInternalServerError,
			mockReturn: []interface{}{
				nil,
				&errors.CustomError{Type: errors.InternalServerError},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockAuthUseCase MockAuthUseCase
			mockAuthUseCase.On("CreateUser", mock.Anything).Return(test.mockReturn...)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(test.input))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			authHandler := &AuthHandler{AuthUseCase: &mockAuthUseCase}
			authHandler.SignUp(c)
			assert.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}

func TestSignIn(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedStatus int
		mockReturn     []interface{}
	}{
		{
			name:           "success",
			input:          `{"email": "test@test.com", "password": "password"}`,
			expectedStatus: http.StatusOK,
			mockReturn: []interface{}{
				&usecase.UserResponse{
					ID:   1,
					Name: "test",
				},
				nil,
			},
		},
		{
			name:           "invalid request for binding error",
			input:          `{"email": , "password": "password"}`,
			expectedStatus: http.StatusBadRequest,
			mockReturn:     []interface{}{nil, nil},
		},
		{
			name:           "not found",
			input:          `{"email": "notfound@test.com", "password": "password"}`,
			expectedStatus: http.StatusNotFound,
			mockReturn: []interface{}{
				nil,
				&errors.CustomError{Type: errors.NotFound},
			},
		},
		{
			name:           "invalid credentials",
			input:          `{"email": "test@test.com", "password": "invalid"}`,
			expectedStatus: http.StatusUnauthorized,
			mockReturn: []interface{}{
				nil,
				&errors.CustomError{Type: errors.InvalidCredentials},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockAuthUseCase MockAuthUseCase
			mockAuthUseCase.On("AuthenticateUser", mock.Anything).Return(test.mockReturn...)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/signin", strings.NewReader(test.input))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			authHandler := &AuthHandler{AuthUseCase: &mockAuthUseCase}
			authHandler.SignIn(c)
			assert.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}
