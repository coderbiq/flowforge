description: Optional review gate that can move a valid proposal into approved state before task application.
allowed-tools: Skill(flowforge)
---

Use the `FlowForge` skill to approve a proposal.
Before proceeding, follow `workflow/guides/rule-loading.md`.

## 执行流程

1. Validate `meta.yaml`, `design.md`, and `task-map.md`
2. Require proposal status `draft` or `proposed`
3. Run `.flowforge/scripts/flowforge-approve-proposal.js`
4. Confirm proposal state changed to `approved`

This command is optional. `/flowforge:apply` can approve a proposed proposal inline and move straight into execution.

## 参数

- 提案编号：`CRYYMMDDNN`

Arguments: $ARGUMENTS
