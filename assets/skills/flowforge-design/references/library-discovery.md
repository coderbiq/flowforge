# Library and Knowledge Source Discovery

FlowForge design work must discover context through CLI only. Do not grep markdown, do not read `02-library/` directly, and do not batch-load whole library trees.

## Three-layer exploration model

Explore in priority order:

1. **FlowForge Library** (highest priority) — curated project conventions, decisions, findings via `library suggest` / `card search --scope library`.
2. **External Knowledge Sources** (medium priority) — configured `knowledge_sources` in `.flowforge/config.yaml`. See "External Knowledge Source Discovery" below.
3. **Project Source Code** (informational only) — read to understand current state; source code is fact, not normative constraint.

## Library Discovery

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

If library returns no useful result and external sources also produce no match, do not expand the search scope without a concrete new query. Record the gap in a log and leave the open question visible.

## External Knowledge Source Discovery

FlowForge supports external knowledge sources configured in `.flowforge/config.yaml` under `knowledge_sources`. Each source has:

```yaml
knowledge_sources:
  - name: team-architecture
    path: /docs/architecture/
    type: file                # access mechanism (file | jira | confluence | url); MVP supports file only
    category: team_knowledge  # content nature (official_docs | team_knowledge | community | experimental | legacy)
    trust: high
    description: Team architecture decision records
```

### Discovery workflow

1. Read the configuration to identify registered external sources (`flowforge config get knowledge_sources` or read `.flowforge/config.yaml`).
2. For `type: file` sources, use file tools (Glob, Grep, Read) to search the configured path for relevant content.
3. Extract and evaluate relevant content from external sources.
4. Decide whether to ingest into library (Strategy A) or embed in a card (Strategy B) — see below.

### Credibility

External source citations must include the source name and trust level. Trust levels: `high`, `medium`, `low`, `unknown`.

## External Content Integration: Strategy A vs B

Decide based on **reusability** (primary) and **length** (auxiliary):

| Condition | Strategy | Action |
|-----------|----------|--------|
| High reusability (will be used by multiple designs/requirements) | A: Ingest into library | Use `library import` to create a library card, then reference by card ID in design/analysis cards. |
| Low reusability, ≤500 chars | B: Embed in card | Direct quotation with source attribution. |
| Low reusability, 500–2000 chars | B: Embed in card | Summary plus key excerpt with source attribution. |
| Low reusability, >2000 chars | Suggest A | Strongly prefer library import; if still embedding, summarize only with source reference. |

### Strategy B format

Embed external content in any card body using the following format:

```markdown
## 外部参考

> **来源**: [source-name] path/to/file
> **可信度**: high
>
> Quoted or summarized content...
```

Rules:
- Source format: `[source-name] relative/path`
- Trust level is mandatory.
- For long content, only embed a summary; the full content should be ingested into the library.
- Do not embed content without source attribution.

### Strategy A: library import

Use existing `library import` command. The `Card.Source` field records the external source provenance.

```bash
flowforge library import \
  --title "External knowledge unit" \
  --type finding \
  --source "[source-name] path/to/file" \
  --source-card <related-card-id>
```

After import, reference the library card by ID in design/analysis cards via normal `card link` relations.

No new library card type is needed. Existing types (`finding`, `convention`, `decision`, `module`, `design`) cover external knowledge units. Add a new card type only when a clear schema-level distinction is required.
