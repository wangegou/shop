package models

import (
	"gorm.io/gorm"
)

// User represents the admin user of the system.
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}
