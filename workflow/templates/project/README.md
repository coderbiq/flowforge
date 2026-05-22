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
├── docs/
│   ├── intake/
│   ├── flowforge/
│   │   ├── _rules/
│   │   └── _templates/
│   ├── explorations/
│   ├── proposals/
│   ├── modules/
│   ├── architecture/
│   └── decisions/
├── .flowforge/
│   ├── config.json
│   ├── workflow/
│   ├── agents/
│   ├── scripts/
│   └── adapters/
├── .claude/
├── .codex/
└── .flowforge/state/
```

Use this template when bootstrapping a new repository or normalizing an existing one.

The installer seeds `docs/flowforge/_rules/` as the initial project-editable
workflow rules bundle. `docs/flowforge/_templates/` remains the place for
workspace-local template copies when a project needs to tailor document shapes.
For model documents, the default template is the single-file `model.md`
document rather than a split parts tree.
`docs/intake/` is the recommended entry point for pre-exploration input
packages.

This layout is also the expected Codex project shape, because Codex reads project instructions from `AGENTS.md` and operates directly on the installed workflow scripts.
