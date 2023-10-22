package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null; size:255; check:name <> ''"`
	Email    string `gorm:"not null; size:255; unique; check:email <> ''"`
	Password string `gorm:"not null; size:255; check:password <> ''"`
}
