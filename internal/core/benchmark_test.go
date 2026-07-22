package core_test

import (
	"testing"

	"github.com/Muhail01/Open-Telemety-System/internal/core"
	"github.com/Muhail01/Open-Telemety-System/internal/demo"
	"github.com/Muhail01/Open-Telemety-System/internal/memory"
)

func BenchmarkDecisionPath(b *testing.B) {
	engine := core.Engine{
		Provider: demo.NewCatalogProvider(),
		Scorer: core.WeightedScorer{Weights: map[string]float64{
			"relevance": 0.55,
			"quality":   0.30,
			"freshness": 0.15,
		}},
		Policy:      core.DefaultPolicy{},
		Store:       memory.NewStore(),
		MaxPerGroup: 2,
	}
	req := core.DecisionRequest{Surface: "home", SessionID: "benchmark", Limit: 4}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := engine.Decide(req); err != nil {
			b.Fatal(err)
		}
	}
}
