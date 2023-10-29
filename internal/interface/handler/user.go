package handler

import (
	"net/http"

	"chatapp/internal/usecase"
	"chatapp/pkg/errors"

	"github.com/labstack/echo/v4"
)

type UserUseCase interface {
	ReadUser(userID string) (*usecase.UserResponse, *errors.CustomError)
	ReadAllUsers() (*usecase.UsersResponse, *errors.CustomError)
	UpdateUser(input *usecase.UpdateUserInput) *errors.CustomError
	DestroyUser(userID string) *errors.CustomError
}

type UserHandler struct {
	UserUseCase UserUseCase
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUserHandler(userService UserUseCase) *UserHandler {
	return &UserHandler{
		UserUseCase: userService,
	}
}

func (h *UserHandler) RetrieveUser(c echo.Context) error {
	userID := c.Param("id")
	user, customErr := h.UserUseCase.ReadUser(userID)
	if customErr != nil {
		return customErr.ErrorResponse(c)
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	users, customErr := h.UserUseCase.ReadAllUsers()
	if customErr != nil {
		return customErr.ErrorResponse(c)
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUserInfo(c echo.Context) error {
	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		customError := errors.NewCustomError(http.StatusBadRequest, err)
		return customError.ErrorResponse(c)
	}

	userID := c.Param("id")
	inputToUseCase := usecase.NewUpdateUserInput(userID, req.Name, req.Email)
	if customErr := h.UserUseCase.UpdateUser(inputToUseCase); customErr != nil {
		return customErr.ErrorResponse(c)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID := c.Param("id")
	customErr := h.UserUseCase.DestroyUser(userID)
	if customErr != nil {
		return customErr.ErrorResponse(c)
	}

	return c.JSON(http.StatusNoContent, nil)
}
