# Multi-Tier Working Memory & Living Documentation

> FlowForge v4 多级工作记忆体系与活文档合成规范。

---

## 1. 背景与核心价值

在复杂的真实工程中，开发过程往往面临两大挑战：
1. **长周期讨论与会话断裂（Memory Decay）**：复杂需求探讨动辄数天，跨越多次会话和上下文压缩。缺乏结构化的外部工作记忆会导致调查过的信息和已达成的共识频繁丢失。
2. **粒度失真与轻量模型跑偏（Information Void vs Cognitive Overload）**：大范围跨模块重构如果只用简陋的单文件描述，轻量模型执行切片时极易因缺乏物理类名、包路径和接口签名而严重跑偏；若一次性全量加载，又会耗尽上下文。

FlowForge v4 建立了**多级工作记忆体系（Multi-Tier Working Memory System）**与**活文档合成机制（Living Doc Synthesis）**，实现跨会话无缝断点续传与系统知识的持续自动合流。

---

## 2. 多级工作记忆体系 (Multi-Tier Working Memory)

```text
┌─────────────────────────────────────────────────────────────┐
│  Tier 1: 项目级全局记忆 (Project Living Memory)             │
│  路径: docs/CONTEXT.md                                      │
│  内容: 系统架构常识、长期技术约束、全局统一语言、活跃提案索引│
└──────────────────────────────┬──────────────────────────────┘
                               │ 注入全局背景
┌──────────────────────────────▼──────────────────────────────┐
│  Tier 2: 提案级活笔记 (Proposal Working Scratchpad)          │
│  路径: 01-workspace/<proposal_id>/                          │
│  模式: Flat 模式 (单文件 README.md)                          │
│        Hierarchical 模式 (README.md + modules/*.md)         │
└──────────────────────────────┬──────────────────────────────┘
                               │ 提取切片最小上下文
┌──────────────────────────────▼──────────────────────────────┐
│  Tier 3: 任务级执行切片 (Slice Execution Context)           │
│  机制: 动态装配 (主控 README + 对应子模块规格)              │
│  内容: 单个 Step 目标 + 涉及代码接缝 + 唯一绑定的测试命令    │
└─────────────────────────────────────────────────────────────┘
```

### 2.1 Tier 1: 项目级全局记忆 (`docs/CONTEXT.md`)
* **定位**：常驻项目根目录的极简统一知识，控制在 300~500 tokens。
* **主要内容**：
  * **Ubiquitous Language**：全局通用的领域术语定义。
  * **Architectural Constraints**：不可逾越的技术选型与工程规范（如：分层依赖单向性、模块可见性、TDD 纪律）。
  * **Active Proposals Index**：当前进行中的提案索引。
* **生命周期**：全局唯一，任何 Agent 会话启动时默认加载。

### 2.2 Tier 2: 提案级活笔记 (`01-workspace/<proposal_id>/`)
* **定位**：长周期讨论的**外脑（External Scratchpad）**。在需求探讨的第一分钟即创建，随着人机对话不断精炼追加。
* **自适应规模与渐进式披露（Adaptive Scale & Progressive Disclosure）**：
  * **单文件极简模式（Flat Mode）**：适用于 Level 2 标准特性任务（单一业务域、<= 5 个切片），所有目标、共识、事实与切片全部内聚在 `README.md` 中。
  * **多层总分模式（Hierarchical Mode）**：当命中 Level 3/4 复杂特征（跨 $\ge 2$ 个独立子模块、大型架构重构、长篇调研事实 > 500字）时，Agent 自动平滑升级为总分结构：
    ```text
    01-workspace/<proposal_id>/
    ├── README.md           # 总控大纲与架构拓扑 (Hub)
    ├── modules/            # 子模块物理规格 (Exact Interface & Move Matrix)
    │   ├── contracts.md    # 确切的 Public Interface 签名与包名
    │   └── domain-xxx.md   # 类物理搬迁表、依赖配置与可见性约束
    ├── references/         # 深度调研与外部知识 (Web/API/Schema，带 file:line)
    └── JOURNAL.md          # 关键演进流水
    ```
* **核心章节与准入标准（Definition of Ready）**：
  1. `Objective`：核心业务/技术目标（1~2 句话）。
  2. `Ubiquitous Language`：本提案专属的名词定义（消除歧义）。
  3. `Explored Facts`：调查代码/数据得到的权威事实（必须附 `path/file:line`）。
  4. `Key Decisions & Consensus`：已达成共识的技术方案（记录 Why，复杂子模块拆解至 `modules/`）。
  5. `Open Questions`：待对齐问题清单（跨会话优先追问）。
  6. `Actionable Slices`：确认后的切片清单及绑定的自动化测试。业务功能采用 Tracer Bullets，架构重构采用 Expand-Contract 批次。
* **断点续传机制**：跨会话时，Agent 优先读取此文件，无需用户重新复述背景，即可 100% 恢复讨论状态。

### 2.3 Tier 3: 任务级执行切片 (Slice Context)
* **定位**：执行具体代码改造时的瞬时上下文。
* **机制**：由主控 `README.md` 与该切片关联的 `modules/<module>.md` 组合动态生成，仅包含当前切片的目标、物理接缝文件与测试命令，确保最大、最干净且信息确定的代码上下文窗口。

---

## 3. 活文档合成与归档机制 (Living Doc Synthesis)

传统的卡片库做完即孤岛。FlowForge v4 引入了**“像代码 Git Merge 一样持续合流”**的活文档机制。

```text
提案交付完成 (Proposal Done)
      │
      ▼
【第一步：提取 (Extract)】 
  ├─ 核心决策 ──▶ 提炼为 ADR (docs/architecture/decisions/0042-xxx.md)
  └─ 业务/接口变更 ──▶ 提炼为 Domain Patch
      │
      ▼
【第二步：定位与比对 (Match & Diff)】
  定位到对应的领域主干文档: docs/domains/policy/README.md
  Agent 生成待合并的 Markdown Diff
      │
      ▼
【第三步：人工确认 (User Confirm)】
  用户确认变更无误后，写入领域主干文档，Proposal 移入 archive/
```

### 3.1 系统级 Living Docs 组织结构
系统级知识按**业务领域（Domain）与架构决策（ADR）**组织，而非按历史提案碎片堆叠：
```text
docs/
├── CONTEXT.md              # Tier 1 全局统一语言与约束
├── architecture/
│   ├── overview.md         # 系统全景图与技术栈
│   └── decisions/          # ADR 归档库 (0001-xxx.md, 0002-xxx.md)
└── domains/                # 领域模型与业务现状主干
    ├── policy/
    │   └── README.md       # 保单生命周期、核心规则、代码接缝
    └── data-migration/
        └── README.md       # 迁移管道、权威映射口径、基线位置
```

### 3.2 归档合成流程（`flowforge-curate`）
1. **Delta 提取**：丢弃排查流水与中间讨论，仅提取关键 ADR 决策与 Domain 现状变更。
2. **Diff 审查**：生成对 `docs/domains/<domain>/README.md` 的差分 Patch。
3. **确认合流**：人工确认后合入主干文档，提案归档至 `00-archive/`。
