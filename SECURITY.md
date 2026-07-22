# Security Policy

## Reporting a vulnerability

Please do not open a public issue for vulnerabilities that could expose credentials, personal data, authorization material, or production systems.

For the initial public phase, contact the project maintainer privately through the GitHub profile associated with this repository. A dedicated security contact may be published later.

## Security principles

Open Telemetry System is designed around strict separation between reusable decisioning infrastructure and private marketplace data.

The public project must never contain:

- production credentials or tokens;
- authentication cookies, JWTs, or authorization headers;
- raw private messages;
- payment or KYC payloads;
- personally identifiable information in examples or fixtures;
- private production endpoints;
- proprietary ranking coefficients copied from a commercial deployment.

Recommendation telemetry should use opaque decision, session, user, and item identifiers. Deployers are responsible for applying their own retention, consent, privacy, and regional compliance requirements.
