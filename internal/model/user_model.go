package model

import "gorm.io/gorm"

// User represents a registered user in the system.
type User struct {
	gorm.Model
	Name         string `gorm:"not null"             json:"name"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null"             json:"-"`
}

func init() {
	Register(&User{})
}