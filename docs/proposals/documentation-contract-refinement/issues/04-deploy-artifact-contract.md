# 04: Deploy one shared artifact contract for Skills

**Blocked by:** None

**Status:** open

**What to build:** Every production Skill that authors or consumes proposal artifacts reaches one deployed, versioned contract reference, and packaging tests fail when a Skill pointer is missing.

See [production packaging risks](../skill-coordination-design-v0.1.md#production-packaging-risks) and [Catalog/CLI responsibility](../skill-coordination-design-v0.1.md#catalogcli-responsibility).

## Touch points

Skill assets, deployment/upgrade manifest, shared references, asset deployment tests.

## Changes

1. Package the approved role, handoff, information-value, diagnostic, and schema rules as progressive-disclosure references.
2. Give each affected Skill a precise context pointer to only the contract branches it consumes.
3. Include missing mechanics/prototype references or remove invalid pointers at their source.
4. Add a deployment test that resolves every packaged Skill reference from the installed layout.

## Constraints

- Do not duplicate the full contract into every Skill or AGENTS.md.
- Do not place undeployed-only material under `assets/`.
- Preserve upstream methodology behavior outside the requirement-to-delivery flow.

## Done and verify

- A clean generated/deployed Skill tree contains every referenced resource and no broken relative pointer.
- `go test ./internal/command/... ./internal/update/...`
