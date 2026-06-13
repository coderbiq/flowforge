# Workflow Rules

Use one primary mode per turn.

## Turn Loop

1. Inspect project, proposal, and context.
2. Pick one mode: index, clarify, analyze, discover library, design, or split tasks.
3. Write or update cards through CLI only.
4. Record a log when the turn changes proposal state.
5. Report changed cards, relations, gaps, and one next step.

## Mode Selection

| Mode | Use When | Main Commands |
|------|----------|---------------|
| index | new demand is not in STR, or STR has too many entries | `structure add/remove/list` |
| clarify | requirement lacks acceptance, scope, or open questions | `card read/update`, `log create` |
| analyze | uncertainty blocks design or task split | `task create --type a --status ready/not_ready`, `task ready --type a` |
| discover library | design needs conventions, modules, decisions, or findings | `library suggest`, `card search`, `card read` |
| design | a stable conclusion exists but no design card exists | `card create --type design` |
| split tasks | requirement and design are stable enough for implementation | `task create --type i --status ready/not_ready` |

## Card Decisions

- New user-visible behavior or constraint: create/update requirement and index it from STR.
- Unknown boundary or impact: create analysis task, usually `--status not_ready` until inputs are clear.
- Reusable fact or risk: create finding.
- Stable design conclusion: create design card.
- Implementation work: create task only when requirement, design, constraints, acceptance, and out-of-scope are clear.

## Ready Rules

Ready implementation tasks require linked requirement, linked design, acceptance, deliverables, out-of-scope, and confirmed library/context constraints. Otherwise create them with `--status not_ready` or keep analysis tasks open.

Ready analysis tasks require goal, inputs, investigation plan, expected outputs, and done-when.

## Library Rules

Do not read library files. First use `library suggest --for <card-id>` or `card search <query> --scope library`; read only confirmed candidates with `card read --summary` or `--section`.

## Output Rules

End with:

- cards created or updated
- relations added
- unresolved gaps
- one recommended next step
