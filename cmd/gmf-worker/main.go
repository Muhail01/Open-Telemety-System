package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Muhail01/Open-Telemety-System/internal/outbox"
	pgstore "github.com/Muhail01/Open-Telemety-System/internal/postgres"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	store, err := pgstore.Open(connectCtx, databaseURL)
	cancel()
	if err != nil {
		log.Fatalf("connect PostgreSQL: %v", err)
	}
	defer store.Close()

	worker := outbox.Worker{
		Store: store,
		Handler: func(_ context.Context, message outbox.Message) error {
			// Reference sink: record only non-sensitive metadata. Real deployments
			// should replace this handler with Kafka/NATS/webhook/etc. delivery.
			log.Printf("outbox deliver topic=%s aggregate_key=%s payload_bytes=%d", message.Topic, message.AggregateKey, len(message.Payload))
			return nil
		},
		BatchSize: 100,
	}

	interval := time.Second
	if raw := os.Getenv("OUTBOX_POLL_INTERVAL"); raw != "" {
		parsed, err := time.ParseDuration(raw)
		if err != nil || parsed <= 0 {
			log.Fatalf("invalid OUTBOX_POLL_INTERVAL %q", raw)
		}
		interval = parsed
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	log.Printf("outbox worker started poll_interval=%s", interval)

	for {
		delivered, err := worker.RunOnce(ctx)
		if err != nil && ctx.Err() == nil {
			log.Printf("outbox delivery error: %v", err)
		}
		if delivered > 0 {
			log.Printf("outbox batch delivered=%d", delivered)
		}

		select {
		case <-ctx.Done():
			log.Printf("outbox worker stopped")
			return
		case <-ticker.C:
		}
	}
}
