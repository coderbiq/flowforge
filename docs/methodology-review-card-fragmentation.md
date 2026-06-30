# FlowForge 设计方法论问题分析：卡片碎片化

> 日期：2026-06-30
>
> 基于提案 `CR26063001_customeragent-portal-sales-模块` 的实证分析，诊
> 断卡片笔记法在实践中的碎片化问题，并提出方法论层面的优化建议。

## 1. 问题背景

提案 `CR26063001_customeragent-portal-sales-模块` 的卡片结构在形式上符合
FlowForge 设计规范（3 层索引树、14 个 EPIC、28 张原子 REQ、类型化链接齐全），
但实际阅读体验极差——读者不停在卡片间跳转，每次进入一张卡片只有 3-4 行有用信
息，然后返回索引再看下一张。从卡片树结构看"合理"，从阅读者体验看"碎片化"。

本文档分析造成这个矛盾的根本原因，并提出方法论优化方向。

## 2. 数据画像

| 维度 | 数值 |
|------|------|
| 总卡片数 | 44 个 .md 文件 |
| REQ 卡片 | 28 张（每条 44-46 行，~900 bytes） |
| STR 卡片 | 15 张（1 根索引 + 14 EPIC 子索引） |
| LOG 卡片 | 1 张 |
| DESIGN 卡片 | **0 张** |
| 分析任务卡片 | **0 张** |
| ROOT 卡片 | **缺失** |
| 总行数 | 1,826 行 |
| 目录深度 | 3 层（根 → 90-cards/ → 43 个 .md 文件） |

## 3. 碎片化的五个根本原因

### 3.1 模板负担危机：73% 结构开销

以一张典型 REQ 卡（45 行）为例：

```
┌─────────────────────────────────────────────┐
│  frontmatter:  15 行  (33%)  ← 元数据开销    │
│  模板结构:     13 行  (29%)  ← # Summary 等  │
│  自动导航:      5 行  (11%)  ← Links/Outgoing │
│  ═══════════════════════════════════════════ │
│  实际业务内容:  12 行  (27%)  ← 只有这些！    │
└─────────────────────────────────────────────┘
```

真正表达业务语义的只有 4-5 行：
- **Summary**: 1 行 — "按配置控制产品在 Portal 是否展示"
- **Source**: 1 行 — "PRD item #1"
- **Acceptance**: 2 行 — 两条验收条件
- **Scope**: 1 行 — 边界描述

**根因**：设计 SKILL 的 `card-templates.md` 定义了 5 段式 REQ 模板作为"最低要
求"，但当需求点本身只有一句话时，强制执行这个模板导致每张卡片被模板"吃掉"
了 3/4 的空间。这违背了 Zettelkasten 的核心理念——**模板应该适配内容，而非
内容填充模板**。

### 3.2 合成层缺失：DESIGN 卡片的"设计鸿沟"

这是最关键的问题。当前卡片网络：

```
STR → EPIC STR → REQ  (只有需求层，没有设计层)
                     ↓
              REQ 是"什么需要做"
              但没有 DESIGN 来回答"怎么做"和"为什么这样做"
```

Zettelkasten 的核心价值在于**卡片之间的对话产生新洞察**。Structure Note 应
该回答 "what do I know about X?"，而非仅仅 "which cards exist about X?"。

当前的 EPIC STR 卡片：

```markdown
## Purpose
Portal 侧产品配置能力：产品可见性、展示配置...

## Entries
- [REQ-xxx] 产品可见性控制
- [REQ-xxx] 产品展示配置
```

这是**目录**，不是**合成**。读者打开了 14 个 EPIC 的 STR 卡，看到的全是"同上"
的结构——无法从任何一个 STR 卡中获得"这个 EPIC 要解决什么核心问题，各需求之
间如何协作，关键设计约束是什么"的理解。

**设计 SKILL 的缺陷**：walkthrough 示例中包含 DESIGN 卡片的创建步骤，但在实际
流程中 DESIGN 被描述为 "when a stable conclusion exists"，导致实践者倾向于先
"把需求都拆完再说"。但需求拆完后如果不立即合成，就变成了 28 个孤立的碎片。

### 3.3 PRD 编号驱动的机械原子化

28 张 REQ 卡与 PRD 项目编号 1:1 映射（PRD #1 ~ #30, #81）。这不是语义驱动
的拆分，而是编号驱动的拆分。

在原始 PRD 中，第 22 项 "Quotation 创建与编辑" 只是一句话的总结，但它隐含依
赖了：
- #17 数据模型（quotation 的结构是什么）
- #19 数据范围（能编辑哪些数据）
- #21 数据转换（数据如何转换）
- #23 保费试算（编辑完要做什么）

当这些被拆成独立 REQ 卡后，每张卡都声称自己是"自足的"（有 Summary +
Acceptance + Scope），但实际上非常依赖上下文。**问题是卡片之间没有显式的
`requires` 链接**——这些依赖关系存在于 STR 的链接列表中，而不在 REQ 卡自身
体内。

**设计 SKILL 的缺陷**：`card-templates.md` 的 REQ 模板没有 `## Dependencies`
或 `## See Also` 段落。模板假设每张 REQ 卡可以独立存在，但没有引导实践者标注
跨卡依赖。

### 3.4 STR 卡片违背 Structure Note 核心理念

Zettelkasten 文献明确定义：

> "A structure note answers 'what do I know about X?'; a hub note answers
> 'how is my knowledge of this whole domain organized?'"

当前 STR 卡实际回答的是：**"which cards have I created about X?"**

具体表现：
- 11/14 个 EPIC STR 卡片是**完全相同的 26 行空壳**：仅 YAML + 一行 Purpose +
  一个指向根索引的链接
- 活跃的 STR（EPIC 1/4/14）虽然包含了 Entries 列表，但也没有合成段落
- 所有 STR 卡都不包含 `## Key Decisions`、`## Design Rationale` 或
  `## Relationships` 等合成性内容

### 3.5 渐进式细化流程被跳过

设计 SKILL 定义的流程：

```
index → clarify → analyze → discover library → design → split tasks
```

实际执行：

```
index (创建 STR 树 + 从 PRD 导入 REQ) → 结束了
         ↑
    clarify / analyze / design 全部被跳过
```

LOG 卡片明确写道："下一步进入 clarify/design 模式：为 EPIC 1/4/14 中 active
requirement 卡创建对应 design 卡"——但这个"下一步"尚未执行，提案就以 28 张
REQ 卡 + 0 张 DESIGN 卡的状态呈现了出来。

**设计 SKILL 的缺陷**：流程被定义为"可回退的循环"，但没有检查点来确保
"index" 阶段之后必须进入 "design" 阶段才能被视为一个可交付的提案。

## 4. 方法论层面的根因

### 根因 A：Zettelkasten "原子性"原则被误用

> Zettelkasten 的原子性是 "one complete thought"，实践变成了 "one small
> fragment"

| Zettelkasten 原文 | FlowForge 实际解释 |
|---|---|
| "每张卡片只包含一个**完整**想法" | 每张卡片只包含一个 PRD 条目编号 |
| "脱离上下文仍能被理解" | 脱离上下文只有 1 行 Summary |
| "可以被其他卡片引用、嵌入到任何上下文" | 可以被 STR 索引，但不能独立被理解 |

真正的问题不是卡片太多，而是卡片**信息密度太低**。如果一个 REQ 只有 4 行有
效内容，它不应该是一张独立卡片——要么合并到更粗粒度的卡片中，要么等待更充
分的分析后再决定是否拆分。

### 根因 B：模板被当作"必须完成的表单"而非"思考的脚手架"

`card-templates.md` 定义了最小模板结构，但设计 SKILL 没有说明：

1. **何时模板不适用**：如果 Summary 只有 1 行、Acceptance 只有 2 条，说明这
   个需求点可能不需要独立成卡
2. **何时应该合并**：如果 3-5 个 REQ 卡共同描述一个完整功能，它们可能应该是
   一张更丰富的 REQ 卡 + 补充性子卡片
3. **密度阈值**：卡片应该有一个"足够丰富"的下限判断

### 根因 C：流程是建议性的，不是强制性的

设计 SKILL 的 7 步流程每一步都是 "use when" 的建议，没有阻塞机制。一个 Agent
可以在 "index" 模式下创建全部 28 张 REQ 卡而不触发任何违规——因为
`workflow-rules.md` 中的 Mode Selection 表只是描述了"什么情况下用什么模式"，
没有定义"什么情况下必须停止当前模式进入下一个模式"。

## 5. 优化建议

### 建议 1：引入"内容密度"概念

在 `card-templates.md` 或 `workflow-rules.md` 中增加：

```markdown
## 内容密度检查

卡片创建后，应在下一个 turn 检查内容密度：

| 密度等级 | 判断标准 | 建议动作 |
|---------|---------|---------|
| **过薄** | 有效内容 < 5 行，Summary 仅 1 句 | 合并到父卡或关联卡，不独立成卡 |
| **合适** | 有效内容 5-20 行，每段有实质信息 | 保持独立 |
| **过厚** | 单段超过 15 行，或总内容 > 50 行 | 考虑拆分 |

注："有效内容"指扣除 frontmatter、模板标题、自动导航后的业务内容。
```

### 建议 2：STR 卡片必须包含合成段落

修改 STR 卡片的模板结构：

```markdown
# EPIC N - Title

## Purpose
（一句话：这个 EPIC 解决什么核心问题）

## Synthesis
（3-8 行：这个 EPIC 下的需求如何协作？
  关键的设计约束是什么？
  与哪些其他 EPIC 有强依赖？）

## Key Decisions
- 关键决策 1
- 关键决策 2

## Entries
- [REQ-xxx] ...
```

STR 卡不能只是链接列表。它必须提供合成信息——这是读者理解整个 EPIC 的核心
入口。

### 建议 3：DESIGN 卡片不应是可选的

将 DESIGN 卡片的创建从 "when a stable conclusion exists" 改为强制规则：

```markdown
## 设计阶段强制规则

当 proposal 中存在 ≥ 3 张 active REQ 卡时，必须至少创建 1 张 DESIGN 卡，
合成这些 REQ 卡的设计意图、架构决策和关键约束。

当 proposal 中存在 ≥ 10 张 active REQ 卡但没有 DESIGN 卡时，
proposal inspect 应报告为 "design_gap" 健康问题。
```

### 建议 4：增加"跨卡依赖"模板段落

在 REQ 模板中增加：

```markdown
## Dependencies
- 依赖的其他 REQ 卡及其原因
- 依赖的外部系统或模块
- 被哪些 REQ 卡依赖

## See Also
- 相关的 DESIGN 或 DEC 卡
```

### 建议 5：区分"播种"与"成熟"提案状态

为 proposal 增加健康检查规则：

```markdown
## Proposal 成熟度检查

proposal inspect 应报告以下健康问题：

| 问题类型 | 触发条件 |
|---------|---------|
| `missing_root` | 缺少 ROOT 卡 |
| `design_gap` | active REQ ≥ 3 但 DESIGN = 0 |
| `str_too_shallow` | STR 卡的 "Synthesis" 段为空或 < 2 行 |
| `card_too_thin` | REQ 卡有效内容 < 5 行 |
| `missing_cross_links` | REQ 卡只有 belongs_to + indexes 链接，缺少功能级链接 |
```

### 建议 6：渐进密度的创建策略

将 index 模式改为两阶段：

```
Phase 1 - Coarse Seeding（粗粒度播种）
  → 创建 EPIC 级别的 REQ 卡（每个 EPIC 1-3 张），内容更丰富
  → 在 EPIC STR 中写入合成段落

Phase 2 - Progressive Refinement（渐进细化）
  → 当某张 REQ 卡的内容增长到 30+ 行时，才对其进行拆分
  → 拆分后的子卡必须标注与父卡的双向依赖
```

读者的初始体验变为：打开 EPIC STR → 看到有深度的合成 → 进入 1-3 张内容丰富
的 REQ 卡 → 如果某张太厚，它已被拆分为子卡。而非当前：打开 STR → 打开
REQ #1（4 行）→ 返回 → 打开 REQ #2（4 行）→ 返回 → ...×28。

## 6. 总结

这个提案的问题是 FlowForge 设计方法论中**模板工程化过早**和**合成机制缺失**
的典型案例。表面上结构完美（3 层索引树、14 个 EPIC、28 张原子 REQ、类型化链
接齐全），但实质上：

1. **模板吞噬内容** — 73% 的行是结构开销
2. **STR 只是目录** — 没有合成，没有洞察
3. **DESIGN 层缺失** — 需求碎片没有聚合点
4. **机械原子化** — 按 PRD 编号拆分，而非按语义拆分
5. **流程跳过** — index 之后直接结束，clarify/analyze/design 全被跳过

核心矛盾在于：**FlowForge 借鉴了 Zettelkasten 的"形式"（原子卡片、类型化链
接、Structure Note），但没有充分内化其"实质"（卡片间的对话产生洞察、
Structure Note 是合成而非目录、原子性指的是思想的完整性而非文本的短小）。**

修复方向不是减少卡片数量，而是在卡片**内部增加密度**、在 STR 层面**增加合
成**、在流程层面**强制关键步骤不可跳过**。让读者在任何一张卡片中停留时都能
获得足够的信息，而不是被迫在碎片之间反复跳转。

## 7. 参考文档

- [Zettelkasten 卡片笔记法](references/zettelkasten.md)
- [知识卡片系统设计](knowledge-system.md)
- [卡片架构不变量与修正方案](card-architecture-invariants.md)
- [flowforge-design 正式草案](flowforge-design-skill.md)
- [Card Templates](assets/skills/flowforge-design/references/card-templates.md)
