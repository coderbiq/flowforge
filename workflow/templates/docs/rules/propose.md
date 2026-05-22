---
doc_type: note
title: Workflow Rules - Propose
status: draft
workspace: default
module_scope: []
system_scope: []
convention_scope: []
ownership: []
information_class: note
topics: []
related_docs: []
archive_target: none
created: <ISO-8601 timestamp>
updated: <ISO-8601 timestamp>
---

# Propose Rules

## Ownership summary

- Primary module: none
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: this proposal rule file

## Proposal defaults

- `meta.yaml` is the proposal bundle manifest.
- Each Markdown artifact carries its own YAML frontmatter for doc-local
  routing.
- `proposal.md` answers why and what, and surfaces `size_class`, `ownership`,
  and any promoted `reusable_rules`.
- Human-readable proposal docs must summarize the ownership graph explicitly.
- Proposals should begin by reviewing the canonical corpus and then describe
  the delta from that corpus.
- Keep `task-map.md` authoritative for task decomposition.
- `task-map.md` must follow deliverable-first splitting, milestone boundaries,
  and explicit checkpoint rules.

## Design surface defaults

- `small`: single-file `design.md`.
- `medium`: single-file `design.md` by default, with `design/` used when the
  work spans multiple concerns.
- `large`: `design/` is mandatory and must include the canonical design
  subdocs required by the workflow core.
- When the proposal introduces multiple core business models, keep the model
  docs readable and explicit rather than compressing them into prose only.

## Baseline and merge expectations

- Record the canonical corpus reviewed for the proposal.
- When changing existing final docs, describe what is updated in place, what is
  appended, and what is superseded.
- If the proposal redistributes knowledge across docs, update the linked docs
  in one archive pass so the reader path stays coherent.

