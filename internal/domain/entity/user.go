package entity

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null; size:255; check:name <> ''"`
	Email    string `gorm:"not null; size:255; unique; check:email <> ''"`
	Password string `gorm:"not null; size:255; check:password <> ''"`
}

func NewUser(name, email, password string) (*User, error) {
	if name == "" || email == "" || password == "" {
		return nil, fmt.Errorf("name, email, and password must not be empty")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	user := &User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hashedPassword), nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
