package httpapi_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Muhail01/Open-Telemety-System/internal/core"
	"github.com/Muhail01/Open-Telemety-System/internal/demo"
	"github.com/Muhail01/Open-Telemety-System/internal/httpapi"
	"github.com/Muhail01/Open-Telemety-System/internal/memory"
)

func testServer() http.Handler {
	engine := core.Engine{
		Provider: demo.NewCatalogProvider(),
		Scorer: core.WeightedScorer{Weights: map[string]float64{
			"relevance": 1,
		}},
		Policy: core.DefaultPolicy{},
		Store:  memory.NewStore(),
	}
	return httpapi.Server{Engine: engine}.Handler()
}

func TestRecommendationsEndpoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/recommendations", bytes.NewBufferString(`{"surface":"home","session_id":"s1","limit":2}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	testServer().ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), `"decision_id"`) || !strings.Contains(res.Body.String(), `"items"`) {
		t.Fatalf("unexpected response: %s", res.Body.String())
	}
}

func TestRecommendationTelemetryRequiresTrace(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/events", bytes.NewBufferString(`{"event_id":"e1","event_type":"recommendation_click"}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	testServer().ServeHTTP(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), "INVALID_RECOMMENDATION_EVENT") {
		t.Fatalf("unexpected response: %s", res.Body.String())
	}
}
