---
doc_type: note
title: Workflow Rules - Explore
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

# Explore Rules

## Ownership summary

- Primary module: none
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: this exploration rule file

## Exploration defaults

- `index.md` is the reading surface.
- `journal/` preserves chronology.
- `findings/` contains atomic statements worth reusing.
- `decisions/` contains candidate decisions with status.
- Declare `ownership`, `expected_size_class`, and `reusable_rules` once the
  question is scoped.
- Mirror the ownership graph in human-readable form so the reader can see the
  owning modules, system or architecture targets, and reusable conventions
  without parsing metadata only.
- Use the archived knowledge base as the first source of truth for research,
  including conventions.
- Treat proposals and explorations as delta records against the existing final
  corpus.
- Do not mix implementation logs into exploration files.

## Analysis expectations

- Exploration should read the input package and infer the structure and
  pressure points before it writes a skeleton.
- The generated skeleton should note which input-package files or sections were
  used as evidence.
- If project rules define special priorities or review order, apply them before
  drafting the skeleton.

