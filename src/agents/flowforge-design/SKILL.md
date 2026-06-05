---
name: flowforge-design
description: |
  FlowForge 需求树驱动渐进式探索设计引擎。在用户表达需求/想法/变更意图，
  需要将需求转化为可执行方案时激活。以需求树为地图、任务为引擎，
  需求树在探索中逐渐长全——初始草案逐步完善为完整的需求全景。

  必须在以下场景激活：
  - 用户描述新需求、想法或变更意图（"我想做..."、"需要..."、"打算..."）
  - 用户明确表达"探索"、"分析"、"设计"、"创建提案"、"拆分任务"
  - 已有 active 状态 proposal 需要继续完善设计或撰写
  - flowforge-implement 在实施中发现设计缺陷需要回退

  不要在以下情况激活：
  - 用户只是询问已有设计内容（只读查询，不修改）
  - 用户在执行 proposal 中的任务（应交给 flowforge-implement）
  - 用户要求归档或沉淀（应交给 flowforge-archive）
  - 仅用于更新进度索引——那是 flowforge-progress 的职责
  - 实施中发现需要补充探索内容——由 flowforge-feedback 结构化捕获后路由，不要直接激活本 SKILL 写 findings
---

# FlowForge Design

负责将需求转化为可执行的 proposal 和任务。

## 工作流

```
  定位上下文 → 理解问题 → 确定 project → 加载 rules → 需求树驱动探索设计 → 撰写 proposal → 细化实施任务
```

---

### 阶段 1：定位上下文

运行 `flowforge design-context [CR-id]` 加载上下文。不指定 CR-id 时自动查找当前 active/draft 状态的 proposal；指定时加载目标 proposal 的上下文（用于跨 proposal 场景）。输出包含：

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
flowforge design-context <projectRoot> --project <projectId>
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

### 阶段 5：需求树驱动的渐进式探索设计

**核心思路**：以需求树为地图、任务为引擎，驱动整个探索与设计过程。需求树在分析中逐渐长全——初始只是草案，随着探索发现新需求点而不断完善。任务和需求树同步生长，树给任务提供结构，任务让树变得更完整。

#### 5.1 构建需求草案

进入阶段 5 后，**先不急于创建任务**——先快速草绘需求树，建立对提案需求的全景认知（不追求完整）：

1. 根据 `naming.proposal_id` 的模板生成 CR-id，在 `<project.wikiRoot>/workspace/proposals/active/<CR-id>/` 下创建 proposal 目录

2. 基于阶段 2-4 收集的信息（用户诉求、intake 材料、project 策略、library 知识），将提案涉及的需求按功能域/模块快速拆解为需求树草案，写入 `notes.md` 的 `## 需求树` 章节：

```
## 需求树

- 认证模块
  - 登录流程调整 [?]
  - Token 刷新机制变更
  - 权限模型扩展
- 会话管理
  - Redis 集群方案
  - 多设备登录策略 [?]
- 数据迁移
  - 旧版本兼容
```

`[?]` 标记当前不确定的需求点，待探索后确认或移除。

**原则**：
- **不求完整**——初始草案有遗漏是正常的，探索中逐步补充
- **叶子节点可分析**——每个叶子应是一个可独立分析和验证的需求单元
- **作为活地图使用**——探索中发现新需求点立即补充到树中，树和任务同步生长

#### 5.2 初始化任务空间

基于需求树草案创建初始任务：

1. 初始化 proposal 的任务空间：

```bash
flowforge task init --proposal <CR-id> "<proposal标题>"
```

 2. 为需求树中**已知的叶子节点**批量创建 analysis 任务：

```bash
flowforge task add-tasks --proposal <CR-id> '[
  {"title":"分析登录流程调整需求","type":"analysis"},
  {"title":"分析Token刷新机制变更需求","type":"analysis"},
  {"title":"分析Redis集群方案需求","type":"analysis"}
]'
```

每个 analysis 任务的 `title` 直接对应需求树叶子节点的描述。任务 ID 由后端自动生成（beads issue ID）。
#### 5.3 探索完善循环

进入持续循环——**执行任务、完善需求树、拆解新任务交替进行**：

```
┌───────────────────────────────────────────────────────────┐
│  选取一个 pending 的 analysis/design 任务                 │
│                     ↓                                     │
│  analysis 任务：探索代码/library → 记录发现               │
│        ↓ 发现新需求点                                     │
│  补充到需求树 + flowforge task add 创建新 analysis 任务  │
│        ↓ 该需求分析充分                                   │
│  flowforge task add 创建 design 任务（sourceTasks 指向分析任务）│
│        ↓ 需求树该节点确定（移除 [?]）                      │
│                     ↓                                     │
│  design 任务：撰写设计文档 → 标记 done                    │
│        ↓ 设计中发现新问题 → 补充需求树 + 新 analysis 任务  │
│                     ↓                                     │
│  当前任务 done → 回到选取下一个 pending 任务              │
└───────────────────────────────────────────────────────────┘
```

**需求树的维护**：
- 每次发现新需求点 → 立即更新 notes.md 中的需求树
- 分析确认某个 `[?]` 节点确实不需要 → 从树中移除
- 分析确认某个 `[?]` 节点确实需要 → 移除 `[?]`，创建对应任务

**选取任务的优先级**：
1. `pending` 的 analysis 任务优先于 design 任务（先分析再设计）
2. 同类任务中，先创建的先处理
3. 有 `dependencies` 的任务，依赖全部完成后才可选

**查询就绪任务**：

```bash
# 按类型查询就绪的分析任务
flowforge task ready --proposal <CR-id> --type analysis

# 就绪的设计任务
flowforge task ready --proposal <CR-id> --type design
```

**处理 analysis 任务**：

1. 认领任务：`flowforge task claim --proposal <CR-id> <taskId>`
2. 根据任务标题确定探索方向，按以下策略执行：
   - 如有 `## Exploration Strategy`，按其探索策略进行
   - 探索代码库、library 文档、外部资料
3. 记录发现（直接写入 library 对应路径）：
   - 系统架构事实 → `library/architecture/<topic>.md`
   - 模块设计事实 → `library/modules/<name>/`
   - 可复用决策 → `library/decisions/`
   - 可复用约定 → `library/conventions/<topic>.md`
4. 在 proposal 的 `notes.md` 中记录探索过程和发现
5. **每个发现携带 `domain` frontmatter**
6. 分析过程中发现新的子探索点 → 创建子 analysis 任务：

```bash
flowforge task add --proposal <CR-id> analysis "<分析子任务标题>" --desc "<描述>" --parent <父任务id>
```

7. 该模块分析充分 → 创建对应的 design 任务：

```bash
flowforge task add --proposal <CR-id> design "<设计任务标题>" --desc "<描述>" --dep <依赖的分析任务id>
```

**analysis 完成标准**（满足以下条件才标记 done）：
- ✅ 需求树中对应的 `[?]` 节点已确认或移除
- ✅ 探索发现已写入 library（architecture / modules / decisions / conventions）
- ✅ 所有子 analysis 任务已完成（没有 pending 的子任务）
- ✅ 该模块的 domain 归属已判定（scope + module + type）
- ✅ 没有遗留的开放问题（或已创建新的 analysis 任务）

8. 任务完成 → `flowforge task done --proposal <CR-id> <taskId> --summary "<完成摘要>"`

**处理 design 任务**：

1. 认领任务：`flowforge task claim --proposal <CR-id> <taskId>`
2. 如有 `## Design Strategy`，参照其项目级设计策略
3. 在 proposal 目录的 `design/` 下撰写设计文档：
   - 参照 `flowforge-docs` 获取 design 写作指南
   - 多模块 proposal 按模块用子目录组织
   - 跨模块架构设计 → `scope: system`，放在 `design/` 根目录
4. 设计中发现新的未探索问题 → 创建新的 analysis 任务：

```bash
flowforge task add --proposal <CR-id> analysis "<分析子任务标题>" --desc "<描述>" --dep <发现问题的design任务id>
```

**design 完成标准**（满足以下条件才标记 done）：
- ✅ 设计文档已写入 design/ 目录（含 frontmatter + domain）
- ✅ 设计文档通过 `flowforge validate-doc` 校验
- ✅ 设计方案覆盖了对应 analysis 任务的所有发现
- ✅ 接口/架构/数据模型等关键决策已在文档中记录

5. 设计文档完成 → `flowforge task done --proposal <CR-id> <taskId> --summary "<完成摘要>"`

**判定 domain 的方法**（同之前，不变）：
```
scope:  文档引用的源文件在哪个模块下？设计最终在哪个模块落地？
        单模块边界内 → module；跨模块 / 全局架构 → system
module: scope=module 时，模块名是什么（如 auth、payment）
type:   架构/接口/数据模型/技术选型 → design
        关键决策+理由+备选方案评估 → decision
        编码规范/命名约定/反例 → convention
```

#### 5.4 终止条件

阶段 5 完成需同时满足：

1. 需求树中不再有 `[?]` 标记的节点（所有不确定性已消除）
2. 需求树的每个叶子节点都有对应的 analysis 任务且状态为 `done`
3. 所有 `design` 类型任务状态为 `done`
4. 没有未解决的开放问题（或已记录为 finding）

验证方式：

```bash
# 检查分析设计阶段是否完成
flowforge task status --proposal <CR-id> --type analysis  # analysis 全部 done？
flowforge task status --proposal <CR-id> --type design    # design 全部 done？
```

达到终止条件后，向用户简要说明当前方案，确认可以进入撰写 proposal 阶段。

---

### 阶段 6：撰写 proposal

将设计决策提炼为 `<project.wikiRoot>/workspace/proposals/active/<CR-id>/proposal.md`。

任务已在阶段 5 初始化并包含完整的 analysis/design 任务追踪记录。proposal.md 应引用已完成的分析和设计任务，确保方案可追溯到具体的分析发现。

参照 `flowforge-docs` SKILL 获取 proposal 的写作指南。

同时参照 `flowforge-docs` 的 proposal meta 契约创建 `meta.yaml`。**必须**写入以下字段：
- `project: <id>`（来自阶段 3 选定的 project.id）——这是下游 SKILL/脚本定位 wikiRoot 的唯一依据
- `modules: [<name>, ...]`（轻量列表，仅包含涉及的模块名，用于 INDEX.md 展示）

`meta.yaml` 中**不再需要** `ownership` 和 `archive_targets`——归档路径由各文档的 `domain` frontmatter 自动推导。

proposal 创建完成后，运行 `flowforge validate-proposal <proposal路径>` 确保结构完整。

---

### 阶段 7：细化实施任务

已完成 analysis 和 design 任务（阶段 5 创建并完成）。本阶段将已完成的设计转化为可执行的 implementation 任务。

#### 任务层级约束

Proposal 的任务空间遵循 4 层结构（详见 `guides/task-hierarchy.md`）：

```
Main Epic → Type Sub-Epic (分析/设计/实施) → Task → Child Task
```

- 独立小任务直接挂在类型子 epic 下（3 层）
- 需要多步骤的大任务拆为父子任务（4 层），不超过 4 层
- 子任务通过 `--parent <parentTaskId>` 挂载

#### 首次创建实施任务

1. 基于已完成的 design 任务，将设计方案拆分为可执行的 implementation 任务
2. 大任务（预计涉及多个文件/步骤）创建为父任务，再拆解子任务：

```bash
# 创建父任务（挂在实施子 epic 下）
flowforge task add --proposal <CR-id> implementation "实现核心配置链路" --desc "DDL、DO/PO、Mapper、Repository、DomainService、Application、API"

# 创建子任务（挂在父任务下）
flowforge task add --proposal <CR-id> implementation "DDL + Liquibase" --parent <父任务id>
flowforge task add --proposal <CR-id> implementation "DO/PO + Mapper" --parent <父任务id>
flowforge task add --proposal <CR-id> implementation "Repository + DomainService" --parent <父任务id>
```

3. 每个任务耗时控制在合理范围，按依赖关系排序

#### 回退细化（从 implement 回退）

部分 implementation 任务可能已执行。根据修改后的设计方案，调整任务列表：

1. 检查已完成的任务是否仍然有效。已完成的任务不应回退——修改方案应基于已有成果继续推进
2. 废弃不再需要的任务：

```bash
flowforge task cancel --proposal <CR-id> <taskId> --reason "<废弃原因>"
```

3. 重开需要修改的已完成任务：

```bash
flowforge task reopen --proposal <CR-id> <taskId>
```

4. 新增实施任务：

```bash
```

5. 完成后通过 `flowforge-progress` 更新状态摘要。

---

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `flowforge design-context` | 加载 projects 列表、intake、naming、当前 proposal 状态 |
| `flowforge task init --proposal <id> <title>` | 初始化 proposal 任务空间 |
| `flowforge task add-tasks --proposal <id> '<json>'` | 批量创建初始任务 |
| `flowforge task add --proposal <id> <type> <title> [flags]` | 增量添加单个任务 |
| `flowforge task cancel --proposal <id> <taskId> [--reason "..."]` | 废弃不再需要的任务 |
| `flowforge task reopen --proposal <id> <taskId>` | 重开已完成任务（回退修改） |
| `flowforge task status --proposal <id>` | 查看任务状态 |

## 引用的 SKILL

| SKILL | 引用场景 |
|-------|---------|
| `flowforge-docs` | 创建文档时获取 frontmatter 契约、meta.yaml 字段约束、任务数据结构 |
