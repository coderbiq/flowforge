# v3 Design SKILL 工作流

> 当前工作流。SKILL 通过 FlowForge CLI 维护结构化事实，v3 规范见 [`proposal-v3/skill-spec.md`](./proposal-v3/skill-spec.md)。

## 工作流

1. 明确 PROP 与 FEATURE 范围。
2. 在 FEATURE 内补充 Summary、Motivation、Design、Constraints 和 Implementation Plan。
3. 需要跨 FEATURE 复用时创建或引用 CONV、DEC、MOD、FIND。
4. 将 FEATURE 演进至 `designed`、`planned`，再由 Executor 按 Step 实施。
5. 用 `context feature` 获取执行上下文，用 `card steps`、`card log` 和 Verification 记录事实。

PROP 的 Feature Map 是人写的语义控制面；`proposal inspect` 负责机械聚合。Journal 只记录协作调度，不新增卡片类型。

## 当前入口

```bash
flowforge card init --type feature --title "..." --proposal <id>
flowforge card evolve <feature-id> --stage designed
flowforge context feature --feature <feature-id> --step <n>
flowforge proposal inspect <proposal-id>
```

旧 task、structure、log create、requirement 工作流直接删除。STR 仅作为 PROP control-plane metadata，不作为普通用户卡。

## 设计边界

设计只写已确认行为；Implementation Plan 是计划，不是实现证明。不能从设计推导迁移、UI、历史 wiki 或旧 ID/links 兼容承诺。遇到未决行为应停止并记录设计缺口，而不是自行扩展范围。
