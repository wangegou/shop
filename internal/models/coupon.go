package models

import (
	"gorm.io/gorm"
)

// Coupon represents a discount code.
type Coupon struct {
	gorm.Model
	Code     string  `gorm:"unique;not null"`
	Discount float64 `gorm:"not null"`
	IsUsed   bool    `gorm:"default:false"`
}
