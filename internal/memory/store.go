package memory

import (
	"fmt"
	"sync"

	"github.com/Muhail01/Open-Telemety-System/internal/core"
)

type Store struct {
	mu        sync.RWMutex
	decisions map[string]core.Decision
	events    map[string]core.Event
}

func NewStore() *Store {
	return &Store{
		decisions: map[string]core.Decision{},
		events:    map[string]core.Event{},
	}
}

func (s *Store) SaveDecision(decision core.Decision) error {
	if decision.DecisionID == "" {
		return fmt.Errorf("decision_id is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.decisions[decision.DecisionID] = decision
	return nil
}

func (s *Store) SaveEvent(event core.Event) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[event.EventID]; exists {
		return true, nil
	}
	s.events[event.EventID] = event
	return false, nil
}

func (s *Store) Decision(id string) (core.Decision, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	decision, ok := s.decisions[id]
	return decision, ok
}
