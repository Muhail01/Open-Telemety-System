package httpapi_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkRecommendationsEndpointParallel(b *testing.B) {
	handler := testServer()
	payload := []byte(`{"surface":"home","session_id":"load-benchmark","limit":4}`)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodPost, "/v1/recommendations", bytes.NewReader(payload))
			res := httptest.NewRecorder()
			handler.ServeHTTP(res, req)
			if res.Code != http.StatusOK {
				b.Fatalf("status=%d body=%s", res.Code, res.Body.String())
			}
		}
	})
}
