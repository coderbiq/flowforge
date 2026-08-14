# v3 卡片模型

> 当前、权威规范。实现状态以源码和测试为准。

## 类型

| 类型 | 含义 | 范围 |
|---|---|---|
| `PROP` | Proposal control-plane metadata | 提案 |
| `FEATURE` | 功能完整生命周期 | 提案 |
| `CONV` | 编码或流程约定 | library |
| `DEC` | 架构决策 | library |
| `MOD` | 模块职责与边界 | library |
| `FIND` | 发现、证据与异常行为 | library |

这是完整的 v3 类型集合。REQ、DES、TASK、LOG、ROOT 和普通 STR 不属于当前模型。

## FEATURE

FEATURE 以一个 Markdown 文件承载 Summary、Motivation、Design、Constraints、Implementation Plan、Verification 和 History。阶段演进为：

```text
draft -> designed -> planned -> in_progress -> done
```

步骤是 FEATURE 内的计划单元。Implementation Plan 描述计划，不是实现事实；实现事实必须来自代码、测试和 Verification 证据。

## PROP

PROP 记录提案目标、Feature Map、架构概览和跨 FEATURE 约束。Feature Map 是人写的语义关系，`proposal inspect` 生成状态、依赖和健康检查。用于聚合的 STR metadata 不作为普通卡片暴露。

## 库卡片

CONV、DEC、MOD、FIND 只在跨 FEATURE 复用时独立成卡。单 FEATURE 的决策和约束写在该 FEATURE 中。卡片通过 `card link` 管理 typed links；当前入口不引入旧 ID 或旧 links。

## 历史边界

旧模型仅可在明确标记 `historical` 的背景材料中出现。v3 不承诺迁移旧 wiki、旧 ID/links 或提供 UI。
