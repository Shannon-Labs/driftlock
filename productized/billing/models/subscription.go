package models

import (
	"time"

	"gorm.io/gorm"
)

type SubscriptionPlan struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"` // Price in cents
	Currency    string `json:"currency"`
	Interval    string `json:"interval"` // month, year
	Features    []Feature `json:"features" gorm:"foreignKey:PlanID"`
}

type Feature struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	PlanID    string `json:"plan_id"`
	Name      string `json:"name"`
	Quota     int    `json:"quota"` // -1 for unlimited
	Available bool   `json:"available"`
}

type UserSubscription struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	UserID         uint           `json:"user_id"`
	PlanID         string         `json:"plan_id"`
	SubscriptionID string         `json:"subscription_id"` // Stripe subscription ID
	CustomerID     string         `json:"customer_id"`     // Stripe customer ID
	Status         string         `json:"status"`          // active, canceled, past_due, etc.
	CurrentPeriodStart time.Time  `json:"current_period_start"`
	CurrentPeriodEnd   time.Time  `json:"current_period_end"`
	TrialEnd       *time.Time     `json:"trial_end,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type UsageRecord struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         uint      `json:"user_id"`
	SubscriptionID string    `json:"subscription_id"`
	FeatureName    string    `json:"feature_name"`
	Quantity       int64     `json:"quantity"`
	Timestamp      time.Time `json:"timestamp"`
}