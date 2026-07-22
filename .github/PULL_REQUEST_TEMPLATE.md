## What changed

Describe the change and the user/developer problem it solves.

## Why this belongs in GMF Core

Explain how the change stays vendor-neutral and reusable outside any one marketplace implementation.

## Validation

- [ ] `go test ./...`
- [ ] `go vet ./...`
- [ ] TypeScript checks pass when SDK code changes
- [ ] Docker Compose smoke still works when runtime/storage code changes
- [ ] New behavior has focused tests or a reproducible example
- [ ] Documentation/OpenAPI updated when public behavior changes

## Public/private boundary

- [ ] No production credentials, tokens, private endpoints, or deployment topology
- [ ] No real user/marketplace production data
- [ ] No proprietary ranking coefficients or seller/fraud-risk formulas
- [ ] No GILLZY-specific private integrations or commercial billing logic

## Notes for reviewers

Call out compatibility concerns, follow-up work, performance tradeoffs, or areas where you specifically want feedback.
