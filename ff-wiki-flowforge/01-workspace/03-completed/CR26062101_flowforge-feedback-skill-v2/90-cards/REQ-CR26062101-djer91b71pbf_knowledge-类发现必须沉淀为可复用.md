---
id: REQ-CR26062101-djer91b71pbf
title: knowledge 类发现必须沉淀为可复用 library 内容
type: requirement
status: active
importance: should
tags:
    - feedback
    - knowledge-ingestion
    - library
links:
    - target: PROP-CR26062101
      relation: belongs_to
created: 2026-06-21T13:17:57.802747963Z
updated: 2026-06-21T13:17:57.802800263Z
source: CR26062101
domain: flowforge
---

# knowledge 类发现必须沉淀为可复用 library 内容

## Summary
knowledge 类反馈是增量知识。feedback SKILL 必须在分类后引导 library import 或 library promote，避免知识只在 proposal 内部循环。

## Source
docs/business-layer-outline.md §7.3

## Acceptance
- knowledge → 构造成 library candidate → library import 或 promote
- library 卡片必须保持 `--source-card` 来源链接
- library 卡片类型为 convention / decision / module / finding / design 之一

## Scope
- 仅导入可独立复用的知识单元，不批量搬运 proposal 内容

## Links

### Outgoing

- [PROP-CR26062101](../../../../03-proposal/CR26062101_flowforge-feedback-skill-v2.md) [proposal] - flowforge-feedback-skill-v2

### Incoming

#### requires
- [TASK-CR26062101-djer9d3vyo2m](TASK-CR26062101-djer9d3vyo2m_编写-flowforge-feedback-skillmd-主文件.md) [task] - 编写 flowforge-feedback SKILL.md 主文件
- [TASK-CR26062101-djer9lgnm8cp](TASK-CR26062101-djer9lgnm8cp_编写-feedback-工作流-rules-reference.md) [task] - 编写 feedback 工作流 rules reference

## Open Questions
None

