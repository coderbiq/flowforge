# Notes 写作指南

## 位置

`workspace/proposals/<CR-id>/notes.md`

## 结构（单文件）

按日期追加记录。同一天多次更新追加到同一段落下。

## 每条记录应包含

- **时间戳**：ISO-8601
- **状态**：in_progress / done / blocked（对应 `rules.implement.task_states`）
- **摘要**：完成了什么、遇到了什么问题、做了什么决策

不需要长篇——每条约 2-3 行。

## 示例

```
## 2026-05-28

16:30 | done | 完成 auth middleware 的 token 校验。测试通过。

15:00 | blocked | JWT 库版本不兼容，暂时回退到 v8。需在下一个提案中升级。
```

## Frontmatter

```yaml
---
doc_type: notes
title: <提案标题> 实施日志
status: active
note_kind: progress
---
```
