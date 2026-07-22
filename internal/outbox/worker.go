package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	ID           int64           `json:"id"`
	Topic        string          `json:"topic"`
	AggregateKey string          `json:"aggregate_key"`
	Payload      json.RawMessage `json:"payload"`
	CreatedAt    time.Time       `json:"created_at"`
}

type Store interface {
	FetchPending(ctx context.Context, limit int) ([]Message, error)
	MarkDelivered(ctx context.Context, id int64) error
}

type Handler func(context.Context, Message) error

type Worker struct {
	Store     Store
	Handler   Handler
	BatchSize int
}

func (w Worker) RunOnce(ctx context.Context) (int, error) {
	if w.Store == nil {
		return 0, fmt.Errorf("outbox store is not configured")
	}
	if w.Handler == nil {
		return 0, fmt.Errorf("outbox handler is not configured")
	}
	limit := w.BatchSize
	if limit <= 0 {
		limit = 100
	}
	messages, err := w.Store.FetchPending(ctx, limit)
	if err != nil {
		return 0, fmt.Errorf("fetch pending outbox: %w", err)
	}
	delivered := 0
	for _, message := range messages {
		if err := w.Handler(ctx, message); err != nil {
			return delivered, fmt.Errorf("deliver outbox message %d: %w", message.ID, err)
		}
		if err := w.Store.MarkDelivered(ctx, message.ID); err != nil {
			return delivered, fmt.Errorf("mark outbox message %d delivered: %w", message.ID, err)
		}
		delivered++
	}
	return delivered, nil
}
