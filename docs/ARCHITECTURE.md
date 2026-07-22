# Architecture

Open Telemetry System is deliberately split into a small decisioning kernel and replaceable adapters.

## Online path

```text
Request
  -> CandidateProvider
  -> Scorer
       -> deterministic baseline
       -> optional model contribution
       -> deterministic rollout / fail-safe fallback
  -> Policy
  -> Diversity rerank
  -> bounded exploration
  -> Decision
  -> Store
       -> decision persistence
       -> transactional outbox
```

## Feedback path

```text
Client surface
  -> impression / click telemetry SDK
  -> POST /v1/events
  -> idempotent event persistence
  -> transactional outbox
  -> downstream analytics / feature pipelines
```

A stable `decision_id` joins the serving decision to downstream impression, click, and outcome events.

## Core packages

- `internal/core`: contracts, engine, ranking, hybrid model scoring, reranking, policy, exploration, registries.
- `internal/demo`: synthetic candidate source used by the runnable example.
- `internal/memory`: zero-dependency development adapter.
- `internal/postgres`: durable event/decision storage and transactional outbox.
- `internal/httpapi`: minimal HTTP transport.
- `sdk/typescript`: browser telemetry client.

## Extension points

### CandidateProvider

Connect databases, search engines, vector stores, or domain-specific candidate generators without changing the decision engine.

### Scorer

The default scorer is a deterministic weighted feature scorer. `HybridScorer` can add optional model scores while preserving baseline behavior when the model is unavailable or outside rollout.

### ModelScoreProvider

Model integrations are provider-neutral. An implementation can call a local model, remote model, learned-to-rank service, or another scoring system.

### Policy

Policies run after scoring and before final serving. This is the intended place for deployment-specific eligibility and safety rules. Proprietary marketplace rules should live outside the public core.

### Store

The in-memory implementation is suitable for tests and demos. PostgreSQL provides durable decisions/events and writes outbox messages in the same transaction.

## Determinism

The baseline ranking path is deterministic:

- explicit feature weights;
- stable tie-breaking by candidate key;
- bounded per-group concentration;
- deterministic rollout bucketing;
- deterministic exploration selection for a stable request identity.

This makes shadow evaluation, replay, debugging, and regression testing practical.

## Privacy boundary

The core expects opaque identifiers and generic feature values. It does not require email addresses, phone numbers, payment details, KYC material, private messages, access tokens, cookies, or raw provider responses.

Deployers are responsible for their own legal basis, consent, retention, deletion, and regional compliance policies.
