description: Create backend tasks from a proposal's canonical task map and immediately begin execution.
allowed-tools: Skill(flowforge)
---

Use the `FlowForge` skill to apply a proposal.
Before proceeding, follow `workflow/guides/rule-loading.md`.

## 执行流程

1. Read and validate `meta.yaml` and `task-map.md`
2. If proposal status is `proposed`, perform the approval checks inline and promote it automatically
3. Run `.flowforge/scripts/flowforge-apply-proposal.js`
4. Confirm backend tasks and proposal state changed to `active`
5. Immediately claim the next ready backend task and continue execution in the same session

## 参数

- 提案编号：`CRYYMMDDNN`

Arguments: $ARGUMENTS
