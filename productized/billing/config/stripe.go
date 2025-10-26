package config

import (
	"log"
	"os"

	"github.com/stripe/stripe-go/v75"
)

type StripeConfig struct {
	SecretKey       string
	PublishableKey  string
	WebhookSecret   string
	SuccessURL      string
	CancelURL       string
	Currency        string
	FreePlanID      string
	ProPlanID       string
	BusinessPlanID  string
	UsageEnabled    bool
}

func LoadStripeConfig() *StripeConfig {
	return &StripeConfig{
		SecretKey:       os.Getenv("STRIPE_SECRET_KEY"),
		PublishableKey:  os.Getenv("STRIPE_PUBLISHABLE_KEY"),
		WebhookSecret:   os.Getenv("STRIPE_WEBHOOK_SECRET"),
		SuccessURL:      os.Getenv("STRIPE_SUCCESS_URL"),
		CancelURL:       os.Getenv("STRIPE_CANCEL_URL"),
		Currency:        os.Getenv("STRIPE_CURRENCY"),
		FreePlanID:      os.Getenv("STRIPE_FREE_PLAN_ID"),
		ProPlanID:       os.Getenv("STRIPE_PRO_PLAN_ID"),
		BusinessPlanID:  os.Getenv("STRIPE_BUSINESS_PLAN_ID"),
		UsageEnabled:    os.Getenv("STRIPE_USAGE_ENABLED") == "true",
	}
}

func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		log.Fatal("Stripe secret key is not set")
	}
}