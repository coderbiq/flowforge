# Journal 写作指南

## 位置

探索日志已合并到 proposal 的 `notes.md`。不再使用独立的 journal 目录。

## 替代方式

探索过程中的记录统一写入 proposal 的 `notes.md`：
- 按日期分段（`## YYYY-MM-DD`）
- 记录 key findings、decisions、下一步计划
- 使用 `note_kind: progress` 格式

## Frontmatter

notes.md 的 frontmatter：
```yaml
---
doc_type: notes
title: <proposal 标题> 实施日志
status: active
note_kind: progress
---
```
