package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Muhail01/GMF-Core/internal/core"
)

type Server struct {
	Engine core.Engine
}

func (s Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
	})
	mux.HandleFunc("POST /v1/recommendations", s.handleRecommendations)
	mux.HandleFunc("GET /v1/decisions/{decision_id}", s.handleDecisionLookup)
	mux.HandleFunc("POST /v1/events", s.handleEvents)
	return mux
}

func (s Server) handleRecommendations(w http.ResponseWriter, r *http.Request) {
	var req core.DecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse("INVALID_JSON", err.Error()))
		return
	}
	decision, err := s.Engine.Decide(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse("DECISION_FAILED", err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, decision)
}

func (s Server) handleDecisionLookup(w http.ResponseWriter, r *http.Request) {
	decisionID := strings.TrimSpace(r.PathValue("decision_id"))
	if decisionID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse("INVALID_DECISION_ID", "decision_id is required"))
		return
	}
	decision, found, err := s.Engine.LookupDecision(decisionID)
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, errorResponse("DECISION_LOOKUP_UNAVAILABLE", err.Error()))
		return
	}
	if !found {
		writeJSON(w, http.StatusNotFound, errorResponse("DECISION_NOT_FOUND", "decision was not found"))
		return
	}
	writeJSON(w, http.StatusOK, decision)
}

func (s Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	var event core.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse("INVALID_JSON", err.Error()))
		return
	}
	if strings.HasPrefix(event.EventType, "recommendation_") && (event.DecisionID == "" || event.ItemID == "" || event.Surface == "") {
		writeJSON(w, http.StatusBadRequest, errorResponse("INVALID_RECOMMENDATION_EVENT", "decision_id, item_id and surface are required"))
		return
	}
	duplicate, err := s.Engine.Ingest(event)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse("EVENT_REJECTED", err.Error()))
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"accepted": !duplicate, "duplicate": duplicate})
}

func errorResponse(code, message string) map[string]any {
	return map[string]any{"error": map[string]string{"code": code, "message": message}}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
