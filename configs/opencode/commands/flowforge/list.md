---
description: List proposals grouped by canonical lifecycle status.
allowed-tools: Skill(flowforge)
---

Use the `FlowForge` skill to list proposals and explorations.

## 输出格式

按状态分组，默认同时展示 proposal 和 exploration：
- 进行中的工作
- 已完成的工作
- 草稿

## 执行方式

- 默认运行 `.flowforge/scripts/flowforge-list-proposals.js`
- 可选：`--kind proposals|explorations|all`

Arguments: $ARGUMENTS
