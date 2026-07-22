CREATE TABLE IF NOT EXISTS gmf_events (
    event_id TEXT PRIMARY KEY,
    event_type TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    session_id TEXT,
    user_id TEXT,
    decision_id TEXT,
    item_id TEXT,
    surface TEXT,
    properties JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS gmf_decisions (
    decision_id TEXT PRIMARY KEY,
    surface TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS gmf_outbox (
    id BIGSERIAL PRIMARY KEY,
    topic TEXT NOT NULL,
    aggregate_key TEXT NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    delivered_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS gmf_outbox_pending_idx
    ON gmf_outbox (id)
    WHERE delivered_at IS NULL;

CREATE INDEX IF NOT EXISTS gmf_events_decision_idx
    ON gmf_events (decision_id)
    WHERE decision_id IS NOT NULL;
