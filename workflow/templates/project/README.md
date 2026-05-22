---
doc_type: note
title: Minimal Project Template
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

# Minimal Project Template

## Ownership summary

- Primary module: none
- System / architecture targets: none
- Convention targets: none
- Canonical reading path: this project scaffold

This directory contains the smallest recommended project skeleton for adopting `FlowForge`.

Suggested target layout:

```text
your-project/
├── AGENTS.md
├── .flowforge/
│   ├── config.json
│   ├── workflow/
│   ├── agents/
│   ├── scripts/
│   └── adapters/
├── .claude/
├── .codex/
├── docs/
│   ├── explorations/
│   ├── proposals/
│   ├── modules/
│   ├── architecture/
│   └── decisions/
└── .flowforge/state/
```

Use this template when bootstrapping a new repository or normalizing an existing one.

This layout is also the expected Codex project shape, because Codex reads project instructions from `AGENTS.md` and operates directly on the installed workflow scripts.
