# Plan Production Validation v0.1

**Date:** 2026-08-25

## Procedure

Apply the production `flowforge-plan` instructions to each fixture. Each preserved fixture contains requirement/design authorities and an `issues-example/` directory of exact generated tickets using the non-catalog `.example` suffix. Copy one fixture to an isolated root, rename `issues-example/` to `issues/`, and copy each `.example` ticket to the same basename with a `.md` suffix, then run:

```bash
./bin/flowforge check --dir <temporary-root>
./bin/flowforge frontier --dir <temporary-root>
```

The source fixtures sit outside an `issues/` directory, so they are documentation examples rather than executable tickets in this proposal.

## Simple existing-seam fixture

**Input:** Add a warning count to the parser's existing diagnostics result. Requirement behavior, parser seam, and verification command are settled by current authorities.

**Expected packaging:** One compact ticket with no blockers.

**Generated output:** [Expose parser warning count](scenario-fixtures/v4/plan-simple/issues-example/01-warning-count.example).

**Checklist:**

- PASS — Delivery states warning-count behavior once.
- PASS — Design context links the current parser-seam authority and revision.
- PASS — Changes preserve the approved seam and contain no architecture choice.
- PASS — Done and verify pair observable warning output with the parser test command.
- PASS — no empty Goal, Summary, Acceptance, or separate Verification repetition.

## Cross-module fixture

**Input:** Move configured provider selection into application composition while concrete construction remains in the web adapter. Authorities settle the interface, migration order, compatibility interval, and verification seams.

**Expected packaging:** Expand ticket; independently parallel consumer-migration tickets blocked by expand; contract ticket blocked by every migration.

**Generated outputs:** [expand composition](scenario-fixtures/v4/plan-complex/issues-example/01-expand-composition.example), [migrate CLI](scenario-fixtures/v4/plan-complex/issues-example/02-migrate-cli.example), [migrate server](scenario-fixtures/v4/plan-complex/issues-example/03-migrate-server.example), and [contract legacy path](scenario-fixtures/v4/plan-complex/issues-example/04-contract-legacy.example).

**Checklist:**

- PASS — every ticket delivers a green, independently observable increment.
- PASS — blocking edges encode only compatibility and removal prerequisites.
- PASS — each Design context contains a local summary and semantic authority link rather than copied rationale.
- PASS — stable modules and test seams appear in Touch points; speculative affected-file inventories do not.
- PASS — contract waits for all migrations and verifies absence of legacy callers.

## Design-gap fixture

**Input:** “Choose whether provider construction belongs in Application or Web, then make tickets.”

**Expected result:** Return the ownership choice to `flowforge-solution-design`; publish no implementation action that selects it.

## Result

The simple fixture check reported a healthy one-ticket graph and frontier `#01`. The complex fixture check reported a healthy four-ticket graph; with expand closed, frontier reported parallel tickets `#02` and `#03`, while `#04` remained blocked by both. Exact commands were run by copying each linked fixture set into a temporary proposal's `issues/` directory and invoking the two commands above.

The unsettled ownership case returns to Solution Design. Across the linked ticket text, every retained paragraph contributes delivery, local design context, location, decided action, constraint, or paired verification information.
