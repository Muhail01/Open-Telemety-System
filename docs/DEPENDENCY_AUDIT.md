# v0.1 Dependency and License Audit

Audit scope: dependencies declared by the committed `go.mod` and `sdk/typescript/package-lock.json` for the standalone public project.

Audit date: 2026-07-22.

This is an engineering release audit, not legal advice.

## Go runtime dependency inventory

| Module | Version | Relationship | Upstream license |
| --- | --- | --- | --- |
| `github.com/jackc/pgx/v5` | `v5.7.2` | direct | MIT |
| `github.com/jackc/pgpassfile` | `v1.0.0` | indirect | MIT |
| `github.com/jackc/pgservicefile` | `v0.0.0-20240606120523-5a60cdf6a761` | indirect | MIT |
| `github.com/jackc/puddle/v2` | `v2.2.2` | indirect | MIT |
| `golang.org/x/crypto` | `v0.31.0` | indirect | BSD 3-Clause |
| `golang.org/x/sync` | `v0.10.0` | indirect | BSD 3-Clause |
| `golang.org/x/text` | `v0.21.0` | indirect | BSD 3-Clause |

Upstream license references reviewed:

- `jackc/pgx`: https://github.com/jackc/pgx/blob/master/LICENSE
- `jackc/pgpassfile`: https://github.com/jackc/pgpassfile/blob/master/LICENSE
- `jackc/pgservicefile`: https://github.com/jackc/pgservicefile/blob/master/LICENSE
- `jackc/puddle`: https://github.com/jackc/puddle/blob/master/LICENSE
- `golang/crypto`: https://github.com/golang/crypto/blob/master/LICENSE
- `golang/sync`: https://github.com/golang/sync/blob/master/LICENSE
- `golang/text`: https://github.com/golang/text/blob/master/LICENSE

## TypeScript package inventory

The telemetry SDK currently has no JavaScript runtime dependency.

| Package | Version | Relationship | Upstream license |
| --- | --- | --- | --- |
| `typescript` | `5.9.3` | development-only | Apache-2.0 |

Upstream license reference reviewed:

- `microsoft/TypeScript`: https://github.com/microsoft/TypeScript/blob/main/LICENSE.txt

## Automated checks

CI performs:

- committed Go checksum verification with `go mod verify`;
- deterministic npm installation from `package-lock.json` using `npm ci`;
- `npm audit --audit-level=high`;
- a repository boundary scan for common credential formats, private-key material, VLESS URIs, known private GILLZY infrastructure identifiers, and private monolith import paths.

## Audit result

No copyleft/restrictive license was identified in the current declared runtime dependency set. Current declared dependencies are under MIT, BSD 3-Clause, or Apache-2.0 licenses.

No JavaScript runtime package is currently shipped by the TypeScript telemetry SDK; TypeScript itself is a development dependency.

Any new dependency must update this inventory before a tagged release. Dependencies that introduce GPL/AGPL/SSPL-style obligations or incompatible commercial restrictions require explicit maintainer review before merge.
