package core

import "fmt"

// DecisionReader is an optional storage capability used by explainability and
// debugging endpoints. It is intentionally separate from Store so write-only
// adapters remain valid.
type DecisionReader interface {
	DecisionByID(id string) (Decision, bool, error)
}

func (e Engine) LookupDecision(id string) (Decision, bool, error) {
	if id == "" {
		return Decision{}, false, fmt.Errorf("decision_id is required")
	}
	reader, ok := e.Store.(DecisionReader)
	if !ok {
		return Decision{}, false, fmt.Errorf("decision lookup is not supported by the configured store")
	}
	return reader.DecisionByID(id)
}
