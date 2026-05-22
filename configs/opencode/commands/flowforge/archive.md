---
description: Archive a completed proposal and update its declared archive targets.
allowed-tools: Skill(flowforge)
---

Use the `FlowForge` skill to archive a proposal.
Before proceeding, follow `workflow/guides/rule-loading.md`.

## 执行流程

1. Confirm the proposal is implemented
2. Confirm backend tasks are closed
3. Run `.flowforge/scripts/flowforge-check-archive.js`
4. Run `.flowforge/scripts/flowforge-archive-proposal.js`
5. Confirm archive targets were updated and the proposal is `archived`

## 参数

- 提案编号：`CRYYMMDDNN`

Arguments: $ARGUMENTS
