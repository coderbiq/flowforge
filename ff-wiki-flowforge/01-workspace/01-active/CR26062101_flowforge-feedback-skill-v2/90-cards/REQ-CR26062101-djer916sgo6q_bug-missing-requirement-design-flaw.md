---
id: REQ-CR26062101-djer916sgo6q
title: bug / missing-requirement / design-flaw 必须生成追踪任务卡
type: requirement
status: active
importance: should
tags:
    - classification
    - feedback
    - task-creation
links:
    - target: PROP-CR26062101
      relation: belongs_to
created: 2026-06-21T13:17:57.536372194Z
updated: 2026-06-21T13:17:57.536437094Z
source: CR26062101
domain: flowforge
---

# bug / missing-requirement / design-flaw 必须生成追踪任务卡

## Summary
这三类发现属于可行动问题，不能只记录日志。每类必须立即生成一张带状态的任务卡，确保后续迭代能消费它。

## Source
docs/business-layer-outline.md §7.3

## Acceptance
- bug → 生成 `task` 卡，--status not_ready，关联来源发现卡
- missing-requirement → 生成 `requirement` 卡，--status draft，关联来源发现卡
- design-flaw → 生成 `requirement` 卡（设计变更需求），--status draft，关联来源发现卡
- 每张生成的任务卡必须通过 `card link` 与发现卡建立 `records` 或 `requires` 关系

## Scope
- 仅生成追踪入口卡，不执行修复
- 修复由后续 implement / design 阶段处理

## Links

### Outgoing

- [PROP-CR26062101](../../../../03-proposal/CR26062101_flowforge-feedback-skill-v2.md) [proposal] - flowforge-feedback-skill-v2

### Incoming

- [TASK-CR26062101-djer9d3vyo2m](TASK-CR26062101-djer9d3vyo2m_编写-flowforge-feedback-skillmd-主文件.md) [task] - 编写 flowforge-feedback SKILL.md 主文件

## Open Questions
None

