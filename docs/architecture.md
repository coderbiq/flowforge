# FlowForge 架构设计 (v4)

## 1. 核心定位

FlowForge 是专为工程级 AI 协作打造的**多级工作记忆（Multi-Tier Working Memory）与活文档合流（Living Docs Synthesis）中枢**。

它坚决反对“形式主义的死板模板”与“限制开发活力的强状态机”，主张通过：
1. **前置复杂度分诊与分流（Complexity Triage & Dynamic Routing）**：杜绝一刀切，根据任务规模（Bug / Feature / Refactor / Greenfield）自动路由到最匹配的工作流；
2. **多级工作记忆体系（Multi-Tier Working Memory）**：支持单文件 Flat 模式与多层 Hierarchical 模式，让 Agent 和人类跨会话、跨天讨论不再遗忘，并为轻量模型提供确定的物理锚点；
3. **敏捷工程方法论（Conversation-first + Expand-Contract + TDD）**：用 Grilling 对齐需求、用 Tracer Bullets 穿透业务特性、用 Expand-Contract 批次治理大面积重构、用自动化测试作为唯一质量门禁；
4. **活文档持续合流（Living Docs Merge）**：让系统级领域知识随提案交付持续演进，永不腐烂。

---

## 2. 三层工作记忆架构

```text
┌─────────────────────────────────────────────────────────────┐
│  Tier 1: 项目级全局记忆 (docs/CONTEXT.md)                    │
│  - 全局统一术语表 (Ubiquitous Language)                      │
│  - 核心架构约束与编码原则                                   │
│  - 当前活跃提案索引                                         │
└──────────────────────────────┬──────────────────────────────┘
                               │ 注入全局背景
┌──────────────────────────────▼──────────────────────────────┐
│  Tier 2: 提案级活笔记 (01-workspace/<proposal_id>/)          │
│  - Flat 模式: README.md (单文件承载目标、决策、切片)         │
│  - Hierarchical 模式:                                       │
│      README.md (总控大纲与架构拓扑)                          │
│      modules/<module>.md (子模块物理规格、接口与迁移表)      │
│      references/ (长篇调研与源码行级事实)                   │
└──────────────────────────────┬──────────────────────────────┘
                               │ 提取切片最小上下文
┌──────────────────────────────▼──────────────────────────────┐
│  Tier 3: 任务级执行切片 (Slice Execution Context)           │
│  - 当前 Slice 目标 + 涉及接缝 (Seams) + 唯一测试命令        │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. 八大原子 Skill 协作流

FlowForge 将工程流程解耦为 8 个职责专一的 Skill，严禁单一 Skill 越权或黑盒脑补：

```text
                               ┌──────────────────────────────────────────────┐
                               │             User Prompt / Task               │
                               └──────────────────────┬───────────────────────┘
                                                      │
                                                      ▼
                                       ┌──────────────────────────────┐
                                       │       flowforge-triage       │ (复杂度分诊与路由)
                                       └──────────────┬───────────────┘
                                                      │
               ┌──────────────────────┬───────────────┴───────────────┬──────────────────────┐
               │ [Level 1: Patch/Bug] │ [Level 2: Feature]            │ [Level 3: Refactor]  │ [Level 4: Greenfield/Fog]
               ▼                      ▼                               ▼                      ▼
        flowforge-diagnose     flowforge-align                 flowforge-align        flowforge-wayfinder
               │                      │ (Flat Mode)                   │ (Hierarchical Mode)  │ (绘制决策探索地图)
               │                      │                               │                      │
               ▼                      ▼                               ▼                      ▼
        flowforge-plan         flowforge-plan                  flowforge-plan         Decision Tickets / Spikes
        (Minimal TDD)          (Tracer Bullets)               (Expand-Contract Batches)      │
               │                      │                               │                      │
               └──────────────────────┴───────────────┬───────────────┴──────────────────────┘
                                                      │
                                                      ▼
                                             flowforge-implement
                                         (Seam-Bound Red-Green TDD)
                                                      │
                                                      ▼
                                              flowforge-review
                                           (Adversarial Invariants)
                                                      │
                                                      ▼
                                              flowforge-curate
                                         (Living Docs & ADR Synthesis)
```

1. **`flowforge-triage`**：复杂度评估与门禁分诊，将任务按影响面和不确定性路由到最适合的工作流；
2. **`flowforge-align`**：领域建模与 Grilling 烤问，在 Flat 或 Hierarchical 模式下建立清晰的边界共识；
3. **`flowforge-wayfinder`**：高迷雾场景下的决策航标，维护决策图（`MAP.md`）推进探索前沿；
4. **`flowforge-explore`**：针对代码库探查事实，将证据（`file:line`）注入工作记忆；
5. **`flowforge-plan`**：多态切片规划（业务特性拆解为 Tracer Bullets，跨模块重构采用 Expand-Contract 批次）；
6. **`flowforge-implement`**：严格绑定 Seam 接口执行红-绿-重构（TDD）循环；
7. **`flowforge-diagnose`**：假说驱动的缺陷隔离与根因分析；
8. **`flowforge-review`**：对抗性审查，拦截架构漂移与认知负荷膨胀；
9. **`flowforge-curate`**：提取核心决策为 ADR，将领域业务变更合流至 `docs/domains/<domain>/README.md`，归档提案。

---

## 4. 活文档与 ADR 归档体系

文档不再是孤岛式卡片，而是系统级的活资产：

* `docs/architecture/decisions/`：记录全局架构决策（ADRs）；
* `docs/domains/<domain>/README.md`：记录各领域当前的真实业务行为、术语与接口接缝；
* `archive/<proposal_id>/`：封存历史提案工作过程。
