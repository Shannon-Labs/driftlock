package kafka

import (
	"context"
	"sync"
	"time"
)

// InMemoryBroker is a simple broker implementation that keeps messages in-process.
type InMemoryBroker struct {
	mu          sync.RWMutex
	subscribers map[string][]chan Message
}

// NewInMemoryBroker creates a new broker instance.
func NewInMemoryBroker() *InMemoryBroker {
	return &InMemoryBroker{
		subscribers: make(map[string][]chan Message),
	}
}

// Publish sends the message to all subscribers of the topic.
func (b *InMemoryBroker) Publish(ctx context.Context, msg Message) error {
	if msg.Topic == "" {
		return nil
	}

	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	b.mu.RLock()
	subs := append([]chan Message(nil), b.subscribers[msg.Topic]...)
	b.mu.RUnlock()

	for _, ch := range subs {
		select {
		case ch <- msg:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// Subscribe registers a subscriber for the given topic.
func (b *InMemoryBroker) Subscribe(topic string, buffer int) (<-chan Message, func()) {
	if buffer <= 0 {
		buffer = 1
	}

	ch := make(chan Message, buffer)
	b.mu.Lock()
	b.subscribers[topic] = append(b.subscribers[topic], ch)
	b.mu.Unlock()

	var once sync.Once
	unsub := func() {
		once.Do(func() {
			b.mu.Lock()
			subs := b.subscribers[topic]
			for i, sub := range subs {
				if sub == ch {
					subs = append(subs[:i], subs[i+1:]...)
					break
				}
			}
			if len(subs) == 0 {
				delete(b.subscribers, topic)
			} else {
				b.subscribers[topic] = subs
			}
			b.mu.Unlock()
			close(ch)
		})
	}

	return ch, unsub
}

// InMemorySubscriber implements Subscriber backed by a channel.
type InMemorySubscriber struct {
	topic   string
	broker  *InMemoryBroker
	buffer  int
	msgs    <-chan Message
	closeFn func()
}

// NewInMemorySubscriber creates a subscriber for a topic.
func NewInMemorySubscriber(broker *InMemoryBroker, topic string, buffer int) *InMemorySubscriber {
	msgs, cancel := broker.Subscribe(topic, buffer)
	return &InMemorySubscriber{
		topic:   topic,
		broker:  broker,
		buffer:  buffer,
		msgs:    msgs,
		closeFn: cancel,
	}
}

// Messages exposes the underlying message channel.
func (s *InMemorySubscriber) Messages() <-chan Message {
	return s.msgs
}

// Close unsubscribes from the topic.
func (s *InMemorySubscriber) Close(ctx context.Context) error {
	_ = ctx
	if s.closeFn != nil {
		s.closeFn()
	}
	return nil
}
