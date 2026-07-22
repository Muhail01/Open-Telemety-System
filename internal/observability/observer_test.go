package observability_test

import (
	"testing"
	"time"

	"github.com/Muhail01/Open-Telemety-System/internal/core"
	"github.com/Muhail01/Open-Telemety-System/internal/observability"
)

type record struct {
	name       string
	attributes map[string]string
}

type fakeSink struct {
	records []record
}

func (s *fakeSink) AddCounter(name string, _ int64, attributes map[string]string) {
	s.records = append(s.records, record{name: name, attributes: attributes})
}
func (s *fakeSink) RecordDuration(name string, _ time.Duration, attributes map[string]string) {
	s.records = append(s.records, record{name: name, attributes: attributes})
}
func (s *fakeSink) RecordGauge(name string, _ int64, attributes map[string]string) {
	s.records = append(s.records, record{name: name, attributes: attributes})
}

func TestObserverUsesOnlyApprovedLowCardinalityAttributes(t *testing.T) {
	sink := &fakeSink{}
	observer := observability.Observer{Sink: sink}
	observer.RecordDecision(core.DecisionObservation{Surface: "home", Fallback: false, ExplorationUsed: true})
	observer.RecordEvent(core.EventObservation{EventType: "recommendation_click", Duplicate: false})

	allowed := map[string]bool{
		"surface": true, "fallback": true, "exploration": true,
		"event_type": true, "duplicate": true,
	}
	for _, item := range sink.records {
		for key := range item.attributes {
			if !allowed[key] {
				t.Fatalf("unexpected potentially high-cardinality attribute %q in metric %s", key, item.name)
			}
		}
	}
}
