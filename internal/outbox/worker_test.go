package outbox_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Muhail01/Open-Telemety-System/internal/outbox"
)

type fakeStore struct {
	messages  []outbox.Message
	delivered []int64
}

func (s *fakeStore) FetchPending(context.Context, int) ([]outbox.Message, error) {
	return append([]outbox.Message(nil), s.messages...), nil
}

func (s *fakeStore) MarkDelivered(_ context.Context, id int64) error {
	s.delivered = append(s.delivered, id)
	return nil
}

func TestWorkerDeliversAndAcknowledgesInOrder(t *testing.T) {
	store := &fakeStore{messages: []outbox.Message{{ID: 1}, {ID: 2}}}
	seen := make([]int64, 0, 2)
	worker := outbox.Worker{
		Store: store,
		Handler: func(_ context.Context, message outbox.Message) error {
			seen = append(seen, message.ID)
			return nil
		},
	}
	delivered, err := worker.RunOnce(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if delivered != 2 || len(seen) != 2 || len(store.delivered) != 2 {
		t.Fatalf("unexpected delivery state: delivered=%d seen=%v ack=%v", delivered, seen, store.delivered)
	}
}

func TestWorkerDoesNotAcknowledgeFailedDelivery(t *testing.T) {
	store := &fakeStore{messages: []outbox.Message{{ID: 1}, {ID: 2}}}
	worker := outbox.Worker{
		Store: store,
		Handler: func(_ context.Context, message outbox.Message) error {
			if message.ID == 2 {
				return errors.New("sink unavailable")
			}
			return nil
		},
	}
	delivered, err := worker.RunOnce(context.Background())
	if err == nil {
		t.Fatal("expected delivery error")
	}
	if delivered != 1 {
		t.Fatalf("expected one delivered message, got %d", delivered)
	}
	if len(store.delivered) != 1 || store.delivered[0] != 1 {
		t.Fatalf("unexpected acknowledgements: %v", store.delivered)
	}
}
