# Roadmap

## v0.1.0 release gate

- CI green for Go and TypeScript.
- Commit dependency checksum/lock files.
- Run secret, license, and dependency audits.
- Add property and fuzz tests for deterministic ranking, event idempotency, and unsafe telemetry input.
- Add a basic load-test profile for the decision endpoint.
- Add an outbox worker/reference delivery loop.
- Add decision lookup for local debugging and explainability.
- Add OpenTelemetry-compatible metrics/tracing hooks without user PII in labels.
- Verify Docker Compose quick start from a clean checkout.
- Publish release notes and tag `v0.1.0`.

## Post-v0.1.0

- Candidate provider registry and reference search/vector adapters.
- Feature enrichment pipeline and feature snapshot interface.
- Provider-neutral remote model scorer adapters.
- Experiment assignment and evaluation helpers.
- React telemetry bindings.
- Go client SDK.
- Replay/shadow-evaluation tooling.
- Prometheus/OpenTelemetry observability package.
- Reference example storefront using synthetic data.

The core must remain useful without proprietary model providers or access to GILLZY systems.
