# FlowForge 卡片模型 v3：阶段演进替代类型拆分

> 日期：2026-07-09
>
> 定义 FEATURE 卡片模型、生命周期、模板、拆分策略、横切类型和 PROP 提案全景。
> CLI 改造方案见 `proposal-cli-spec-v3.md`，SKILL 方法论见 `proposal-skill-spec-v3.md`。

---

## 1. 问题诊断

基于 CR26063001 提案（44 张卡片，28 REQ，0 DESIGN，1826 行）的实证分析和 CR26070801 的真实设计过程，
确认 FlowForge 卡片系统存在以下五个根本问题：

### 1.1 信息分散而非递进

REQ、DES、TASK 是同一功能的三个抽象层级，但被定义为三种独立的卡片类型：

```
STR → REQ (要做什么) → DES (怎么做) → TASK (执行步骤)
```

理解一个功能需要打开 4 张卡片，每张卡片都要重复"这是什么功能"的上下文。这不是原子性，这是拆碎了。

### 1.2 TASK 卡片沦为"指向器"

当前 TASK 模板仅要求 Goal、Deliverables、Acceptance 等元信息，不要求包含文件路径、方法签名、边界条件。
结果 TASK 变成 `"Read Before Work: 参考 DES-xxx"` —— 一个指向其他卡片的指针，不提供任何实现级别的信息。

### 1.3 STR 卡片是目录而非合成

Zettelkasten 的 Structure Note 应该回答 "what do I know about X?"，FlowForge 的 STR 回答的是
"which cards have I created about X?"。11/14 个 EPIC STR 是相同空壳，不含合成段落。

### 1.4 模板负担：73% 结构开销

典型 REQ 卡（45 行）：frontmatter 15 行 + 模板结构 13 行 + 自动导航 5 行 = 33 行结构开销，
只有 12 行业务内容（27%）。

### 1.5 流程是建议性的，不是强制性的

设计 SKILL 定义的 `index → clarify → analyze → discover → design → split tasks` 流程
是 "Use When" 建议表，无阻塞机制。Agent 可以在 index 模式创建 28 张 REQ 后直接结束。

---

## 2. 为什么补丁式修复不够

`remediation-card-fragmentation.md` 的修复方案已在代码中全部实现：
密度检查、合成检查、跨卡链接检查、design_gap 检测、Mode Gating Rules、SKILL 硬规则。

但这些修复的共同前提是 **"10 种卡片类型的拆分模型是正确的，只是执行不到位"**。
真实数据（CR26070801）表明，执行层面之外存在结构性缺陷：

| 现象 | 根本原因（非执行问题） |
|------|----------------------|
| REQ 的 Acceptance 和 DES 的 Constraints 表述几乎一致 | REQ 和 DES 不是两类独立信息，是同一信息的两个抽象层级 |
| TASK 的 Goal 就是 REQ 的 Summary 换了动词 | TASK 没有在 REQ/DES 的基础上增加实现级别的信息 |
| 同一决策在 5 类卡片中出现 5 次 | 按类型拆分导致信息在层级间重复，而非递进 |
| 完整理解一个功能需要跨 6+ 卡片跳转 | "原子性"被误解为"机械拆分"而非"思想完整性" |

**核心问题不是卡片太多，而是类型拆分本身在制造碎片。** 当 WHAT、HOW、DO 被定义为三种独立卡片类型，
系统就隐含承诺了每种类型有足够的内容密度来证明其独立存在——但这个承诺在很多场景下是假的。

### 2.1 方法论范式转换声明

本方案代表一个**有意的范式转换**，与 v1/v2 基于 Zettelkasten 原子卡片网络的方法论有根本区别：

| 维度 | v1/v2 (卡片网络) | v3 (阶段化文档) |
|------|-----------------|----------------|
| **信息单元** | 原子卡片（一张卡 = 一个思想） | 结构化文档（一张卡 = 一个功能的完整生命周期） |
| **知识产生** | 卡片之间的链接和对话产生洞察 | 同一张卡内的阶段演进和设计推理产生洞察 |
| **导航方式** | 跨卡跳转（STR→REQ→DES→TASK） | 段落渐进填充 + CLI 自动聚合视图 |
| **不可变性** | 卡片创建后应保持稳定，通过链接扩展 | 卡片是活文档，随阶段演进持续修改 |
| **密度要求** | 每张卡自足但可以很薄 | 每个阶段有内容密度门控 |

Zettelkasten 的"原子性"在软件工程上下文中被证明会退化为碎片化——当理解的单元是一个功能时，
按 WHAT/HOW/DO 拆分为三张卡片产生的维护成本远超洞察收益。
v3 选择结构化文档 + 阶段演进 + CLI 门控的路径，承认软件功能设计的自然单元是一个**完整功能**，
而非按类型拆分的片段。

这不是否定 Zettelkasten——CONV/DEC/MOD/FIND 仍然遵循原子横切原则。
这是将"一张卡片"的定义从"一个思想片段"修正为软件工程语境下的"一个完整可交付功能"。

---

## 3. 新模型：阶段演进替代类型拆分

### 3.1 核心转变

```
当前：按内容类别拆分为不同类型的卡片
  STR (索引) → REQ (要做什么) → DES (怎么做) → TASK (执行)

新：按功能单元组织，同一张卡片随认知深入而演进
  FEATURE (draft → designed → planned → in_progress → done)
```

一张 FEATURE 卡片承载一个**用户可感知的功能单元**的完整生命周期，从模糊需求到具体实现计划。
它不是 REQ+DES+TASK 的简单合并——它要求内容在卡片**内部**递进，而不是在卡片**之间**拆分。

### 3.2 卡片类型体系

```
提案 (proposal)

  ├── PROP    提案全景——回答"这个提案要达成什么、由哪些 FEATURE 组成、它们如何协作"

  ├── FEATURE 一个功能的完整全景，随阶段演进

  └── 横切类型 跨多个 FEATURE 生效，独立于功能演进
        ├── CONV  编码约定——一条可执行的规则
        ├── DEC   架构决策——影响多个功能的技术选择
        ├── MOD   模块知识——一个模块的定位和职责
        └── FIND  探索发现——一个意外行为或认知
```

**类型从 10 种精简为 6 种（含 PROP）。**

### 3.3 被移除的类型

| 移除类型 | 去向 |
|---------|------|
| REQ | 合并到 FEATURE 的 `## Summary` + `## Motivation` |
| DES | 合并到 FEATURE 的 `## Design` + `## Constraints` |
| TASK | 替换为 FEATURE 的 `## Implementation Plan` |
| STR | 替换为 CLI 自动生成的 `proposal inspect` 聚合视图 |
| LOG | 替换为 FEATURE 的 `## History` 段落（CLI 追加） |
| ROOT | 合并到 PROP，PROP 本身承载提案全景 |

### 3.4 为什么保留横切类型

CONV、MOD、DEC、FIND 与 FEATURE 有本质区别：

- **FEATURE 是纵向的**：一个功能的完整生命周期，创建→演进→完成→归档
- **横切类型是横向的**：一条 CONV 约束 5 个 FEATURE，一个 MOD 描述模块定位不绑定任何单一功能

如果把它们也塞进 FEATURE——比如把 CONV 写入每个被约束的 FEATURE——就会产生维护灾难：
修改一条约定需要更新 5 个文件。

保留横切类型的原则是 **"跨功能生效才独立成卡"**。如果一个决策只影响一个功能，
它应该写入那个 FEATURE 的 Design 段落，而不是创建独立的 DEC 卡。

### 3.5 PROP 卡片的提案全景职责

在旧模型中，PROP 是几乎空白的根卡片（仅一行 Summary），STR 承担了导航职责但无合成。
新模型中 PROP 升级为**提案全景入口**——它回答"这个提案要达成什么、FEATURE 如何分工协作"。

```markdown
---
id: PROP-CR26070801
title: FileProcessor Clone 能力
type: proposal
status: active
links:
  - target: FEAT-001
    relation: indexes
  - target: FEAT-002
    relation: indexes
  - target: FEAT-003
    relation: indexes
created: ...
updated: ...
---

# FileProcessor Clone 能力

## Goal

<!-- 1-3 句话：这个提案完成后，用户能做什么以前不能做的事？ -->

## Feature Map

<!-- 每个 FEATURE 在这个提案中承担什么职责？它们如何分工协作？ -->
<!-- 不是链接列表，是语义描述 -->

| Feature | 职责 | 协作关系 |
|---------|------|---------|
| FEAT-001 Clone 后端 API | 提供 clone 服务接口和领域编排 | FEAT-002 和 FEAT-003 的基础依赖 |
| FEAT-002 子对象树复制 | 实现深层复制和规则重建逻辑 | 依赖 FEAT-001 的 Service 层 |
| FEAT-003 Clone 前端页面 | 用户触发 clone 操作的交互界面 | 依赖 FEAT-001 的 API 契约 |

## Architecture Overview

<!-- 3-8 行：FEATURE 之间的架构关系、共享的技术决策、关键约束 -->

## Key Constraints

<!-- 跨所有 FEATURE 的共同约束（来自业务规则、系统架构） -->
<!-- 单 FEATURE 的约束写在该 FEATURE 的 Constraints 段落 -->
```

**PROP 与 `proposal inspect` 的关系：**
- PROP 的 `## Feature Map` 是**人写的语义描述**——为什么这些 FEATURE 组成一个提案、它们如何协作
- `proposal inspect` 是 **CLI 自动生成的机械聚合**——状态进度、阻塞关系、统计数据
- 两者互补不重复：人写"为什么"，机器算"怎么样"

**PROP 不应在 `proposal create` 时就写完。** 它随 FEATURE 的创建和演进逐步填充——
创建第一个 FEATURE 后不急着写 Architecture Overview，等 Design 逐步明确后再补充。

### 3.6 Proposal 生命周期：状态驱动而非目录迁移

**问题：** v2 的 wiki 目录结构为 `01-active/` 和 `03-completed/`，归档时将 proposal 目录
物理移动到 `03-completed/`。开放 Agent 直接编辑卡片内容后，卡片 body 中的文件路径引用会因
目录迁移而失效。

**v3 方案：** 所有 proposal 创建后不再移动目录。生命周期通过 PROP 卡片的 `status` 字段管理。
CLI 通过状态过滤提供 active/completed/all 视图。

**wiki 目录结构变更：**

```
v2:                                 v3:
01-workspace/                       01-workspace/
├── 01-active/                      ├── CR26061201/          # status: active
│   └── CR26061201/                 ├── CR26060101/          # status: completed
├── 02-intake/                      └── CR26053001/          # status: completed
└── 03-completed/
    ├── CR26060101/

02-library/                         02-library/              # 不变
```

**`02-intake/` 的处理：** 该目录在 v2 中规划为"待处理需求入口"，但从未被实现。
v3 中移除——新需求直接创建 proposal，不经过中间状态。

**迁移方式：** `flowforge upgrade` 命令内置版本迁移逻辑。升级时自动检测
版本跨越边界，按顺序执行所需的数据迁移步骤。

```bash
flowforge upgrade   # 下载新版本 → 安装 → 检测 v2→v3 迁移 → 执行
```

用户无需感知迁移细节——`upgrade` 是唯一的升级入口。

**生命周期状态流：**

```
proposal create → status: active
    │
    ├── 所有 FEATURE done + 验证通过
    │   → proposal archive → status: completed
    │       → 提取可复用知识到 library（curate Mode B）
    │
    └── proposal delete → status: deleted（或物理删除）
```

**`proposal archive` 的新行为：**
1. 将 PROP 卡片的 `status` 从 `active` 改为 `completed`
2. 运行 `proposal inspect` 生成最终报告
3. 提示：`"使用 flowforge-curate 将可复用知识提取到 library"`
4. **不移动任何目录或文件**

**`proposal list` 增强：**
```bash
flowforge proposal list                    # 默认 active
flowforge proposal list --status all       # 全部
flowforge proposal list --status completed # 已完成
```

---

## 4. FEATURE 卡片模板与生命周期

### 4.1 阶段定义

```
draft ──────→ designed ──────→ planned ──────→ in_progress ──────→ done
  │               │                │                │
  WHAT only       WHAT + HOW      WHAT + HOW       WHAT + HOW
                                  + 实现计划        + 实现计划
                                                    + 部分完成
```

| 阶段 | 含义 | 必填段落 | 门控条件 |
|------|------|---------|---------|
| `draft` | 刚识别出的功能，只有基本描述 | Summary, Motivation, Open Questions | - |
| `designed` | 设计方案已明确 | + Design, Constraints | Design 至少 1 个关键决策 + 理由 + 1 个约束；Open Questions 已清零或标注为假设 |
| `planned` | 实现计划已拆解完毕 | + Implementation Plan | 每个步骤至少写明：文件路径、方法签名或伪代码、边界条件；禁止出现"参考其他卡片" |
| `in_progress` | 正在实现 | + History (进行中) | 至少 1 个 Implementation Plan 步骤标记为完成 |
| `done` | 全部完成并验证 | + History (完成记录) | 所有步骤完成；Verification 各验收项有对应验证结果 |

### 4.2 阶段演进规则

```
draft → designed：满足 designed 门控 → card evolve <id> --stage designed
designed → planned：满足 planned 门控 → card evolve <id> --stage planned
planned → in_progress：至少 1 个步骤开始执行 → card steps <id> --start 1
in_progress → done：所有步骤完成 → card evolve <id> --stage done
```

**不能跳过阶段。** `card evolve` 在执行状态变更前必须验证门控条件。

### 4.3 模板

```markdown
---
id: FEAT-<proposal>-<ts>
title: <功能名称>
type: feature
status: draft
importance: should
links:
  - target: PROP-<proposal-id>
    relation: belongs_to
  - target: CONV-xxx
    relation: constrains
  - target: MOD-xxx
    relation: references
  - target: DEC-xxx
    relation: references
  - target: FEAT-xxx
    relation: depends_on
created: ...
updated: ...
source: <proposal-id>
---

# <功能名称>

## Summary

<!-- 1-3 句话：这个功能要解决什么用户问题或实现什么能力 -->
<!-- 脱离 proposal 上下文仍能被理解。draft 阶段即可填写 -->

## Motivation

<!-- 为什么需要这个功能？不做的后果是什么？谁需要它？ -->

## Design

<!-- 设计方案。draft 阶段可为空或写 TBD -->
<!-- designed 阶段必须填写 -->

### Key Decisions

<!-- 采用的关键设计决策，每条附带理由 -->

### Architecture

<!-- 涉及的模块、类、接口、数据流 -->

### Alternatives Considered

<!-- 考虑过但未采用的方案及原因。没有则写 None -->

## Constraints

<!-- 必须遵守的约束（来源：CONV/MOD/DEC/业务规则） -->
<!-- 也包含明确不做的事情（Out of Scope） -->

## Implementation Plan

<!-- 卡片进入 planned 阶段后填充 -->
<!-- 每个 ### Step N: 是一个可独立提取的单元（用于分段读取控制 Token） -->

### Step N: <步骤目标>

<!-- step-status: not_started -->

- **Goal**: 本步骤交付的可验证结果
- **Files**: 创建或修改的文件路径（相对于项目根目录）
- **Approach**: 实现策略——关键方法签名、算法选择、数据结构、状态流转
- **Edge Cases**: 边界条件和处理方式（空数据、重复、并发、失败回滚）
- **Dependencies**: 依赖的其他 FEATURE 或本卡步骤（标注 FEAT ID 或 Step 编号，写明等待策略如 mock）
- **Parallel**: yes 或 no——是否可以与本 FEATURE 的其他步骤并行执行
- **Verification**: 本步骤如何验证（测试场景、关键断言）

## Verification

<!-- 整体验证方案：端到端测试场景、关键断言、验收检查项 -->

## History

<!-- CLI 通过 card log 命令自动追加的事件记录 -->
<!-- Agent 和开发者不手写此段落 -->
<!-- 格式：- <ISO时间> | <kind> | <摘要> -->

## Open Questions

<!-- draft 阶段可以有多个；designed 阶段应该清零或明确标注为假设 -->

## Dependencies

<!-- 本功能依赖的其他 FEATURE 卡片，及被哪些 FEATURE 依赖 -->
```

### 4.4 模板原则

| 原则 | 说明 |
|------|------|
| **自足性** | 任何段落的描述不得依赖"参考 XXX 卡片"——必要信息必须在当前卡片内 |
| **渐进密度** | draft 阶段只有 2-3 个段落有内容；designed 增加 Design/Constraints；planned 增加 Implementation Plan |
| **去模板化** | 段落是思考的脚手架，不是必须填满的表格。没有内容写 `None` |
| **禁止跨卡引用** | Implementation Plan 步骤不写"参考 DES-xxx"——设计决策已写入当前 Design 段落 |
| **步骤可分段提取** | 每个 `### Step N:` 是自足的独立单元；步骤状态用 HTML 注释标记（`<!-- step-status: ... -->`），CLI 解析 |

---

## 5. FEATURE 拆分策略

### 5.1 何时拆分

**原则：拆分只发生在功能单元无法被一个开发者在一个迭代内完成时。** 拆分是万不得已的手段，不是默认操作。

| 场景 | 应拆分 | 不应拆分 |
|------|--------|---------|
| Implementation Plan 超过 ~10 个步骤 | 是 | - |
| 存在明显独立的子功能，可独立交付和验证 | 是 | - |
| 不同子功能由不同开发者并行实现 | 是 | - |
| 不同子功能有不同的验收方或交付时间线 | 是 | - |
| 卡片只是"变长了" | - | 否（用 `--section` 分段读取） |
| Design 段落比较丰富 | - | 否（丰富是好事，不是拆分理由） |
| "感觉应该按类型拆分" | - | 否（这是旧的 REQ/DES 思维） |

### 5.2 拆分机制

```
父 FEATURE (epic 级)
  │
  ├── 子 FEATURE-1 (独立子功能)
  ├── 子 FEATURE-2 (独立子功能)
  └── 子 FEATURE-3 (独立子功能)
```

**拆分过程（由 `card split` 命令执行）：**

1. 识别父 FEATURE 中可独立交付的子能力
2. 为每个子能力创建子 FEATURE 卡片（draft 阶段）
3. 父 FEATURE 保留 Design/Constraints/Motivation，移除 Implementation Plan
4. 子 FEATURE 继承父 FEATURE 的 Constraints，加上自己的 Design 和 Implementation Plan
5. 链接关系：子 FEATURE `part_of → 父 FEATURE`，父 FEATURE `decomposes → 各子 FEATURE`

**父 FEATURE 的角色（容器 FEATURE）：** 父 FEATURE 是"合成点"——让读者在不进入子卡片的情况下理解这个功能家族的全貌。
它不是目录（那是旧 STR 的问题），而是**有实质内容的概览**。如果拆分后父 FEATURE 只剩链接列表，
说明拆错了——应合并回一张卡。

**容器 FEATURE 与叶子 FEATURE 的区分：**

拆分后存在两种 FEATURE 角色，各自有不同的生命周期：

| 属性 | 叶子 FEATURE | 容器 FEATURE |
|------|------------|------------|
| **有 Implementation Plan** | 是 | 否 |
| **可被 `card evolve`** | 是，正常演进 | 否，stage 由子 FEATURE 聚合决定 |
| **Stage 的计算方式** | 由 `card evolve` 显式升级 | 自动聚合：所有子 FEATURE done → done；任一 in_progress → in_progress；否则 maintained |
| **可被 `card steps`** | 是 | 否 |
| **可被 `context feature --step`** | 是 | 否（返回子 FEATURE 列表） |
| **可再次拆分** | 是（如果 Implementation Plan > 10 步） | 否 |
| **可被归档** | 是 | 当所有子 FEATURE archived 后 |

`card split` 执行后自动将原 FEATURE 标记为容器角色（在 frontmatter 中增加 `role: container`）。
容器 FEATURE 不参与 `card evolve` 门控验证，其 `## Sub-Features` 段落由 CLI 自动维护。

### 5.3 拆分反模式

| 反模式 | 为什么错 |
|--------|---------|
| "REQ 部分太长，拆成 REQ+DES" | 回到旧模型——按类型拆分而非按功能拆分 |
| "Implementation Plan 有 5 步，拆成 5 个子 FEATURE" | 5 个步骤可能是一个功能的不同阶段，拆分后太薄 |
| "前端和后端拆成两个 FEATURE" | 如果同一个功能的不同实现面，拆分会失去端到端可验证性 |
| 拆分后父 FEATURE 只剩链接列表 | 等于旧 STR——纯目录，无合成价值 |

---

## 6. 横切卡片类型

### 6.1 CONV（约定）

```markdown
---
id: CONV-<proposal>-<ts> 或 CONV-NNN (library 中为全局编号)
title: <规则名称>
type: convention
status: active
importance: must | should | may
links:
  - target: FEAT-xxx
    relation: constrains
---
# <规则名称>

## Rule
<!-- 这条约定要求什么？用 must / should / may 明确力度 -->

## Rationale
<!-- 为什么有这条约定？ -->

## Applies To
<!-- 适用场景——什么时候必须遵守？什么时候不适用？ -->

## Examples
<!-- 正例 + 反例 -->
```

### 6.2 DEC（决策）

保持 ADR 格式，与现行 DEC 模板一致。

### 6.3 MOD（模块知识）

保持现行 MOD 格式。增加"被哪些 FEATURE 引用"的反向链接视图（CLI 自动生成）。

### 6.4 FIND（发现）

保持现行 FIND 格式。增加"影响哪些 FEATURE"的显式链接——旧模型下 FIND 经常缺少这个链接导致发现被遗忘。
