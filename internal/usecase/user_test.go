package usecase

import (
	"errors"
	"testing"

	"chatapp/internal/domain/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepo) FindByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepo) FindByID(id string) (*entity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockUserRepo) FindAll() ([]*entity.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *mockUserRepo) Update(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepo) Delete(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		in         *CreateUserInput
		mockReturn []interface{}
		wantErr    bool
	}{
		{
			name: "success",
			in: &CreateUserInput{
				Name:     "test",
				Email:    "test@test.com",
				Password: "password",
			},
			mockReturn: []interface{}{
				&entity.User{
					Name:     "test",
					Email:    "test@test.com",
					Password: "password",
				},
				nil,
			},
			wantErr: false,
		},
		{
			name: "error when creating user",
			in: &CreateUserInput{
				Name:     "test",
				Email:    "test@test.com",
				Password: "password",
			},
			mockReturn: []interface{}{nil, errors.New("error")},
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockRepo mockUserRepo
			mockRepo.On("Create", mock.Anything).Return(test.mockReturn...)

			u := &UserUseCase{UserRepo: &mockRepo}
			userResponse, err := u.CreateUser(test.in)
			if test.wantErr && test.mockReturn[0] == nil {
				assert.Error(t, err)
			} else {
				user, _ := test.mockReturn[0].(*entity.User)
				assert.NoError(t, err)
				assert.Equal(t, user.Name, userResponse.Name)
			}
		})
	}
}

func TestAuthenticateUser(t *testing.T) {
	plainPassword := "password"
	user, _ := entity.NewUser("test", "test@test.com", plainPassword)

	tests := []struct {
		name       string
		in         *AuthenticateUserInput
		mockReturn []interface{}
		wantErr    bool
	}{
		{
			name: "success",
			in: &AuthenticateUserInput{
				Email:    "test@test.com",
				Password: plainPassword,
			},
			mockReturn: []interface{}{
				user,
				nil,
			},
			wantErr: false,
		},
		{
			name: "error when authenticating user",
			in: &AuthenticateUserInput{
				Email:    "test@test.com",
				Password: "invalidPassword",
			},
			mockReturn: []interface{}{
				user,
				nil,
			},
			wantErr: true,
		},
		{
			name: "error when finding user",
			in: &AuthenticateUserInput{
				Email:    "test@test.com",
				Password: plainPassword,
			},
			mockReturn: []interface{}{nil, errors.New("error")},
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockRepo mockUserRepo
			mockRepo.On("FindByEmail", mock.Anything).Return(test.mockReturn...)

			u := &UserUseCase{UserRepo: &mockRepo}
			userResponse, err := u.AuthenticateUser(test.in)
			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, userResponse)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userResponse)
				assert.Equal(t, test.mockReturn[0].(*entity.User).Name, userResponse.Name)
			}
		})
	}
}

func TestReadUser(t *testing.T) {
	tests := []struct {
		name       string
		in         string
		mockReturn []interface{}
		wantErr    bool
	}{
		{
			name: "success",
			in:   "1",
			mockReturn: []interface{}{
				&entity.User{
					Name:     "test",
					Email:    "test@test.com",
					Password: "password",
				},
				nil,
			},
			wantErr: false,
		},
		{
			name:       "not found",
			in:         "1",
			mockReturn: []interface{}{nil, nil},
			wantErr:    true,
		},
		{
			name:       "error when finding user",
			in:         "1",
			mockReturn: []interface{}{nil, errors.New("error")},
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockRepo mockUserRepo
			mockRepo.On("FindByID", mock.Anything).Return(test.mockReturn...)

			u := &UserUseCase{UserRepo: &mockRepo}
			userResponse, err := u.ReadUser(test.in)
			if test.wantErr {
				assert.Error(t, err)
			} else if test.mockReturn[0] == nil {
				assert.NoError(t, err)
				assert.Nil(t, userResponse)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.mockReturn[0].(*entity.User).Name, userResponse.Name)
			}
		})
	}
}

func TestReadAllUsers(t *testing.T) {
	tests := []struct {
		name       string
		mockReturn []interface{}
		wantErr    bool
	}{
		{
			name: "success",
			mockReturn: []interface{}{
				[]*entity.User{
					{
						Name:     "test",
						Email:    "test@test.com",
						Password: "password",
					},
					{
						Name:     "test2",
						Email:    "test2test.com",
						Password: "password2",
					},
				},
				nil,
			},
			wantErr: false,
		},
		{
			name:       "not found",
			mockReturn: []interface{}{nil, nil},
			wantErr:    false,
		},
		{
			name:       "error when finding users",
			mockReturn: []interface{}{nil, errors.New("error")},
			wantErr:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockRepo mockUserRepo
			mockRepo.On("FindAll").Return(test.mockReturn...)

			u := &UserUseCase{UserRepo: &mockRepo}
			usersResponse, err := u.ReadAllUsers()
			if test.wantErr {
				assert.Error(t, err)
			} else if test.mockReturn[0] == nil {
				assert.NoError(t, err)
				assert.Equal(t, 0, len(usersResponse.Users))
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, usersResponse)
				assert.Equal(t, 2, len(usersResponse.Users))
				assert.Equal(t, test.mockReturn[0].([]*entity.User)[0].Name, usersResponse.Users[0].Name)
				assert.Equal(t, test.mockReturn[0].([]*entity.User)[1].Name, usersResponse.Users[1].Name)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name             string
		inUserID         string
		inUserInput      *UpdateUserInput
		findMockReturn   []interface{}
		updateMockReturn error
		wantErr          bool
	}{
		{
			name:     "success",
			inUserID: "1",
			inUserInput: &UpdateUserInput{
				Name:  "test",
				Email: "update@test.com",
			},
			findMockReturn: []interface{}{
				&entity.User{
					Name:     "test",
					Email:    "test@test.com",
					Password: "password",
				},
				nil,
			},
			updateMockReturn: nil,
			wantErr:          false,
		},
		{
			name:     "error when finding user",
			inUserID: "1",
			inUserInput: &UpdateUserInput{
				Name:  "test",
				Email: "update@test.com",
			},
			findMockReturn:   []interface{}{nil, nil},
			updateMockReturn: errors.New("error"),
			wantErr:          true,
		},
		{
			name:     "error when user not found",
			inUserID: "1",
			inUserInput: &UpdateUserInput{
				Name:  "test",
				Email: "test@test.com",
			},
			findMockReturn:   []interface{}{nil, nil},
			updateMockReturn: nil,
			wantErr:          true,
		},
		{
			name:     "error when updating user",
			inUserID: "1",
			inUserInput: &UpdateUserInput{
				Name:  "test",
				Email: "update@test.com",
			},
			findMockReturn: []interface{}{
				&entity.User{
					Name:     "test",
					Email:    "test@test.com",
					Password: "password",
				},
				nil,
			},
			updateMockReturn: errors.New("error"),
			wantErr:          true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockRepo mockUserRepo
			mockRepo.On("FindByID", test.inUserID).Return(test.findMockReturn...)
			mockRepo.On("Update", mock.Anything).Return(test.updateMockReturn)

			u := &UserUseCase{UserRepo: &mockRepo}
			err := u.UpdateUser(test.inUserID, test.inUserInput)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name             string
		in               string
		findMockReturn   []interface{}
		deleteMockReturn error
		wantErr          bool
	}{
		{
			name: "success",
			in:   "1",
			findMockReturn: []interface{}{
				&entity.User{
					Name:     "test",
					Email:    "test@test.com",
					Password: "password",
				},
				nil,
			},
			deleteMockReturn: nil,
			wantErr:          false,
		},
		{
			name:             "error when finding user",
			in:               "1",
			findMockReturn:   []interface{}{nil, errors.New("error")},
			deleteMockReturn: nil,
			wantErr:          true,
		},
		{
			name:             "error when user not found",
			in:               "1",
			findMockReturn:   []interface{}{nil, nil},
			deleteMockReturn: nil,
			wantErr:          true,
		},
		{
			name: "error when deleting user",
			in:   "1",
			findMockReturn: []interface{}{
				&entity.User{
					Name:  "test",
					Email: "test@test.com",
				},
				nil,
			},
			deleteMockReturn: errors.New("error"),
			wantErr:          true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var mockRepo mockUserRepo
			mockRepo.On("FindByID", test.in).Return(test.findMockReturn...)
			mockRepo.On("Delete", mock.Anything).Return(test.deleteMockReturn)

			u := &UserUseCase{UserRepo: &mockRepo}
			err := u.DeleteUser(test.in)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
