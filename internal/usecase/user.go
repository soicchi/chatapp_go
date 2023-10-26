package usecase

import (
	"fmt"

	"chatapp/internal/domain/entity"
	"chatapp/pkg/errors"
)

// UserRepository is a repository for the user entity
type UserRepository interface {
	Create(user *entity.User) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
	FindAll() ([]*entity.User, error)
	Update(user *entity.User) error
	Delete(user *entity.User) error
}

// UserUseCase is a use case for the user entity
type UserUseCase struct {
	UserRepo UserRepository
}

// UserResponse is a response for the user entity
type UserResponse struct {
	ID   uint
	Name string
}

// UsersResponse is a response for the user entity
type UsersResponse struct {
	Users []UserResponse
}

// CreateUserInput is an input for creating a user
type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

// AuthenticateUserInput is an input for authenticating a user
type AuthenticateUserInput struct {
	Email    string
	Password string
}

// UpdateUserInput is an input for updating a user
type UpdateUserInput struct {
	Name  string
	Email string
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{UserRepo: repo}
}

// Create creates a new user
func (u *UserUseCase) CreateUser(input *CreateUserInput) (*UserResponse, error) {
	user, err := entity.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	newUser, err := u.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}

	responseUser := &UserResponse{
		ID:   newUser.ID,
		Name: newUser.Name,
	}

	return responseUser, nil
}

// AuthenticateUser authenticates a user
func (u *UserUseCase) AuthenticateUser(input *AuthenticateUserInput) (*UserResponse, error) {
	user, err := u.UserRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		entityErr := errors.NewEntityError("user")
		return nil, entityErr.NotFoundError()
	}

	if !user.CheckPassword(input.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	responseUser := &UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}

	return responseUser, nil
}

// ReadUser reads a user
func (u *UserUseCase) ReadUser(userID string) (*UserResponse, error) {
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		entityErr := errors.NewEntityError("user")
		return nil, entityErr.NotFoundError()
	}

	responseUser := &UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}

	return responseUser, nil
}

// FindAllUsers finds all users
func (u *UserUseCase) ReadAllUsers() (*UsersResponse, error) {
	users, err := u.UserRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responseUsers := make([]UserResponse, len(users))
	for i, user := range users {
		responseUsers[i] = UserResponse{
			ID:   user.ID,
			Name: user.Name,
		}
	}

	return &UsersResponse{Users: responseUsers}, nil
}

// UpdateUser updates a user
func (u *UserUseCase) UpdateUser(userID string, input *UpdateUserInput) error {
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		entityErr := errors.NewEntityError("user")
		return entityErr.NotFoundError()
	}

	user.Name = input.Name
	user.Email = input.Email
	if err := u.UserRepo.Update(user); err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user
func (u *UserUseCase) DeleteUser(userID string) error {
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		entityErr := errors.NewEntityError("user")
		return entityErr.NotFoundError()
	}

	if err := u.UserRepo.Delete(user); err != nil {
		return err
	}

	return nil
}
