---
doc_type: note
title: Exploration Template
status: draft
workspace: default
module_scope: []
system_scope: []
convention_scope: []
ownership: []
information_class: exploration
topics: []
related_docs: []
archive_target: none
created: <ISO-8601 timestamp>
updated: <ISO-8601 timestamp>
---

# Exploration Template

Canonical structure for a new exploration:

```text
docs/explorations/<slug>/
├── index.md
├── journal/
│   └── 2026-05-20-entry.md
├── findings/
│   └── F-001-example.md
├── decisions/
│   └── D-001-candidate.md
└── artifacts/
```

Use `index.md` as the main reading surface.

## Ownership summary

- Primary module: <type:target or none>
- System / architecture targets: <type:target or none>
- Convention targets: <type:target or none>
- Canonical reading path: this exploration overview

When starting a new exploration, review the existing canonical corpus first:

- relevant module docs
- relevant architecture docs
- relevant ADRs
- earlier proposals or explorations in the same area
