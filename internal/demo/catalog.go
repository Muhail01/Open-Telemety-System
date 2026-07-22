package demo

import "github.com/Muhail01/GMF-Core/internal/core"

type CatalogProvider struct {
	Items []core.Candidate
}

func NewCatalogProvider() CatalogProvider {
	return CatalogProvider{Items: []core.Candidate{
		{Key: "game-1", ItemID: "game-1", Group: "games", Features: map[string]float64{"relevance": 0.94, "quality": 0.88, "freshness": 0.52}},
		{Key: "game-2", ItemID: "game-2", Group: "games", Features: map[string]float64{"relevance": 0.90, "quality": 0.91, "freshness": 0.41}},
		{Key: "skin-1", ItemID: "skin-1", Group: "skins", Features: map[string]float64{"relevance": 0.84, "quality": 0.95, "freshness": 0.79}},
		{Key: "skin-2", ItemID: "skin-2", Group: "skins", Features: map[string]float64{"relevance": 0.78, "quality": 0.86, "freshness": 0.92}},
		{Key: "gift-1", ItemID: "gift-1", Group: "gift-cards", Features: map[string]float64{"relevance": 0.73, "quality": 0.93, "freshness": 0.65}},
	}}
}

func (p CatalogProvider) Candidates(_ core.DecisionRequest) ([]core.Candidate, error) {
	return append([]core.Candidate(nil), p.Items...), nil
}
