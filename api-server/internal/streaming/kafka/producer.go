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

// ProducerConfig describes how to connect to Kafka for publishing.
type ProducerConfig struct {
	Brokers      []string
	ClientID     string
	TLSConfig    *tls.Config
	BatchSize    int
	BatchTimeout time.Duration
}

type writer interface {
	WriteMessages(context.Context, ...segment.Message) error
	Close() error
}

// Producer implements the Publisher interface backed by Kafka.
type Producer struct {
	dialer      *segment.Dialer
	cfg         ProducerConfig
	mu          sync.Mutex
	writers     map[string]writer
	writerMaker func(topic string) (writer, error)
}

// NewProducer constructs a Kafka-backed Publisher.
func NewProducer(cfg ProducerConfig) (*Producer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka producer: at least one broker is required")
	}

	dialer := &segment.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		ClientID:  cfg.ClientID,
		TLS:       cfg.TLSConfig,
	}

	p := &Producer{
		dialer:  dialer,
		cfg:     cfg,
		writers: make(map[string]writer),
	}

	p.writerMaker = func(topic string) (writer, error) {
		if topic == "" {
			return nil, errors.New("kafka producer: topic is required")
		}

		wCfg := segment.WriterConfig{
			Brokers:      cfg.Brokers,
			Topic:        topic,
			Dialer:       dialer,
			BatchSize:    cfg.BatchSize,
			BatchTimeout: cfg.BatchTimeout,
		}

		if wCfg.BatchSize == 0 {
			wCfg.BatchSize = 100
		}
		if wCfg.BatchTimeout == 0 {
			wCfg.BatchTimeout = 5 * time.Millisecond
		}

		return segment.NewWriter(wCfg), nil
	}

	return p, nil
}

// Publish sends the provided message to Kafka.
func (p *Producer) Publish(ctx context.Context, msg Message) error {
	if msg.Topic == "" {
		return errors.New("kafka producer: message topic is required")
	}

	writer, err := p.getWriter(msg.Topic)
	if err != nil {
		return fmt.Errorf("kafka producer: get writer: %w", err)
	}

	headers := make([]segment.Header, 0, len(msg.Headers))
	for k, v := range msg.Headers {
		headers = append(headers, segment.Header{Key: k, Value: []byte(v)})
	}

	ts := msg.Timestamp
	if ts.IsZero() {
		ts = time.Now()
	}

	return writer.WriteMessages(ctx, segment.Message{
		Topic:   msg.Topic,
		Key:     msg.Key,
		Value:   msg.Value,
		Time:    ts,
		Headers: headers,
	})
}

// Close releases all underlying writer resources.
func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var firstErr error
	for topic, w := range p.writers {
		if err := w.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("close topic %s: %w", topic, err)
		}
		delete(p.writers, topic)
	}
	return firstErr
}

func (p *Producer) getWriter(topic string) (writer, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if existing, ok := p.writers[topic]; ok {
		return existing, nil
	}

	w, err := p.writerMaker(topic)
	if err != nil {
		return nil, err
	}

	p.writers[topic] = w
	return w, nil
}
