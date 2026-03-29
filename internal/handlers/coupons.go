package handlers

import (
	"auto-vending-system/internal/database"
	"auto-vending-system/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ShowCouponsPage renders the page that lists all coupons with pagination.
func ShowCouponsPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 10

	var coupons []models.Coupon
	var total int64

	database.DB.Model(&models.Coupon{}).Count(&total)
	offset := (page - 1) * limit
	database.DB.Offset(offset).Limit(limit).Find(&coupons)

	c.HTML(http.StatusOK, "coupons.html", gin.H{
		"title":      "Coupons",
		"active_nav": "coupons",
		"coupons":    coupons,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"prev_page":  page - 1,
		"next_page":  page + 1,
		"has_prev":   page > 1,
		"has_next":   (int64(page) * int64(limit)) < total,
	})
}

// ShowNewCouponForm renders the form to create a new coupon.
func ShowNewCouponForm(c *gin.Context) {
	c.HTML(http.StatusOK, "coupon_form.html", gin.H{
		"title":      "New Coupon",
		"active_nav": "coupons",
		"coupon":     models.Coupon{},
	})
}

// ShowEditCouponForm renders the form to edit an existing coupon.
func ShowEditCouponForm(c *gin.Context) {
	id := c.Param("id")
	var coupon models.Coupon
	if err := database.DB.First(&coupon, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Error",
			"error": "Coupon not found",
		})
		return
	}
	c.HTML(http.StatusOK, "coupon_form.html", gin.H{
		"title":      "Edit Coupon",
		"active_nav": "coupons",
		"coupon":     coupon,
	})
}


type CouponInput struct {
	Code     string  `json:"code" binding:"required"`
	Discount float64 `json:"discount" binding:"required"`
}

// CreateCoupon adds a new coupon.
func CreateCoupon(c *gin.Context) {
	var input CouponInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coupon := models.Coupon{Code: input.Code, Discount: input.Discount}
	if err := database.DB.Create(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coupon"})
		return
	}

	c.JSON(http.StatusOK, coupon)
}

// GetCoupons lists all coupons.
func GetCoupons(c *gin.Context) {
	var coupons []models.Coupon
	database.DB.Find(&coupons)
	c.JSON(http.StatusOK, coupons)
}

// GetCoupon finds a single coupon.
func GetCoupon(c *gin.Context) {
	id := c.Param("id")
	var coupon models.Coupon
	if err := database.DB.First(&coupon, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}
	c.JSON(http.StatusOK, coupon)
}

// UpdateCoupon updates a coupon.
func UpdateCoupon(c *gin.Context) {
	id := c.Param("id")
	var coupon models.Coupon
	if err := database.DB.First(&coupon, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	var input CouponInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&coupon).Updates(models.Coupon{Code: input.Code, Discount: input.Discount})
	c.JSON(http.StatusOK, coupon)
}

// DeleteCoupon deletes a coupon.
func DeleteCoupon(c *gin.Context) {
	id := c.Param("id")
	var coupon models.Coupon
	if err := database.DB.First(&coupon, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	database.DB.Delete(&coupon)
	c.JSON(http.StatusOK, gin.H{"message": "Coupon deleted"})
}
