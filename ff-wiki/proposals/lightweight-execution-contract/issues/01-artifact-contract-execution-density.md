---
flowforge:
  schema: 1
  role: ticket
---

# 01: Add execution-density sections and write set to artifact contract

**Blocked by:** None

**Status:** closed

## Delivery

The shared artifact contract defines the three new ticket sections (Execution detail, Implementation note, Review rounds), the write-set convention in Constraints, and the checkbox format for Changes, so that Plan, Implement, and Review skills can reference them.

## Design context

The artifact contract is the single source of truth for what a ticket contains. Plan, Implement, and Review skills all link to it. Adding the new sections here lets the skill changes be purely instructional—each skill says "follow the contract" rather than re-defining structure.

See [spec](../spec.md).

## Touch points

- `assets/skills/_shared/ARTIFACT-CONTRACT.md` — ticket section definitions, Constraints, Changes format

## Changes

- [x] 1. Add `Write set:` to the Constraints definition: an explicit line listing the only directories or files the implementer may modify
- [x] 2. Change Changes definition to require `- [ ]` checkbox format and mechanical-level action description
- [x] 3. Add `Execution detail` section definition: placed after a `---` separator, containing Settled decisions, Expected tests, and Conventions subsections; visible but skippable by humans
- [x] 4. Add `Implementation note` section definition: written by implementer after execution; records completed Changes, commands and results, files modified, write-set compliance; not evidence
- [x] 5. Add `Review rounds` section definition: accumulates review history per round; each round records fixed point (commit SHA), Standards/Spec findings, fix Changes created, disposition
- [x] 6. Update the ticket section order to include the new sections after Done and verify, separated by `---`

## Constraints

- Write set: only `assets/skills/_shared/ARTIFACT-CONTRACT.md`
- Do not change the requirement, design, or evidence role definitions
- Do not add machine-parsed YAML fields for the new sections; they are human-visible Markdown
- Keep the information-value test unchanged; the new sections must pass it

## Done and verify

- `flowforge check --dir docs/proposals` passes with no new diagnostics
- ARTIFACT-CONTRACT.md defines all three new sections, write set, and checkbox format
- Existing `documentation-contract-refinement` proposal tickets remain valid under the updated contract
- `go build ./cmd/flowforge` compiles without error

---

## Execution detail

### Settled decisions

- Execution detail, Implementation note, and Review rounds are human-visible Markdown sections, not YAML metadata
- Section order in ticket: Delivery, Design context, Blocked by, Touch points, Changes, Constraints, Done and verify, then `---`, then Execution detail, Implementation note (written after execution), Review rounds (written after review)
- Changes use `- [ ]` for unchecked and `- [x]` for completed; fix Changes from review use `Fix:` prefix in the item text
- Write set is one line starting with `Write set:` inside the Constraints section
- Implementation note is not evidence; it is a status report for the review agent
- Review rounds accumulate; they are not overwritten—each round is a new subsection

### Expected tests

No Go tests. Verification is structural consistency:
- ARTIFACT-CONTRACT.md internally references all new sections it defines
- No contradiction between the new section definitions and the existing hand-offs or diagnostics sections

## Completion evidence

Delivered: ARTIFACT-CONTRACT.md now defines three information tiers (human-priority, shared execution contract, agent execution detail), checkbox Changes with `Fix:` prefix convention, Write set in Constraints, Done and verify with expected results, Execution detail section, Implementation note section, and Review rounds section. A new "Execution and review loop" section describes the Plan → Implement (lightweight) → Review → fix → loop → convergence path.
Verification: `flowforge check --dir docs/proposals` — healthy DAG, no new diagnostics; `go build -trimpath -o /tmp/opencode/flowforge ./cmd/flowforge` — OK.
Review: skipped (documentation contract change, no Go code; structural consistency verified by check).
Modified: `assets/skills/_shared/ARTIFACT-CONTRACT.md` only.
