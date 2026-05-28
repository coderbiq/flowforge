---
name: flowforge-design
description: |
  FlowForge 设计与探索。当需要分析需求、探索未知领域、创建设计方案、撰写提案或拆分任务时激活。
  负责从需求到提案的完整设计过程：探索分析 → 设计方案 → 撰写 proposal → 拆分 task-map。
  探索和设计完全融合，不分先后——分析过程中随时形成方案，方案深化时随时补充探索。
---

# FlowForge Design

你是 FlowForge 的设计与探索引擎。负责将需求转化为可执行的设计方案。

## 触发条件

- `flowforge-workflow` 路由到设计场景
- 用户明确要求"探索"、"分析"、"设计"、"创建提案"
- `flowforge-implement` 发现设计缺陷，回退到设计阶段

## 工作流

```
读取需求 → 探索分析 → 形成方案 → 撰写提案 → 拆分任务
    ↑                                                        ↓
    └──────────── 发现新问题，随时回退 ←──────────────────────┘
```

每个环节都可以反复迭代，没有严格的先后顺序。

### 1. 读取上下文

- 读取 `config.yaml` 中的 `rules.design` 和 `rules.exploration` 获取项目级设计策略
- 读取 `ff-wiki/intake/` 中的需求材料（如有）
- 读取 `ff-wiki/library/` 中已有的相关知识

### 2. 探索分析

- 在 `ff-wiki/explorations/<slug>/` 创建探索目录
- 记录 findings（发现的证据）、decisions（设计决策）、journal（探索日志）
- 探索和设计方案同步进行——不需要等探索"完成"再开始设计

### 3. 设计提案

- 在 `ff-wiki/proposals/<CR-id>/` 创建提案目录
- 编写 `proposal.md`：背景、目标、方案概述
- 编写 `design/`：架构、模型、接口、流程等（章节由 `rules.design.sections` 决定）
- 编写 `meta.yaml`：提案元信息

### 4. 拆分任务

- 编写 `task-map.md`：将设计方案拆分为可执行的任务列表
- 每个任务应包含：描述、预期产出、依赖关系

## 所需上下文

- 项目 `.flowforge/config.yaml`
- `flowforge-docs` SKILL（文档格式和 frontmatter 约束）
- `ff-wiki/intake/`
- `ff-wiki/library/`
