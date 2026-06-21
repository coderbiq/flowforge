---
id: LOG-CR26062101-djes35cgu87t
title: 'Library discover: no conventions found for feedback skill tasks'
type: log
status: active
importance: should
tags:
    - discovery
    - feedback
    - library
    - v2
links:
    - target: PROP-CR26062101
      relation: belongs_to
    - target: TASK-CR26062101-djer9jcjxe2n
      relation: records
    - target: TASK-CR26062101-djer9lgnm8cp
      relation: records
    - target: TASK-CR26062101-djer9d3vyo2m
      relation: records
    - target: DES-CR26062101-djer91do74qp
      relation: references
created: 2026-06-21T13:57:17.511710669Z
updated: 2026-06-21T13:57:17.511723369Z
source: CR26062101
domain: flowforge
---

# Library discover: no conventions found for feedback skill tasks

## Kind
library-discovery

## Event
对三个 ready 任务执行 `library suggest --for <task> --types convention,module`，均返回空结果。

## Context
- TASK-CR26062101-djer9jcjxe2n: 编写 feedback 分类规则 reference
- TASK-CR26062101-djer9lgnm8cp: 编写 feedback 工作流 rules reference
- TASK-CR26062101-djer9d3vyo2m: 编写 flowforge-feedback SKILL.md 主文件

Library 当前无任何 facets、无 convention/module 卡片。
健康检查因此提示 "ready implementation task has no linked convention constraints"。

## Result
三个任务均无法从 library 获得约束。
设计卡中的 "Constraints" 节目前依赖的是文档推理而非 library 证据。
待任务实施后，产出的 SKILL/reference 文件应通过 `library import` 或 `library promote` 回流为 library 内容。

## Links

### Outgoing

- [PROP-CR26062101]() [proposal] - flowforge-feedback-skill-v2

