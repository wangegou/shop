package handlers

import (
	"auto-vending-system/internal/database"
	"auto-vending-system/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// ShowStorefront renders the main shop page.
func ShowStorefront(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// ShowOrderPage renders the page to place an order for a specific product.
func ShowOrderPage(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.String(http.StatusNotFound, "Product not found")
		return
	}
	c.HTML(http.StatusOK, "order.html", gin.H{
		"product": product,
	})
}

// ShowStatusPage renders the order status page.
func ShowStatusPage(c *gin.Context) {
	orderNo := c.Param("order_no")
	c.HTML(http.StatusOK, "status.html", gin.H{
		"order_no": orderNo,
	})
}

// GetPublicProducts lists all categories with their products for the storefront.
func GetPublicProducts(c *gin.Context) {
	var categories []models.Category
	if err := database.DB.Preload("Products").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

type CreateOrderInput struct {
	ProductID  uint   `json:"product_id" binding:"required"`
	Email      string `json:"email" binding:"required"`
	CouponCode string `json:"coupon_code"`
}

// CreateOrder handles the creation of a new order.
func CreateOrder(c *gin.Context) {
	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	total := product.Price
	var coupon *models.Coupon

	if input.CouponCode != "" {
		var tempCoupon models.Coupon
		if err := database.DB.Where("code = ? AND is_used = ?", input.CouponCode, false).First(&tempCoupon).Error; err == nil {
			total -= tempCoupon.Discount
			if total < 0 {
				total = 0
			}
			coupon = &tempCoupon
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or used coupon"})
			return
		}
	}

	var stock models.ProductStock
	if err := database.DB.Where("product_id = ? AND is_sold = ?", product.ID, false).First(&stock).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Out of stock"})
		return
	}

	order := models.Order{
		OrderNo:      uuid.New().String(),
		ProductID:    product.ID,
		Email:        input.Email,
		Total:        total,
		Status:       "pending",
		PurchasedKey: stock.Key,
	}

	if coupon != nil {
		order.CouponID = &coupon.ID
	}

	tx := database.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	if err := tx.Model(&stock).Update("is_sold", true).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}

	if coupon != nil {
		if err := tx.Model(coupon).Update("is_used", true).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to use coupon"})
			return
		}
	}

	if err := tx.Model(&order).Update("status", "paid").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrderStatus retrieves the status of an order and the purchased key if paid.
func GetOrderStatus(c *gin.Context) {
	orderNo := c.Param("order_no")
	var order models.Order
	if err := database.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	response := gin.H{
		"order_no": order.OrderNo,
		"status":   order.Status,
	}

	if order.Status == "paid" {
		response["key"] = order.PurchasedKey
	}

	c.JSON(http.StatusOK, response)
}

// VerifyCoupon checks if a coupon is valid.
func VerifyCoupon(c *gin.Context) {
	code := c.Query("code")
	var coupon models.Coupon
	if err := database.DB.Where("code = ? AND is_used = ?", code, false).First(&coupon).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or used coupon"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Coupon is valid", "discount": coupon.Discount})
}
