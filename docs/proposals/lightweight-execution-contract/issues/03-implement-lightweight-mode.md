---
flowforge:
  schema: 1
  role: ticket
---

# 03: Add lightweight execution mode to Implement skill

**Blocked by:** 01

**Status:** closed

## Delivery

The Implement skill supports a lightweight mode where a flash model executes unchecked Changes in order, runs mechanical self-checks, checks off completed Changes, writes an Implementation note, and stops—without reviewing, closing the ticket, or making design decisions.

## Design context

The current Implement skill assumes one capable agent owns the entire turn: implement, invoke review, resolve findings, write evidence, and close. When the implementer is a flash model, it cannot review, and the review runs as a separate session. The Implement skill needs a mode that produces a structured status (Implementation note) and stops, so the review agent can pick up.

See [spec](../spec.md) and [artifact contract](../../../assets/skills/_shared/ARTIFACT-CONTRACT.md).

## Touch points

- `assets/skills/flowforge-implement/SKILL.md` — the implementation process instructions

## Changes

- [x] 1. Add a "Determine execution mode" step at the start: if the ticket has unchecked Changes with mechanical steps and a write set, use lightweight mode; otherwise use the full mode (existing process)
- [x] 2. In lightweight mode, add "Execute unchecked Changes" step: work through each `- [ ]` item in order, implement the described action at the named file and symbol, run `go build` after each Change, run focused tests after relevant Changes
- [x] 3. Add "Mechanical self-check" step: run all Done and verify commands; record pass/fail for each; if any fail, stop and report the failure in Implementation note—do not attempt to debug
- [x] 4. Add "Check off completed Changes" step: change `- [ ]` to `- [x]` for each completed Change; leave unchecked any Change that could not be completed
- [x] 5. Add "Write Implementation note" step: record which Changes were completed, commands run and results, files modified, and write-set compliance statement
- [x] 6. Add "Stop" instruction: do not invoke review, do not write Completion evidence, do not close the ticket, do not make design decisions; if a Change requires a design decision, stop and note it as a design return in the Implementation note
- [x] 7. Keep the full mode (existing steps 1-5) for capable agents that own the entire turn

## Constraints

- Write set: only `assets/skills/flowforge-implement/SKILL.md`
- Do not remove the full mode; add lightweight mode alongside it
- In lightweight mode, the implementer must not change ticket Status
- In lightweight mode, the implementer must not modify files outside the ticket's Write set
- If a Done and verify command fails, the implementer stops and reports; it does not attempt fixes

## Done and verify

- `flowforge check --dir docs/proposals` passes
- Implement SKILL.md defines both lightweight and full modes
- Lightweight mode explicitly says: do not review, do not close, do not make design decisions
- `go build ./cmd/flowforge` compiles without error

---

## Execution detail

### Settled decisions

- Mode selection is implicit: if the ticket has checkbox Changes with mechanical steps and a Write set, lightweight mode applies
- Lightweight mode does not invoke flowforge-tdd independently; it follows the ticket's mechanical steps which may include test-writing actions
- Implementation note is the handoff artifact to the review agent; it is not evidence
- If a Change cannot be completed (e.g. file doesn't exist, symbol not found, build fails after the change), the implementer stops, leaves the Change unchecked, and describes the blocker in the Implementation note
- The implementer runs Done and verify commands but does not judge whether the results are "correct"—it records pass/fail and observed output

### Expected tests

No Go tests. Verification:
- Both modes are described in the skill
- Lightweight mode has a clear "stop" boundary: no review, no close, no evidence
- Implementation note section is defined in the process

## Completion evidence

Delivered: Implement SKILL.md now defines two execution modes. Lightweight mode (step 3) executes unchecked Changes mechanically, runs build and focused tests per Change, runs full Done/verify commands, checks off completed Changes, writes an Implementation note, and stops—without reviewing, closing, or making design decisions. Full mode (step 4) preserves the existing capable-agent flow. Mode selection (step 2) is implicit from ticket structure: checkbox Changes + Write set + Execution detail → lightweight; otherwise full.
Verification: `flowforge check --dir docs/proposals` — healthy DAG; `go build -trimpath -o /tmp/opencode/flowforge ./cmd/flowforge` — OK.
Review: skipped (skill instruction change, no Go code; structural consistency verified by check).
Modified: `assets/skills/flowforge-implement/SKILL.md` only.
