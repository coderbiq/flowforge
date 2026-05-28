---
doc_type: "finding"
title: "F-004 平台 commands 只是脚本 surface 的适配层"
status: "validated"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "architecture/installed-flowforge-upgrade-policy.md"
    role: "primary"
  - type: "module"
    target: "modules/workflow-core"
    role: "secondary"
information_class: "exploration"
topics: []
related_docs:
  - "default:explorations/installed-flowforge-safe-upgrade/index.md"
archive_target: "default:architecture/installed-flowforge-upgrade-policy.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
exploration_slug: "installed-flowforge-safe-upgrade"
finding_id: "F-004-platform-commands-track-the-script-surface"
evidence_sources: []
---

# F-004 平台 commands 只是脚本 surface 的适配层

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/installed-flowforge-upgrade-policy.md
- Convention targets: none
- Canonical reading path: installed-flowforge-safe-upgrade/findings/F-004-platform-commands-track-the-script-surface.md

## Statement

Claude Code 和 OpenCode 的 `flowforge` commands 不是独立业务实现，而是围绕 `.flowforge/scripts/` 的平台包装层。它们的职责是把平台命令映射到同一套 FlowForge 脚本和工作流语义。

## Why it matters

如果脚本升级了，但平台 command 没有同步更新，用户在平台入口看到的行为、参数和说明会与实际工具版本脱节。

## References

- [Adapter contract](../../../workflow/guides/adapter-contract.md)
- `configs/claude/commands/flowforge/`
- `configs/opencode/commands/flowforge/`
