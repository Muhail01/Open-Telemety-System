package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Muhail01/Open-Telemety-System/internal/core"
	"github.com/Muhail01/Open-Telemety-System/internal/demo"
	"github.com/Muhail01/Open-Telemety-System/internal/httpapi"
	"github.com/Muhail01/Open-Telemety-System/internal/memory"
	pgstore "github.com/Muhail01/Open-Telemety-System/internal/postgres"
)

func main() {
	var store core.Store = memory.NewStore()
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		postgresStore, err := pgstore.Open(ctx, databaseURL)
		cancel()
		if err != nil {
			log.Fatalf("connect PostgreSQL: %v", err)
		}
		defer postgresStore.Close()
		store = postgresStore
		log.Printf("Open Telemetry System storage: PostgreSQL")
	} else {
		log.Printf("Open Telemetry System storage: in-memory")
	}

	engine := core.Engine{
		Provider: demo.NewCatalogProvider(),
		Scorer: core.WeightedScorer{Weights: map[string]float64{
			"relevance": 0.55,
			"quality":   0.30,
			"freshness": 0.15,
		}},
		Policy:      core.DefaultPolicy{},
		Store:       store,
		MaxPerGroup: 2,
	}

	addr := os.Getenv("GMF_CORE_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	server := httpapi.Server{Engine: engine}
	log.Printf("Open Telemetry System listening on %s", addr)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
