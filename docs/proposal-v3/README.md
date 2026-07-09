# FlowForge 卡片模型 v3 重构方案

> 日期：2026-07-09
>
> 基于 CR26063001 和 CR26070801 两个真实提案的实证分析，提出从 10 种卡片类型到 FEATURE 阶段演进模型的根本性重构。

---

## 方案概述

**核心变化**：将 REQ+DES+TASK 三种按类型拆分的卡片合并为一张 FEATURE 卡片，用阶段演进（draft → designed → planned → in_progress → done）替代跨卡跳转。卡片类型从 10 种精简为 5 种（+ PROP）。同时放宽 Agent 直接读写卡片的约束，CLI 专注于不变式保护。

## 文档索引

| 文档 | 内容 | 读它当你想... |
|------|------|-------------|
| [card-model.md](./card-model.md) | 问题诊断、FEATURE 生命周期/模板、PROP 全景、拆分策略、横切类型 | 理解"为什么改"和"改成什么样" |
| [cli-spec.md](./cli-spec.md) | 约束放宽、废弃清单（16 个命令）、新增命令（init/evolve/log/steps/split/context feature）、修改命令、迁移策略 | 实现 CLI 改造或了解命令变更 |
| [skill-spec.md](./skill-spec.md) | Token 消耗控制设计、需求分解方法论、信息探索/设计推理/任务拆分/约束构建的 Agent 思维链 | 编写或审查 SKILL 文件 |

## 关键决策记录

| 决策 | 理由 |
|------|------|
| REQ+DES+TASK → FEATURE (阶段演进) | 三者信息高度重叠，拆分只产生维护成本 |
| 保留 CONV/MOD/DEC/FIND | 横切关注点，天然跨功能生效 |
| STR → `proposal inspect` 自动聚合 | 手动维护的合成价值不抵维护成本 |
| 门控在 CLI 层 (`card evolve`) 强制执行 | 不依赖 SKILL 自律 |
| Agent 可直接读写 .md 文件 | CLI 为人类设计，Agent 需要直接文件访问效率 |
| CLI 保留链接/阶段/进展操作 | 多文件一致性和复杂解析逻辑 |

## 相关文档

- [methodology-review-card-fragmentation.md](../methodology-review-card-fragmentation.md) — 原始问题诊断
- [remediation-card-fragmentation.md](../remediation-card-fragmentation.md) — 补丁方案（已实现但未触及根因）
- [methodology-card-model-simplification.md](../methodology-card-model-simplification.md) — 早期模型简化方向探索
