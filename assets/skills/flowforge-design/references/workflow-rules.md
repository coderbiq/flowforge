# Workflow Rules

Use one primary mode per turn.

## Turn Loop

1. Inspect project, proposal, and context.
2. Review `proposal inspect` Health Issues.
3. Pick one mode: index, clarify, analyze, discover library, design, split tasks, or refresh navigation.
4. Write or update cards through CLI only.
5. Record a log when the turn changes proposal state.
6. Run `flowforge validate all` after creating or changing proposal structure.
7. Report changed cards, relations, validation result, gaps, and one next step.

## Mode Selection

| Mode | Use When | Main Commands |
|------|----------|---------------|
| index | new demand not in STR, or STR has too many entries. First analyze the demand: identify core goal, scope boundary, implicit assumptions, and risks; then decompose into themes and write STR entries | `structure add/remove/list` |
| clarify | requirement lacks acceptance, scope, or open questions | `card read/update`, `log create` |
| analyze | uncertainty blocks design or task split | `task create --type a --status ready/not_ready`, `task ready --type a` |
| discover library | design needs conventions, modules, decisions, or findings | `library suggest`, `card search`, `card read` |
| design | a stable conclusion exists but no design card exists | `card create --type design --status draft` |
| split tasks | requirement and design are stable enough for implementation | `task create --type i --status ready/not_ready` |
| refresh navigation | requirement/design relationships changed | `card refresh <REQ/DES id>` |

## Mode Gating Rules

The following rules are mandatory and block further work:

| Condition | Required Action | Rationale |
|-----------|----------------|-----------|
| Proposal has >= 3 active REQ cards but 0 DESIGN cards | Enter **design** mode next turn. Create at least 1 DESIGN card synthesizing design intent, architecture decisions, and key constraints. | Prevents requirement fragmentation without synthesis (design_gap). |
| An STR card has only `## Purpose` + `## Entries` (no `## Synthesis` / `## Key Decisions`) | Enter **refresh navigation** mode. Write synthesis content into the STR card via `card update`. | STR cards must be synthesis, not just directories. |
| A REQ card has < 5 lines of effective content (excluding frontmatter, section headings, auto-generated nav) | Do not create this card independently. Merge into parent or related card. | Prevents high structural overhead (template burden). |
| After index mode creates new REQ cards, the next turn must be design or clarify (not another index pass) | Report as a required next step; index mode may not repeat without design. | Prevents "index -> done" syndrome (progressive refinement skipped). |

## Batch Card Creation

When creating multiple cards at once (e.g., seeding a proposal with requirements, designs, and tasks), generate a YAML manifest and use `card batch --manifest "cards:\n  - ..."`. Use `ref` for cross-references within the same batch. Use `-o json` to capture created card IDs.

## Link Invariants

- Do not write wiki files or frontmatter manually.
- Do not create internal card links in card bodies by hand.
- Do not create `[[wikilink]]`.
- FlowForge CLI is the only component that formats and inserts internal card navigation links.
- Hand-written Markdown links are allowed only for external source references.
- Every non-root card needs at least one outbound frontmatter link.
- Proposal-scoped cards should belong to `ROOT-<proposal>`; the CLI adds this when creating cards in a proposal.
- Requirement indexes may only `indexes` requirement cards or child structure cards.
- Use `decomposes` for subtask -> parent task. Do not use `related` for parent/child, ownership, or indexing.
- Use `records` for logs, `constrains` for convention cards, `requires` for required inputs, `implements` for design implementation, and `satisfies` for requirement satisfaction.

## Card Decisions

- New user-visible behavior or constraint: create/update requirement and index it from STR.
- Unknown boundary or impact: create analysis task, usually `--status not_ready` until inputs are clear.
- Reusable fact or risk: create finding.
- Stable design conclusion: create design card.
- Implementation work: create task only when requirement, design, constraints, acceptance, and out-of-scope are clear.
- After adding analysis/design/task links for a requirement or implementation links for a design, run `card refresh <id>`.

## Card and Task Co-evolution

Analysis tasks manage complex analysis processes. Cards and tasks evolve concurrently (not sequentially):

- Creating an analysis task for a requirement does not block the requirement card from being updated.
- Write partial analysis findings to the requirement/design card immediately via `card update`, even while the analysis task remains in progress.
- The analysis task tracks remaining sub-questions; its `Done When` field governs completion.
- A card transitions from `draft` to `stable` when: all linked analysis tasks are `done` AND all Open Questions in the card body are resolved.

## Ready Rules

Ready implementation tasks require linked requirement, linked design, acceptance, deliverables, out-of-scope, and confirmed library/context constraints. Otherwise create them with `--status not_ready` or keep analysis tasks open.

Ready analysis tasks require goal, inputs, investigation plan, expected outputs, and done-when.

## Design Completion Rules

A design card is complete when all linked analysis tasks are `done` AND all open questions in the design card body are resolved.

A proposal is ready for implementation when:
- Every active REQ card is linked to at least one DESIGN card via `designs` or `satisfies`.
- Every STR card has non-placeholder `## Synthesis` content.

## Library Rules

Do not read library files. First use `library suggest --for <card-id>` or `card search <query> --scope library`; read only confirmed candidates with `card read --summary` or `--section`.

Use `library import` only after source material has already been decomposed into a structured candidate. Use `library promote <card-id>` to copy a stable proposal card into library. Do not write library card files directly.

## Exploration Depth Criteria

Design work must explore three sources in priority order before forming design conclusions:

1. **FlowForge Library** (highest priority): `library suggest`, `card search --scope library` — curated project conventions, decisions, findings.
2. **External Knowledge Sources** (medium priority): configured `knowledge_sources` in `.flowforge/config.yaml` — agent reads files from configured paths using file tools.
3. **Project Source Code** (informational only): read to understand current state; source code is fact, not normative constraint.

### When exploration is sufficient

Hard rules (must satisfy):
- No open question on a requirement card blocks the current design decision.
- All analysis task specified input sources have been checked.

Heuristic guidance (should satisfy):
- Library, external sources, and source code have all been checked; all returned no actionable new result.
- Two consecutive exploration rounds yielded no new information and no new search terms.

### Recording exploration

Record what was searched, what matched, and what was missed in a log card (`log create --kind progress`). Logs serve as audit trail, not as replacement for requirement/design/finding cards.

## Output Rules

End with:

- cards created or updated
- relations added
- unresolved gaps
- one recommended next step
