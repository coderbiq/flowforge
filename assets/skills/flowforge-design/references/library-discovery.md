# Library Discovery

FlowForge design work must discover library context through CLI only.
Do not grep markdown, do not read `02-library/` directly, and do not batch-load whole library trees.

## Roles of the CLI commands

- `flowforge library facets`: discover the project's library facet vocabulary before guessing tags or dimensions.
- `flowforge library classify --for <card-id>`: classify a requirement, design, or task against discovered facets.
- `flowforge library suggest --for <card-id>`: primary recommendation pass for a requirement, analysis task, design, or implementation task.
- `flowforge library suggest --for <card-id> --facet key:value`: recommendation pass constrained by confirmed project facets.
- `flowforge card search <query> --scope library`: targeted keyword and type search when you need to narrow the candidate set.
- `flowforge card read <id> --summary`: quick validation of a candidate.
- `flowforge card read <id> --section <name>`: deep read for a confirmed candidate.
- `flowforge library import`: write an already-structured source-material candidate into library.
- `flowforge library promote <card-id>`: copy a stable proposal card into library while keeping source traceability.

## Candidate discovery

Start from the focus card:
- title
- summary
- tags
- domain
- status
- direct relationships

Run `library facets` first if you do not already know the project's facet vocabulary.
Run `library classify --for <card-id>` before facet-constrained suggestions.
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

## Knowledge ingestion

Use `library import` only after source material has already been decomposed into an atomic candidate with type, title, body, tags, and source evidence.

Use `library promote <card-id>` when proposal work produces reusable knowledge. Promote stable findings, designs, decisions, or conventions; keep process logs as trace evidence unless they contain a reusable conclusion that has first been captured as a finding or design.

Do not write library files directly. Do not import title-only candidates. Every imported or promoted library card must keep an outbound source link, either through `--source-card` or `--links`.

## No-hit handling

If the library does not return a useful candidate:
- do not invent a rule
- record the search result in a log
- leave the open question visible in the analysis task or design card
- if the new knowledge is reusable, create a finding card

If the gap still blocks implementation, keep the task `not_ready` or create a new analysis task.
