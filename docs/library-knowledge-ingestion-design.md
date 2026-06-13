# Library Knowledge Ingestion Design

> Status: draft
> Scope: FlowForge v2 library knowledge supply chain

## Problem

FlowForge needs more than task-to-library search. A project library must accept knowledge from two sources:

- existing knowledge bases such as project SKILLs, engineering guides, and legacy docs
- discoveries produced during proposal design, implementation, and feedback

The library also has to expose how knowledge is organized so an Agent can classify a task and compose the right convention, module, decision, design, and finding cards without hard-coding a taxonomy.

## Principles

- Cards are the source of truth. sqlite indexes are derived.
- Agents access library knowledge through CLI commands, not direct file tree reads.
- FlowForge defines the facet mechanism, not a fixed project taxonomy.
- Imports and promotions produce reviewable plans before changing active library knowledge.
- Stable constraints are linked from task/design cards; process evidence links back from log/finding cards.
- Library cards must remain human-readable, not frontmatter-only records.

## Library Card Roles

| Type | Role | Example |
|------|------|---------|
| `convention` | Rules that constrain future work | service pagination query rule |
| `module` | Existing system/module knowledge | customer module, auth module |
| `decision` | Accepted architectural or product decisions | shared error envelope |
| `design` | Reusable historical design | import pipeline design |
| `finding` | Reusable facts and caveats | legacy API compatibility issue |
| `structure` | Navigational indexes and maps | convention map, module map |

## Facet Tags

Facets are project-defined dimensions encoded as card tags.

Preferred form:

```yaml
tags:
  - layer:service
  - scenario:page-query
  - domain:customer
```

Compatible form:

```yaml
tags:
  - facet:layer:service
```

FlowForge treats both as the same facet `layer:service`.

Facet keys and values are discovered from existing library cards. A backend project might use `layer`, `scenario`, and `framework`; a data platform project might use `pipeline`, `storage`, and `quality`.

## Discovery Commands

### `flowforge library facets`

Summarizes the library's available facet vocabulary.

Output:

```markdown
## Library Facets

| Facet | Value | Cards |
|-------|-------|-------|
| layer | service | 8 |
| scenario | page-query | 5 |

## Common Combinations

| Facets | Cards |
|--------|-------|
| layer:service + scenario:page-query | 4 |
```

### `flowforge library classify --for <card-id>`

Classifies a requirement, design, or task against discovered library facets.

The command does not write tags or links. It reports extracted candidate facets and evidence so the Agent can decide whether to use them.

Output:

```markdown
## Library Classification

| Facet | Source | Evidence | LibraryCards |
|-------|--------|----------|--------------|
| layer:service | tag | layer:service | 8 |
| scenario:page-query | text | page-query | 5 |

## Suggested Commands

- flowforge library suggest --for TASK... --facet layer:service --facet scenario:page-query
```

### `flowforge library suggest --facet key:value`

Uses explicit facet filters in addition to keyword scoring.

Example:

```bash
flowforge library suggest \
  --for TASK-... \
  --types convention,module,decision \
  --relation constrains \
  --facet layer:service \
  --facet scenario:page-query
```

Facet matches must be exact. The Agent still validates each candidate with `card read --summary` or a targeted section read before linking.

## External Knowledge Import

External import is for existing SKILLs, engineering guides, or legacy docs.

Target workflow:

```text
scan source
  -> extract candidate knowledge units
  -> propose facets
  -> propose cards and structure indexes
  -> review plan
  -> apply as draft/active library cards
  -> rebuild indexes
```

Future command shape:

```bash
flowforge library import scan <path>
flowforge library import plan <scan-id>
flowforge library import apply <plan-id>
```

The plan should contain:

- proposed facet vocabulary
- proposed cards with type, title, summary, tags, importance, and source evidence
- proposed structure cards
- duplicate or merge candidates
- warnings for oversized or vague cards

Imports should not directly flood active library. The first implementation can create `status: draft` cards and require explicit promotion to `active` or `accepted`.

## Proposal Knowledge Promotion

Proposal work generates logs, findings, decisions, and design cards. Not all of them are reusable.

Target workflow:

```text
proposal log/finding/design
  -> reusable candidate
  -> promotion plan
  -> duplicate/merge check
  -> create or update library card
  -> link source evidence
```

Future command shape:

```bash
flowforge library promote --from FIND-...
flowforge library promote --from LOG-...
flowforge library promote --proposal CR...
flowforge library promote apply <plan-id>
```

Promotion actions:

| Action | Meaning |
|--------|---------|
| `create` | Create a new library card |
| `merge` | Add a section or evidence link to an existing card |
| `supersede` | Mark old knowledge as superseded and link replacement |
| `skip` | Keep proposal-local only |

## Link Ownership

Stable links:

```text
TASK -> CONV  constrains
TASK -> MOD   references
TASK -> DEC   references
DES  -> CONV  constrains
DES  -> MOD   references
STR  -> CONV  indexes
```

Evidence links:

```text
LOG  -> TASK  records
FIND -> TASK  records
CONV -> FIND  derived-from
DEC  -> LOG   derived-from
```

Tasks should not accumulate every process evidence link. Evidence cards link to tasks and are shown through backlinks.

## MVP Scope

Implemented first:

- `library facets`
- `library classify --for`
- `library suggest --facet`

Deferred:

- bulk import scan/plan/apply
- promotion plan/apply
- sqlite FTS/BM25-backed ranking
- embedding/vector retrieval

## Validation Scenario

1. Create convention cards with project-specific facets.
2. Create an implementation task with title/body/tags that imply those facets.
3. Run `library facets`.
4. Run `library classify --for <task>`.
5. Run `library suggest --for <task> --facet ...`.
6. Link confirmed convention cards to the task.
7. Run `context task --task <task>` and verify the linked conventions appear as stable context.
