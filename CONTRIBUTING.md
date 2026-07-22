# Contributing

Thank you for contributing to Open Telemetry System.

## Development

```bash
go test ./...
go vet ./...
```

Keep the core vendor-neutral. New functionality should be expressed through reusable contracts, interfaces, configuration, or adapters rather than marketplace-specific conditions.

## Pull requests

A pull request should include:

- a clear problem statement;
- tests for behavioral changes;
- deterministic behavior where possible;
- documentation for public API changes;
- no secrets, production data, or private commercial rules.

## Architecture boundaries

Suitable contributions include candidate providers, generic feature sources, ranking plugins, reranking algorithms, policy guardrails, exploration strategies, storage adapters, telemetry SDKs, observability, and developer tooling.

Do not contribute copied proprietary business rules, private risk formulas, payment flows, KYC data models, production credentials, or customer data.

## Compatibility

The core runtime should remain usable without a proprietary model provider. Model-assisted scoring belongs behind optional interfaces and adapters.
