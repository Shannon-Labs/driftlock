package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"time"

	segment "github.com/segmentio/kafka-go"
)

// ConsumerConfig describes how to connect to Kafka for consuming.
type ConsumerConfig struct {
	Brokers   []string
	GroupID   string
	Topic     string
	ClientID  string
	TLSConfig *tls.Config
	MinBytes  int
	MaxBytes  int
	Buffer    int
}

// Consumer implements the Subscriber interface backed by Kafka.
type Consumer struct {
	reader   *segment.Reader
	cancel   context.CancelFunc
	msgs     chan Message
	wg       sync.WaitGroup
	closeErr error
	once     sync.Once
	errMu    sync.Mutex
}

// NewConsumer constructs a Kafka-backed Subscriber.
func NewConsumer(cfg ConsumerConfig) (*Consumer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka consumer: at least one broker is required")
	}
	if cfg.Topic == "" {
		return nil, errors.New("kafka consumer: topic is required")
	}

	if cfg.MinBytes == 0 {
		cfg.MinBytes = 1
	}
	if cfg.MaxBytes == 0 {
		cfg.MaxBytes = 10 << 20 // 10MiB
	}
	if cfg.Buffer <= 0 {
		cfg.Buffer = 128
	}

	dialer := &segment.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		ClientID:  cfg.ClientID,
		TLS:       cfg.TLSConfig,
	}

	reader := segment.NewReader(segment.ReaderConfig{
		Brokers:         cfg.Brokers,
		GroupID:         cfg.GroupID,
		Topic:           cfg.Topic,
		MinBytes:        cfg.MinBytes,
		MaxBytes:        cfg.MaxBytes,
		StartOffset:     segment.LastOffset,
		Dialer:          dialer,
		ReadLagInterval: 0,
	})

	c := &Consumer{
		reader: reader,
		msgs:   make(chan Message, cfg.Buffer),
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.wg.Add(1)
	go c.run(ctx)

	return c, nil
}

// Messages returns a channel that streams Kafka records.
func (c *Consumer) Messages() <-chan Message {
	return c.msgs
}

// Close stops consumption and releases resources.
func (c *Consumer) Close(ctx context.Context) error {
	c.once.Do(func() {
		c.cancel()

		done := make(chan struct{})
		go func() {
			c.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
		case <-ctx.Done():
		}

		c.setErr(c.reader.Close())
	})
	return c.error()
}

func (c *Consumer) run(ctx context.Context) {
	defer c.wg.Done()
	defer close(c.msgs)

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || ctx.Err() != nil {
				return
			}
			// Transient errors: record and continue
			c.setErr(fmt.Errorf("kafka consumer read: %w", err))
			return
		}

		msg := Message{
			Topic:     m.Topic,
			Key:       m.Key,
			Value:     m.Value,
			Timestamp: m.Time,
			Headers:   make(map[string]string, len(m.Headers)),
		}
		for _, h := range m.Headers {
			msg.Headers[h.Key] = string(h.Value)
		}

		select {
		case c.msgs <- msg:
		case <-ctx.Done():
			return
		}
	}
}

func (c *Consumer) setErr(err error) {
	if err == nil {
		return
	}
	c.errMu.Lock()
	if c.closeErr == nil {
		c.closeErr = err
	}
	c.errMu.Unlock()
}

func (c *Consumer) error() error {
	c.errMu.Lock()
	defer c.errMu.Unlock()
	return c.closeErr
}
