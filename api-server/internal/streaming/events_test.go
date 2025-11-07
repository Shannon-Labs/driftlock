package streaming

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/shannon-labs/driftlock/api-server/internal/models"
	"github.com/shannon-labs/driftlock/api-server/internal/streaming/kafka"
)

func TestKafkaEventPublisher_AnomalyCreated(t *testing.T) {
	broker := kafka.NewInMemoryBroker()
	publisher := NewKafkaEventPublisher(broker, "anomalies")

	anomaly := &models.Anomaly{ID: uuid.New()}
	if err := publisher.AnomalyCreated(context.Background(), anomaly); err != nil {
		t.Fatalf("AnomalyCreated returned error: %v", err)
	}

	msgs, cancel := broker.Subscribe("anomalies", 1)
	defer cancel()

	if err := publisher.AnomalyCreated(context.Background(), anomaly); err != nil {
		t.Fatalf("AnomalyCreated returned error: %v", err)
	}

	select {
	case msg := <-msgs:
		if string(msg.Key) != anomaly.ID.String() {
			t.Fatalf("unexpected key: %s", string(msg.Key))
		}
		if msg.Headers["driftlock-event-type"] != "anomaly.created" {
			t.Fatalf("unexpected header: %#v", msg.Headers)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestKafkaEventPublisher_NoProducer(t *testing.T) {
	publisher := NewKafkaEventPublisher(nil, "anomalies")
	if err := publisher.AnomalyCreated(context.Background(), &models.Anomaly{}); err != nil {
		t.Fatalf("expected nil error when producer missing, got %v", err)
	}
}
