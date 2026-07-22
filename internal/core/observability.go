package core

import "time"

// Observer exposes low-cardinality, privacy-safe runtime measurements. Public
// implementations should not add user IDs, session IDs, item IDs, raw queries,
// messages, tokens, or other high-cardinality/private data as metric labels.
type Observer interface {
	RecordDecision(DecisionObservation)
	RecordEvent(EventObservation)
}

type DecisionObservation struct {
	Surface          string
	Duration         time.Duration
	CandidateCount   int
	ServedItemCount  int
	Fallback         bool
	ExplorationUsed  bool
}

type EventObservation struct {
	EventType string
	Duplicate bool
}
