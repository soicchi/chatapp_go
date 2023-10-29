package handler

import (
	"fmt"
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

type mockUserUseCase struct {
	mock.Mock
}

func (m *mockUserUseCase) ReadUser(userID string) (*usecase.UserResponse, *errors.CustomError) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.CustomError)
	}
	return args.Get(0).(*usecase.UserResponse), nil
}

func (m *mockUserUseCase) ReadAllUsers() (*usecase.UsersResponse, *errors.CustomError) {
	args := m.Called()
	if args.Get(0) == nil && args.Get(1) == nil {
		return nil, nil
	}

	if args.Get(0) == nil {
		return nil, args.Get(1).(*errors.CustomError)
	}

	return args.Get(0).(*usecase.UsersResponse), nil
}

func (m *mockUserUseCase) UpdateUser(input *usecase.UpdateUserInput) *errors.CustomError {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*errors.CustomError)
}

func (m *mockUserUseCase) DestroyUser(userID string) *errors.CustomError {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*errors.CustomError)
}

func TestRetrieveUser(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		mockReturn []interface{}
		wantStatus int
	}{
		{
			name:  "success",
			input: "1",
			mockReturn: []interface{}{
				&usecase.UserResponse{
					ID:   uint(1),
					Name: "test",
				},
				nil,
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "error when reading user",
			input: "1",
			mockReturn: []interface{}{
				nil,
				errors.NewCustomError(errors.InternalServerError, fmt.Errorf("error")),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockUserUseCase mockUserUseCase
			mockUserUseCase.On("ReadUser", test.input).Return(test.mockReturn...)

			userHandler := NewUserHandler(&mockUserUseCase)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/users/:id", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(test.input)

			userHandler.RetrieveUser(c)
			assert.Equal(t, test.wantStatus, rec.Code)
		})
	}
}

func TestListUsers(t *testing.T) {
	tests := []struct {
		name       string
		mockReturn []interface{}
		wantStatus int
	}{
		{
			name: "success",
			mockReturn: []interface{}{
				&usecase.UsersResponse{
					Users: []usecase.UserResponse{
						{
							ID:   uint(1),
							Name: "test",
						},
						{
							ID:   uint(2),
							Name: "test2",
						},
					},
				},
				nil,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error when reading users",
			mockReturn: []interface{}{
				nil,
				errors.NewCustomError(errors.InternalServerError, fmt.Errorf("error")),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockUserUseCase mockUserUseCase
			mockUserUseCase.On("ReadAllUsers").Return(test.mockReturn...)

			userHandler := NewUserHandler(&mockUserUseCase)
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			userHandler.ListUsers(c)
			assert.Equal(t, test.wantStatus, rec.Code)
		})
	}
}

func TestUpdateUserInfo(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		mockReturn []interface{}
		wantStatus int
	}{
		{
			name:       "success",
			input:      `{ "name": "test", "email": "test@test.com"}`,
			mockReturn: []interface{}{nil},
			wantStatus: http.StatusNoContent,
		},
		{
			name:  "error when updating user",
			input: `{ "name": "test", "email": "test@test.com"}`,
			mockReturn: []interface{}{
				errors.NewCustomError(errors.InternalServerError, fmt.Errorf("error")),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockUserUseCase mockUserUseCase
			mockUserUseCase.On("UpdateUser", mock.Anything).Return(test.mockReturn...)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/users/:id", strings.NewReader(test.input))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("1")

			userHandler := NewUserHandler(&mockUserUseCase)
			userHandler.UpdateUserInfo(c)
			assert.Equal(t, test.wantStatus, rec.Code)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name       string
		mockReturn []interface{}
		wantStatus int
	}{
		{
			name:       "success",
			mockReturn: []interface{}{nil},
			wantStatus: http.StatusNoContent,
		},
		{
			name: "error when deleting user",
			mockReturn: []interface{}{
				errors.NewCustomError(errors.InternalServerError, fmt.Errorf("error")),
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockUserUseCase mockUserUseCase
			mockUserUseCase.On("DestroyUser", mock.Anything).Return(test.mockReturn...)

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/users/:id", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("1")

			userHandler := NewUserHandler(&mockUserUseCase)
			userHandler.DeleteUser(c)
			assert.Equal(t, test.wantStatus, rec.Code)
		})
	}
}
