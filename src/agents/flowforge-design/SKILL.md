---
name: flowforge-design
description: |
  FlowForge 设计与探索。当需要分析需求、探索未知领域、创建设计方案、撰写提案或拆分任务时激活。
  负责从需求到可执行方案的全部工作。探索和设计完全融合——分析过程中随时形成方案，方案深化时随时补充探索。
---

# FlowForge Design

你是 FlowForge 的设计引擎。你的工作是：**理解问题 → 探索+设计融合循环 → 产出 proposal 和 task-map**。

## 触发条件

- `flowforge-workflow` 路由到 `continue-proposal` 或 `new-requirement` 场景
- `flowforge-implement` 实施中发现设计缺陷，回退到设计阶段
- 用户明确要求"探索"、"分析"、"设计"、"创建提案"

## 工作流

```
定位上下文 → 理解问题 → [探索 ⇄ 设计] 循环 → 撰写 proposal → 拆分任务
```

---

### 阶段 1：定位上下文

运行 `scripts/design-context.js` 加载全部上下文。

**如果场景是 `new-requirement`**：
- 检查输出中是否有 intake 材料和 `intake.steps`
- 有则逐条执行步骤中的 `action`，将结论作为后续输入
- 无则进入阶段 2 向用户提问

**如果场景是 `continue-proposal`**：
- 上下文中已包含当前 proposal 的状态和已有文件
- 当前 proposal 已经有 `CR-id` 和目录结构
- 跳过阶段 2，直接进入阶段 3 继续设计

---

### 阶段 2：理解问题

仅在 `new-requirement` 场景且没有 intake 材料时执行。

向用户确认核心诉求、影响范围和已知约束。信息不足就提问，不跳过。

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

**如果是 `new-requirement` 且首次进入**：
1. 根据 `naming.exploration_slug` 生成 slug，创建 exploration 目录
2. 根据 `naming.proposal_id` 的模板生成 CR-id，创建 proposal 目录

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

将设计决策提炼为 `ff-wiki/workspace/proposals/<CR-id>/proposal.md`。

参照 `flowforge-docs` SKILL 获取 proposal 的写作指南。

同时参照 `flowforge-docs` 的 proposal meta 契约创建 `meta.yaml`。

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
| `scripts/design-context.js` | 加载 project rules、intake、naming、当前 proposal 状态 |

## 引用的 SKILL

| SKILL | 引用场景 |
|-------|---------|
| `flowforge-docs` | 创建文档时获取 frontmatter 契约、meta.yaml 字段约束 |

项目级策略（探索方法、设计章节、命名规则、任务粒度）均通过脚本从 `config.yaml` 加载，不在此 SKILL 硬编码。
