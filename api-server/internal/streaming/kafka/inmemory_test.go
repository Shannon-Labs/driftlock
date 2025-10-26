package kafka

import (
	"context"
	"testing"
	"time"
)

func TestInMemoryBrokerPublishAndSubscribe(t *testing.T) {
	broker := NewInMemoryBroker()
	ctx := context.Background()

	msgs, cancel := broker.Subscribe("events", 1)
	defer cancel()

	msg := Message{Topic: "events", Key: []byte("key"), Value: []byte("value")}
	if err := broker.Publish(ctx, msg); err != nil {
		t.Fatalf("Publish returned error: %v", err)
	}

	select {
	case received := <-msgs:
		if string(received.Value) != "value" {
			t.Fatalf("expected value 'value', got %q", string(received.Value))
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestInMemoryBrokerUnsubscribe(t *testing.T) {
	broker := NewInMemoryBroker()
	ctx := context.Background()

	msgs, cancel := broker.Subscribe("events", 1)
	select {
	case <-msgs:
	default:
	}

	cancel()

	msg := Message{Topic: "events", Value: []byte("value")}
	if err := broker.Publish(ctx, msg); err != nil {
		t.Fatalf("Publish returned error: %v", err)
	}

	select {
	case _, ok := <-msgs:
		if ok {
			t.Fatal("expected channel to be closed")
		}
	default:
	}
}

func TestInMemorySubscriberClose(t *testing.T) {
	broker := NewInMemoryBroker()
	ctx := context.Background()

	subscriber := NewInMemorySubscriber(broker, "events", 1)
	if err := broker.Publish(ctx, Message{Topic: "events", Value: []byte("hello")}); err != nil {
		t.Fatalf("Publish returned error: %v", err)
	}

	select {
	case <-subscriber.Messages():
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for message before close")
	}

	if err := subscriber.Close(ctx); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	select {
	case _, ok := <-subscriber.Messages():
		if ok {
			t.Fatal("expected subscriber channel to be closed")
		}
	default:
	}
}
