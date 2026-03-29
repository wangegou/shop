package models

import (
	"gorm.io/gorm"
)

// Product represents an item for sale.
type Product struct {
	gorm.Model
	Name        string
	Description string
	Price       float64
	CategoryID  uint
	Category    Category
	Stock       []ProductStock
}
