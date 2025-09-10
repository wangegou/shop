package handlers

import (
	"auto-vending-system/internal/database"
	"auto-vending-system/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ShowOrdersPage renders the page that lists all orders with pagination.
func ShowOrdersPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 10 // Items per page

	var orders []models.Order
	var total int64

	database.DB.Model(&models.Order{}).Count(&total)
	offset := (page - 1) * limit
	database.DB.Preload("Product").Order("created_at desc").Offset(offset).Limit(limit).Find(&orders)

	c.HTML(http.StatusOK, "orders.html", gin.H{
		"title":      "Orders",
		"active_nav": "orders",
		"orders":     orders,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"prev_page":  page - 1,
		"next_page":  page + 1,
		"has_prev":   page > 1,
		"has_next":   (int64(page) * int64(limit)) < total,
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
