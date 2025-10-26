package kafka

import (
	"context"
	"errors"
	"testing"

	segment "github.com/segmentio/kafka-go"
)

type mockWriter struct {
	messages []segment.Message
	err      error
	closed   bool
}

func (m *mockWriter) WriteMessages(_ context.Context, msgs ...segment.Message) error {
	if m.err != nil {
		return m.err
	}
	m.messages = append(m.messages, msgs...)
	return nil
}

func (m *mockWriter) Close() error {
	m.closed = true
	return nil
}

func TestProducerPublish(t *testing.T) {
	producer, err := NewProducer(ProducerConfig{
		Brokers: []string{"localhost:9092"},
	})
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}

	mock := &mockWriter{}
	producer.writerMaker = func(topic string) (writer, error) {
		if topic != "anomaly-events" {
			t.Fatalf("expected topic anomaly-events, got %s", topic)
		}
		return mock, nil
	}

	err = producer.Publish(context.Background(), Message{
		Topic:   "anomaly-events",
		Key:     []byte("key"),
		Value:   []byte("value"),
		Headers: map[string]string{"foo": "bar"},
	})
	if err != nil {
		t.Fatalf("expected no error publishing: %v", err)
	}

	if len(mock.messages) != 1 {
		t.Fatalf("expected one message, got %d", len(mock.messages))
	}

	msg := mock.messages[0]
	if msg.Topic != "anomaly-events" {
		t.Errorf("expected topic anomaly-events, got %s", msg.Topic)
	}
	if string(msg.Key) != "key" {
		t.Errorf("expected key 'key', got %s", msg.Key)
	}
	if string(msg.Value) != "value" {
		t.Errorf("expected value 'value', got %s", msg.Value)
	}
	if len(msg.Headers) != 1 {
		t.Fatalf("expected one header, got %d", len(msg.Headers))
	}
	if msg.Headers[0].Key != "foo" || string(msg.Headers[0].Value) != "bar" {
		t.Errorf("unexpected header payload: %+v", msg.Headers[0])
	}
	if msg.Time.IsZero() {
		t.Errorf("expected timestamp to be set")
	}
}

func TestProducerPublishError(t *testing.T) {
	producer, err := NewProducer(ProducerConfig{
		Brokers: []string{"localhost:9092"},
	})
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}

	producer.writerMaker = func(topic string) (writer, error) {
		return nil, errors.New("boom")
	}

	err = producer.Publish(context.Background(), Message{
		Topic: "events",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestProducerClose(t *testing.T) {
	producer, err := NewProducer(ProducerConfig{
		Brokers: []string{"localhost:9092"},
	})
	if err != nil {
		t.Fatalf("expected no error: %v", err)
	}

	mock := &mockWriter{}
	producer.writerMaker = func(topic string) (writer, error) {
		return mock, nil
	}

	if err := producer.Publish(context.Background(), Message{Topic: "events"}); err != nil {
		t.Fatalf("unexpected publish error: %v", err)
	}

	if err := producer.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}

	if !mock.closed {
		t.Fatalf("expected writer to be closed")
	}
}

func TestNewProducerRequiresBroker(t *testing.T) {
	_, err := NewProducer(ProducerConfig{})
	if err == nil {
		t.Fatalf("expected error for missing brokers")
	}
}
