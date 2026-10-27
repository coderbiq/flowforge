---
prototype: true
artifact_role: design
flowforge:
  id: tangram-backend-refinement-design
  revision: 3
requirement_dependencies:
  tangram-backend-refinement-requirements: 2
design_areas:
  agent-composition:
    status: gap
  configured-web-construction:
    status: covered
  rpc-route-execution:
    status: covered
  trigger-lifecycle:
    status: covered
---

# EXAMPLE — Tangram Backend Refinement Design

> Scenario fixture, not approved Tangram design authority or an approved schema.

## Requirement coverage

- [Application composes without implementation imports](complex-requirements.example.md#outcomes) through application-owned composition interfaces.
- [Contracts contain declarations without runtime policy](complex-requirements.example.md#outcomes) by keeping construction, reflection, and invocation inside adapters.
- [Configured providers honor enabled and disabled states](complex-requirements.example.md#scenarios-and-acceptance) through an adapter factory driven by application-bound settings.
- [RPC and Trigger behavior remains observable](complex-requirements.example.md#scenarios-and-acceptance) through route-execution and lifecycle seams.
- [Completion exceeds static checks](complex-requirements.example.md#requirement-constraints) through compile and behavior verification at those seams.

## Responsibility map

| Module | Owns | Excludes |
|---|---|---|
| Contracts | Framework-neutral declarations and passive cross-seam values | Builders, beans, reflection handles, construction policy |
| Agent domain | Agent orchestration and an application-facing composition interface | Deployment binding and concrete web adapters |
| Task domain | Trigger registration and lifecycle coordination | Scheduler polling implementation |
| Implementation adapters | Concrete RPC, scheduler, persistence, and web behavior | Public domain interfaces and application assembly policy |
| Spring/container adapter | Static discovery and runtime route invocation | Domain and provider-selection decisions |
| Application | Deployment binding, adapter selection, and composition | Imports of implementation or domain-internal types |

## Selected seams

### Configured web construction

Application owns provider selection and configuration binding. A contract factory accepts passive provider settings; the web adapter owns concrete construction.

- Disabled providers MUST NOT be constructed or prevent startup.
- Application MUST NOT import implementation classes.

### RPC route execution

Application dispatch crosses a framework-neutral route-execution interface. The Spring RPC adapter owns bean discovery, reflection, and invocation.

- Contract values MUST NOT carry beans or reflection handles.

### Trigger lifecycle

Task domain owns Trigger coordination; scheduler supplies a Trigger adapter. Registration may occur early, but start and stop follow application lifecycle transitions.

- A Trigger MUST NOT start before the application is running.
- Every started Trigger MUST stop during shutdown.

## Migration and verification

1. Expand the selected interfaces beside existing paths.
2. Prove each seam with focused behavior tests.
3. Migrate consumers in coherent batches.
4. Remove old imports, runtime contract values, and obsolete construction only after migration.
5. Run cross-module compilation and integration behavior before close-out.

Static dependency checks remain additional evidence, not the completion proof.

## Credible rejected alternatives

- Annotation scanning constructs configured adapters: rejected because scanning cannot derive runtime provider values.
- RPC contracts expose reflection handles: rejected because it moves runtime invocation implementation into contract knowledge.

## Open design question

### What is the Agent application-composition interface?

Responsibility is settled, but the smallest framework-neutral interface shape is not. This gap affects Agent composition and application Agent migration only.

