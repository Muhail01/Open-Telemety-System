# Open-source boundary

This document defines what may move from the private GILLZY codebase into Open Telemetry System.

## Public

The public project may contain reusable implementations of:

- generic event envelopes and validation;
- idempotent event ingestion;
- transactional outbox patterns;
- session-safe opaque identifiers;
- candidate provider interfaces;
- generic feature maps and feature definitions;
- deterministic ranking algorithms;
- configurable example weights created specifically for the public project;
- model scorer interfaces and provider-neutral adapters;
- diversity and concentration controls;
- policy engine primitives;
- bounded exploration strategies;
- decision logs and explainability contracts;
- feature/model/experiment/exploration registries;
- impression/click/outcome telemetry contracts;
- in-memory and generic PostgreSQL adapters;
- synthetic fixtures and demo catalogs;
- metrics, tracing, tests, benchmarks, documentation, and SDKs.

## Private

The following remain in GILLZY or another private repository:

- production user, seller, supplier, order, payment, wallet, and dispute data;
- credentials, tokens, private endpoints, infrastructure addresses, and deployment topology;
- KYC, sanctions, fraud, abuse, and enforcement implementations;
- private seller-risk formulas and thresholds;
- production recommendation coefficients and commercial tuning;
- paid advertising auction, billing, reserve, charging, payout, or money-movement logic;
- private model artifacts, training data, prompts, features, and production model configuration;
- internal campaign, advertiser, and commercial targeting rules;
- raw customer messages, support content, provider responses, stack traces, and operational incident data;
- GILLZY-specific database joins or adapters that expose private schema assumptions.

## Extraction rule

Code that is cleanly generic may be adapted into the public project after a secret/privacy/license review. Code that mixes generic decisioning with GILLZY commercial behavior must be reimplemented behind a neutral interface rather than copied verbatim.

The public repository must be independently understandable and runnable without access to GILLZY.
