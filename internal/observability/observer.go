package observability

import (
	"strconv"
	"time"

	"github.com/Muhail01/GMF-Core/internal/core"
)

// MetricSink is deliberately minimal so deployments can adapt it to an
// OpenTelemetry Meter, Prometheus collector, StatsD client, or test recorder
// without making the core depend on a specific observability vendor.
type MetricSink interface {
	AddCounter(name string, value int64, attributes map[string]string)
	RecordDuration(name string, value time.Duration, attributes map[string]string)
	RecordGauge(name string, value int64, attributes map[string]string)
}

type Observer struct {
	Sink MetricSink
}

func (o Observer) RecordDecision(observation core.DecisionObservation) {
	if o.Sink == nil {
		return
	}
	attributes := map[string]string{
		"surface":     observation.Surface,
		"fallback":    strconv.FormatBool(observation.Fallback),
		"exploration": strconv.FormatBool(observation.ExplorationUsed),
	}
	o.Sink.AddCounter("decision.total", 1, attributes)
	o.Sink.RecordDuration("decision.duration", observation.Duration, attributes)
	o.Sink.RecordGauge("decision.candidates", int64(observation.CandidateCount), attributes)
	o.Sink.RecordGauge("decision.items", int64(observation.ServedItemCount), attributes)
}

func (o Observer) RecordEvent(observation core.EventObservation) {
	if o.Sink == nil {
		return
	}
	o.Sink.AddCounter("event.total", 1, map[string]string{
		"event_type": observation.EventType,
		"duplicate":  strconv.FormatBool(observation.Duplicate),
	})
}
