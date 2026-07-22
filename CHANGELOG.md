# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2026-07-22

### Added

- Vendor-neutral decisioning and organic recommendation engine.
- Pluggable candidate provider contract and synthetic marketplace provider.
- Deterministic weighted baseline ranking with stable tie-breaking.
- Optional model-assisted scoring with deterministic rollout and fail-safe baseline fallback.
- Diversity/concentration reranking and bounded deterministic exploration.
- Policy guardrails and safe empty-result fallback behavior.
- Feature, model, experiment, and exploration registry primitives.
- Stable decision IDs, score breakdowns, reason codes, and persisted decision lookup.
- Idempotent event ingestion for recommendation impression/click telemetry.
- In-memory and PostgreSQL storage adapters.
- Transactional outbox persistence and reference PostgreSQL delivery worker.
- Privacy-safe observability hooks with provider-neutral metric sink adapter.
- TypeScript telemetry SDK.
- OpenAPI 3.1 contract.
- Dockerfile and Docker Compose reference stack with PostgreSQL and outbox worker.
- Apache-2.0 licensing, security policy, contributor guidance, code of conduct, architecture documentation, public/private extraction boundary, and dependency/license audit.
- Deterministic/property tests, telemetry fuzz coverage, API contract tests, and parallel decision endpoint benchmark.
- CI gates for Go tests/vet, locked module verification, TypeScript typecheck, npm vulnerability audit, public secret/private-boundary scan, and clean Docker Compose end-to-end smoke validation.

### Security and privacy

- Public repository boundary excludes production GILLZY data, credentials, payment/wallet/KYC/fraud/dispute/support logic, paid-ad billing and money movement, proprietary seller-risk formulas, production ranking coefficients, private model artifacts/configuration, private endpoints, and deployment topology.
- Observability contract excludes user IDs, session IDs, item IDs, raw queries, messages, tokens, and other private/high-cardinality values from metric attributes.
