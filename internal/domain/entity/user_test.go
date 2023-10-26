package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]string
		wantErr bool
	}{
		{
			name: "success to create a new user",
			input: map[string]string{
				"name":     "test",
				"email":    "test@test.com",
				"password": "password",
			},
			wantErr: false,
		},
		{
			name: "fail to create a new user because name is empty",
			input: map[string]string{
				"name":     "",
				"email":    "test@test.com",
				"password": "password",
			},
			wantErr: true,
		},
		{
			name: "fail to create a new user because email is empty",
			input: map[string]string{
				"name":     "test",
				"email":    "",
				"password": "password",
			},
			wantErr: true,
		},
		{
			name: "fail to create a new user because password is empty",
			input: map[string]string{
				"name":     "test",
				"email":    "test@test.com",
				"password": "",
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := NewUser(test.input["name"], test.input["email"], test.input["password"])
			if test.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, test.input["name"], user.Name)
				assert.Equal(t, test.input["email"], user.Email)
				// Password must not equal because the password is hashed
				assert.NotEqual(t, test.input["password"], user.Password)
			}
		})
	}
}
