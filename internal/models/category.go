package models

import (
	"gorm.io/gorm"
)

// Category represents a product category.
type Category struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	Products []Product
}
