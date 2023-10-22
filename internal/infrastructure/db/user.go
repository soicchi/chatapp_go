package db

import (
	"errors"
	"fmt"

	"chatapp/internal/domain/entity"

	"gorm.io/gorm"
)

// UserRepository is a repository for the user entity
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *entity.User) error {
	result := r.DB.Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	return nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

// FindAll finds all users
func (r *UserRepository) FindAll() ([]*entity.User, error) {
	var users []*entity.User
	err := r.DB.Find(&users).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find all users: %w", err)
	}

	return users, nil
}

// Update updates a user
func (r *UserRepository) Update(user *entity.User) error {
	if err := r.DB.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(user *entity.User) error {
	if err := r.DB.Delete(user).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
