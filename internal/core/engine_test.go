package core_test

import (
	"testing"

	"github.com/Muhail01/GMF-Core/internal/core"
	"github.com/Muhail01/GMF-Core/internal/demo"
	"github.com/Muhail01/GMF-Core/internal/memory"
)

func TestDecisionIsDeterministicForSameCandidates(t *testing.T) {
	store := memory.NewStore()
	engine := core.Engine{
		Provider: demo.NewCatalogProvider(),
		Scorer: core.WeightedScorer{Weights: map[string]float64{
			"relevance": 0.55,
			"quality":   0.30,
			"freshness": 0.15,
		}},
		Policy:      core.DefaultPolicy{},
		Store:       store,
		MaxPerGroup: 2,
	}

	first, err := engine.Decide(core.DecisionRequest{Surface: "home", SessionID: "s1", Limit: 4})
	if err != nil {
		t.Fatal(err)
	}
	second, err := engine.Decide(core.DecisionRequest{Surface: "home", SessionID: "s1", Limit: 4})
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Items) != len(second.Items) {
		t.Fatalf("item count changed: %d != %d", len(first.Items), len(second.Items))
	}
	for i := range first.Items {
		if first.Items[i].ItemID != second.Items[i].ItemID {
			t.Fatalf("rank %d changed: %s != %s", i, first.Items[i].ItemID, second.Items[i].ItemID)
		}
	}
	if first.DecisionID == second.DecisionID {
		t.Fatal("decision ids must be unique")
	}
}

func TestEventIngestIsIdempotent(t *testing.T) {
	store := memory.NewStore()
	engine := core.Engine{Store: store}
	event := core.Event{EventID: "evt-1", EventType: "recommendation_click", DecisionID: "d1", ItemID: "i1", Surface: "home"}

	duplicate, err := engine.Ingest(event)
	if err != nil || duplicate {
		t.Fatalf("first ingest duplicate=%v err=%v", duplicate, err)
	}
	duplicate, err = engine.Ingest(event)
	if err != nil || !duplicate {
		t.Fatalf("second ingest duplicate=%v err=%v", duplicate, err)
	}
}
