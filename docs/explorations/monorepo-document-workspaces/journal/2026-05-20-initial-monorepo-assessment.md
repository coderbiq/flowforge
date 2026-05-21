# 过程记录

- Timestamp: 2026-05-20T00:00:00Z
- Actor: Codex

## 本次变化

针对“根文档目录和子项目文档目录同时存在”的 monorepo 场景审视了当前 `tg-workflow` 模型。确认当前实现建立在单一 docs root 之上，而 monorepo 支持需要更强的抽象层。

## 证据

- `workflow/guides/configuration.md`
- `workflow/schema/proposal.schema.yaml`
- `scripts/lib/flowforge.js`
- `docs/PROPOSAL-WORKFLOW.md`

## 新问题

- 文档工作区应该如何在 config 和 metadata 中表达？
- 生命周期中的哪些操作必须显式带上 workspace？
