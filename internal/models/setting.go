package models

import (
	"gorm.io/gorm"
)

// Setting represents a key-value pair for system configuration.
type Setting struct {
	gorm.Model
	Key   string `gorm:"unique;not null"`
	Value string
}
