package core

type DefaultPolicy struct {
	MinimumScore *float64
}

func (p DefaultPolicy) Apply(_ DecisionRequest, candidates []RankedCandidate) ([]RankedCandidate, []string) {
	out := make([]RankedCandidate, 0, len(candidates))
	reasons := make([]string, 0, 2)
	for _, candidate := range candidates {
		if candidate.Candidate.Blocked {
			if !containsReason(reasons, "BLOCKED_CANDIDATE_SUPPRESSED") {
				reasons = append(reasons, "BLOCKED_CANDIDATE_SUPPRESSED")
			}
			continue
		}
		if p.MinimumScore != nil && candidate.Score < *p.MinimumScore {
			if !containsReason(reasons, "MINIMUM_SCORE_FILTERED") {
				reasons = append(reasons, "MINIMUM_SCORE_FILTERED")
			}
			continue
		}
		out = append(out, candidate)
	}
	return out, reasons
}

func containsReason(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
