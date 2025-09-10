package handlers

import (
	"auto-vending-system/internal/database"
	"auto-vending-system/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ShowOrdersPage renders the page that lists all orders.
func ShowOrdersPage(c *gin.Context) {
	var orders []models.Order
	database.DB.Preload("Product").Order("created_at desc").Find(&orders)
	c.HTML(http.StatusOK, "orders.html", gin.H{
		"title":      "Orders",
		"active_nav": "orders",
		"orders":     orders,
	})
}

// GetOrders lists all orders.
func GetOrders(c *gin.Context) {
	var orders []models.Order
	database.DB.Preload("Product").Preload("Coupon").Order("created_at desc").Find(&orders)
	c.JSON(http.StatusOK, orders)
}

// GetOrder finds a single order.
func GetOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Preload("Product").Preload("Coupon").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}
