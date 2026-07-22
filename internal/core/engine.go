package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type Engine struct {
	Provider          CandidateProvider
	Scorer            Scorer
	Policy            Policy
	Store             Store
	Observer          Observer
	MaxPerGroup       int
	ExplorationPolicy *ExplorationPolicy
}

func (e Engine) Decide(req DecisionRequest) (Decision, error) {
	startedAt := time.Now()
	if e.Provider == nil || e.Scorer == nil || e.Policy == nil || e.Store == nil {
		return Decision{}, fmt.Errorf("engine dependencies are not configured")
	}
	if req.Surface == "" {
		return Decision{}, fmt.Errorf("surface is required")
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	candidates, err := e.Provider.Candidates(req)
	if err != nil {
		return Decision{}, fmt.Errorf("build candidates: %w", err)
	}
	ranked, err := e.Scorer.Score(req, candidates)
	if err != nil {
		return Decision{}, fmt.Errorf("score candidates: %w", err)
	}
	filtered, reasons := e.Policy.Apply(req, ranked)
	diversified := RerankDiversity(filtered, len(filtered), e.MaxPerGroup)

	var final []RankedCandidate
	explorationUsed := false
	if e.ExplorationPolicy != nil {
		final, explorationUsed = ApplyExploration(req, diversified, req.Limit, *e.ExplorationPolicy)
		if explorationUsed {
			reasons = append(reasons, "EXPLORATION_APPLIED")
		}
	} else {
		final = RerankDiversity(diversified, req.Limit, 0)
	}

	decision := Decision{
		DecisionID: newDecisionID(),
		Surface:    req.Surface,
		Items:      make([]DecisionItem, 0, len(final)),
		CreatedAt:  time.Now().UTC(),
		Reasons:    reasons,
	}
	if len(final) == 0 {
		decision.Fallback = true
		decision.Reasons = append(decision.Reasons, "NO_ELIGIBLE_ITEMS")
	}
	for i, candidate := range final {
		decision.Items = append(decision.Items, DecisionItem{
			ItemID:     candidate.Candidate.ItemID,
			Rank:       i + 1,
			Score:      candidate.Score,
			ReasonCode: "RANKED_ORGANIC",
			Breakdown:  candidate.Breakdown,
		})
	}
	if err := e.Store.SaveDecision(decision); err != nil {
		return Decision{}, fmt.Errorf("save decision: %w", err)
	}
	if e.Observer != nil {
		e.Observer.RecordDecision(DecisionObservation{
			Surface:         req.Surface,
			Duration:        time.Since(startedAt),
			CandidateCount:  len(candidates),
			ServedItemCount: len(decision.Items),
			Fallback:        decision.Fallback,
			ExplorationUsed: explorationUsed,
		})
	}
	return decision, nil
}

func (e Engine) Ingest(event Event) (bool, error) {
	if e.Store == nil {
		return false, fmt.Errorf("event store is not configured")
	}
	if event.EventID == "" || event.EventType == "" {
		return false, fmt.Errorf("event_id and event_type are required")
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now().UTC()
	}
	duplicate, err := e.Store.SaveEvent(event)
	if err != nil {
		return false, err
	}
	if e.Observer != nil {
		e.Observer.RecordEvent(EventObservation{EventType: event.EventType, Duplicate: duplicate})
	}
	return duplicate, nil
}

func newDecisionID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return fmt.Sprintf("decision-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b[:])
}
