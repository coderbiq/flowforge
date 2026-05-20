---
description: Create or revise a proposal from a validated exploration.
allowed-tools: Skill(tg-workflow)
---

Use the `tg-workflow` skill to create or revise a proposal.

## 执行流程

1. Collect title, source exploration, archive targets, and task backend
2. Run `scripts/tg-create-proposal.js`
3. Review and refine `proposal.md`, `design.md`, and `task-map.md`
4. Validate with `scripts/tg-validate-proposal.js`

## 触发场景

- 用户要创建提案
- 探索笔记已定稿

## 约束

- 提案编号格式：`CRYYMMDDNN`
- 至少一个 `archive target`

Arguments: $ARGUMENTS
