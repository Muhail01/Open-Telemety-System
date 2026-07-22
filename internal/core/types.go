package core

import "time"

type Event struct {
	EventID    string         `json:"event_id"`
	EventType  string         `json:"event_type"`
	OccurredAt time.Time      `json:"occurred_at,omitempty"`
	SessionID  string         `json:"session_id,omitempty"`
	UserID     string         `json:"user_id,omitempty"`
	DecisionID string         `json:"decision_id,omitempty"`
	ItemID     string         `json:"item_id,omitempty"`
	Surface    string         `json:"surface,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}

type Candidate struct {
	Key      string             `json:"key"`
	ItemID   string             `json:"item_id"`
	Group    string             `json:"group,omitempty"`
	Features map[string]float64 `json:"features,omitempty"`
	Metadata map[string]any     `json:"metadata,omitempty"`
	Blocked  bool               `json:"blocked,omitempty"`
}

type RankedCandidate struct {
	Candidate Candidate          `json:"candidate"`
	Score     float64            `json:"score"`
	Breakdown map[string]float64 `json:"breakdown"`
}

type DecisionRequest struct {
	Surface   string         `json:"surface"`
	SessionID string         `json:"session_id,omitempty"`
	UserID    string         `json:"user_id,omitempty"`
	Limit     int            `json:"limit,omitempty"`
	Context   map[string]any `json:"context,omitempty"`
}

type DecisionItem struct {
	ItemID     string             `json:"item_id"`
	Rank       int                `json:"rank"`
	Score      float64            `json:"score"`
	ReasonCode string             `json:"reason_code"`
	Breakdown  map[string]float64 `json:"breakdown,omitempty"`
}

type Decision struct {
	DecisionID string         `json:"decision_id"`
	Surface    string         `json:"surface"`
	Items      []DecisionItem `json:"items"`
	Fallback   bool           `json:"fallback"`
	Reasons    []string       `json:"reason_codes"`
	CreatedAt  time.Time      `json:"created_at"`
}

type CandidateProvider interface {
	Candidates(req DecisionRequest) ([]Candidate, error)
}

type Scorer interface {
	Score(req DecisionRequest, candidates []Candidate) ([]RankedCandidate, error)
}

type Policy interface {
	Apply(req DecisionRequest, candidates []RankedCandidate) ([]RankedCandidate, []string)
}

type Store interface {
	SaveDecision(decision Decision) error
	SaveEvent(event Event) (duplicate bool, err error)
}
