---
name: flowforge-design
description: |
  FlowForge 任务驱动渐进式探索设计引擎。在用户表达需求/想法/变更意图，
  需要将需求转化为可执行方案时激活。通过 analysis/design 任务驱动
  探索和设计过程，使每一步探索都有明确目标、每一步设计都有来源追溯。

  必须在以下场景激活：
  - 用户描述新需求、想法或变更意图（"我想做..."、"需要..."、"打算..."）
  - 用户明确表达"探索"、"分析"、"设计"、"创建提案"、"拆分任务"
  - 已有 active 状态 proposal 需要继续完善设计或撰写
  - flowforge-implement 在实施中发现设计缺陷需要回退

  不要在以下情况激活：
  - 用户只是询问已有设计内容（只读查询，不修改）
  - 用户在执行 task-map 中的任务（应交给 flowforge-implement）
  - 用户要求归档或沉淀（应交给 flowforge-archive）
  - 仅用于更新进度索引——那是 flowforge-progress 的职责
  - 实施中发现需要补充探索内容——由 flowforge-feedback 结构化捕获后路由，不要直接激活本 SKILL 写 findings
---

# FlowForge Design

负责将需求转化为可执行的 proposal 和 task-map。

## 工作流

```
定位上下文 → 理解问题 → 确定 project → 加载 rules → 渐进式探索设计 → 撰写 proposal → 细化实施任务
```

---

### 阶段 1：定位上下文

运行 `scripts/design-context.js` 加载上下文。输出包含：

- `## Projects`：本仓库配置的所有 project（id、name、wikiRoot、srcDirs、description、keywords）
- `## Intake Material`：待处理 intake 文件，按所属 project 分组
- `## Current Proposal`（如有）：当前 proposal 路径 + 已锁定的 project/wikiRoot

根据上下文判断当前是**新需求**还是**继续已有 proposal**：

**新需求**（无活跃 proposal，或用户表达全新意图）：
- 检查输出中是否有 intake 材料
- 有则根据 intake 内容辅助判断 project
- 进入阶段 2 收集更多信息（如需要）

**继续已有 proposal**（`## Current Proposal` 中有 project 和 wikiRoot）：
- project 和 wikiRoot 已锁定，**永不重新决策归属**
- 跳过阶段 2 和阶段 3，**直接进入阶段 4** 加载 rules

---

### 阶段 2：理解问题

仅在**新需求**且没有 intake 材料时执行。

向用户确认核心诉求、影响范围和已知约束。信息不足就提问，不跳过。

如有 intake 材料，参照 `design-context.js` 输出的 `## Intake Strategy` 分析策略，按策略中描述的优先级和方法提取需求信息。

---

### 阶段 3：确定 project 归属

仅在**新需求**时执行（修改已有 proposal 时跳过——`## Current Proposal` 已锁定 project）。

**目标**：从 `## Projects` 列出的候选中选定一个 project.id，决定后续所有文档落点的 wikiRoot。一次性决策，写入 meta.yaml 后永不变更。

**决策算法**：

1. **若 projects.length == 1** → 直接选定该唯一 project，跳过后续步骤。

2. **若 projects.length > 1**，对每个 project 计算两个得分：
   - **srcDir_score**：列出本次设计涉及/将引用的源文件路径，统计落在该 project.srcDirs 下的文件数
   - **keyword_score**：分析 intake 材料 + 用户描述 + 阶段 2 收集的核心诉求，统计命中该 project.keywords 的次数

3. **总分** = `srcDir_score * 2 + keyword_score`（源码引用权重更高，因为它是硬证据）

4. **判定**：
   - 若所有 project 得分均为 0 → **必须列出所有 projects 让用户选**
   - 否则按总分排序：
     - `top1 / (top1 + top2) >= 0.7` → 选 top1（无歧义）
     - 否则 → **必须问用户**，列出 top3 + 每个的得分依据

5. **询问用户的格式**：
   ```
   本次新建文档可能属于以下 project，请选择：
   1. <id1> (<name>) — srcDir 命中 N 个, keyword 命中 M 个
   2. <id2> (<name>) — srcDir 命中 N 个, keyword 命中 M 个
   ...
   ```

**跨项目影响处理**：若设计明显涉及多个 project（例如同时改前端和后端），仅选定一个**主归属** project（按上述算法），将跨项目依赖记录为 finding（`type: cross-module`），但 proposal 主体只能落在一个 wikiRoot 下。

**输出**：选定的 `project.id` 和对应的 `project.wikiRoot`。

---

### 阶段 4：加载 project 规则

确定 projectId 后，再次运行 design-context.js：

```bash
scripts/design-context.js <projectRoot> --project <projectId>
```

此次输出包含：

- `## Exploration Strategy`：该 project 的探索策略（如存在）
- `## Design Rules`：命名规则（proposal_id 模板）、任务规则（字段、时间估计）
- `## Design Strategy`：指导 Agent 如何进行方案分析和设计决策的项目级策略（如存在）
- `## Implement Rules`：任务状态机、日志字段
- `## Library Rules`：归档行为（requireReview、autoUpdateHistory）
- `## Domain 分类指引`：如何为文档设置 domain 字段（scope、module、type 的判定规则）
- `## Intake Strategy`：intake 分析策略（如存在）

**重要**：project 配置中不再有 `modules` 注册表。模块判定基于源文件路径和设计落地位置，不依赖预注册。

---

### 阶段 5：任务驱动的渐进式探索设计

**核心思路**：从阶段 4 结束时即创建 task-map.yaml，以 analysis 和 design 两类任务驱动整个探索与设计过程。每个任务有明确的目标和产出，探索和设计不再是模糊的循环，而是可追踪、可拆分的工作单元。

#### 5.1 初始化任务地图

进入阶段 5 后立即执行（仅首次进入时）：

1. 根据 `naming.proposal_id` 的模板生成 CR-id，在 `<project.wikiRoot>/workspace/proposals/active/<CR-id>/` 下创建 proposal 目录

2. 创建 `task-map.yaml`，写入一个根级 analysis 任务作为分析的起点：

```yaml
proposal_id: <CR-id>
tasks:
  - id: "1"
    title: 分析识别提案的需求边界
    type: analysis
    status: pending
    dependencies: []
    sourceTasks: []
    epic: []
```

3. 运行脚本将任务同步到存储层：

```bash
node scripts/task-create.js <projectRoot> <CR-id>
```

#### 5.2 任务驱动的探索设计循环

进入持续循环，每次选取一个 `pending` 的 analysis 或 design 任务执行：

```
┌─────────────────────────────────────────────────────┐
│  选取一个 pending 任务（优先 analysis → design）    │
│                     ↓                               │
│  analysis 任务：探索代码/library → 记录发现         │
│        ↓ 发现子分析点                               │
│  task-add.js 添加子 analysis 任务                    │
│        ↓ 分析充分                                   │
│  task-add.js 添加 design 任务（sourceTasks 指向分析）│
│                     ↓                               │
│  design 任务：撰写设计文档 → 完成设计               │
│        ↓ 设计中发现新问题                           │
│  task-add.js 添加新 analysis 任务                    │
│                     ↓                               │
│  标记当前任务 done → 回到选取下一个 pending 任务     │
└─────────────────────────────────────────────────────┘
```

**选取任务的优先级**：
1. `pending` 的 analysis 任务优先于 design 任务（先分析再设计）
2. 同类任务中，id 较小的优先（先创建的先处理）
3. 有 `dependencies` 的任务，依赖全部完成后才可选

**处理 analysis 任务**：

1. 认领任务：标记 status 为 `in_progress`
2. 根据任务标题确定探索方向，按以下策略执行：
   - 如有 `## Exploration Strategy`，按其探索策略进行
   - 探索代码库、library 文档、外部资料
3. 记录发现（直接写入 library 对应路径）：
   - 系统架构事实 → `library/architecture/<topic>.md`
   - 模块设计事实 → `library/modules/<name>/`
   - 可复用决策 → `library/decisions/`
   - 可复用约定 → `library/conventions/<topic>.md`
   - 参照 `flowforge-docs` 获取各 doc_type 的写作指南
4. 在 proposal 的 `notes.md` 中记录探索过程（探索了哪些 library 文档、发现了什么）
5. **每个发现携带 `domain` frontmatter**，标注 scope/module/type：
   ```yaml
   ---
   domain:
     scope: module
     module: auth
     type: design
   ---
   ```
6. 分析过程中发现新的子探索点 → 通过 `task-add.js` 创建子 analysis 任务：
   ```bash
   node scripts/task-add.js <projectRoot> <CR-id> "<分析子任务标题>" "<描述>" <父任务id>
   ```
   新任务的 `epic` 指向根 analysis 任务（id: "1"）
7. 该模块分析充分 → 创建对应的 design 任务：
   ```bash
   node scripts/task-add.js <projectRoot> <CR-id> "<设计任务标题>" "<描述>" <依赖的分析任务id>
   ```
   design 任务的 `sourceTasks` 指向其依据的 analysis 任务，`epic` 指向根 analysis 任务
8. 任务完成 → 标记为 `done`

**处理 design 任务**：

1. 认领任务：标记 status 为 `in_progress`
2. 如有 `## Design Strategy`，参照其项目级设计策略指导方案分析、架构决策和设计文档的撰写方向
3. 在 proposal 目录的 `design/` 下撰写设计文档：
   - 参照 `flowforge-docs` SKILL 获取 design 的写作指南
   - **多模块拆分规则**：如果一个 proposal 涉及多个模块（如同时改 auth 和 session），按模块用子目录组织。例如：
     ```
     design/auth/architecture.md    → domain: { scope: module, module: auth, type: design }
     design/auth/api.md             → domain: { scope: module, module: auth, type: design }
     design/session/architecture.md → domain: { scope: module, module: session, type: design }
     ```
     单模块 proposal 可省略子目录，直接在 `design/` 下平铺。
   - 跨模块的通用架构设计（如整体数据流变更）→ `scope: system`，放在 `design/` 根目录下。
   - **每个设计文档只设一个 domain**，确保归档时能精确路由。
4. 设计中发现新的未探索问题 → 通过 `task-add.js` 创建新的 analysis 任务，status 为 `pending`
5. 设计文档完成 → 标记任务为 `done`

**判定 domain 的方法**（同之前，不变）：
```
scope:  文档引用的源文件在哪个模块下？设计最终在哪个模块落地？
        单模块边界内 → module；跨模块 / 全局架构 → system
module: scope=module 时，模块名是什么（如 auth、payment）
type:   架构/接口/数据模型/技术选型 → design
        关键决策+理由+备选方案评估 → decision
        编码规范/命名约定/反例 → convention
```

#### 5.3 终止条件

阶段 5 在以下**所有**条件满足时结束：

1. task-map.yaml 中所有 `analysis` 类型任务状态为 `done`
2. task-map.yaml 中所有 `design` 类型任务状态为 `done`
3. 没有未解决的开放问题（或已记录为 finding）

达到终止条件后，向用户简要说明当前方案，确认可以进入撰写 proposal 阶段，收到确认后再进入阶段 6。

---

### 阶段 6：撰写 proposal

将设计决策提炼为 `<project.wikiRoot>/workspace/proposals/active/<CR-id>/proposal.md`。

task-map.yaml 已在阶段 5 创建并包含完整的 analysis/design 任务追踪记录。proposal.md 应引用已完成的分析和设计任务，确保方案可追溯到具体的分析发现。

参照 `flowforge-docs` SKILL 获取 proposal 的写作指南。

同时参照 `flowforge-docs` 的 proposal meta 契约创建 `meta.yaml`。**必须**写入以下字段：
- `project: <id>`（来自阶段 3 选定的 project.id）——这是下游 SKILL/脚本定位 wikiRoot 的唯一依据
- `modules: [<name>, ...]`（轻量列表，仅包含涉及的模块名，用于 INDEX.md 展示）

`meta.yaml` 中**不再需要** `ownership` 和 `archive_targets`——归档路径由各文档的 `domain` frontmatter 自动推导。

proposal 创建完成后，运行 `scripts/validate-proposal.js <proposal路径>` 确保结构完整。

---

### 阶段 7：细化实施任务

task-map.yaml 中已有 analysis 和 design 任务（阶段 5 创建并完成）。本阶段将已完成的设计转化为可执行的 implementation 任务，或对已存在的 implementation 任务进行细化。

#### 首次创建实施任务

如果 task-map.yaml 中尚无 implementation 类型的任务：

1. 基于已完成的 design 任务，将设计方案拆分为可执行的 implementation 任务
2. 每个 implementation 任务包含字段：`id`、`title`、`description`、`deliverable`、`dependencies`、`status`（初始 `pending`）
3. 每个 implementation 任务的 `sourceTasks` 指向其依据的 design 任务，`epic` 指向根 analysis 任务
4. 每个任务耗时控制在 `time_estimate` 范围内，按依赖关系排序
5. 通过 `task-add.js` 逐个添加：

```bash
node scripts/task-add.js <projectRoot> <CR-id> "<标题>" "<描述>" <依赖任务id> <依赖任务id> ...
```

#### 回退细化（从 implement 回退）

部分 implementation 任务可能已执行（status 为 done 或 in_progress）。根据修改后的设计方案，调整任务列表：

1. 检查当前 task-map.yaml 中已完成的任务是否仍然有效。已完成的任务不应回退为 pending——它们已完成，修改方案应基于已有成果继续推进。

2. 标记需要废弃的任务（方案调整后不再需要的任务）：

```bash
node scripts/task-cancel.js <projectRoot> <CR-id> <taskId> "<废弃原因>"
```

3. 新增实施任务，每个单独创建，依赖字段关联到已有任务：

```bash
node scripts/task-add.js <projectRoot> <CR-id> "<标题>" "<描述>" <依赖任务id> <依赖任务id> ...
```

新增任务的 `sourceTasks` 指向对应的 design 任务，`epic` 继承原任务的 epic 归属。
依赖关系应考虑：被替代任务的原有依赖需继承、新方案引入的前置任务需声明。

4. 完成后通过 `flowforge-progress` 更新状态摘要。

---

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `scripts/design-context.js` | 加载 projects 列表、intake、naming、当前 proposal 状态 |
| `scripts/task-create.js <root> <id>` | 首次创建：初始化 task-map.yaml 中的任务到存储层 |
| `scripts/task-add.js <root> <id> <title> <desc> [depIds...]` | 增量添加：创建单个 analysis/design/implementation 任务，支持指定依赖 |
| `scripts/task-cancel.js <root> <id> <taskId> [reason]` | 回退修改：废弃不再需要的任务 |

## 引用的 SKILL

| SKILL | 引用场景 |
|-------|---------|
| `flowforge-docs` | 创建文档时获取 frontmatter 契约、meta.yaml 字段约束、task-map 格式 |
