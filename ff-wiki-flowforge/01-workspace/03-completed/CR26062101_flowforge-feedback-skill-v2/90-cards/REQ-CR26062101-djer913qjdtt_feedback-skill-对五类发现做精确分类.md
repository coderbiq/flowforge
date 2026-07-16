---
id: REQ-CR26062101-djer913qjdtt
title: feedback SKILL 对五类发现做精确分类
type: requirement
status: active
importance: should
tags:
    - classification
    - feedback
    - v2
links:
    - target: PROP-CR26062101
      relation: belongs_to
created: 2026-06-21T13:17:57.351750277Z
updated: 2026-06-21T13:17:57.351765577Z
source: CR26062101
domain: flowforge
---

# feedback SKILL 对五类发现做精确分类

## Summary
当任务执行或设计审查中发现偏差时，feedback SKILL 必须将每个发现精确归入以下五类之一，不能混用。

## Source
docs/business-layer-outline.md §7.3

## Acceptance
- 每个反馈条目被分类为 bug / finding / knowledge / missing-requirement / design-flaw 之一
- 分类结果写入卡片 frontmatter `type` 字段
- 分类逻辑可在 card body 中复现

## Scope
- 覆盖实施阶段发现的所有反馈入口（测试失败、行为偏差、认知更新、需求缺口、设计缺陷）
- 不处理库卡片归并（由 curate SKILL 负责）

## Links

### Outgoing

- [PROP-CR26062101](../../../../03-proposal/CR26062101_flowforge-feedback-skill-v2.md) [proposal] - flowforge-feedback-skill-v2

### Incoming

#### 
- [TASK-CR26062101-djer9d3vyo2m](TASK-CR26062101-djer9d3vyo2m_编写-flowforge-feedback-skillmd-主文件.md) [task] - 编写 flowforge-feedback SKILL.md 主文件
- [TASK-CR26062101-djer9jcjxe2n](TASK-CR26062101-djer9jcjxe2n_编写-feedback-分类规则-reference.md) [task] - 编写 feedback 分类规则 reference
- [DES-CR26062101-djer91do74qp](DES-CR26062101-djer91do74qp_feedback-skill-v2-采用五类分类路由器-cli.md) [design] - feedback SKILL v2 采用五类分类路由器 + CLI 原子操作
#### requires
- [TASK-CR26062101-djer9d3vyo2m](TASK-CR26062101-djer9d3vyo2m_编写-flowforge-feedback-skillmd-主文件.md) [task] - 编写 flowforge-feedback SKILL.md 主文件
- [TASK-CR26062101-djer9jcjxe2n](TASK-CR26062101-djer9jcjxe2n_编写-feedback-分类规则-reference.md) [task] - 编写 feedback 分类规则 reference

## Open Questions
None

