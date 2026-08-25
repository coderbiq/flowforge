---
prototype: true
artifact_role: requirements
flowforge:
  id: tangram-backend-refinement-requirements
  revision: 2
---

# EXAMPLE — Tangram Backend Refinement Requirements

> Scenario fixture, not approved Tangram requirement authority or an approved schema.

## Problem and evidence

Backend implementation adapters, domain orchestration, contracts, and application composition have crossed their intended responsibilities. The previous refinement could pass a static architecture script while application compilation, configured-provider construction, and runtime seams remained broken.

## Outcomes

- Application code composes capabilities without importing implementation-internal types.
- Contract modules expose stable declarations without runtime policy, beans, reflection handles, or mutable builders.
- Configured providers can be enabled or disabled without invalid construction or startup failure.
- RPC dispatch and Trigger lifecycle remain observable through stable interfaces.
- Completion requires compilation and behavior proof in addition to static dependency checks.

## Scope

The refinement covers Agent composition, configured web construction, RPC execution, Trigger lifecycle, their cross-module contracts, and application assembly.

It does not redesign unrelated domain behavior, credentials policy, provider algorithms, or user-facing Agent features.

## Scenarios and acceptance

- An enabled configured provider is constructed from its deployment configuration and makes its capability available.
- A disabled provider is not constructed and does not prevent startup.
- Application dispatch invokes a registered RPC route without receiving reflection or container implementation values.
- Registered Triggers start only after the application is running and every started Trigger stops during shutdown.
- Application compilation contains no imports of implementation-internal or Agent-domain-internal types.

## Requirement constraints

- Application MUST NOT import or instantiate implementation adapters.
- Contract modules MUST NOT contain runtime policy or framework implementation handles.
- Static dependency checks MUST NOT be the only completion evidence.

## Unknowns

None that change the solution space in this scenario fixture. Interface shapes remain solution-design decisions.

