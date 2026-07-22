package core

import "hash/fnv"

func ApplyExploration(req DecisionRequest, in []RankedCandidate, limit int, policy ExplorationPolicy) ([]RankedCandidate, bool) {
	if limit <= 0 || limit > len(in) {
		limit = len(in)
	}
	base := append([]RankedCandidate(nil), in[:limit]...)
	if policy.Status != "active" || policy.Percent <= 0 || len(in) <= limit || limit == 0 {
		return base, false
	}
	if policy.Surface != "" && policy.Surface != req.Surface {
		return base, false
	}
	if !selectedForRollout(req, policy.Percent) {
		return base, false
	}

	key := policy.Key + "|" + req.Surface + "|" + req.SessionID + "|" + req.UserID
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	extraIndex := limit + int(h.Sum32()%uint32(len(in)-limit))
	base[limit-1] = in[extraIndex]
	return base, true
}
