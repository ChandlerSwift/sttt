package models

import (
	"gorm.io/gorm"
)

// A User is a user. You probably don't need this explained.
type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string // TODO: hash with an auth lib; set `gorm:"size:255"`
}
