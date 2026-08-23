# FlowForge Skills Specification & Engineering Methodology

> FlowForge v4 原子 Skill 矩阵、复杂度分诊与敏捷工程方法论。

---

## 1. 设计哲学：复杂度自适应与分流治理

过去的版本将所有职责强塞给通用的扁平流程，导致复杂架构重构与新模块开发在轻量模型下出现“信息真空”和严重跑偏。
FlowForge v4 借鉴 **`mattpocock/skills`** 的敏捷分诊与渐进式展开哲学，构建了自适应的原子 Skill 矩阵：

```text
┌─────────────────────────────────────────────────────────────────────────────┐
│                            FlowForge Skill 矩阵                             │
└─────────────────────────────────────────────────────────────────────────────┘
  0. flowforge-triage (Complexity Classification & Routing)
     └─ 目标：评估任务足迹、不确定性与爆炸半径，路由至适配的工作流
     └─ 借鉴：/triage + Phase Boundaries

  1. flowforge-align (Grill & Domain Alignment)
     └─ 目标：在对话中深挖需求、统一领域语言（CONTEXT.md）、支持 Flat / Hierarchical 工作记忆
     └─ 借鉴：/grill-me + Domain Modeling + Working Memory

  2. flowforge-wayfinder (Fog-of-War Map & Decision Frontier)
     └─ 目标：高不确定性场景下绘制决策树地图（MAP.md），以前沿票据推进
     └─ 借鉴：/wayfinder + Frontier Exploration

  3. flowforge-explore (Fact Finding)
     └─ 目标：针对特定疑点深入代码库调研，将事实回填至 Proposal
     └─ 借鉴：/explore + Code Seam Discovery

  4. flowforge-plan (Polymorphic Decomposition & Slicing)
     └─ 目标：共识达成后进行多态拆解（Tracer Bullets 功能穿透 vs Expand-Contract 宽幅重构批次）
     └─ 借鉴：/to-tickets + Expand-Contract Refactoring

  5. flowforge-implement (Seam-Bound TDD Cycle)
     └─ 目标：专注执行一个切片：写红测试 → 写代码 → 绿测试 → 消除认知负荷
     └─ 借鉴：/tdd + Red-Green-Refactor

  6. flowforge-diagnose (Hypothesis-Driven Root Cause Analysis)
     └─ 目标：针对缺陷构建最小隔离复现，证伪假设并修复根因
     └─ 借鉴：/diagnosing-bugs

  7. flowforge-review (Non-blocking Adversarial Review)
     └─ 目标：对抗性审查架构漂移、包可见性泄漏与认知负荷膨胀
     └─ 借鉴：/code-review

  8. flowforge-curate (Living Doc Synthesis)
     └─ 目标：提案交付后提取 ADR 与领域补丁，合流并归档提案
     └─ 借鉴：Living Documentation Merge
```

---

## 2. 原子 Skill 详细规范

### 2.0 `flowforge-triage` (前置分诊与复杂度路由)
* **触发场景**：接收用户新的工程任务输入。
* **复杂度级别判定**：
  * **Level 1 (Direct Patch / Local Bug)**: 单一文件/简单 Bug，直接进入 `diagnose` / `implement`。
  * **Level 2 (Standard Feature - Tracer Bullet)**: 纵向穿越单一模块分层，需求明确，走 `align (Flat)` $\rightarrow$ `plan (Tracer Bullets)`。
  * **Level 3 (Architecture Refactor / Wide Impact)**: 跨 $\ge 2$ 个子模块的大规模重构/解耦，走 `align (Hierarchical)` $\rightarrow$ `plan (Expand-Contract Batches)`。
  * **Level 4 (Greenfield / High Uncertainty)**: 全新孵化或高迷雾场景，走 `wayfinder (MAP.md)` $\rightarrow$ `Spikes` $\rightarrow$ 迷雾消退后降级规划。

### 2.1 `flowforge-align` (对齐与分层工作记忆捕获)
* **触发场景**：用户提出业务需求、架构设计或恢复正在讨论的提案。
* **执行准则**：
  1. **Domain Modeling First**：核对与更新 `docs/CONTEXT.md` 统一术语，杜绝同名异义与歧义；
  2. **Grilling 追问机制**：主动挖掘边界、假设与非目标（Non-goals）；
  3. **自适应存储结构**：
     - 单模块/标准任务：输出单个 `README.md`。
     - 跨模块/复杂重构：生成 `README.md`（总控与拓扑） + `modules/<module>.md`（物理迁移矩阵、接口签名与可见性约束）。

### 2.2 `flowforge-wayfinder` (迷雾航标与决策地图)
* **触发场景**：存在多个互斥架构方案、高不确定性技术选型或迷雾任务。
* **执行准则**：
  1. 维护 `01-workspace/<proposal_id>/MAP.md` 决策拓扑图；
  2. 将未决分歧转化为可探索的 Decision Tickets；
  3. 沿未阻塞的决策前沿（Decision Frontier）有序推进，严禁盲目大面积切片。

### 2.3 `flowforge-explore` (代码与数据事实探查)
* **触发场景**：讨论过程中遇到不确定的代码逻辑、遗留表结构、接口边界时。
* **执行准则**：
  1. **只查事实，不写业务代码**：快速定位具体函数、数据样本、异常堆栈；
  2. **结构化回填**：必须输出带有文件坐标（`file:line`）的事实链，回填至活笔记的 `[Explored Facts]`。

### 2.4 `flowforge-plan` (多态切片与批次规划)
* **触发场景**：用户确认对齐共识（"looks good", "开始规划"）。
* **执行准则**：
  1. **业务功能（Feature）**：采用垂直 Tracer Bullets（3~6 片，包含端到端行为）；
  2. **架构重构（Refactor）**：采用 Expand-Contract 批次（Batch 1: Expand 骨架 $\rightarrow$ Batch 2..N: 按模块迁移 $\rightarrow$ Batch Final: Contract 清理与可见性收敛）；
  3. **明确物理锚点**：每个切片必须具备明确的 Seams/文件路径与自动化测试验证命令。

### 2.5 `flowforge-implement` (缝隙约束 TDD 增量实现)
* **触发场景**：用户要求执行某个具体 Slice。
* **执行准则**：
  1. **红-绿-重构循环 (Red-Green-Refactor)**：
     - **Red**：先编写/更新失败的单元测试；
     - **Green**：编写最小代码使测试通过；
     - **Refactor (Cognitive Load Reduction)**：优化代码结构，消除死代码，内联浅封装，早返回，保持测试全绿。
  2. **测试作为唯一裁判**：只有绑定的测试命令 100% 运行通过，该切片才算完成。

### 2.6 `flowforge-diagnose` (假设驱动排错)
* **触发场景**：遇到复杂 Bug、偶发测试失败、生产回归或死锁排查。
* **执行准则**：
  1. 提出可证伪的假设；
  2. 构建最小隔离复现用例；
  3. 追踪变量与状态跃迁；
  4. 修复根因而非表面症状。

### 2.7 `flowforge-review` (非阻断式对抗性审查)
* **触发场景**：提案切片交付后、合并前进行红队安全与架构漂移审查。
* **执行准则**：
  1. 架构漂移检测（对比 `docs/CONTEXT.md`）；
  2. 认知负荷削减建议（死代码、浅封装、不必要的多层抽象）；
  3. 模块可见性（是否把 internal 的类泄漏为了 public）。

### 2.8 `flowforge-curate` (活文档合流与归档)
* **触发场景**：Proposal 所有切片完成并通过测试，准备归档。
* **执行准则**：
  1. **Delta 提取**：提取关键架构决策为 ADR，提取业务现状变更生成 Domain Patch；
  2. **Diff 审查与合流**：展示对 `docs/domains/<domain>/README.md` 的变更 Diff，用户确认后合并，并将 Proposal 移入 `00-archive/`。
