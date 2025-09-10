package handlers

import (
	"auto-vending-system/internal/database"
	"auto-vending-system/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CategoryInput struct {
	Name string `json:"name" binding:"required"`
}

// ShowCategoriesPage renders the page that lists all categories with pagination.
func ShowCategoriesPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 10

	var categories []models.Category
	var total int64

	database.DB.Model(&models.Category{}).Count(&total)
	offset := (page - 1) * limit
	database.DB.Offset(offset).Limit(limit).Find(&categories)

	c.HTML(http.StatusOK, "categories.html", gin.H{
		"title":      "Categories",
		"active_nav": "categories",
		"categories": categories,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"prev_page":  page - 1,
		"next_page":  page + 1,
		"has_prev":   page > 1,
		"has_next":   (int64(page) * int64(limit)) < total,
	})
}

// ShowNewCategoryForm renders the form to create a new category.
func ShowNewCategoryForm(c *gin.Context) {
	c.HTML(http.StatusOK, "category_form.html", gin.H{
		"title":      "New Category",
		"active_nav": "categories",
		"category":   models.Category{},
	})
}

// ShowEditCategoryForm renders the form to edit an existing category.
func ShowEditCategoryForm(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Error",
			"error": "Category not found",
		})
		return
	}

	c.HTML(http.StatusOK, "category_form.html", gin.H{
		"title":      "Edit Category",
		"active_nav": "categories",
		"category":   category,
	})
}

// CreateCategory handles the API request to create a new category.
func CreateCategory(c *gin.Context) {
	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{Name: input.Name}
	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// GetCategories lists all categories. (API)
func GetCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

// GetCategory finds a single category. (API)
func GetCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	c.JSON(http.StatusOK, category)
}


// UpdateCategory handles the API request to update a category.
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&category).Update("name", input.Name)
	c.JSON(http.StatusOK, category)
}

// DeleteCategory deletes a category.
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	database.DB.Delete(&category)
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
