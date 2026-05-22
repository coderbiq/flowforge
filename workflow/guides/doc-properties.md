# Document Properties

FlowForge attaches a property block (YAML frontmatter) to every Markdown
artifact so each document is independently identifiable, indexable, and
routable. The frontmatter is intentionally minimal: it answers
"what is this document, where does it belong, and how does it flow",
not "what does the document say". All narrative and design content
stays in the Markdown body.

## Design principles

- Each `.md` artifact carries its own property block.
- `meta.yaml` only describes the proposal container (lifecycle, links,
  archive targets, task backend). It never replaces per-doc properties.
- Property blocks are standard YAML frontmatter, compatible with
  Obsidian and any other tool that indexes Markdown by frontmatter.
- Fields are small and stable. Narrative content must stay in the body.
- Frontmatter must not contain long paragraphs, design rationale,
  or any prose that needs more than a short phrase.

## Minimal example

Use the same shape for every document, then add only the
type-specific fields that apply:

```yaml
---
doc_type: exploration
title: Data Service Config Management
status: draft
workspace: default
module_scope:
  - docs/modules/data-service
system_scope: []
convention_scope: []
ownership:
  - type: module
    target: docs/modules/data-service
    role: primary
information_class: exploration
topics:
  - data-service
  - config-management
related_docs: []
archive_target: default:modules/data-service/README.md
created: 2026-05-22T00:00:00Z
updated: 2026-05-22T00:00:00Z
exploration_slug: data-service-config-management
question: How should config data be modeled and validated?
expected_size_class: large
reusable_rules:
  - config data should be validated before persistence
---
```

The body below the frontmatter should explain the actual problem,
evidence, and design content in prose.

## Universal fields

Every FlowForge document must declare these fields:

- `doc_type`: which artifact this is (see types below).
- `title`: human-readable title.
- `status`: lifecycle state of this document.
- `workspace`: document workspace name from `flowforge.config.yaml`.
- `module_scope`: list of `docs/modules/<name>` paths the document is
  scoped to. Use an empty list when the doc is not module-scoped.
- `system_scope`: list of `docs/architecture/<topic>` paths the
  document affects. Empty when not architectural.
- `convention_scope`: list of `docs/conventions/<topic>` paths the
  document promotes or governs. Empty when not convention-related.
- `ownership`: the canonical ownership graph for this document. Use
  the same `type`, `target`, and `role` shape as proposal ownership.
- `information_class`: one of `exploration`, `proposal`, `design`,
  `model`, `finding`, `decision`, `journal`, `note`, `task-map`,
  `convention`, `module`, `architecture`, `adr`.
- `topics`: free-form tag list for cross-cutting topics. Used for
  Obsidian search and graph view. Keep entries short.
- `related_docs`: list of `workspace:ref` references to other
  FlowForge documents that share context with this one.
- `archive_target`: where this document's knowledge is expected to
  land after archive. Use `workspace:ref` form, or `none` for
  transient docs such as journal entries.
- `created`: ISO-8601 timestamp.
- `updated`: ISO-8601 timestamp.

## Routing cheat sheet

When deciding where a document belongs, read the fields in this order:

- `doc_type` says what kind of document this is.
- `information_class` says which workflow family it belongs to.
- `ownership` says what final knowledge corpus it contributes to.
- `archive_target` says where this document's knowledge should land.
- `module_scope`, `system_scope`, and `convention_scope` say what
  areas the document is about, even before archive.

In practice:

- `module` ownership usually archives to `docs/modules/<module>/`.
- `system` ownership usually archives to `docs/architecture/<topic>.md`.
- `cross-module` ownership usually archives to architecture plus
  supporting module history docs.
- `convention` ownership usually archives to `docs/conventions/<topic>.md`.
- `none` is valid only for transient docs such as journals or scratch notes.

## Type-specific extensions

Each `doc_type` adds a small number of routing fields. None of these
fields should be used to store body content.

### exploration

- `exploration_slug`: directory slug of the exploration.
- `question`: short question the exploration tries to answer.
- `reusable_rules`: candidate convention-level rules surfaced during
  the exploration. Each entry should stay short and may later be
  promoted into `docs/conventions/`.
- `expected_size_class`: predicted size class for the resulting
  proposal (`small | medium | large`).

### proposal

- `proposal_id`: CRYYMMDDNN id.
- `size_class`: `small | medium | large`.
- `ownership_primary`: `type:target` of the primary ownership entry.
- `design_layout`: `single | split`.

### design

- `design_section`: section name, e.g. `architecture`, `lifecycle`,
  `flow`, `api`, `constraints`, `tradeoffs`, `model-overview`,
  `entry`.
- `proposal_id`: CRYYMMDDNN id.
- `canonical_entry_point`: path to the canonical doc that remains
  the reader entry point after archive.

### model

- `proposal_id`: CRYYMMDDNN id.
- `model_name`: model identifier.
- `model_role`: `core | lifecycle | view-facing | shared`.
- `data_scope`: `single-record | master-table | event | derived`.
- `model_status_in_proposal`: `new | modified | retained`.

### finding

- `exploration_slug`: parent exploration slug.
- `finding_id`: `F-NNN` id.
- `evidence_sources`: list of repo paths or external references that
  support the finding. Short strings only.

### decision

- `exploration_slug` or `proposal_id`: parent context id.
- `decision_id`: `D-NNN` or `ADR-NNN` id.
- `decision_status`: `candidate | accepted | rejected | superseded`.

### journal

- `exploration_slug` or `proposal_id`: parent context id.
- `journal_date`: ISO date.

### note

- `proposal_id`: parent proposal id.
- `note_kind`: `progress | follow-up | decision-log`.

### task-map

- `proposal_id`: parent proposal id.
- `task_backend`: `beads | github | linear | none`.

### convention

- `convention_status`: `active | superseded | deprecated`.
- `enforcement`: `must | should | may`.
- `applies_to`: short list of artifact / layer names the rule covers.
- `origin_proposal`: CRYYMMDDNN id.

### module

- `module_name`: module identifier.
- `module_status`: `active | deprecated`.

### architecture

- `architecture_topic`: short topic name.
- `architecture_status`: `active | deprecated`.

### adr

- `adr_id`: `ADR-NNN`.
- `adr_status`: `proposed | accepted | superseded | deprecated`.

## What stays in the body

The following kinds of content must remain in the document body and
must never be encoded as frontmatter values:

- problem statement, context, motivation
- design rationale and alternatives
- model data structure tables and constraints
- rule explanations and counter-examples
- decision reasoning
- implementation history and follow-ups

If a piece of information needs explanation, a citation, or more than
a short phrase, it belongs in the body, not in the property block.

## Relationship to `meta.yaml`

- `meta.yaml` describes the proposal bundle: lifecycle status,
  ownership graph, archive targets, source explorations, canonical
  corpus, links between proposal files, and task backend state.
- Document frontmatter describes the document itself.
- The two layers must stay consistent but do not duplicate each
  other. `meta.yaml` is the proposal-level contract; frontmatter is
  the document-level contract.
- Validators check both layers and surface mismatches.

## Obsidian compatibility

- Frontmatter must be standard YAML enclosed by `---` lines.
- Property keys use snake_case so they map cleanly to Obsidian
  property names.
- List-typed properties are written as YAML sequences so Obsidian
  can render them as multi-value properties.
- `topics`, `module_scope`, `system_scope`, `convention_scope`, and
  `related_docs` are the primary fields used to build the Obsidian
  graph view and dataview-style queries.
