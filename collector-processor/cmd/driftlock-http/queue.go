package main

import (
	"context"

	"github.com/google/uuid"
)

type queueJob struct {
	ID       uuid.UUID
	TenantID uuid.UUID
	Type     string
	Payload  []byte
}

type queueStats struct {
	Mode     string
	Pending  int
	Capacity int
}

type jobQueue interface {
	Enqueue(ctx context.Context, job queueJob) error
	Stats() queueStats
}

type memoryQueue struct {
	name string
	ch   chan queueJob
}

func newMemoryQueue(capacity int) *memoryQueue {
	if capacity <= 0 {
		capacity = 256
	}
	return &memoryQueue{
		name: "memory",
		ch:   make(chan queueJob, capacity),
	}
}

func (mq *memoryQueue) Enqueue(ctx context.Context, job queueJob) error {
	select {
	case mq.ch <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (mq *memoryQueue) Stats() queueStats {
	return queueStats{
		Mode:     mq.name,
		Pending:  len(mq.ch),
		Capacity: cap(mq.ch),
	}
}
