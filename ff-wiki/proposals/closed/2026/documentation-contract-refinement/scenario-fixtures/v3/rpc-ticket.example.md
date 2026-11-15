---
prototype: true
artifact_role: ticket
flowforge:
  id: migrate-rpc-route-execution
  revision: 1
blocked_by: []
requirement_dependencies:
  tangram-backend-refinement-requirements: 2
design_dependencies:
  tangram-backend-refinement-design: 3
---

# EXAMPLE — Migrate application dispatch through route execution

> Scenario fixture, not an executable Tangram ticket or an approved schema.

## Delivery

Application dispatch resolves and invokes registered routes without receiving Spring beans or reflection handles.

## Design context

Application crosses one framework-neutral route-execution seam; Spring owns discovery and runtime invocation.

See [RPC route execution](complex-design.example.md#rpc-route-execution).

## Touch points

- RPC extension contract route-execution interface and passive values.
- Spring RPC adapter discovery and invocation implementation.
- Application dispatcher and its focused behavior tests.

## Changes

1. Expand the framework-neutral route-execution contract beside the current path.
2. Adapt Spring discovery/invocation behind it.
3. Migrate application dispatch and focused tests.
4. Remove exposed bean/reflection values after no consumer remains.

## Constraints

- [Contract values contain no runtime implementation handles](complex-requirements.example.md#requirement-constraints).
- Application MUST NOT import the Spring RPC implementation.
- Route behavior and error observation MUST remain stable through migration.

## Done and verify

- A focused dispatcher test resolves and invokes a registered route through the contract seam.
- Contract source contains no bean or reflection-handle type.
- The RPC contract, Spring adapter, and application compile together.
- Exact focused Gradle commands remain an environment warning; the observable seam is settled.

## Completion evidence

See [RPC migration evidence](rpc-evidence.example.md) after execution.

