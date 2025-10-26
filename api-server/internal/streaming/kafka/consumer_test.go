package kafka

import "testing"

func TestNewConsumerValidation(t *testing.T) {
	if _, err := NewConsumer(ConsumerConfig{}); err == nil {
		t.Fatalf("expected error when brokers are missing")
	}

	if _, err := NewConsumer(ConsumerConfig{
		Brokers: []string{"localhost:9092"},
	}); err == nil {
		t.Fatalf("expected error when topic is missing")
	}
}
