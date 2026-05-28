---
description: Create or revise a proposal from a validated exploration.
allowed-tools: Skill(flowforge)
---

Use the `FlowForge` skill to create or revise a proposal.
Before proceeding, follow `workflow/guides/rule-loading.md`.

## 执行流程

1. Collect title, source exploration, archive targets, and task backend
2. Run `.flowforge/scripts/flowforge-create-proposal.js`
3. Review and refine `proposal.md`, `design.md`, and `task-map.md`
4. Validate with `.flowforge/scripts/flowforge-validate-proposal.js`

## 触发场景

- 用户要创建提案
- 探索笔记已定稿

## 约束

- 提案编号格式：`CRYYMMDDNN`
- 至少一个 `archive target`

Arguments: $ARGUMENTS
