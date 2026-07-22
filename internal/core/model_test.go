package core

import (
	"errors"
	"testing"
)

type failingModel struct{}

func (failingModel) Scores(DecisionRequest, []Candidate) (map[string]float64, error) {
	return nil, errors.New("model unavailable")
}

func TestHybridScorerFallsBackToBaseline(t *testing.T) {
	baseline := WeightedScorer{Weights: map[string]float64{"relevance": 1}}
	scorer := HybridScorer{Baseline: baseline, Model: failingModel{}, ModelWeight: 1, RolloutPercent: 100}
	items, err := scorer.Score(DecisionRequest{Surface: "home", SessionID: "s1"}, []Candidate{
		{Key: "a", ItemID: "a", Features: map[string]float64{"relevance": 2}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Score != 2 {
		t.Fatalf("unexpected baseline result: %#v", items)
	}
	if items[0].Breakdown["model_fallback"] != 1 {
		t.Fatal("expected model fallback marker")
	}
}

func TestExplorationIsBoundedToLastVisibleSlot(t *testing.T) {
	in := []RankedCandidate{
		{Candidate: Candidate{Key: "a", ItemID: "a"}, Score: 4},
		{Candidate: Candidate{Key: "b", ItemID: "b"}, Score: 3},
		{Candidate: Candidate{Key: "c", ItemID: "c"}, Score: 2},
		{Candidate: Candidate{Key: "d", ItemID: "d"}, Score: 1},
	}
	out, explored := ApplyExploration(
		DecisionRequest{Surface: "home", SessionID: "stable-session"},
		in,
		3,
		ExplorationPolicy{Key: "demo", Surface: "home", Percent: 100, Status: "active"},
	)
	if !explored {
		t.Fatal("expected exploration")
	}
	if len(out) != 3 {
		t.Fatalf("unexpected result length: %d", len(out))
	}
	if out[0].Candidate.Key != "a" || out[1].Candidate.Key != "b" {
		t.Fatal("exploration must not disturb leading deterministic slots")
	}
	if out[2].Candidate.Key != "d" {
		t.Fatalf("expected lower-ranked candidate in exploration slot, got %s", out[2].Candidate.Key)
	}
}
