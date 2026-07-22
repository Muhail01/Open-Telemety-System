package core

import (
	"fmt"
	"sort"
)

type WeightedScorer struct {
	Weights map[string]float64
}

func (s WeightedScorer) Score(_ DecisionRequest, candidates []Candidate) ([]RankedCandidate, error) {
	if len(candidates) == 0 {
		return nil, nil
	}

	out := make([]RankedCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.Key == "" || candidate.ItemID == "" {
			return nil, fmt.Errorf("candidate key and item_id are required")
		}
		breakdown := make(map[string]float64, len(s.Weights))
		score := 0.0
		for feature, weight := range s.Weights {
			contribution := candidate.Features[feature] * weight
			breakdown[feature] = contribution
			score += contribution
		}
		out = append(out, RankedCandidate{Candidate: candidate, Score: score, Breakdown: breakdown})
	}

	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Score == out[j].Score {
			return out[i].Candidate.Key < out[j].Candidate.Key
		}
		return out[i].Score > out[j].Score
	})
	return out, nil
}

func RerankDiversity(in []RankedCandidate, limit, maxPerGroup int) []RankedCandidate {
	if limit <= 0 || limit > len(in) {
		limit = len(in)
	}
	if maxPerGroup <= 0 {
		return append([]RankedCandidate(nil), in[:limit]...)
	}

	counts := map[string]int{}
	out := make([]RankedCandidate, 0, limit)
	for _, candidate := range in {
		group := candidate.Candidate.Group
		if group != "" && counts[group] >= maxPerGroup {
			continue
		}
		out = append(out, candidate)
		if group != "" {
			counts[group]++
		}
		if len(out) == limit {
			break
		}
	}
	return out
}
