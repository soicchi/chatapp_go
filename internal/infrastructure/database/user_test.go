package database

import (
	"testing"

	"chatapp/internal/domain/entity"
	helper "chatapp/tests"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		input   *entity.User
		wantErr bool
	}{
		{
			name: "success",
			input: &entity.User{
				Name:     "test",
				Email:    "test@test.com",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "error empty name",
			input: &entity.User{
				Name:     "",
				Email:    "test2@test.com",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "error empty email",
			input: &entity.User{
				Name:     "test",
				Email:    "",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "error duplicated email",
			input: &entity.User{
				Name:     "test",
				Email:    "initial@test.com",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "error empty password",
			input: &entity.User{
				Name:     "test",
				Email:    "test3@test.com",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		// Create transaction
		tx := testDB.Begin()
		helper.CreateTestUser(tx, "initial", "initial@test.com", "password")

		t.Run(test.name, func(t *testing.T) {
			repo := &UserRepository{DB: tx}
			user, err := repo.Create(test.input)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, user.ID)
				assert.Equal(t, test.input.Name, user.Name)
				assert.Equal(t, test.input.Email, user.Email)
			}
		})
		tx.Rollback()
	}
}

func TestFindByEmail(t *testing.T) {
	tests := []struct {
		name       string
		inputEmail string
		wantErr    bool
	}{
		{
			name:       "success",
			inputEmail: "test@test.com",
			wantErr:    false,
		},
		{
			name:       "not found",
			inputEmail: "not_found@test.com",
			wantErr:    false,
		},
		{
			name:       "error empty email",
			inputEmail: "",
			wantErr:    false,
		},
	}

	for _, test := range tests {
		// Create transaction
		tx := testDB.Begin()
		helper.CreateTestUser(tx, "test", "test@test.com", "password")

		t.Run(test.name, func(t *testing.T) {
			repo := &UserRepository{DB: tx}
			user, err := repo.FindByEmail(test.inputEmail)
			if test.wantErr {
				assert.Error(t, err)
			} else if user == nil {
				assert.NoError(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "test", user.Name)
				assert.Equal(t, "test@test.com", user.Email)
			}
		})
		tx.Rollback()
	}
}

func TestFindAll(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
		{
			name:    "not found",
			wantErr: false,
		},
	}

	for _, test := range tests {
		// Create transaction
		tx := testDB.Begin()
		helper.CreateTestUser(tx, "test", "test@test.com", "password")
		helper.CreateTestUser(tx, "test2", "test2@test.com", "password")

		t.Run(test.name, func(t *testing.T) {
			repo := &UserRepository{DB: tx}
			users, err := repo.FindAll()
			if test.wantErr {
				assert.Error(t, err)
			} else if users == nil {
				assert.NoError(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 2, len(users))
				assert.Equal(t, "test", users[0].Name)
				assert.Equal(t, "test@test.com", users[0].Email)
				assert.Equal(t, "test2", users[1].Name)
				assert.Equal(t, "test2@test.com", users[1].Email)
			}
		})
		tx.Rollback()
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
	}

	for _, test := range tests {
		// Create transaction
		tx := testDB.Begin()
		helper.CreateTestUser(tx, "test", "test@test.com", "password")

		t.Run(test.name, func(t *testing.T) {
			var targetUser entity.User
			tx.Where("email = ?", "test@test.com").First(&targetUser)
			targetUser.Name = "updated"
			targetUser.Email = "updated@test.com"

			repo := &UserRepository{DB: tx}
			err := repo.Update(&targetUser)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				var updatedUser entity.User
				tx.Where("email = ?", "updated@test.com").First(&updatedUser)
				assert.Equal(t, "updated", updatedUser.Name)
				assert.Equal(t, "updated@test.com", updatedUser.Email)
			}
		})
		tx.Rollback()
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "success",
			wantErr: false,
		},
	}

	for _, test := range tests {
		// Create transaction
		tx := testDB.Begin()
		helper.CreateTestUser(tx, "test", "test@test.com", "password")

		t.Run(test.name, func(t *testing.T) {
			var targetUser entity.User
			tx.Where("email = ?", "test@test.com").First(&targetUser)

			repo := &UserRepository{DB: tx}
			err := repo.Delete(&targetUser)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				var deletedUser *entity.User
				tx.Where("email = ?", "test@test.com").First(deletedUser)
				assert.Nil(t, deletedUser)
			}
		})
		tx.Rollback()
	}
}
