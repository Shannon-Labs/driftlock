package streaming

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/Shannon-Labs/driftlock/api-server/internal/models"
	"github.com/Shannon-Labs/driftlock/api-server/internal/streaming/kafka"
)

// EventPublisher publishes domain events to downstream systems.
type EventPublisher interface {
	AnomalyCreated(ctx context.Context, anomaly *models.Anomaly) error
}

// KafkaEventPublisher publishes events via Kafka.
type KafkaEventPublisher struct {
	producer kafka.Publisher
	topic    string
}

// NewKafkaEventPublisher returns a Kafka-backed event publisher.
func NewKafkaEventPublisher(producer kafka.Publisher, topic string) *KafkaEventPublisher {
	return &KafkaEventPublisher{producer: producer, topic: topic}
}

// AnomalyCreated emits an anomaly-created event.
func (p *KafkaEventPublisher) AnomalyCreated(ctx context.Context, anomaly *models.Anomaly) error {
	if p == nil || p.producer == nil || anomaly == nil {
		return nil
	}

	payload := struct {
		Type    string          `json:"type"`
		Version string          `json:"version"`
		Anomaly *models.Anomaly `json:"anomaly"`
	}{
		Type:    "anomaly.created",
		Version: "v1",
		Anomaly: anomaly,
	}

	value, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal anomaly event: %w", err)
	}

	msg := kafka.Message{
		Topic: p.topic,
		Key:   []byte(anomaly.ID.String()),
		Value: value,
		Headers: map[string]string{
			"driftlock-event-type": payload.Type,
		},
	}

	return p.producer.Publish(ctx, msg)
}

// InMemoryPublisher provides an in-process publisher primarily for tests.
type InMemoryPublisher struct {
	broker    *kafka.InMemoryBroker
	topic     string
	subscribe <-chan kafka.Message
	closeFn   func()
}

// NewInMemoryPublisher creates an event publisher backed by InMemoryBroker.
func NewInMemoryPublisher(broker *kafka.InMemoryBroker, topic string) *InMemoryPublisher {
	msgs, cancel := broker.Subscribe(topic, 64)
	return &InMemoryPublisher{
		broker:    broker,
		topic:     topic,
		subscribe: msgs,
		closeFn:   cancel,
	}
}

// Publish forwards to the underlying broker.
func (p *InMemoryPublisher) Publish(ctx context.Context, msg kafka.Message) error {
	if msg.Topic == "" {
		msg.Topic = p.topic
	}
	if len(msg.Key) == 0 {
		msg.Key = []byte(uuid.New().String())
	}
	return p.broker.Publish(ctx, msg)
}

// Messages returns the subscription channel for tests.
func (p *InMemoryPublisher) Messages() <-chan kafka.Message {
	return p.subscribe
}

// Close releases the subscription.
func (p *InMemoryPublisher) Close() {
	if p.closeFn != nil {
		p.closeFn()
	}
}

// AnomalyCreated implements EventPublisher interface by publishing anomaly events.
func (p *InMemoryPublisher) AnomalyCreated(ctx context.Context, anomaly *models.Anomaly) error {
	if p == nil || anomaly == nil {
		return nil
	}

	payload := struct {
		Type    string          `json:"type"`
		Version string          `json:"version"`
		Anomaly *models.Anomaly `json:"anomaly"`
	}{
		Type:    "anomaly.created",
		Version: "v1",
		Anomaly: anomaly,
	}

	value, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal anomaly event: %w", err)
	}

	msg := kafka.Message{
		Topic: p.topic,
		Key:   []byte(anomaly.ID.String()),
		Value: value,
		Headers: map[string]string{
			"driftlock-event-type": payload.Type,
		},
	}

	return p.Publish(ctx, msg)
}
