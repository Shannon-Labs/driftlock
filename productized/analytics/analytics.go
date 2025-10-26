package analytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"driftlock/productized/analytics/config"
)

type AnalyticsService struct {
	config *config.AnalyticsConfig
}

type Event struct {
	ClientID   string                 `json:"client_id"`
	UserID     string                 `json:"user_id"`
	EventName  string                 `json:"event_name"`
	Parameters map[string]interface{} `json:"parameters"`
	Timestamp  time.Time              `json:"timestamp"`
}

func NewAnalyticsService(config *config.AnalyticsConfig) *AnalyticsService {
	return &AnalyticsService{
		config: config,
	}
}

func (a *AnalyticsService) TrackEvent(event Event) error {
	if !a.config.EnableTracking {
		return nil // Analytics disabled
	}

	// Send to Google Analytics 4 if configured
	if a.config.GA4MeasurementID != "" && a.config.GA4APIKey != "" {
		if err := a.sendToGA4(event); err != nil {
			log.Printf("Failed to send event to GA4: %v", err)
			// Continue even if GA4 fails
		}
	}

	// Send to custom analytics backend if needed
	if err := a.sendToCustomBackend(event); err != nil {
		log.Printf("Failed to send event to custom backend: %v", err)
	}

	return nil
}

func (a *AnalyticsService) sendToGA4(event Event) error {
	// GA4 Measurement Protocol
	ga4URL := fmt.Sprintf("https://www.google-analytics.com/mp/collect?measurement_id=%s&api_secret=%s", 
		a.config.GA4MeasurementID, a.config.GA4APIKey)

	// GA4 Event payload
	ga4Payload := map[string]interface{}{
		"client_id": event.ClientID,
		"user_id":   event.UserID,
		"events": []map[string]interface{}{
			{
				"name": event.EventName,
				"params": event.Parameters,
			},
		},
	}

	payloadBytes, err := json.Marshal(ga4Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal GA4 payload: %v", err)
	}

	resp, err := http.Post(ga4URL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send GA4 event: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GA4 returned status: %d", resp.StatusCode)
	}

	return nil
}

func (a *AnalyticsService) sendToCustomBackend(event Event) error {
	// In a real implementation, this would send to your custom analytics backend
	// For now, we'll just log the event
	log.Printf("Custom analytics event: %s for user %s", event.EventName, event.UserID)
	return nil
}

// TrackUserAction tracks user actions in the application
func (a *AnalyticsService) TrackUserAction(userID, action string, metadata map[string]interface{}) error {
	event := Event{
		UserID:    userID,
		EventName: "user_action",
		Parameters: map[string]interface{}{
			"action":   action,
			"metadata": metadata,
		},
		Timestamp: time.Now(),
	}

	return a.TrackEvent(event)
}

// TrackAPIUsage tracks API usage for analytics
func (a *AnalyticsService) TrackAPIUsage(userID, endpoint string, responseTime time.Duration) error {
	event := Event{
		UserID:    userID,
		EventName: "api_usage",
		Parameters: map[string]interface{}{
			"endpoint":      endpoint,
			"response_time": responseTime.Milliseconds(),
		},
		Timestamp: time.Now(),
	}

	return a.TrackEvent(event)
}

// TrackAnomalyDetection tracks anomaly detection events
func (a *AnalyticsService) TrackAnomalyDetection(userID string, anomalyID uint, severity string) error {
	event := Event{
		UserID:    userID,
		EventName: "anomaly_detected",
		Parameters: map[string]interface{}{
			"anomaly_id": anomalyID,
			"severity":   severity,
		},
		Timestamp: time.Now(),
	}

	return a.TrackEvent(event)
}