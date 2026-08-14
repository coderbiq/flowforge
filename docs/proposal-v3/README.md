# FlowForge v3 规范

> `authoritative` / `current`。本文档族是 v3 唯一规范；实现状态以源码、测试和 `docs/index-management.md` 为准。

## 模型

| 类型 | 作用 |
|---|---|
| `PROP` | Proposal control-plane metadata；提案范围、FEATURE 关系和聚合入口。 |
| `FEATURE` | 一个可交付功能的完整生命周期。 |
| `CONV` | 可执行的跨功能约定。 |
| `DEC` | 跨功能架构决策。 |
| `MOD` | 可复用模块知识。 |
| `FIND` | 探索发现和证据。 |

FEATURE 阶段为 `draft`、`designed`、`planned`、`in_progress`、`done`。PROP 的 Feature Map 是语义控制面；`proposal inspect` 是机械聚合视图。

STR 仅表示 Proposal control-plane metadata，不是普通用户卡。旧类型、旧 ID/links、历史 wiki 和迁移能力不属于当前 v3 规范。

## CLI 原则

当前入口是 `card init`、`card evolve`、`card log`、`card steps`、`card link`、`context feature`、`proposal inspect` 和库查询命令。旧 `task`、`structure`、`log create`、`requirement` CLI 直接删除，不提供 deprecated 兼容入口。

## 文档边界

- `card-model.md` 定义模型与生命周期。
- `cli-spec.md` 定义当前 CLI 契约。
- `skill-spec.md` 定义 Agent 工作流。
- `implementation-plan.md` 只记录计划，不能证明实现已完成。
