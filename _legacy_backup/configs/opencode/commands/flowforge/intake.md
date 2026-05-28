---
description: Create or revise a durable intake package before exploration begins.
allowed-tools: Skill(flowforge)
---

Use the `FlowForge` skill to create or revise an intake package.

## 执行流程

1. Create or update `docs/intake/<slug>/`
2. Populate `index.md`, `references.md`, `questions.md`, and any supporting assets
3. Run `.flowforge/scripts/flowforge-intake-context.js` or `.flowforge/scripts/flowforge-explore-context.js` as needed to review the current package
4. Keep the intake package current while exploration is in progress

## 触发场景

- 用户要先提供需求材料再进入探索
- 用户提到“输入包”、“起始目录”、“brief”
- 用户希望先整理截图、链接、约束和未决问题

Arguments: $ARGUMENTS
