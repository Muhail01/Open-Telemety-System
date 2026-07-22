package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Muhail01/GMF-Core/internal/core"
	"github.com/Muhail01/GMF-Core/internal/demo"
	"github.com/Muhail01/GMF-Core/internal/httpapi"
	"github.com/Muhail01/GMF-Core/internal/memory"
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

func TestDecisionLookupReturnsPersistedDecision(t *testing.T) {
	handler := testServer()
	createReq := httptest.NewRequest(http.MethodPost, "/v1/recommendations", bytes.NewBufferString(`{"surface":"home","session_id":"lookup-session","limit":2}`))
	createRes := httptest.NewRecorder()
	handler.ServeHTTP(createRes, createReq)
	if createRes.Code != http.StatusOK {
		t.Fatalf("create status=%d body=%s", createRes.Code, createRes.Body.String())
	}
	var decision core.Decision
	if err := json.Unmarshal(createRes.Body.Bytes(), &decision); err != nil {
		t.Fatal(err)
	}
	if decision.DecisionID == "" {
		t.Fatal("expected decision_id")
	}

	lookupReq := httptest.NewRequest(http.MethodGet, "/v1/decisions/"+decision.DecisionID, nil)
	lookupRes := httptest.NewRecorder()
	handler.ServeHTTP(lookupRes, lookupReq)
	if lookupRes.Code != http.StatusOK {
		t.Fatalf("lookup status=%d body=%s", lookupRes.Code, lookupRes.Body.String())
	}
	if !strings.Contains(lookupRes.Body.String(), decision.DecisionID) {
		t.Fatalf("lookup response missing decision id: %s", lookupRes.Body.String())
	}
}

func TestDecisionLookupReturnsNotFound(t *testing.T) {
	res := httptest.NewRecorder()
	testServer().ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/v1/decisions/missing", nil))
	if res.Code != http.StatusNotFound {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
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
