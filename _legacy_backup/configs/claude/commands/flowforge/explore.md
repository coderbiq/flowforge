---
description: Create or extend a durable exploration using the canonical FlowForge structure.
allowed-tools: Skill(flowforge)
---

Use the `FlowForge` skill to create or extend an exploration.
Before proceeding, follow `workflow/guides/rule-loading.md` and materialize
the current intake package with `scripts/flowforge-explore-context.js` when
one exists. If no intake package exists yet, use the helper without an intake
argument so the exploration still starts from the rule bundle.

## 执行流程

1. Create `docs/explorations/<slug>/`
2. Initialize `index.md`, `journal/`, `findings/`, `decisions/`, and `artifacts/`
3. Record durable findings before implementation
4. Keep `index.md` readable and current

## 触发场景

- 用户描述新需求
- 用户询问"如何实现..."
- 用户提到"探索"、"调研"、"分析"

Arguments: $ARGUMENTS
