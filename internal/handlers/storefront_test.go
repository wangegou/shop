package handlers

import (
	"auto-vending-system/internal/database"
	"auto-vending-system/internal/models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupTestDB initializes an in-memory SQLite database for tests.
func setupTestDB() {
	var err error
	database.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to in-memory database")
	}
	err = database.DB.AutoMigrate(
		&models.User{}, &models.Category{}, &models.Product{},
		&models.ProductStock{}, &models.Order{}, &models.Coupon{}, &models.Setting{},
	)
	if err != nil {
		panic("Failed to migrate test database")
	}
}

func TestCreateOrder_OutOfStock(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	// Setup: Create a product but no stock
	product := models.Product{Name: "Test Product", Price: 10.0}
	database.DB.Create(&product)

	// Prepare request
	orderInput := CreateOrderInput{ProductID: product.ID, Email: "test@example.com"}
	body, _ := json.Marshal(orderInput)
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/orders", CreateOrder)
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "Out of stock")
}

func TestVerifyCoupon(t *testing.T) {
	setupTestDB()
	gin.SetMode(gin.TestMode)

	// Setup: Create a valid, unused coupon
	validCoupon := models.Coupon{Code: "VALID10", Discount: 10, IsUsed: false}
	database.DB.Create(&validCoupon)

	// Setup: Create a used coupon
	usedCoupon := models.Coupon{Code: "USED5", Discount: 5, IsUsed: true}
	database.DB.Create(&usedCoupon)

	r := gin.Default()
	r.GET("/coupons/verify", VerifyCoupon)

	// Test Case 1: Valid coupon
	reqValid, _ := http.NewRequest(http.MethodGet, "/coupons/verify?code=VALID10", nil)
	wValid := httptest.NewRecorder()
	r.ServeHTTP(wValid, reqValid)

	assert.Equal(t, http.StatusOK, wValid.Code)
	assert.Contains(t, wValid.Body.String(), "Coupon is valid")

	// Test Case 2: Used coupon
	reqUsed, _ := http.NewRequest(http.MethodGet, "/coupons/verify?code=USED5", nil)
	wUsed := httptest.NewRecorder()
	r.ServeHTTP(wUsed, reqUsed)

	assert.Equal(t, http.StatusNotFound, wUsed.Code)
	assert.Contains(t, wUsed.Body.String(), "Invalid or used coupon")

	// Test Case 3: Non-existent coupon
	reqInvalid, _ := http.NewRequest(http.MethodGet, "/coupons/verify?code=BOGUS", nil)
	wInvalid := httptest.NewRecorder()
	r.ServeHTTP(wInvalid, reqInvalid)

	assert.Equal(t, http.StatusNotFound, wInvalid.Code)
	assert.Contains(t, wInvalid.Body.String(), "Invalid or used coupon")
}
