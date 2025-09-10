package email

import (
	"auto-vending-system/internal/database"
	"auto-vending-system/internal/models"
	"gopkg.in/gomail.v2"
	"strconv"
)

// SendEmail constructs and sends an email using SMTP settings from the database.
func SendEmail(to, subject, body string) error {
	// Fetch settings from the database
	settings, err := getEmailSettings()
	if err != nil {
		return err
	}

	port, _ := strconv.Atoi(settings["mail_port"])

	m := gomail.NewMessage()
	m.SetHeader("From", settings["mail_from_address"])
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(settings["mail_host"], port, settings["mail_username"], settings["mail_password"])

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// getEmailSettings retrieves all email-related settings from the database.
func getEmailSettings() (map[string]string, error) {
	var settings []models.Setting
	keys := []string{"mail_host", "mail_port", "mail_username", "mail_password", "mail_from_address"}
	if err := database.DB.Where("key IN ?", keys).Find(&settings).Error; err != nil {
		return nil, err
	}

	settingsMap := make(map[string]string)
	for _, s := range settings {
		settingsMap[s.Key] = s.Value
	}
	return settingsMap, nil
}
