---
description: Start execution for an approved proposal using its canonical task map.
allowed-tools: Skill(tg-workflow)
---

Use the `tg-workflow` skill to apply a proposal.

## 执行流程

1. Read and validate `meta.yaml` and `task-map.md`
2. Require proposal status `approved`
3. Run `scripts/tg-apply-proposal.js`
4. Confirm backend tasks and proposal state changed to `active`

## 参数

- 提案编号：`CRYYMMDDNN`

Arguments: $ARGUMENTS
