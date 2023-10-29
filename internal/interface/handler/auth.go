package handler

import (
	"fmt"
	"net/http"

	"chatapp/internal/usecase"
	"chatapp/pkg/errors"

	"github.com/labstack/echo/v4"
)

type AuthUseCase interface {
	CreateUser(input *usecase.CreateUserInput) (*usecase.UserResponse, *errors.CustomError)
	AuthenticateUser(input *usecase.AuthenticateUserInput) (*usecase.UserResponse, *errors.CustomError)
}

type AuthHandler struct {
	AuthUseCase AuthUseCase
}

type SignUpInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(authUseCase AuthUseCase) *AuthHandler {
	return &AuthHandler{
		AuthUseCase: authUseCase,
	}
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	var input SignUpInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("invalid request"))
	}

	inputToUsecase := usecase.NewCreateUserInput(input.Name, input.Email, input.Password)
	user, err := h.AuthUseCase.CreateUser(inputToUsecase)
	if err != nil {
		return err.ErrorResponse(c)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *AuthHandler) SignIn(c echo.Context) error {
	var input SignUpInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("invalid request"))
	}

	inputToUsecase := usecase.NewAuthenticateUserInput(input.Email, input.Password)
	user, err := h.AuthUseCase.AuthenticateUser(inputToUsecase)
	if err != nil {
		return err.ErrorResponse(c)
	}

	return c.JSON(http.StatusOK, user)
}
