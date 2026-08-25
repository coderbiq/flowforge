# 07: Plan high-information tracer tickets

**Blocked by:** 05, 06

**Status:** open

**What to build:** Plan publishes one independently verifiable increment per ticket with genuine DAG edges, locally sufficient design context, ordered decided changes, scoped constraints, and paired completion/verification evidence requirements.

See [Ticket interface](../artifact-handoff-interface-v0.2.md#ticket-interface) and [Plan coordination contract](../skill-coordination-design-v0.1.md#flowforge-plan).

## Touch points

Plan Skill instructions, ticket template/reference, semantic dependency metadata, frontier validation loop.

## Changes

1. Replace overlapping goal/summary/acceptance prose with the approved Delivery, Design context, Blocked by, Touch points, Changes, Constraints, and Done and verify semantics.
2. Record consumed authority revisions and human semantic links without making IDs the reading interface.
3. Return any hidden architecture choice to solution design instead of publishing it as an implementation action.
4. Run Catalog/DAG checks after publishing and explain warnings without treating them as stale workflow state.

## Constraints

- Tickets must remain vertical, independently verifiable, and small enough for one fresh context.
- Empty sections, minimum word counts, and repeated upstream prose must not be required.

## Done and verify

- Simple work remains one compact ticket; complex work produces linked tracer tickets with valid blocking edges.
- Information-value review can delete no retained paragraph without losing a required semantic category.
- `go run ./cmd/flowforge check --dir docs/proposals`
