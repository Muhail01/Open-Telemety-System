# Open Telemetry System

Open Telemetry System is a vendor-neutral, open-source decisioning, recommendation, telemetry, and organic ranking engine extracted from production-oriented marketplace architecture.

It is designed around an explainable online decision pipeline:

`event ingest -> candidate providers -> feature enrichment -> deterministic baseline ranking -> optional model scoring -> diversity rerank -> policy/guardrails -> controlled exploration -> explainable decision -> decision log + transactional outbox -> impression/click/outcome telemetry`

## Why it exists

Many marketplace recommendation systems become tightly coupled to one product, one database schema, and private commercial rules. This project separates reusable decisioning primitives from those proprietary concerns.

The public core intentionally excludes:

- production user and marketplace data;
- payment and wallet logic;
- paid-ad billing and money movement;
- fraud, KYC, dispute and support systems;
- private seller-risk formulas;
- production ranking coefficients;
- private model artifacts and production model configuration;
- credentials, private endpoints and deployment topology;
- GILLZY-specific domain integrations.

## Included in the initial scaffold

- neutral recommendation and telemetry contracts;
- pluggable candidate providers;
- deterministic baseline scoring;
- configurable feature weights;
- optional model-assisted scoring with deterministic rollout and fail-safe baseline fallback;
- diversity-aware reranking;
- policy guardrails and safe fallback behavior;
- bounded deterministic exploration;
- feature/model/experiment/exploration registry primitives;
- in-memory and PostgreSQL event/decision storage;
- transactional outbox;
- HTTP API for event ingest and recommendations;
- decision IDs and score breakdowns for explainability and impression/click attribution;
- synthetic marketplace catalog;
- TypeScript telemetry SDK;
- OpenAPI 3.1 contract;
- Docker Compose reference stack;
- tests, benchmarks and CI.

## Quick start

Requires Go 1.24+.

```bash
go test ./...
go run ./cmd/gmf-core
```

Then:

```bash
curl -s http://localhost:8080/healthz
```

Request recommendations:

```bash
curl -s -X POST http://localhost:8080/v1/recommendations \
  -H 'Content-Type: application/json' \
  -d '{
    "surface":"home",
    "session_id":"demo-session",
    "limit":4,
    "context":{"category":"games"}
  }'
```

Ingest a telemetry event:

```bash
curl -s -X POST http://localhost:8080/v1/events \
  -H 'Content-Type: application/json' \
  -d '{
    "event_id":"evt-1",
    "event_type":"recommendation_impression",
    "session_id":"demo-session",
    "decision_id":"<decision-id>",
    "item_id":"game-1",
    "surface":"home"
  }'
```

## Architecture

The core is deliberately modular. Candidate sources, feature enrichment, baseline scoring, model scoring, reranking, policy and exploration are independent stages. Each decision is persisted with a stable `decision_id`, reason codes and score breakdowns so downstream impression, click and outcome telemetry can be joined back to the decision that produced the slate.

See `docs/ARCHITECTURE.md` and `docs/OPEN_SOURCE_BOUNDARY.md`.

## Status

This repository is the standalone public extraction of reusable recommendation/decisioning infrastructure originally developed while building GILLZY. The project is being prepared for a first `v0.1.0` release.

## License

Apache-2.0. See `LICENSE`.
