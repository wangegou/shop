package models

import (
	"gorm.io/gorm"
)

// ProductStock represents a single item in the product's inventory, e.g., a license key.
type ProductStock struct {
	gorm.Model
	ProductID uint
	Product   Product
	Key       string `gorm:"not null"`
	IsSold    bool   `gorm:"default:false"`
}
