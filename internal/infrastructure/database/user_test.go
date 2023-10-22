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

	helper.CreateTestUser(testDB, "initial", "initial@test.com", "password")

	for _, test := range tests {
		tx := testDB.Begin()
		t.Run(test.name, func(t *testing.T) {
			repo := &UserRepository{DB: tx}
			err := repo.Create(test.input)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				var user entity.User
				tx.Where("email = ?", test.input.Email).First(&user)
				assert.NotZero(t, user.ID)
				assert.Equal(t, test.input.Name, user.Name)
				assert.Equal(t, test.input.Email, user.Email)
			}
			tx.Rollback()
		})
	}
}
