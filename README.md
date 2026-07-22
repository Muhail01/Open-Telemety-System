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

## Included in the current v0.1 scaffold

- neutral recommendation and telemetry contracts;
- pluggable candidate providers;
- deterministic baseline scoring and stable tie-breaking;
- configurable feature weights;
- optional model-assisted scoring with deterministic rollout and fail-safe baseline fallback;
- diversity-aware reranking;
- policy guardrails and safe fallback behavior;
- bounded deterministic exploration;
- feature/model/experiment/exploration registry primitives;
- in-memory and PostgreSQL event/decision storage;
- transactional outbox plus a provider-neutral reference delivery worker;
- HTTP API for event ingest, recommendations, health and persisted decision lookup;
- decision IDs and score breakdowns for explainability and impression/click attribution;
- privacy-safe observability hooks that can be adapted to OpenTelemetry, Prometheus, StatsD, or another metric sink;
- synthetic marketplace catalog;
- TypeScript telemetry SDK;
- OpenAPI 3.1 contract;
- Docker Compose reference stack with API, PostgreSQL and outbox worker;
- deterministic/property/fuzz tests, parallel endpoint benchmark and CI;
- automated public-boundary/secret-pattern checks.

## Quick start

Requires Go 1.24+.

```bash
go test ./...
go run ./cmd/gmf-core
```

Or launch the reference PostgreSQL-backed stack:

```bash
docker compose up --build
```

Health check:

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

A successful response contains a `decision_id`. Retrieve the persisted decision and its score breakdowns:

```bash
curl -s http://localhost:8080/v1/decisions/<decision-id>
```

Ingest a telemetry event linked to that decision:

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

The PostgreSQL adapter writes decisions/events and their outbox records transactionally. `cmd/gmf-worker` is a reference delivery loop; deployments can replace its sink with Kafka, NATS, webhooks, queues, or another transport.

Observability deliberately exposes only low-cardinality runtime dimensions such as `surface`, `event_type`, fallback and exploration state. User IDs, session IDs, item IDs, raw queries, messages, tokens, and other private/high-cardinality values are excluded from the public metric contract.

See `docs/ARCHITECTURE.md`, `docs/OPEN_SOURCE_BOUNDARY.md`, and `docs/ROADMAP.md`.

## Status

This repository is the standalone public extraction of reusable recommendation/decisioning infrastructure originally developed while building GILLZY. The project is being hardened for its first `v0.1.0` release.

## License

Apache-2.0. See `LICENSE`.
