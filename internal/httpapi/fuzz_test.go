package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func FuzzRecommendationTelemetryValidation(f *testing.F) {
	f.Add("recommendation_click", "decision-1", "item-1", "home")
	f.Add("recommendation_impression", "", "item-1", "home")
	f.Add("page_view", "", "", "")

	f.Fuzz(func(t *testing.T, eventType, decisionID, itemID, surface string) {
		payload, err := json.Marshal(map[string]any{
			"event_id":    "fuzz-event",
			"event_type":  eventType,
			"decision_id": decisionID,
			"item_id":     itemID,
			"surface":     surface,
		})
		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()
		testServer().ServeHTTP(res, httptest.NewRequest(http.MethodPost, "/v1/events", bytes.NewReader(payload)))

		if strings.HasPrefix(eventType, "recommendation_") && (decisionID == "" || itemID == "" || surface == "") {
			if res.Code != http.StatusBadRequest {
				t.Fatalf("unsafe recommendation telemetry accepted: status=%d body=%s", res.Code, res.Body.String())
			}
			return
		}
		if res.Code < 200 || res.Code >= 500 {
			t.Fatalf("unexpected status=%d body=%s", res.Code, res.Body.String())
		}
	})
}
