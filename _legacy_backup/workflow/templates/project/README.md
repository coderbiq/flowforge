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

这个目录提供了采用 `FlowForge` 时建议的最小项目骨架。

Suggested target layout:

```text
your-project/
├── AGENTS.md   # 可插入到项目本地文件中的 FlowForge 章节
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

在初始化新仓库，或者整理已有仓库时使用这个模板。

安装器会把 `docs/flowforge/_rules/` 作为项目最初可编辑的 workflow rules bundle。
当项目需要调整文档形状时，`docs/flowforge/_templates/` 仍然是 workspace-local template copies 的位置。
对于 model documents，默认模板是单文件 `model.md`，而不是拆分成 parts tree。
`docs/intake/` 是 pre-exploration input packages 的推荐入口。

这个布局也是 Codex 期望看到的项目形态，因为 Codex 会从 `AGENTS.md` 读取项目说明，并直接操作已安装的 workflow scripts。这里种子的 `AGENTS.md` 内容刻意只是一段 FlowForge 章节片段，不是对项目自有说明的完整替换。
