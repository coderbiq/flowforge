---
name: flowforge-frontend-implement
description: Deliver frontend/UI tickets with design-system context and a mandatory screenshot self-review loop — read the project's design spec slot, respect the component single-entry list, implement, screenshot, critique against the counter-example list, fix, repeat. Use for tickets touching JSX, styles, layout, or component layers. Not for backend logic or API work.
---

# Frontend implement with visual self-review

This skill is the frontend specialization of `flowforge-implement`: it inherits
that skill's ticket contract, evidence rules, STATUS discipline, and closeout.
It adds the visual loop only — do not restate or override the parent contract.

## Project slots (must be declared by the project's AGENTS.md)

1. `design-spec`: path to the authoritative design-system guide (e.g. a
   numbered-clauses visual standard). Missing slot → `STATUS: BLOCKED`.
2. `component-entry-list`: the single-entry component package and its allowed
   surface (e.g. `@tangram/ui` AppLayout/ThemeProvider/ReadableRegistry; the
   app layer must not import the underlying UI library directly).
3. `counter-examples`: the numbered anti-pattern list (violations to check
   screenshots against). Missing slot → `STATUS: BLOCKED`.

## Process

### 1. Load context

Read the three slots above before writing any code. Re-read the
counter-example list immediately before each critique pass — it is task-time
context, not background knowledge.

### 2. Implement (negative constraints first)

- Do not hand-assemble Layout/Header/Sider/Menu/Footer outside the entry list.
- Do not hardcode color values inline; use tokens/theme variables.
- Do not bypass the theme provider or wrap it a second time.
- Do not invent a second component system for something the entry list covers.

### 3. Screenshot

Start (or reuse) the dev server, then capture desktop AND mobile screenshots
of the changed surface with Playwright (or the project's declared screenshot
entry). Interactions route through accessibility snapshots; screenshots are
for visual verification only.

- No screenshot entry declared → `STATUS: BLOCKED` (do not hand-wave visual
  verification).

### 4. Critique (the vision step)

Load the screenshots as images and critique them against the
counter-example list, clause by clause, citing the clause number for every
negative finding. A critique with zero findings must state which clauses were
checked — silence is not a pass.

- If your runtime cannot accept image input, `STATUS: BLOCKED` with that
  exact reason — never critique from memory or DOM assumptions.

### 5. Fix and repeat

Fix findings, then repeat 3–4. Stop after 5 failed repair rounds or when the
iteration budget is nearly exhausted (parent skill's budget-closure rule) and
report.

### 6. Deliver

Attach to the ticket: the final screenshots (paths), the per-clause critique
record, and the normal parent-skill closeout (STATUS, verification, evidence).
