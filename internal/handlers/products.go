package handlers

import (
	"auto-vending-system/internal/database"
	"auto-vending-system/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ShowProductsPage renders the page that lists all products with pagination.
func ShowProductsPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 10 // Items per page

	var products []models.Product
	var total int64

	database.DB.Model(&models.Product{}).Count(&total)
	offset := (page - 1) * limit
	database.DB.Preload("Category").Offset(offset).Limit(limit).Find(&products)

	c.HTML(http.StatusOK, "products.html", gin.H{
		"title":       "Products",
		"active_nav":  "products",
		"products":    products,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"prev_page":   page - 1,
		"next_page":   page + 1,
		"has_prev":    page > 1,
		"has_next":    (int64(page) * int64(limit)) < total,
	})
}

// ShowNewProductForm renders the form to create a new product.
func ShowNewProductForm(c *gin.Context) {
	var categories []models.Category
	database.DB.Find(&categories)
	c.HTML(http.StatusOK, "product_form.html", gin.H{
		"title":      "New Product",
		"active_nav": "products",
		"product":    models.Product{},
		"categories": categories,
	})
}

// ShowEditProductForm renders the form to edit an existing product.
func ShowEditProductForm(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Error",
			"error": "Product not found",
		})
		return
	}
	var categories []models.Category
	database.DB.Find(&categories)
	c.HTML(http.StatusOK, "product_form.html", gin.H{
		"title":      "Edit Product",
		"active_nav": "products",
		"product":    product,
		"categories": categories,
	})
}

// ShowStockPage renders the page to manage stock for a product.
func ShowStockPage(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Error",
			"error": "Product not found",
		})
		return
	}
	var stock []models.ProductStock
	database.DB.Where("product_id = ?", id).Find(&stock)
	c.HTML(http.StatusOK, "stock.html", gin.H{
		"title":      "Manage Stock",
		"active_nav": "products",
		"product":    product,
		"stock":      stock,
	})
}


type ProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

// CreateProduct adds a new product.
func CreateProduct(c *gin.Context) {
	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProducts lists all products.
func GetProducts(c *gin.Context) {
	var products []models.Product
	database.DB.Preload("Category").Find(&products)
	c.JSON(http.StatusOK, products)
}

// GetProduct finds a single product.
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.Preload("Category").Preload("Stock").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// UpdateProduct updates a product.
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&product).Updates(models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		CategoryID:  input.CategoryID,
	})

	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product.
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Also delete associated stock
	database.DB.Where("product_id = ?", id).Delete(&models.ProductStock{})
	database.DB.Delete(&product)

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
