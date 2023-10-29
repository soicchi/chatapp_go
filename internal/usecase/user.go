package usecase

import (
	"fmt"
	"log"

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

// NewCreateUserInput creates a new input for creating a user
func NewCreateUserInput(name, email, password string) *CreateUserInput {
	return &CreateUserInput{
		Name:     name,
		Email:    email,
		Password: password,
	}
}

// NewAuthenticateUserInput creates a new input for authenticating a user
func NewAuthenticateUserInput(email, password string) *AuthenticateUserInput {
	return &AuthenticateUserInput{
		Email:    email,
		Password: password,
	}
}

// Create creates a new user
func (u *UserUseCase) CreateUser(input *CreateUserInput) (*UserResponse, *errors.CustomError) {
	log.Println("CreateUser:", input)

	user, err := entity.NewUser(input.Name, input.Email, input.Password)
	if err != nil {
		return nil, errors.NewCustomError(errors.BadRequest, err)
	}

	newUser, err := u.UserRepo.Create(user)
	if err != nil {
		return nil, errors.NewCustomError(errors.InternalServerError, err)
	}

	responseUser := &UserResponse{
		ID:   newUser.ID,
		Name: newUser.Name,
	}

	return responseUser, nil
}

// AuthenticateUser authenticates a user
func (u *UserUseCase) AuthenticateUser(input *AuthenticateUserInput) (*UserResponse, *errors.CustomError) {
	log.Println("AuthenticateUser:", input)

	user, err := u.UserRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, errors.NewCustomError(errors.InternalServerError, err)
	}

	if user == nil {
		return nil, errors.NewCustomError(errors.NotFound, fmt.Errorf("user not found"))
	}

	if !user.CheckPassword(input.Password) {
		return nil, errors.NewCustomError(errors.InvalidCredentials, fmt.Errorf("invalid credentials"))
	}

	responseUser := &UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}

	return responseUser, nil
}

// ReadUser reads a user
func (u *UserUseCase) ReadUser(userID string) (*UserResponse, *errors.CustomError) {
	log.Println("ReadUser:", userID)

	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return nil, errors.NewCustomError(errors.InternalServerError, err)
	}

	if user == nil {
		return nil, errors.NewCustomError(errors.NotFound, err)
	}

	responseUser := &UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}

	return responseUser, nil
}

// FindAllUsers finds all users
func (u *UserUseCase) ReadAllUsers() (*UsersResponse, *errors.CustomError) {
	log.Println("ReadAllUsers")

	users, err := u.UserRepo.FindAll()
	if err != nil {
		return nil, errors.NewCustomError(errors.InternalServerError, err)
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
func (u *UserUseCase) UpdateUser(userID string, input *UpdateUserInput) *errors.CustomError {
	log.Println("UpdateUser:", input)

	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return errors.NewCustomError(errors.InternalServerError, err)
	}
	if user == nil {
		return errors.NewCustomError(errors.NotFound, fmt.Errorf("user not found"))
	}

	user.Name = input.Name
	user.Email = input.Email
	if err := u.UserRepo.Update(user); err != nil {
		return errors.NewCustomError(errors.InternalServerError, err)
	}

	return nil
}

// DeleteUser deletes a user
func (u *UserUseCase) DeleteUser(userID string) *errors.CustomError {
	user, err := u.UserRepo.FindByID(userID)
	if err != nil {
		return errors.NewCustomError(errors.InternalServerError, err)
	}
	if user == nil {
		return errors.NewCustomError(errors.NotFound, fmt.Errorf("user not found"))
	}

	if err := u.UserRepo.Delete(user); err != nil {
		return errors.NewCustomError(errors.InternalServerError, err)
	}

	return nil
}
