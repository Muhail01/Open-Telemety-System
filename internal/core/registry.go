package core

import "sync"

type ModelVersion struct {
	Key            string         `json:"key"`
	Version        string         `json:"version"`
	Status         string         `json:"status"`
	RolloutPercent int            `json:"rollout_percent"`
	Config         map[string]any `json:"config,omitempty"`
}

type FeatureDefinition struct {
	Key         string `json:"key"`
	EntityType  string `json:"entity_type"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

type Experiment struct {
	ID             string `json:"id"`
	Surface        string `json:"surface"`
	Variant        string `json:"variant"`
	TrafficPercent int    `json:"traffic_percent"`
	Status         string `json:"status"`
}

type ExplorationPolicy struct {
	Key            string `json:"key"`
	Surface        string `json:"surface"`
	Percent        int    `json:"percent"`
	MaxPerDecision int    `json:"max_per_decision"`
	Status         string `json:"status"`
}

type Registry struct {
	mu          sync.RWMutex
	models      map[string]ModelVersion
	features    map[string]FeatureDefinition
	experiments map[string]Experiment
	exploration map[string]ExplorationPolicy
}

func NewRegistry() *Registry {
	return &Registry{
		models:      map[string]ModelVersion{},
		features:    map[string]FeatureDefinition{},
		experiments: map[string]Experiment{},
		exploration: map[string]ExplorationPolicy{},
	}
}

func (r *Registry) PutModel(model ModelVersion) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.models[model.Key] = model
}

func (r *Registry) Model(key string) (ModelVersion, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	model, ok := r.models[key]
	return model, ok
}

func (r *Registry) PutFeature(feature FeatureDefinition) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.features[feature.Key] = feature
}

func (r *Registry) Feature(key string) (FeatureDefinition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	feature, ok := r.features[key]
	return feature, ok
}

func (r *Registry) PutExperiment(experiment Experiment) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.experiments[experiment.ID] = experiment
}

func (r *Registry) Experiment(id string) (Experiment, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	experiment, ok := r.experiments[id]
	return experiment, ok
}

func (r *Registry) PutExploration(policy ExplorationPolicy) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.exploration[policy.Key] = policy
}

func (r *Registry) Exploration(key string) (ExplorationPolicy, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	policy, ok := r.exploration[key]
	return policy, ok
}
