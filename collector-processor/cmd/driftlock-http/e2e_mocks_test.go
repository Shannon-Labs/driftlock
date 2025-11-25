package main

import (
	"sync"
	"time"
)

// mockEmailService captures emails instead of sending them
type mockEmailService struct {
	mu                  sync.Mutex
	verificationEmails  []sentEmail
	welcomeEmails       []sentEmail
	trialEndingEmails   []sentEmail
	paymentFailedEmails []sentEmail
}

type sentEmail struct {
	To          string
	CompanyName string
	Token       string // for verification emails
	APIKey      string // for welcome emails
	Days        int    // for trial ending emails
	GraceEnd    time.Time
	SentAt      time.Time
}

func newMockEmailService() *mockEmailService {
	return &mockEmailService{}
}

// sendVerificationEmail captures the verification token
func (m *mockEmailService) sendVerificationEmail(toEmail, companyName, token string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.verificationEmails = append(m.verificationEmails, sentEmail{
		To:          toEmail,
		CompanyName: companyName,
		Token:       token,
		SentAt:      time.Now(),
	})
}

// sendWelcomeEmail captures the welcome email
func (m *mockEmailService) sendWelcomeEmail(toEmail, companyName, apiKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.welcomeEmails = append(m.welcomeEmails, sentEmail{
		To:          toEmail,
		CompanyName: companyName,
		APIKey:      apiKey,
		SentAt:      time.Now(),
	})
}

// sendTrialEndingEmail captures trial ending emails
func (m *mockEmailService) sendTrialEndingEmail(toEmail, companyName string, daysRemaining int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.trialEndingEmails = append(m.trialEndingEmails, sentEmail{
		To:          toEmail,
		CompanyName: companyName,
		Days:        daysRemaining,
		SentAt:      time.Now(),
	})
}

// sendPaymentFailedEmail captures payment failed emails
func (m *mockEmailService) sendPaymentFailedEmail(toEmail, companyName string, gracePeriodEnd time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.paymentFailedEmails = append(m.paymentFailedEmails, sentEmail{
		To:          toEmail,
		CompanyName: companyName,
		GraceEnd:    gracePeriodEnd,
		SentAt:      time.Now(),
	})
}

// getLastVerificationToken returns the last verification token sent
func (m *mockEmailService) getLastVerificationToken() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.verificationEmails) == 0 {
		return ""
	}
	return m.verificationEmails[len(m.verificationEmails)-1].Token
}

// getVerificationTokenFor returns the verification token for a specific email
func (m *mockEmailService) getVerificationTokenFor(email string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := len(m.verificationEmails) - 1; i >= 0; i-- {
		if m.verificationEmails[i].To == email {
			return m.verificationEmails[i].Token
		}
	}
	return ""
}

// getLastWelcomeAPIKey returns the API key from the last welcome email
func (m *mockEmailService) getLastWelcomeAPIKey() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.welcomeEmails) == 0 {
		return ""
	}
	return m.welcomeEmails[len(m.welcomeEmails)-1].APIKey
}

// clear resets all captured emails
func (m *mockEmailService) clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.verificationEmails = nil
	m.welcomeEmails = nil
	m.trialEndingEmails = nil
	m.paymentFailedEmails = nil
}

// emailerInterface allows both real and mock email services
type emailerInterface interface {
	sendVerificationEmail(toEmail, companyName, token string)
	sendWelcomeEmail(toEmail, companyName, apiKey string)
	sendTrialEndingEmail(toEmail, companyName string, daysRemaining int)
	sendPaymentFailedEmail(toEmail, companyName string, gracePeriodEnd time.Time)
}

// Ensure both types implement the interface
var _ emailerInterface = (*emailService)(nil)
var _ emailerInterface = (*mockEmailService)(nil)
