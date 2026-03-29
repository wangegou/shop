package handlers

import (
	"auto-vending-system/internal/database"
	"auto-vending-system/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type StockInput struct {
	Keys []string `json:"keys" binding:"required"`
}

// AddStock adds new stock (keys) to a product.
func AddStock(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var input StockInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var stockItems []models.ProductStock
	for _, key := range input.Keys {
		stockItems = append(stockItems, models.ProductStock{
			ProductID: uint(productID),
			Key:       key,
		})
	}

	if err := database.DB.Create(&stockItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add stock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock added successfully"})
}

// GetStockForProduct lists all stock for a given product.
func GetStockForProduct(c *gin.Context) {
	productID := c.Param("id")
	var stock []models.ProductStock
	database.DB.Where("product_id = ?", productID).Find(&stock)
	c.JSON(http.StatusOK, stock)
}
