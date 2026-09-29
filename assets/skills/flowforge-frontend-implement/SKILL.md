---
name: flowforge-frontend-implement
description: Deliver frontend/UI tickets with design-system context and a mandatory screenshot self-review loop — follow the ticket's transcribed clauses first, consult slot slices by anchor only, implement, run the declared verification entry, critique screenshots against the counter-example list, fix, repeat. Use for tickets touching JSX, styles, layout, or component layers. Not for backend logic or API work.
---

# Frontend implement with visual self-review (v2)

This skill is the frontend specialization of `flowforge-implement`: it inherits
that skill's ticket contract, evidence rules, STATUS discipline, and closeout.
It adds the visual loop only — do not restate or override the parent contract.

## Hierarchy of authority (in order)

1. **Ticket clauses** — the transcribed `must`/`must not` frontend clauses in
   the ticket's Constraints are the ONLY compliance source you must satisfy.
2. **Slot slices by anchor** — when a clause cites a slot file + clause number,
   read that slice (not the whole file) to resolve the clause's exact meaning.
3. **Counter-example list** — the numbered anti-pattern list used for the
   critique pass (visual comparison material, not prose to internalize).

Never read a full design spec front-to-back during execution; understanding
standards is the design session's job, already distilled into your ticket.

## Project slots (declared by the project's manifest)

`.agents/references/frontend/manifest.md` declares:

1. `design-tokens` — token source (DTCG recommended, not required)
2. `component-entry` — component allowlist + self-check command (project's own
   enforcement mechanism; run it when the ticket cites it)
3. `counter-examples` — numbered anti-pattern list
4. `verification-entry` — repeatable visual verification command that prints
   grep-able PASS/FAIL with numeric values

## Process

### 0. Pre-flight (all three, immediately, before any work)

- Slots complete per manifest? Missing slot → `STATUS: BLOCKED` naming it.
- Runtime accepts image input? No → `STATUS: BLOCKED` with that exact reason.
- verification-entry command exists and is runnable? No → `STATUS: BLOCKED`.

A late BLOCKED after implementation is a defect of this process; fail here,
fail now.

### 1. Implement (negative constraints first)

- Do not hand-assemble Layout/Header/Sider/Menu/Footer outside the entry list.
- Do not hardcode color values inline; use tokens/theme variables.
- Do not bypass or double-wrap the theme provider.
- Do not invent a second component system for what the entry list covers.

### 2. Verify

Run the declared verification-entry. Its grep-able output is the compliance
record. Interactions route through accessibility snapshots; screenshots are
for visual verification only.

### 3. Critique (the vision step)

Load the delivered screenshots as images. Critique against the
counter-example list clause by clause, citing the clause number for every
negative finding. A critique with zero findings must state which clauses were
checked — silence is not a pass. Never critique from memory or DOM
assumptions; if you cannot look at the image, you cannot pass this step.

### 4. Fix and repeat

Fix findings, then repeat 2–3. Stop after 5 failed repair rounds or when the
iteration budget is nearly exhausted (parent skill's budget-closure rule).

### 5. Deliver

Attach to the ticket: the final verification-entry output (verbatim), the
screenshot paths, and the per-clause critique record. Then the parent-skill
closeout (STATUS, verification, evidence). Numbers in the output are hard
evidence — a later reviewer will cross-check them against the screenshots.
