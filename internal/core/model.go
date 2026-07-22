package core

import (
	"hash/fnv"
	"sort"
)

type ModelScoreProvider interface {
	Scores(req DecisionRequest, candidates []Candidate) (map[string]float64, error)
}

type HybridScorer struct {
	Baseline       Scorer
	Model          ModelScoreProvider
	ModelWeight    float64
	RolloutPercent int
}

func (s HybridScorer) Score(req DecisionRequest, candidates []Candidate) ([]RankedCandidate, error) {
	baseline, err := s.Baseline.Score(req, candidates)
	if err != nil {
		return nil, err
	}
	if s.Model == nil || s.ModelWeight == 0 || !selectedForRollout(req, s.RolloutPercent) {
		return baseline, nil
	}

	modelScores, err := s.Model.Scores(req, candidates)
	if err != nil {
		for i := range baseline {
			if baseline[i].Breakdown == nil {
				baseline[i].Breakdown = map[string]float64{}
			}
			baseline[i].Breakdown["model_fallback"] = 1
		}
		return baseline, nil
	}

	for i := range baseline {
		contribution := modelScores[baseline[i].Candidate.Key] * s.ModelWeight
		baseline[i].Score += contribution
		if baseline[i].Breakdown == nil {
			baseline[i].Breakdown = map[string]float64{}
		}
		baseline[i].Breakdown["model_score"] = contribution
	}
	sort.SliceStable(baseline, func(i, j int) bool {
		if baseline[i].Score == baseline[j].Score {
			return baseline[i].Candidate.Key < baseline[j].Candidate.Key
		}
		return baseline[i].Score > baseline[j].Score
	})
	return baseline, nil
}

func selectedForRollout(req DecisionRequest, percent int) bool {
	if percent <= 0 {
		return false
	}
	if percent >= 100 {
		return true
	}
	key := req.Surface + "|" + req.SessionID + "|" + req.UserID
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32()%100) < percent
}
