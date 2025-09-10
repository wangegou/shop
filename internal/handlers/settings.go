package handlers

import (
	"auto-vending-system/internal/database"
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

	// Convert to a map for easier use on the frontend
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
		// Use FirstOrCreate to either update the existing setting or create it if it doesn't exist.
		setting := models.Setting{Key: key}
		database.DB.FirstOrCreate(&setting, models.Setting{Key: key})
		database.DB.Model(&setting).Update("value", value)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

// EmailTest sends a test email (simulation).
func EmailTest(c *gin.Context) {
	// In a real application, this would use the settings from the database
	// to configure an SMTP client and send an email.
	// For this simulation, we just log that the function was called.

	// Example of how you might get a setting:
	// var mailHost models.Setting
	// database.DB.Where("key = ?", "mail_host").First(&mailHost)
	// log.Printf("Simulating sending test email to host: %s", mailHost.Value)

	c.JSON(http.StatusOK, gin.H{"message": "Test email sent (simulation)!"})
}
