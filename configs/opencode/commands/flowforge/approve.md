---
description: Move a valid proposal into approved state before task application.
allowed-tools: Skill(flowforge)
---

Use the `FlowForge` skill to approve a proposal.

## 执行流程

1. Validate `meta.yaml`, `design.md`, and `task-map.md`
2. Require proposal status `draft` or `proposed`
3. Run `scripts/flowforge-approve-proposal.js`
4. Confirm proposal state changed to `approved`

## 参数

- 提案编号：`CRYYMMDDNN`

Arguments: $ARGUMENTS
