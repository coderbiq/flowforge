# FlowForge 架构设计 (v4)

## 1. 核心定位

FlowForge 是专为工程级 AI 协作打造的**多级工作记忆（Multi-Tier Working Memory）与活文档合流（Living Docs Synthesis）中枢**。

它坚决反对“形式主义的死板模板”与“限制开发活力的强状态机”，主张通过：
1. **多级工作记忆体系**：让 Agent 和人类跨会话、跨天讨论不再遗忘；
2. **敏捷工程方法论（Conversation-first + TDD）**：用 Grilling 对齐需求、用 Tracer Bullets 拆解极小切片、用自动化测试作为唯一质量门禁；
3. **活文档持续合流（Living Docs Merge）**：让系统级领域知识随提案交付持续演进，永不腐烂。

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
│  Tier 2: 提案级活笔记 (01-workspace/<proposal_id>/README.md) │
│  - Objective: 核心目标 (1-2 句)                             │
│  - Ubiquitous Language: 本提案专用术语                      │
│  - Explored Facts: 代码/数据探查事实 (带 file:line)         │
│  - Key Decisions: 共识与决策 (Why)                          │
│  - Open Questions: 待对齐疑问                               │
│  - Actionable Slices: 3-6 个 Tracer Bullet 切片与测试绑定   │
└──────────────────────────────┬──────────────────────────────┘
                               │ 提取切片最小上下文
┌──────────────────────────────▼──────────────────────────────┐
│  Tier 3: 任务级执行切片 (Slice Execution Context)           │
│  - 当前 Slice 目标 + 涉及接缝 + 唯一测试命令                │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. 五大原子 Skill 协作流

FlowForge 将工程流程解耦为 5 个职责专一的 Skill，严禁单一 Skill 越权或黑盒脑补：

1. **`flowforge-align`**：通过启发式提问（Grilling）澄清边界、非目标与失败语义，维护提案活笔记；
2. **`flowforge-explore`**：针对代码库与历史数据探查事实，将证据（`file:line`）回填活笔记；
3. **`flowforge-plan`**：在达成共识后，将方案拆解为 3~6 个极小 Tracer Bullet 切片并绑定自动化测试；
4. **`flowforge-implement`**：严格执行红-绿-重构（TDD）循环，测试通过即切片完成；
5. **`flowforge-curate`**：提取核心决策为 ADR，将领域业务变更合流至 `docs/domains/<domain>/README.md`，归档提案。

---

## 4. 活文档与 ADR 归档体系

文档不再是孤岛式卡片，而是系统级的活资产：

* `docs/architecture/decisions/`：记录全局架构决策（ADRs）；
* `docs/domains/<domain>/README.md`：记录各领域当前的真实业务行为、术语与接口接缝；
* `archive/<proposal_id>/`：封存历史提案工作过程。
