package kafka

import (
	"context"
	"time"
)

// Message represents a Kafka record.
type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Headers   map[string]string
	Timestamp time.Time
}

// Publisher publishes messages to Kafka.
type Publisher interface {
	Publish(ctx context.Context, msg Message) error
}

// Subscriber consumes messages from Kafka.
type Subscriber interface {
	Messages() <-chan Message
	Close(ctx context.Context) error
}

// Broker provides topic subscription and publication primitives.
type Broker interface {
	Publisher
	Subscribe(topic string, buffer int) (<-chan Message, func())
}
