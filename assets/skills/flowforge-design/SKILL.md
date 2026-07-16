---
name: flowforge-design
description: Use ONLY when the user wants to design or decompose a FlowForge proposal before implementation. Do NOT use for implementation, feedback, archiving, or single-card creation.
---

# flowforge-design

## Start

Run `flowforge project current`, `flowforge proposal current`, `flowforge proposal inspect <id>`. If project/proposal is missing, ask the user.

## Workflow

Four-phase turn loop: **seed → clarify → enrich → plan**.

1. **seed**: Create FEATURE cards via `card init --type feature`. Fill Summary + Motivation. Identify cross-cutting concerns (CONV/DEC/MOD/FIND).
2. **clarify**: Ask user clarifying questions. Resolve Open Questions.
3. **enrich**: Explore library (`library suggest --for <id>`) and codebase. Fill Design + Constraints. Run `card evolve --stage designed`.
4. **plan**: Fill Implementation Plan steps. Each step needs Files, Approach, Edge Cases. Run `card evolve --stage planned`.

Follow `references/card-templates.md` for card body templates.

## Hard Rules

- Create cards via `card init --type feature`; then edit the `.md` file directly for body content.
- Use `card link`/`card unlink` for all link operations.
- Use `card evolve` for stage transitions — never hand-edit status in frontmatter.
- Run `flowforge validate all` after any `.md` file changes.
- Never create >5 draft FEATURE cards in a single round.
- Never skip stages: draft → designed → planned must be sequential via `card evolve`.
- Before enriching Design, always run `library suggest --for <feature-id>`.
- Each Key Decision must include a "why" (≥1 sentence), not just "what".
- Implementation Plan steps must include Files, Approach, and Edge Cases — no cross-card references.
- All Open Questions must be cleared before `card evolve --stage designed`.
- Run `flowforge proposal inspect <id>` and address all health issues before implementation.

## Output

Report cards created/updated, relations added, unresolved gaps, and one next step.
