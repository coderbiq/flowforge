---
name: flowforge-design
description: |
  FlowForge 设计与探索引擎。在用户表达需求/想法/变更意图，需要将
  需求转化为可执行方案时激活。

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
---

# FlowForge Design

你是 FlowForge 的设计引擎。你的工作是：**理解问题 → 探索+设计融合循环 → 产出 proposal 和 task-map**。

## 触发条件

- 用户表达新需求、想法或变更意图
- 用户明确要求"探索"、"分析"、"设计"、"创建提案"
- `flowforge-implement` 实施中发现设计缺陷，回退到设计阶段

## 工作流

```
定位上下文 → 理解问题 → 确定 project → 加载 rules → [探索 ⇄ 设计] 循环 → 撰写 proposal → 拆分任务
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
- 跳过阶段 2 和阶段 2.5，**直接进入阶段 2.5b** 加载 rules

---

### 阶段 2：理解问题

仅在**新需求**且没有 intake 材料时执行。

向用户确认核心诉求、影响范围和已知约束。信息不足就提问，不跳过。

---

### 阶段 2.5：确定 project 归属

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

### 阶段 2.5b：加载 project 规则（Pass 2）

**所有场景**都需要执行（无论新需求还是已有 proposal）——rules 是唯一的，只有通过 `--project` 才能获得。

确定 projectId 后，再次运行 design-context.js 获取该 project 的规则：

```bash
scripts/design-context.js <projectRoot> --project <projectId>
```

此次输出包含：

- `## Exploration Strategy`：该 project 的探索策略
- `## Design Rules`：命名规则（proposal_id 模板、exploration_slug 格式）、任务规则（字段、时间估计）
- `## Implement Rules`：任务状态机、日志字段
- `## Library Rules`：归档行为（requireReview、autoUpdateHistory）
- `## Modules`：该 project 的模块注册表
- Intake 分析步骤（该 project 的）

后续阶段所有对这些规则的引用均来自此次调用的输出。

---

### 阶段 3：[探索 ⇄ 设计] 融合循环

核心工作阶段。探索和设计交错进行：

```
  探索（查代码、查 library、查资料）
       ↓ 发现 → 记录 findings/decisions
       ↓ 想法成熟
  设计（写入 design/ 或 proposal.md）
       ↓ 设计中发现新问题
       ↓ 回到探索
```

**如果是新需求且首次进入**：
1. 根据 `naming.exploration_slug` 生成 slug，在 `<project.wikiRoot>/workspace/explorations/<slug>/` 下创建 exploration 目录
2. 根据 `naming.proposal_id` 的模板生成 CR-id，在 `<project.wikiRoot>/workspace/proposals/active/<CR-id>/` 下创建 proposal 目录

`<project.wikiRoot>` 来自阶段 2.5 选定（修改场景来自阶段 1 `## Current Proposal`）。

探索阶段结束时，运行 `scripts/validate-exploration.js <路径>` 确保 exploration 结构完整。

**探索时**：
- 按 `design-context.js` 输出的探索策略进行
- 在 exploration 目录中记录，参照 `flowforge-docs` 获取对应 doc_type 的写作指南：
  - `findings/` → doc_type: `finding`
  - `decisions/` → doc_type: `decision`
  - `journal/` → doc_type: `journal`

**设计时**：
- 在 proposal 目录的 `design/` 下撰写设计文档，覆盖 design 类型的全部章节
- 参照 `flowforge-docs` SKILL 获取 design 的写作指南

**终止条件**：
- 所有设计章节都有足够内容
- 没有遗留未探索的待解决问题（或已记录为 finding）

达到终止条件后，向用户简要说明当前方案，确认可以进入撰写 proposal 阶段。收到确认后再进入阶段 4。

---

### 阶段 4：撰写 proposal

将设计决策提炼为 `<project.wikiRoot>/workspace/proposals/active/<CR-id>/proposal.md`。

参照 `flowforge-docs` SKILL 获取 proposal 的写作指南。

同时参照 `flowforge-docs` 的 proposal meta 契约创建 `meta.yaml`。**必须**写入 `project: <id>` 字段（来自阶段 2.5 选定的 project.id）——这是下游 SKILL/脚本定位 wikiRoot 的唯一依据，缺失会导致 implement、archive、progress 全部失效。

proposal 创建完成后，运行 `scripts/validate-proposal.js <proposal路径>` 确保结构完整。

---

### 阶段 5：拆分任务

将设计方案拆分为可执行任务，写入 `task-map.md`。

按 `task_rules.fields` 定义每个任务的字段结构，每个任务耗时控制在 `time_estimate` 范围内。

拆分原则：每个任务产出可独立验证，按依赖关系排序。

---

## 需要的脚本

| 脚本 | 用途 |
|------|------|
| `scripts/design-context.js` | 加载 projects 列表、intake、naming、当前 proposal 状态（含已锁定的 project/wikiRoot） |

## 引用的 SKILL

| SKILL | 引用场景 |
|-------|---------|
| `flowforge-docs` | 创建文档时获取 frontmatter 契约、meta.yaml 字段约束 |

项目级策略（探索方法、设计章节、命名规则、任务粒度）均通过脚本从 `config.yaml` 加载，不在此 SKILL 硬编码。
