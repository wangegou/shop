package models

import (
	"gorm.io/gorm"
)

// Order represents a customer's purchase.
type Order struct {
	gorm.Model
	OrderNo      string `gorm:"unique;not null"`
	ProductID    uint
	Product      Product
	Email        string `gorm:"not null"`
	Total        float64
	Status       string `gorm:"default:'pending'"` // pending, paid, completed
	PurchasedKey string
	CouponID     *uint // Pointer to allow null values
	Coupon       *Coupon
}
