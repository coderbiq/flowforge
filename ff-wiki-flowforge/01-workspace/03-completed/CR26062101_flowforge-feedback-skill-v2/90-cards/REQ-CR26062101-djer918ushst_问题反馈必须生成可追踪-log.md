---
id: REQ-CR26062101-djer918ushst
title: 问题反馈必须生成可追踪 log 卡和任务卡
type: requirement
status: active
importance: should
tags:
    - feedback
    - log-card
    - traceability
links:
    - target: PROP-CR26062101
      relation: belongs_to
created: 2026-06-21T13:17:57.661230473Z
updated: 2026-06-21T13:17:57.661244373Z
source: CR26062101
domain: flowforge
---

# 问题反馈必须生成可追踪 log 卡和任务卡

## Summary
feedback SKILL 不能只把问题写进日志；必须同时生成可追踪的结构化卡片，使后续任务能通过 `card search` 或 `context proposal` 定位到这些发现。

## Source
docs/business-layer-outline.md §7.3

## Acceptance
- 每个有意义的发现生成一张 `log` 卡，kind=feedback，关联来源
- log 卡通过 `records` 关系指向来源发现卡或任务卡
- 后续任务可通过 `Backlink Evidence` 区域找到相关 log

## Scope
- log 卡作为过程证据保留，不作为替代 requirement / design 的摘要

## Links

### Outgoing

- [PROP-CR26062101](../../../../03-proposal/CR26062101_flowforge-feedback-skill-v2.md) [proposal] - flowforge-feedback-skill-v2

### Incoming

- [TASK-CR26062101-djer9lgnm8cp](TASK-CR26062101-djer9lgnm8cp_编写-feedback-工作流-rules-reference.md) [task] - 编写 feedback 工作流 rules reference
#### requires
- [TASK-CR26062101-djer9d3vyo2m](TASK-CR26062101-djer9d3vyo2m_编写-flowforge-feedback-skillmd-主文件.md) [task] - 编写 flowforge-feedback SKILL.md 主文件
- [TASK-CR26062101-djer9lgnm8cp](TASK-CR26062101-djer9lgnm8cp_编写-feedback-工作流-rules-reference.md) [task] - 编写 feedback 工作流 rules reference

## Open Questions
None

