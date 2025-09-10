package database

import (
	"auto-vending-system/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect initializes the database connection and migrates the schema.
func Connect() {
	var err error
	DB, err = gorm.Open(sqlite.Open("vending.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Product{},
		&models.ProductStock{},
		&models.Order{},
		&models.Coupon{},
		&models.Setting{},
	)
	if err != nil {
		panic("failed to migrate database")
	}
}
