package core_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/Muhail01/GMF-Core/internal/core"
)

func TestEqualScoreOrderingIsStableAcrossInputPermutations(t *testing.T) {
	rng := rand.New(rand.NewSource(62026))
	scorer := core.WeightedScorer{Weights: map[string]float64{"score": 1}}

	for iteration := 0; iteration < 100; iteration++ {
		candidates := make([]core.Candidate, 20)
		expected := make([]string, 20)
		for i := range candidates {
			key := fmt.Sprintf("candidate-%02d", i)
			candidates[i] = core.Candidate{Key: key, ItemID: key, Features: map[string]float64{"score": 1}}
			expected[i] = key
		}
		rng.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })

		ranked, err := scorer.Score(core.DecisionRequest{Surface: "home"}, candidates)
		if err != nil {
			t.Fatal(err)
		}
		sort.Strings(expected)
		for i := range ranked {
			if ranked[i].Candidate.Key != expected[i] {
				t.Fatalf("iteration %d rank %d: got %s want %s", iteration, i, ranked[i].Candidate.Key, expected[i])
			}
		}
	}
}
