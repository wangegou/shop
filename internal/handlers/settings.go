package handlers

import (
	"auto-vending-system/internal/database"
	"auto-vending-system/internal/email"
	"auto-vending-system/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ShowSettingsPage renders the settings page.
func ShowSettingsPage(c *gin.Context) {
	var settings []models.Setting
	database.DB.Find(&settings)

	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	c.HTML(http.StatusOK, "settings.html", gin.H{
		"title":      "Settings",
		"active_nav": "settings",
		"settings":   settingsMap,
	})
}

// GetSettings retrieves all settings.
func GetSettings(c *gin.Context) {
	var settings []models.Setting
	database.DB.Find(&settings)

	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}

	c.JSON(http.StatusOK, settingsMap)
}

// UpdateSettings updates a batch of settings.
func UpdateSettings(c *gin.Context) {
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for key, value := range input {
		setting := models.Setting{Key: key}
		database.DB.FirstOrCreate(&setting, models.Setting{Key: key})
		database.DB.Model(&setting).Update("value", value)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

type EmailTestInput struct {
	To string `json:"to" binding:"required,email"`
}

// EmailTest sends a real test email to the specified address.
func EmailTest(c *gin.Context) {
	var input EmailTestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address provided."})
		return
	}

	err := email.SendEmail(input.To, "Test Email", "This is a test email from the vending system.")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send test email: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test email sent successfully to " + input.To + "!"})
}
