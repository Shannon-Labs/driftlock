package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type emailService struct {
	client      *sendgrid.Client
	fromAddress string
	fromName    string
}

func newEmailService() *emailService {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		log.Println("WARNING: SENDGRID_API_KEY not set, email service disabled")
		return nil
	}
	return &emailService{
		client:      sendgrid.NewSendClient(apiKey),
		fromAddress: env("EMAIL_FROM_ADDRESS", "noreply@driftlock.net"),
		fromName:    env("EMAIL_FROM_NAME", "Driftlock"),
	}
}

func (s *emailService) sendWelcomeEmail(toEmail, companyName, apiKey string) {
	if s == nil {
		log.Printf("MOCK EMAIL: Welcome to %s (%s). API Key: %s...", toEmail, companyName, apiKey[:5])
		return
	}

	from := mail.NewEmail(s.fromName, s.fromAddress)
	subject := "Welcome to Driftlock!"
	to := mail.NewEmail(companyName, toEmail)
	
	plainTextContent := fmt.Sprintf(`Welcome to Driftlock, %s!

Your API Key is: %s

You can use this key to start sending events to our API.
Documentation: https://driftlock.net/docs

Happy Detecting!
The Driftlock Team`, companyName, apiKey)

	htmlContent := fmt.Sprintf(`
		<div style="font-family: sans-serif; color: #333;">
			<h2>Welcome to Driftlock!</h2>
			<p>Hi %s,</p>
			<p>Thanks for signing up. Here is your API key to get started:</p>
			<div style="background: #f4f4f4; padding: 15px; border-radius: 5px; font-family: monospace; margin: 20px 0;">
				%s
			</div>
			<p>You can view our <a href="https://driftlock.net/docs">documentation</a> to learn how to integrate.</p>
			<p>Happy Detecting!<br>The Driftlock Team</p>
		</div>
	`, companyName, apiKey)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	
	go func() {
		response, err := s.client.Send(message)
		if err != nil {
			log.Printf("ERROR: Failed to send welcome email to %s: %v", toEmail, err)
		} else if response.StatusCode >= 400 {
			log.Printf("ERROR: SendGrid API returned %d for welcome email to %s: %s", response.StatusCode, toEmail, response.Body)
		} else {
			log.Printf("Sent welcome email to %s", toEmail)
		}
	}()
}

func (s *emailService) sendVerificationEmail(toEmail, companyName, token string) {
	if s == nil {
		log.Printf("MOCK EMAIL: Verification for %s (%s). Token: %s", toEmail, companyName, token)
		return
	}

	verifyLink := fmt.Sprintf("https://driftlock.net/verify?token=%s", token)

	from := mail.NewEmail(s.fromName, s.fromAddress)
	subject := "Verify your Driftlock account"
	to := mail.NewEmail(companyName, toEmail)
	
	plainTextContent := fmt.Sprintf(`Welcome to Driftlock!

Please verify your email address by clicking the link below:
%s

If you didn't sign up for Driftlock, you can ignore this email.`, verifyLink)

	htmlContent := fmt.Sprintf(`
		<div style="font-family: sans-serif; color: #333;">
			<h2>Verify your email</h2>
			<p>Welcome to Driftlock! Please verify your email address to activate your account.</p>
			<p>
				<a href="%s" style="background: #2563eb; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; display: inline-block;">Verify Email</a>
			</p>
			<p style="font-size: 12px; color: #666;">Or paste this link in your browser: %s</p>
		</div>
	`, verifyLink, verifyLink)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	go func() {
		response, err := s.client.Send(message)
		if err != nil {
			log.Printf("ERROR: Failed to send verification email to %s: %v", toEmail, err)
		} else if response.StatusCode >= 400 {
			log.Printf("ERROR: SendGrid API returned %d for verification email to %s: %s", response.StatusCode, toEmail, response.Body)
		} else {
			log.Printf("Sent verification email to %s", toEmail)
		}
	}()
}

