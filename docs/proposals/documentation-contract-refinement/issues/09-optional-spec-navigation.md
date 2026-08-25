# 09: Make to-spec an optional navigation synthesis

**Blocked by:** 05, 06

**Status:** open

**What to build:** `flowforge-to-spec` creates a compact review/navigation entry point over existing requirement and design authorities without becoming a duplicate authority or inventing missing decisions.

See [optional synthesis responsibility](../skill-coordination-design-v0.1.md#flowforge-to-spec-optional-synthesis).

## Touch points

To-spec Skill, spec role metadata, semantic link examples, information-value checks.

## Changes

1. Replace mandatory extensive user-story and implementation-decision generation with concise scope, key decisions, authority links, verification navigation, and remaining gaps.
2. Emit non-executable `role: spec` metadata.
3. Refuse to synthesize through missing requirement or design decisions; link their scoped open items instead.

## Constraints

- The synthesis must not become a second requirement or design authority.
- Human links must state semantic meaning rather than expose IDs as the reading interface.

## Done and verify

- A generated spec contains no copied authority section and never enters issue count or frontier.
- Deleting the optional spec leaves authoritative traceability intact.
