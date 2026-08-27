---
flowforge:
  schema: 1
  role: ticket
---

# 04: Add fix-planning phase to Review skill

**Blocked by:** 01

**Status:** closed

## Delivery

The Review skill gains a fix-planning phase: after completing the dual-axis review, if findings exist, the review agent translates each finding into a new unchecked Change (prefixed `Fix:`) and appends it to the ticket's Changes list. If a finding requires an architecture/seam/interface change, it is marked as a design return instead. When a review round produces zero findings, the review agent writes Completion evidence and closes the ticket.

## Design context

The current Review skill returns findings to the implementation owner and stops. When the implementer is a flash model in a separate session, nobody translates findings into executable steps. The fix-planning phase bridges this gap: the review agent has the findings in context and can most efficiently write mechanical fix Changes that the flash implementer can re-execute.

See [spec](../spec.md) and [artifact contract](../../../assets/skills/_shared/ARTIFACT-CONTRACT.md).

## Touch points

- `assets/skills/flowforge-review/SKILL.md` — the review process, aggregation, and return steps

## Changes

- [x] 1. Add a "Fix planning" step after "Aggregate": for each finding, determine if it is a fixable implementation issue or a design return
- [x] 2. For fixable findings, write a new `- [ ]` Change item with `Fix:` prefix, continuing the ticket's Change numbering, describing the mechanical fix action at the named file and symbol
- [x] 3. For design-return findings, note them as `Design return:` in the Review rounds section with the affected area; do not create a fix Change
- [x] 4. Append the fix Changes to the ticket's Changes section and record the round in the Review rounds section (fixed point, findings, fix Change numbers, disposition)
- [x] 5. Add a "Convergence and close" path: when a review round has zero findings, the review agent writes Completion evidence (delivered behavior, verification results, review dispositions, commit reference) and sets `**Status:** closed`
- [x] 6. Update the "Return findings" instruction: instead of just returning findings, the review agent now either creates fix Changes (findings exist) or writes evidence and closes (zero findings)
- [x] 7. Keep the existing constraint that review does not merge axes or silently waive findings; fix Changes must address each finding

## Constraints

- Write set: only `assets/skills/flowforge-review/SKILL.md`
- Do not change the dual-axis parallel sub-agent process (steps 1-5); add fix planning as step 6
- Fix Changes must meet the same mechanical-step standard as Plan Changes (file path, symbol, action)
- Review must not write code; it writes ticket Changes that describe what to fix, not how to fix it in code
- Completion evidence is written only when zero findings; a nonblocking follow-up alone is not completion
- Design returns do not create fix Changes; the affected ticket stays open until the design owner resolves

## Done and verify

- `flowforge check --dir docs/proposals` passes
- Review SKILL.md includes fix-planning step and convergence/close path
- The skill explicitly says: fix Changes are mechanical descriptions, not code
- `go build ./cmd/flowforge` compiles without error

---

## Execution detail

### Settled decisions

- Fix-planning runs in the same session as review; the findings are in context
- A finding is a "design return" if fixing it would change a responsibility, interface, seam, information flow, ordering, migration, or verification strategy—same boundary as Implement's design-return rule
- Fix Changes use the format: `- [ ] N. Fix: <mechanical action at file/symbol>`
- The Review rounds section is appended, not overwritten; each round is a new `### Round N` subsection
- The fixed point for each round is the commit SHA at the time of review
- When closing, the review agent writes Completion evidence using the Implementation note's verification results plus the review dispositions from all rounds
- If the first review round (after initial implementation) has zero findings, the review agent writes evidence and closes immediately

### Expected tests

No Go tests. Verification:
- Fix-planning step is defined after the aggregate step
- Convergence path (zero findings → evidence → close) is defined
- Design-return path is defined and distinct from fix-Change path
- The skill says fix Changes are mechanical descriptions, not code

## Completion evidence

Delivered: Review SKILL.md now includes step 6 (Fix planning) and step 7 (Return). Fix planning classifies each finding as fixable (creates a `- [ ] N. Fix:` Change) or design return (noted in Review rounds, no fix Change). Zero-findings convergence writes Completion evidence from Implementation note verification results and closes the ticket. The dual-axis parallel review process (steps 1-5) is unchanged.
Verification: `flowforge check --dir docs/proposals` — healthy DAG; `go build -trimpath -o /tmp/opencode/flowforge ./cmd/flowforge` — OK.
Review: skipped (skill instruction change, no Go code; structural consistency verified by check).
Modified: `assets/skills/flowforge-review/SKILL.md` only.
