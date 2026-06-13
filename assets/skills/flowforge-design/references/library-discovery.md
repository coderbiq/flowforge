# Library Discovery

FlowForge design work must discover library context through CLI only.
Do not grep markdown, do not read `02-library/` directly, and do not batch-load whole library trees.

## Roles of the CLI commands

- `flowforge library suggest --for <card-id>`: primary recommendation pass for a requirement, analysis task, design, or implementation task.
- `flowforge card search <query> --scope library`: targeted keyword and type search when you need to narrow the candidate set.
- `flowforge card read <id> --summary`: quick validation of a candidate.
- `flowforge card read <id> --section <name>`: deep read for a confirmed candidate.

## Candidate discovery

Start from the focus card:
- title
- summary
- tags
- domain
- status
- direct relationships

Use the user request and the current proposal context to form the query.
Prefer candidates that match the same project, module, or domain.

## Candidate filtering

Inspect the candidate list before reading anything in full.
Prefer:
- convention, decision, module, and design cards that directly constrain the work
- active or accepted cards over deprecated or superseded cards
- cards that share the same proposal, module, or domain

Do not link every plausible candidate.
Only keep cards that are confirmed relevant to the current analysis or design decision.

## Deep read

Read only the smallest useful surface:
- summary first
- then one or two sections only if the card is still relevant

Use `--section` when you need a rule, constraint, or decision.
Do not expand the whole library just because the first candidate set is small.

## Link writing

Write links only after the candidate is confirmed relevant.

Suggested relationships:
- analysis task `references ->` module / decision / finding cards
- analysis task `constrains ->` convention cards
- design card `references ->` decision / finding cards
- design card `constrains ->` convention / module cards
- implementation task `constrains ->` convention cards

Unselected candidates may be mentioned in a log, but they should not be linked into the center card.

## No-hit handling

If the library does not return a useful candidate:
- do not invent a rule
- record the search result in a log
- leave the open question visible in the analysis task or design card
- if the new knowledge is reusable, create a finding card

If the gap still blocks implementation, keep the task `not_ready` or create a new analysis task.
