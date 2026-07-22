package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Muhail01/GMF-Core/internal/core"
	"github.com/Muhail01/GMF-Core/internal/outbox"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	db *sql.DB
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) SaveDecision(decision core.Decision) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("postgres store is not configured")
	}
	payload, err := json.Marshal(decision)
	if err != nil {
		return fmt.Errorf("marshal decision: %w", err)
	}
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO gmf_decisions (decision_id, surface, payload, created_at)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (decision_id) DO NOTHING`,
		decision.DecisionID, decision.Surface, payload, decision.CreatedAt,
	); err != nil {
		return fmt.Errorf("insert decision: %w", err)
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO gmf_outbox (topic, aggregate_key, payload) VALUES ($1, $2, $3)`,
		"decision.created", decision.DecisionID, payload,
	); err != nil {
		return fmt.Errorf("insert decision outbox: %w", err)
	}
	return tx.Commit()
}

func (s *Store) DecisionByID(id string) (core.Decision, bool, error) {
	if s == nil || s.db == nil {
		return core.Decision{}, false, fmt.Errorf("postgres store is not configured")
	}
	var payload []byte
	err := s.db.QueryRowContext(context.Background(),
		`SELECT payload FROM gmf_decisions WHERE decision_id = $1`, id,
	).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return core.Decision{}, false, nil
	}
	if err != nil {
		return core.Decision{}, false, fmt.Errorf("lookup decision: %w", err)
	}
	var decision core.Decision
	if err := json.Unmarshal(payload, &decision); err != nil {
		return core.Decision{}, false, fmt.Errorf("decode decision: %w", err)
	}
	return decision, true, nil
}

func (s *Store) SaveEvent(event core.Event) (bool, error) {
	if s == nil || s.db == nil {
		return false, fmt.Errorf("postgres store is not configured")
	}
	properties, err := json.Marshal(event.Properties)
	if err != nil {
		return false, fmt.Errorf("marshal event properties: %w", err)
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return false, fmt.Errorf("marshal event: %w", err)
	}
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx,
		`INSERT INTO gmf_events
		 (event_id, event_type, occurred_at, session_id, user_id, decision_id, item_id, surface, properties)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 ON CONFLICT (event_id) DO NOTHING`,
		event.EventID, event.EventType, event.OccurredAt, nullable(event.SessionID), nullable(event.UserID), nullable(event.DecisionID), nullable(event.ItemID), nullable(event.Surface), properties,
	)
	if err != nil {
		return false, fmt.Errorf("insert event: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rows == 0 {
		return true, nil
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO gmf_outbox (topic, aggregate_key, payload) VALUES ($1, $2, $3)`,
		"event.accepted", event.EventID, payload,
	); err != nil {
		return false, fmt.Errorf("insert event outbox: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return false, nil
}

func (s *Store) FetchPending(ctx context.Context, limit int) ([]outbox.Message, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("postgres store is not configured")
	}
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, topic, aggregate_key, payload, created_at
		 FROM gmf_outbox
		 WHERE delivered_at IS NULL
		 ORDER BY id
		 LIMIT $1`, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query pending outbox: %w", err)
	}
	defer rows.Close()

	messages := make([]outbox.Message, 0, limit)
	for rows.Next() {
		var message outbox.Message
		if err := rows.Scan(&message.ID, &message.Topic, &message.AggregateKey, &message.Payload, &message.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan outbox message: %w", err)
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate outbox messages: %w", err)
	}
	return messages, nil
}

func (s *Store) MarkDelivered(ctx context.Context, id int64) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("postgres store is not configured")
	}
	result, err := s.db.ExecContext(ctx,
		`UPDATE gmf_outbox SET delivered_at = now() WHERE id = $1 AND delivered_at IS NULL`, id,
	)
	if err != nil {
		return fmt.Errorf("mark outbox delivered: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read outbox acknowledgement result: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("outbox message %d is missing or already delivered", id)
	}
	return nil
}

func nullable(value string) any {
	if value == "" {
		return nil
	}
	return value
}
