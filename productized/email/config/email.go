package config

import (
	"os"
	"strconv"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	Provider     string // "smtp", "sendgrid", etc.
	SendGridKey  string
}

func LoadEmailConfig() *EmailConfig {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		port = 587 // Default SMTP port
	}

	return &EmailConfig{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     port,
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("EMAIL_FROM"),
		Provider:     os.Getenv("EMAIL_PROVIDER"),
		SendGridKey:  os.Getenv("SENDGRID_API_KEY"),
	}
}