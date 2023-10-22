package tests

import (
	"chatapp/internal/domain/entity"

	"gorm.io/gorm"
)

func CreateTestUser(db *gorm.DB, name, email, password string) error {
	user := &entity.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	if err := db.Create(user).Error; err != nil {
		return err
	}

	return nil
}
