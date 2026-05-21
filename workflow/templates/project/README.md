# Minimal Project Template

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
