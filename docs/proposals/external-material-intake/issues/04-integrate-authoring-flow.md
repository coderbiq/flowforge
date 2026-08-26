---
flowforge:
  schema: 1
  role: ticket
  id: integrate-authoring-flow
  revision: 1
  consumes:
    requirements:
      external-material-intake-requirements: 1
    design:
      external-material-intake-design: 1
---

# 04: Integrate import, rewriting, and publication checks into the Skill flow

**Blocked by:** 01, 02

**Status:** open

**Delivery:** Route, Align, Solution Design, Plan, and Review guide a request that begins with external material through the right authority owner, semantic rewrite, strict publication check, and only then user-approved ticket planning.

## Design context

Apply the [Import ownership boundaries and writing rules](../design.md#external-material-intake-design) across existing Skills to meet [the requirement authority's source-to-authority outcome](../requirements.md#external-material-intake-requirements). The human-readable document remains authoritative; this integration must not create a status gate.

## Touch points

`assets/skills/flowforge-route/`; `assets/skills/flowforge-align/`; `assets/skills/flowforge-solution-design/`; `assets/skills/flowforge-plan/`; `assets/skills/flowforge-review/`; `README.md`; `docs/` workflow documentation; skill packaging tests.

## Changes

1. Route material-first requests to Import, then retain existing direct-request routes.
2. Make authority-owning Skills perform strict feature-directory validation after writing schema authority and repair only diagnostics introduced by their edit.
3. Apply shared language and information-value rules at the authoring and Standards-review boundaries.
4. Document the end-to-end external-material branch in the theory-first, actual-demand workflow documentation.

## Constraints

- Plan does not create tickets before its title, Delivery, and genuine DAG edges are accepted.
- Keep direct simple work compact; do not require source notes or all artifact roles for every request.

## Done and verify

- A documented mixed-source scenario reaches Import → Align → Solution Design → Plan without mechanical document conversion, empty artifact roles, or status transitions.
- Skill reference/package tests and relevant documentation checks pass: `go test ./internal/command/...` and `git diff --check`.
