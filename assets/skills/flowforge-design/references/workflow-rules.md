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
| index | new demand is not in STR, or STR has too many entries | `structure add/remove/list` |
| clarify | requirement lacks acceptance, scope, or open questions | `card read/update`, `log create` |
| analyze | uncertainty blocks design or task split | `task create --type a --status ready/not_ready`, `task ready --type a` |
| discover library | design needs conventions, modules, decisions, or findings | `library suggest`, `card search`, `card read` |
| design | a stable conclusion exists but no design card exists | `card create --type design` |
| split tasks | requirement and design are stable enough for implementation | `task create --type i --status ready/not_ready` |
| refresh navigation | requirement/design relationships changed | `card refresh <REQ/DES id>` |

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

## Ready Rules

Ready implementation tasks require linked requirement, linked design, acceptance, deliverables, out-of-scope, and confirmed library/context constraints. Otherwise create them with `--status not_ready` or keep analysis tasks open.

Ready analysis tasks require goal, inputs, investigation plan, expected outputs, and done-when.

## Library Rules

Do not read library files. First use `library suggest --for <card-id>` or `card search <query> --scope library`; read only confirmed candidates with `card read --summary` or `--section`.

Use `library import` only after source material has already been decomposed into a structured candidate. Use `library promote <card-id>` to copy a stable proposal card into library. Do not write library card files directly.

## Output Rules

End with:

- cards created or updated
- relations added
- unresolved gaps
- one recommended next step
