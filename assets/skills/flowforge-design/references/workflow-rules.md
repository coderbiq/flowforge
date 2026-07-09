# Workflow Rules

Use one primary phase per turn.

## Turn Loop

1. Run `flowforge proposal inspect <id>` to check health.
2. Pick one phase: seed, clarify, enrich, or plan.
3. Write or update FEATURE cards — create via `card init`, edit body directly.
4. Use `card link`/`card unlink` for link operations.
5. Use `card evolve` for stage transitions.
6. Run `flowforge validate all` after changes.
7. Report changed cards, relations, validation result, gaps, and one next step.

## Phase Selection

| Phase | Use When | Main Commands |
|-------|----------|---------------|
| seed | new demand, no FEATURE cards exist | `card init --type feature`, direct edit for Summary/Motivation |
| clarify | Open Questions exist; user input needed | direct edit for Open Questions, `card log` |
| enrich | draft FEATUREs need Design/Constraints | `library suggest`, direct edit for Design/Constraints, `card evolve --stage designed` |
| plan | designed FEATUREs need Implementation Plan | direct edit for Implementation Plan, `card evolve --stage planned` |

## Stage Gates (CLI-enforced)

`card evolve` validates:
- **draft → designed**: Design.Key Decisions ≥1, Constraints ≥1, Open Questions cleared
- **designed → planned**: Implementation Plan steps with Files + Approach + Edge Cases, no cross-card refs
- **in_progress → done**: All steps done, Verification results present

## Link Invariants

- Do not hand-write `links` in frontmatter. Use `card link`/`card unlink`.
- FEATURE cards link to PROP via `belongs_to`.
- FEATURE cards link to other FEATUREs via `depends_on`.
- FEATURE cards link to CONV/DEC via `constrains`/`references`.
- Run `card related <id>` to view links.

## Card Decisions

- New user-visible capability: `card init --type feature`
- Reusable rule across features: `card init --type convention`
- Architecture decision affecting multiple features: `card init --type decision`
- Module knowledge: `card init --type module`
- Observation or risk: `card init --type finding`

## PROP Update Triggers

Update the PROP card's Feature Map and Architecture Overview when:
- First FEATURE created
- Any FEATURE reaches `designed` stage
- `card split` executed
- Any FEATURE reaches `done` stage
- Run `proposal inspect` to detect stale Feature Maps

## Output Rules

End with:
- cards created or updated
- relations added
- unresolved gaps
- one recommended next step
