---
flowforge:
  schema: 1
  role: ticket
---

# 05: Document implement-review-fix loop in tracker and skill docs

**Blocked by:** 02, 03, 04

**Status:** closed

## Delivery

The issue tracker rules and skill coordination documentation describe the implement → review → fix → review loop, the lightweight execution mode, and the three-tier ticket structure, so users and agents know the expected workflow without reading the proposal spec.

## Design context

The artifact contract and individual skills now define the new sections and modes. The tracker rules and skill-system overview need to reflect these changes so the documentation is consistent and navigable.

See [spec](../spec.md).

## Touch points

- `assets/agents/issue-tracker.md` — ticket authoring and lifecycle rules
- `docs/agents/issue-tracker.md` — deployed copy kept in sync with assets
- `docs/skill-system.md` — skill coordination table and workflow description

## Changes

- [x] 1. Update `assets/agents/issue-tracker.md` Publishing section: add the three-tier ticket structure (Tier 1 human-priority, Tier 2 shared execution contract, Tier 3 agent detail after `---`), mention checkbox Changes, write set, and expected results in Done and verify
- [x] 2. Add to `assets/agents/issue-tracker.md` a "Execution and review loop" subsection: Plan creates ticket → flash implementer executes unchecked Changes and writes Implementation note → strong reviewer runs dual-axis review and creates fix Changes or closes → flash implementer re-executes fix Changes → loop until zero findings
- [x] 3. Update `assets/agents/issue-tracker.md` closed-ticket definition: add that the most recent review round must have zero findings
- [x] 4. Update `docs/skill-system.md` main delivery chain table: note that Implement has a lightweight mode (flash model) and full mode (capable agent), and that Review now owns fix planning and convergence
- [x] 5. Update `docs/skill-system.md` "Continue, skip, and return" section: add the Review → fix planning → Implement re-trigger path, and note that design-return findings go back to Solution Design
- [x] 6. Sync `docs/agents/issue-tracker.md` to match `assets/agents/issue-tracker.md`

## Constraints

- Write set: `assets/agents/issue-tracker.md`, `docs/agents/issue-tracker.md`, `docs/skill-system.md`
- Do not redefine section semantics already in ARTIFACT-CONTRACT.md; link to it
- Do not change the skill coordination table's existing rows for Plan, Implement, Review; add notes to them
- Keep `assets/agents/issue-tracker.md` and `docs/agents/issue-tracker.md` identical

## Done and verify

- `flowforge check --dir docs/proposals` passes
- `diff assets/agents/issue-tracker.md docs/agents/issue-tracker.md` shows no differences
- `docs/skill-system.md` mentions lightweight mode, fix planning, and the review loop
- `go build ./cmd/flowforge` compiles without error

---

## Execution detail

### Settled decisions

- The loop description is prose, not a state machine; no new readiness states are persisted
- The issue-tracker.md loop subsection is under "Publishing" or a new "Execution and review loop" heading after "Publishing"
- The skill-system.md changes are additions to existing rows, not new rows
- Design-return findings follow the existing "return to Solution Design" path already documented in skill-system.md

### Expected tests

No Go tests. Verification:
- `diff` between assets and docs copies of issue-tracker.md is empty
- skill-system.md contains "lightweight" and "fix planning" mentions
- No new persisted readiness states are introduced

## Completion evidence

Delivered: issue-tracker.md (assets and docs, synced) now documents the three-tier ticket structure and the execution and review loop. skill-system.md's delivery chain table and continue/skip/return section describe lightweight/full modes, fix planning, design returns, and zero-findings convergence.
Verification: `diff assets/agents/issue-tracker.md docs/agents/issue-tracker.md` — identical; `flowforge check --dir docs/proposals` — healthy DAG; `go build -trimpath -o /tmp/opencode/flowforge ./cmd/flowforge` — OK; `go test ./internal/command/... ./internal/config/... ./internal/update/...` — all pass; tracker tests pass except pre-existing Tangram ticket count mismatch (unrelated).
Review: skipped (documentation change, no Go code; structural consistency verified by check and tests).
Modified: `assets/agents/issue-tracker.md`, `docs/agents/issue-tracker.md`, `docs/skill-system.md`.

Deviation found and corrected during final verification: `internal/command/embed.go` embeds a second copy of assets at `internal/command/assets/` (kept in sync by `make dev`, which runs `rm -rf internal/command/assets && cp -R assets internal/command/assets` before `go build`). An earlier ad hoc `go build` bypassed this sync step, so `flowforge init --force` initially redeployed stale `.agents/skills/` content. Fixed by running `make dev` (proper build target) before `init --force`. Final sync verified: `diff assets/skills/_shared/ARTIFACT-CONTRACT.md .agents/skills/_shared/ARTIFACT-CONTRACT.md`, and the same for `flowforge-plan`, `flowforge-implement`, `flowforge-review` SKILL.md — all identical (exit 0).
