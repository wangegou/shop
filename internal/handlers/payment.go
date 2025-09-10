package handlers

import (
	"auto-vending-system/internal/database"
	"auto-vending-system/internal/email"
	"auto-vending-system/internal/models"
	"auto-vending-system/internal/payment"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// RedirectToPayment generates the payment URL and redirects the user.
func RedirectToPayment(c *gin.Context) {
	orderNo := c.Param("order_no")
	var order models.Order
	if err := database.DB.Preload("Product").Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"title": "Error", "error": "Order not found"})
		return
	}

	var pid, key models.Setting
	database.DB.Where("key = ?", "payment_pid").First(&pid)
	database.DB.Where("key = ?", "payment_key").First(&key)

	if pid.Value == "" || key.Value == "" {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"title": "Error", "error": "Payment gateway not configured."})
		return
	}

	baseURL := "http://localhost:8080"
	params := map[string]string{
		"pid":          pid.Value,
		"out_trade_no": order.OrderNo,
		"notify_url":   baseURL + "/payment/notify",
		"return_url":   baseURL + "/status/" + order.OrderNo,
		"name":         order.Product.Name,
		"money":        fmt.Sprintf("%.2f", order.Total),
		"type":         "alipay",
		"sign_type":    "MD5",
	}
	params["sign"] = payment.GenerateSign(params, key.Value)

	paymentURL := "https://pay.v8jisu.cn/submit.php"
	var queryParts []string
	for k, v := range params {
		queryParts = append(queryParts, url.QueryEscape(k)+"="+url.QueryEscape(v))
	}
	fullURL := paymentURL + "?" + strings.Join(queryParts, "&")

	c.Redirect(http.StatusFound, fullURL)
}

// NotifyHandler handles the asynchronous callback from the payment gateway.
func NotifyHandler(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	params := make(map[string]string)
	for k, v := range queryParams {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}

	receivedSign := params["sign"]
	if receivedSign == "" {
		c.String(http.StatusBadRequest, "fail: missing sign")
		return
	}
	delete(params, "sign")
	delete(params, "sign_type")

	var key models.Setting
	database.DB.Where("key = ?", "payment_key").First(&key)
	if key.Value == "" {
		c.String(http.StatusInternalServerError, "fail: payment key not configured")
		return
	}

	if receivedSign != payment.GenerateSign(params, key.Value) {
		c.String(http.StatusBadRequest, "fail: invalid sign")
		return
	}

	if params["trade_status"] == "TRADE_SUCCESS" {
		orderNo := params["out_trade_no"]
		var order models.Order
		if err := database.DB.Preload("Product").Where("order_no = ?", orderNo).First(&order).Error; err == nil {
			if order.Status == "paid" {
				c.String(http.StatusOK, "success") // Already processed
				return
			}

			tx := database.DB.Begin()
			if err := tx.Model(&order).Update("status", "paid").Error; err != nil {
				tx.Rollback()
				return
			}
			if err := tx.Model(&models.ProductStock{}).Where("key = ?", order.PurchasedKey).Update("is_sold", true).Error; err != nil {
				tx.Rollback()
				return
			}
			if order.CouponID != nil {
				if err := tx.Model(&models.Coupon{}).Where("id = ?", *order.CouponID).Update("is_used", true).Error; err != nil {
					tx.Rollback()
					return
				}
			}
			tx.Commit()

			// Send email notification
			var emailTemplateSetting models.Setting
			database.DB.Where("key = ?", "email_template").First(&emailTemplateSetting)
			emailBody := strings.Replace(emailTemplateSetting.Value, "{{key}}", order.PurchasedKey, -1)
			emailBody = strings.Replace(emailBody, "{{product_name}}", order.Product.Name, -1)

			err := email.SendEmail(order.Email, "Your Purchase from Our Store", emailBody)
			if err != nil {
				log.Printf("Failed to send email for order %s: %v", order.OrderNo, err)
			}

			c.String(http.StatusOK, "success")
			return
		}
	}

	c.String(http.StatusOK, "fail")
}
