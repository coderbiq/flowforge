# 卡片模型简化方案：从"类型拆分"到"阶段演进"

> 日期：2026-07-08
>
> 基于 `methodology-review-card-fragmentation.md` 的诊断、`remediation-card-fragmentation.md` 的补丁方案，
> 以及 CR26063001、CR26070801 两个真实提案的实证分析，提出根本性的模型修正。

---

## 1. 背景：为什么补丁式修复不够

`remediation-card-fragmentation.md` 提出了密度检查、门控规则、合成段落等一系列修复。但这些修复的共同前提是
**"10 种卡片类型的拆分模型是正确的，只是执行不到位"**。

真实数据（CR26070801）表明执行层面之外存在结构性缺陷：

| 现象 | 根本原因 |
|------|---------|
| REQ 的 Acceptance 和 DES 的 Constraints 表述几乎一致 | REQ 和 DES 不是两类独立信息，是同一信息的两个抽象层级 |
| TASK 的 Goal 就是 REQ 的 Summary 换了动词 | TASK 没有在 REQ/DES 的基础上增加任何实现级别的信息 |
| 同一决策在 5 类卡片中出现 5 次 | 按类型拆分导致信息在层级间重复，而非递进 |
| 完整理解一个功能需要跨 6+ 卡片跳转 | "原子性"被误解为"机械拆分"而非"思想完整性" |

**核心问题不是卡片太多，而是类型拆分本身在制造碎片。** 当 WHAT、HOW、DO 被定义为三种独立卡片类型，
系统就隐含承诺了每种类型有足够的内容密度来证明其独立存在——但这个承诺在很多场景下是假的。

---

## 2. 新模型：阶段演进替代类型拆分

### 2.1 核心转变

```
当前：按内容类别拆分为不同类型的卡片
  STR (索引) → REQ (要做什么) → DES (怎么做) → TASK (执行)

新：按功能单元组织，同一张卡片随认知深入而演进
  FEATURE (draft → designed → planned → done)
```

一张 FEATURE 卡片承载一个**用户可感知的功能单元**的完整生命周期，从模糊需求到具体实现计划。
它不是 REQ+DES+TASK 的简单合并——它要求内容在卡片**内部**递进，而不是在卡片**之间**拆分。

### 2.2 卡片类型体系

```
提案 (proposal)

  ├── FEATURE   (一个功能的完整全景，随阶段演进)
  │
  └── 横切类型  (跨多个 FEATURE 生效，独立于功能演进)
        ├── CONV  编码约定——一条可执行的规则
        ├── DEC   架构决策——影响多个功能的技术选择
        ├── MOD   模块知识——一个模块的定位和职责
        └── FIND  探索发现——一个意外行为或认知
```

**类型从 10 种精简为 5 种。** 被移除的类型：

| 移除类型 | 去向 |
|---------|------|
| REQ | 合并到 FEATURE 的 `## Summary` + `## Motivation` |
| DES | 合并到 FEATURE 的 `## Design` + `## Constraints` |
| TASK | 替换为 FEATURE 的 `## Implementation Plan` |
| STR | 替换为 CLI 自动生成的 `proposal inspect` 聚合视图 |
| LOG | 替换为 FEATURE 的 `## History` 段落（CLI 追加） |
| ROOT | 保留但简化为 proposal 入口，不再承载卡片类型语义 |

### 2.3 为什么保留横切类型

CONV、MOD、DEC、FIND 与 FEATURE 有本质区别：

- **FEATURE 是纵向的**：一个功能的完整生命周期，创建→演进→完成→归档
- **横切类型是横向的**：一条 CONV 约束 5 个 FEATURE，一个 MOD 描述模块定位不绑定任何单一功能

如果把它们也塞进 FEATURE——比如把 CONV 写入每个被约束的 FEATURE——就会产生维护灾难：修改一条约定需要更新 5 个文件。

保留横切类型的原则是 **"跨功能生效才独立成卡"**。如果一个决策只影响一个功能，它应该写入那个 FEATURE 的 Design 段落，而不是创建独立的 DEC 卡。

---

## 3. FEATURE 卡片生命周期

### 3.1 阶段定义

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

### 3.2 阶段演进规则

```
draft → designed：满足 designed 门控 → card evolve <id> --stage designed
designed → planned：满足 planned 门控 → card evolve <id> --stage planned
planned → in_progress：至少 1 个步骤开始执行 → card steps <id> --start 1
in_progress → done：所有步骤完成 → card evolve <id> --stage done
```

**不能跳过阶段。** `card evolve` 在执行状态变更前必须验证门控条件。
如果 FEATURE 还在 `draft`，但已经有实现计划的想法，可以预先填充 Implementation Plan 但阶段不升级——等到 Design 段落也满足门控后才能一起升级为 `planned`。

---

## 4. FEATURE 卡片模板

```markdown
---
id: FEAT-<proposal>-<ts>
title: <功能名称>
type: feature
status: draft
importance: should
links:
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
<!-- draft 阶段即可填写 -->

## Design

<!-- 设计方案。draft 阶段可为空或写 TBD -->
<!-- designed 阶段必须填写。包含： -->

### Key Decisions

<!-- 采用的关键设计决策，每条附带理由 -->
<!-- 示例： -->
<!-- 1. 使用独立的 FpFileProcessorCloneCmd 而非复用 CreateCmd — 理由：clone 语义不同，复用会造成字段爆炸和校验冲突 -->

### Architecture

<!-- 涉及的模块、类、接口、数据流 -->
<!-- 示例： -->
<!-- - 新增 FileProcessorCloneService（domain 层），负责根对象克隆、子树复制、rule 克隆的事务编排 -->
<!-- - FpFileProcessorMgmtAPI 新增 POST /{fileProcessorId}/clone -->

### Alternatives Considered

<!-- 考虑过但未采用的方案及原因 -->
<!-- 没有则写 None -->

## Constraints

<!-- 必须遵守的约束 -->
<!-- - 来自 CONV 卡片（标注 CONV ID） -->
<!-- - 来自 MOD 模块边界 -->
<!-- - 来自 DEC 架构决策 -->
<!-- - 来自业务规则 -->
<!-- 也包含明确不做的事情（Out of Scope） -->

## Implementation Plan

<!-- 卡片进入 planned 阶段后填充。替代当前 TASK 卡片 -->
<!-- 每个步骤包含足够的实现细节，开发者不需要跳转其他卡片 -->

### Step N: <步骤目标>

- **Goal**: 本步骤交付的可验证结果
- **Files**: 创建或修改的文件路径（相对于项目根目录）
- **Approach**: 实现策略——关键方法签名、算法选择、数据结构、状态流转
- **Edge Cases**: 边界条件和处理方式（空数据、重复、并发、失败回滚）
- **Dependencies**: 依赖哪些其他 FEATURE 的完成（标注 FEAT ID 和原因）
- **Verification**: 本步骤如何验证（测试场景、关键断言）

## Verification

<!-- 如何验证这个功能整体正确实现 -->
<!-- 包含：端到端测试场景、关键断言、验收检查项 -->
<!-- designed 阶段可以写草案，planned 阶段应该完善 -->

## History

<!-- CLI 通过 card log 命令自动追加的事件记录 -->
<!-- Agent 和开发者不手写此段落 -->
<!-- 格式：- <ISO时间> | <kind> | <摘要> -->

## Open Questions

<!-- draft 阶段可以有多个；designed 阶段应该清零或明确标注为假设 -->
<!-- 每个问题标注：阻塞对象、影响范围、建议的解决方式 -->

## Dependencies

<!-- 本功能依赖的其他 FEATURE 卡片，及被哪些 FEATURE 依赖 -->
<!-- drafted 阶段可以写初步判断；planned 阶段应该准确 -->
```

### 模板原则

| 原则 | 说明 |
|------|------|
| **自足性** | 任何段落的描述不得依赖"参考 XXX 卡片"——必要信息必须在当前卡片内 |
| **渐进密度** | draft 阶段只有 2-3 个段落有内容，designed 阶段增加 Design/Constraints，planned 阶段增加 Implementation Plan |
| **去模板化** | 段落是思考的脚手架，不是必须填满的表格。如果某个段落确实没有内容（如 Open Questions 已清零），写 `None` 即可 |
| **禁止跨卡引用** | Implementation Plan 的步骤不能写"参考 DES-xxx"——因为 DES 的设计决策已经写入当前卡片的 Design 段落 |

---

## 5. FEATURE 拆分策略

### 5.1 设计原则

FEATURE 拆分的核心张力：我们希望避免 REQ→DES→TASK 的碎片化，但 FEATURE 本身在分析设计过程中
可能发现粒度太大需要拆分。如何拆分而不重蹈覆辙？

**原则：拆分只发生在功能单元无法被一个开发者在一个迭代内完成时。** 拆分是万不得已的手段，
不是默认操作。

| 场景 | 应拆分 | 不应拆分 |
|------|--------|---------|
| Implementation Plan 超过 ~10 个步骤 | ✅ | - |
| 存在明显独立的子功能，可独立交付和验证 | ✅ | - |
| 不同子功能由不同开发者并行实现 | ✅ | - |
| 不同子功能有不同的验收方或交付时间线 | ✅ | - |
| 卡片只是"变长了" | - | ❌ 用 `--section` 分段读取 |
| Design 段落比较丰富 | - | ❌ 丰富是好事，不是拆分理由 |
| "感觉应该按类型拆分" | - | ❌ 这是旧的 REQ/DES 思维 |

### 5.2 拆分机制

```
父 FEATURE (epic 级)
  │
  ├── 子 FEATURE-1 (独立子功能, status: designed → planned → done)
  ├── 子 FEATURE-2 (独立子功能)
  └── 子 FEATURE-3 (独立子功能)
```

**拆分过程：**

1. 识别父 FEATURE 中可独立交付的子能力
2. 为每个子能力创建子 FEATURE 卡片
3. 父 FEATURE 保留：
   - `## Summary`：整体功能定位
   - `## Motivation`：为什么需要这个功能家族
   - `## Design`：子功能共享的架构决策和设计约束
   - `## Constraints`：子功能共同遵守的约束
   - **不再保留** `## Implementation Plan`——步骤下沉到子 FEATURE
4. 子 FEATURE 继承父 FEATURE 的 Constraints，加上自己的 Design 和 Implementation Plan
5. 链接关系：子 FEATURE `part_of → 父 FEATURE`

**父 FEATURE 的角色：**

父 FEATURE 是"合成点"——它让读者在不进入子卡片的情况下理解这个功能家族的全貌。
它不是目录（那是旧 STR 的问题），而是**有实质内容的概览**。

```markdown
# FileProcessor Clone 功能家族

## Summary
为 FileProcessor 提供完整的 clone 能力，包括后端 API、领域编排、子树深拷贝、
规则重建，以及前端 clone 页面。

## Motivation
当前 FileProcessor 不支持 clone，用户在创建相似配置时需要手动重新填写所有字段...
[此处应有 3-5 行实质性动机描述]

## Design

### Key Decisions
1. clone 使用独立的 Cmd + Service 模式（参考 GiPrdtScheme）
2. 整个 clone 过程保持单事务
3. rule 不复用 createRule，避免状态被强制改写

### Architecture
- 后端：FileProcessorCloneService（domain）+ clone API（adapter）
- 前端：独立 clone 页面，复用 FileProcessorBasicInfo 组件

## Constraints
- fileProcessorCode 由前端提供且必须唯一
- 不得修改源 FileProcessor 实例
- 子对象复制顺序：sheet → data field → field mapping → rule
- [CONV-xxx] clone 使用独立 cmd + service 模式

## Sub-Features
- [FEAT-001] Clone 后端 API 与领域编排 (depends_on: none)
- [FEAT-002] Clone 子对象树复制与规则重建 (depends_on: FEAT-001)
- [FEAT-003] Clone 前端实现 (depends_on: FEAT-001)

## Verification
- clone 后所有 id 重新生成，源对象不受影响
- 前后端联调通过
```

读者打开父 FEATURE 就能理解全貌，不需要打开子卡片。子卡片是给实施者用的。

### 5.3 拆分反模式

| 反模式 | 为什么错 |
|--------|---------|
| "REQ 部分太长，拆成 REQ+DES" | 这是回到旧模型——按类型拆分而非按功能拆分 |
| "Implementation Plan 有 5 步，拆成 5 个子 FEATURE" | 5 个步骤可能是一个功能的不同阶段，拆分后每张子卡太薄 |
| "前端和后端拆成两个 FEATURE" | 如果前后端是同一个功能的不同实现面，拆分会失去端到端可验证性 |
| 拆分后父 FEATURE 只剩链接列表 | 父 FEATURE 如果没有 Design/Constraints，就等于旧 STR——纯目录 |

---

## 6. FEATURE 依赖模型

### 6.1 依赖类型

一个 proposal 内的 FEATURE 卡片之间可能存在真实的功能依赖。这与旧模型中机械的 REQ→DES→TASK 链接
不同——FEATURE 间的依赖表达的是**功能交付的时序约束**。

| 关系 | 语义 | 示例 | Link 方向 |
|------|------|------|-----------|
| `depends_on` | B 未完成则 A 无法完成 | 前端 clone 页面 **depends_on** 后端 clone API | A → B (A 依赖 B) |
| `extends` | A 构建在 B 提供的能力之上 | Clone 功能 **extends** 现有 FileProcessor create 能力 | A → B |
| `conflicts_with` | A 和 B 的设计或需求冲突，需要协调 | - | A → B (双向) |

`depends_on` 和 `extends` 的区别：
- `depends_on`：B **必须先完成**，A 才能开始。强时序约束。
- `extends`：B 提供了某种基础能力，A **在此基础上增加**。没有强时序，但 A 需要理解 B 的接口契约。

### 6.2 依赖的识别时机

```text
draft 阶段：初步识别。在 Motivation 和 Summary 中提到"需要 XX 能力支撑"。
            Open Questions 中标注"依赖 XX 是否已存在？"

designed 阶段：确认依赖。在 Design 段落中明确"本功能调用 XX 接口，
              由 FEAT-XXX 提供"。在 frontmatter links 中写入 depends_on。

planned 阶段：精确依赖。Implementation Plan 的每个步骤标注
             "此步骤依赖 FEAT-XXX 的完成"。
```

### 6.3 依赖的验证

`proposal inspect` 应报告依赖健康状态：

```markdown
### Dependency Graph
FEAT-003 (Clone 前端) depends_on FEAT-001 (Clone API)
FEAT-002 (子对象复制) depends_on FEAT-001 (Clone API)

### Dependency Health
- ✅ FEAT-001 → done (满足 FEAT-002、FEAT-003 的依赖)
- ⚠️ FEAT-002 → in_progress (不阻塞，FEAT-003 只依赖 FEAT-001)
- ❌ 无循环依赖

### Blocked Features
- 无
```

### 6.4 依赖与 Implementation Plan 的关系

当 Implementation Plan 的某个步骤依赖其他 FEATURE 时：

```markdown
### Step 2: 实现 clone API 路由
- **Goal**: 新增 POST /file-processors/mgmt/{fileProcessorId}/clone
- **Dependencies**: 
  - FEAT-002 (FileProcessorCloneService) — Step 2 完成后才能调用 clone 方法
  - **等待策略**: Step 2 可以先实现路由骨架和参数校验，Service 接口可用 mock
```

这比旧 TASK 的 "参考 DES-xxx" 有实质信息：它不仅说了依赖什么，还说了**如何解耦等待**（mock、骨架、接口契约）。

---

## 7. 横切卡片类型

### 7.1 CONV (约定)

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

### 7.2 DEC (决策)

保持 ADR 格式，与现行 DEC 模板一致。

### 7.3 MOD (模块知识)

保持现行 MOD 格式。但增加 **"被哪些 FEATURE 引用"** 的反向链接视图（CLI 自动生成）。

### 7.4 FIND (发现)

保持现行 FIND 格式。但增加 **"影响哪些 FEATURE"** 的显式链接——旧模型下 FIND 经常缺少这个链接导致
发现被遗忘。

---

## 8. SKILL 改造

### 8.1 flowforge-design SKILL

**当前流程：**
```
index → clarify → analyze → discover library → design → split tasks
  创建STR+REQ     追问        分析TASK      library suggest   创建DES   创建TASK(指向器)
```

**新流程：**
```
seed → clarify → analyze → enrich → plan
 创建      追问      查看     填充     填充
FEATURE   用户     代码/    Design  Implementation
(draft)           library  段落    Plan 段落
```

每轮不是"创建新卡片"，而是**丰富现有 FEATURE 卡片的更多段落**。唯一创建新卡片的时机是：
1. 识别出新的独立功能 → 创建新的 FEATURE
2. 发现需要拆分的过大 FEATURE → 创建子 FEATURE
3. 发现横切知识 → 创建 CONV/DEC/MOD/FIND

### 8.2 flowforge-design 强制门控

```markdown
## Stage Gating Rules (Hard Rules)

### Seed → Clarify
- 新建的 FEATURE (draft) 至少需要 Summary + Motivation
- 如果 Summary 只有 1 句话且 Motivation 为空 → 追问用户补充，不继续创建更多卡片
- 禁止单轮创建 > 5 张 draft FEATURE（防止机械导入 PRD 条目）

### Clarify → Design (draft → designed)
- 所有 Open Questions 必须已清零或明确标注为假设
- 如果标注为假设，Design 段落必须说明"以下设计基于 X 假设，如果假设不成立需要重新设计"
- Design 的 Key Decisions 至少包含 1 个决策 + 理由
- Constraints 至少包含 1 条约束（可来自 library 查询或用户确认）
- 如果 Design 涉及的决策影响多个 FEATURE → 先创建 DEC 卡，再在各 FEATURE 中 references DEC

### Design → Plan (designed → planned)
- Implementation Plan 至少包含 1 个步骤
- 每个步骤必须写明：Files (文件路径)、Approach (方法签名或伪代码)、Edge Cases (至少 1 个边界条件)
- 步骤中禁止出现"参考 DES-xxx"或"参考 REQ-xxx"——所有信息必须自足
- 如果某个步骤的具体实现依赖尚未确认的设计假设 → 该步骤标注为 [blocked]，该 FEATURE 不能升级为 planned
- 所有 Open Questions 必须清空

### 密度门控（适用于所有阶段）
- draft FEATURE 的有效内容必须 ≥ 5 行（Summary + Motivation 合计）
- designed FEATURE 的 Design 段落必须 ≥ 8 行有效内容
- planned FEATURE 的每个 Implementation Plan 步骤必须 ≥ 3 个字段有实质内容
```

### 8.3 flowforge-implement SKILL

**当前流程：**
```
读 TASK → 读 REQ (跳转) → 读 DES (跳转) → 读参考源码 → 实现 → 创建 LOG
```

**新流程：**
```
读 FEATURE → 理解全貌 → 按 Implementation Plan 步骤顺序执行 → CLI 追加 History → 更新步骤状态
```

不再需要跨卡跳转。所有上下文在一张 FEATURE 卡里。进展记录通过 CLI 追加而非创建 LOG 卡。

implement SKILL 硬规则：
```markdown
## Hard Rules
- Start by reading the FEATURE card with --section Implementation Plan
- Execute steps in order; do not skip blocked steps
- After each step, use card log to record progress
- Do NOT create separate LOG cards — all progress goes into FEATURE History
- If implementation reveals a missing detail in the plan, update the FEATURE card first, then continue
```

---

## 9. CLI 改造

### 9.1 新增/修改命令

| 命令 | 说明 |
|------|------|
| `card create --type feature --title "..." [--links ...]` | 创建 FEATURE 卡片 |
| `card evolve <id> --stage designed\|planned\|done` | 升级阶段，触发门控验证 |
| `card log <id> --event "..." [--kind progress\|bug\|blocked]` | 向 History 追加事件 |
| `card steps <id> --status done\|in_progress <step-number>` | 更新 Implementation Plan 步骤状态 |
| `card split <id> --titles "子功能1,子功能2"` | 拆分 FEATURE 为父子结构 |
| `proposal inspect <id> [--view dependency\|timeline]` | 聚合视图，自动生成 |
| `proposal context` | 上下文聚合，适配新卡片类型 |

### 9.2 card evolve 门控验证

`card evolve <id> --stage designed` 在执行前必须验证：

```
1. 卡片当前阶段 = draft
2. ## Design 段落存在且 Key Decisions 条目数 >= 1
3. ## Constraints 段落存在且条目数 >= 1
4. ## Open Questions 为 None 或所有条目以 "[假设]" 开头
5. 有效内容行数 >= 15 (Summary + Motivation + Design + Constraints)
```

任一条件不满足 → 拒绝升级，输出具体缺失项和建议命令。

### 9.3 card split 行为

```bash
flowforge card split FEAT-xxx --titles "Clone API","子对象复制","前端实现"
```

执行：
1. 验证父 FEATURE 处于 designed 或 planned 阶段
2. 验证子功能数量 >= 2
3. 创建子 FEATURE 卡片（draft），各带指定的 title
4. 子 FEATURE 自动添加 `part_of → 父 FEATURE`
5. 父 FEATURE 自动添加 `decomposes → 各子 FEATURE`
6. 父 FEATURE 移除 `## Implementation Plan`（下沉到子卡片）
7. 父 FEATURE 保留 Design/Constraints/Motivation，增加 `## Sub-Features` 链接列表
8. 输出拆分摘要

### 9.4 proposal inspect 聚合视图

```markdown
## Proposal: <proposal-id> - <title>

### Feature Map
| ID | Title | Stage | Steps | Dependencies | Blocked By |
|----|-------|-------|-------|-------------|------------|
| FEAT-001 | Clone API | planned | 0/3 | - | - |
| FEAT-002 | 子对象复制 | designed | - | FEAT-001 | FEAT-001 |
| FEAT-003 | Clone 前端 | draft | - | FEAT-001 | FEAT-001 |

### Dependency Health
- ⚠️ FEAT-002 blocked by FEAT-001 (status: planned, not done)
- ⚠️ FEAT-003 blocked by FEAT-001 (status: planned, not done)

### Cross-cutting Cards
| ID | Type | Title | Constrains/References |
|----|------|-------|----------------------|
| CONV-001 | convention | clone 模式约定 | FEAT-001, FEAT-002 |

### Stage Summary
| Stage | Count |
|-------|-------|
| draft | 1 |
| designed | 1 |
| planned | 1 |
| in_progress | 0 |
| done | 0 |

### Recommendations
- 执行 FEAT-001 的 Step 1-3 以解除 FEAT-002 和 FEAT-003 的阻塞
- FEAT-003 仍为 draft，建议先 clarify 前端需求细节
```

---

## 10. STR 卡片的命运

STR 卡片不再作为手动维护的卡片类型存在。其功能拆分如下：

| STR 原有功能 | 新机制 |
|-------------|--------|
| 导航入口 | `proposal inspect` 自动聚合视图 |
| 需求索引 | FEATURE 卡片的 `depends_on` / `part_of` 关系自动构建依赖图 |
| 内容合成 | FEATURE 卡片的 Design 段落（设计阶段）或父 FEATURE 的 Sub-Features 概览（拆分后） |
| 知识组织 | library 中的主题索引改为自动生成（CLI 扫描 CONV/MOD/DEC/FIND 的 links 关系） |

---

## 11. 迁移策略

### 阶段一：方法论文档更新（不涉及代码）
1. 更新 `card-templates.md`：新增 FEATURE 模板，标记旧类型为 deprecated
2. 更新 `workflow-rules.md`：用阶段门控规则替换模式选择表
3. 更新 `SKILL.md`：重写执行流程

### 阶段二：CLI 增量支持（向后兼容）
1. `card create --type feature` + `card evolve`
2. `proposal inspect` 自动聚合视图
3. `card log` 追加 History
4. `card split` 拆分支持

旧类型（REQ/DES/TASK）保留读取能力，创建能力标记为 deprecated 但继续工作。

### 阶段三：渐进清理
1. 归档工具支持将旧提案中的 REQ+DES+TASK 合并为 FEATURE
2. `card create` 默认类型改为 feature
3. STR 创建命令移除

---

## 12. 风险与缓解

| 风险 | 缓解措施 |
|------|---------|
| FEATURE 卡片过长（> 200 行有效内容） | Implementation Plan 步骤用 `--section` 裁剪读取；门控规则中 single-step 上限 200 行 |
| 失去跨卡片"意外连接" | 跨功能连接通过 `depends_on` / `extends` 显式表达；横切类型（CONV/DEC）保持跨功能链接 |
| 多人并行编辑冲突 | FEATURE 粒度控制在"一个开发者一个迭代"；必要时通过 `card split` 为并行创造条件 |
| SKILL 不执行门控 | `card evolve` 在 CLI 层面强制执行门控，不依赖 SKILL 自律 |
| Implementation Plan 步骤仍然单薄 | 门控明确要求每个步骤的方法签名、文件路径、边界条件；SKILL 硬规则禁止"参考其他卡片" |
| 旧数据兼容 | 只读兼容旧类型；archive 工具渐近合并 |

---

## 13. 与已有修复方案的关系

本方案是 `remediation-card-fragmentation.md` 的**根本性替代**，而非补充：

| 维度 | remediation (补丁) | 本方案 (模型修正) |
|------|-------------------|-------------------|
| 问题诊断 | 执行不到位（没填合成、没做设计、密度太低） | 模型设计不当（类型拆分本身制造碎片） |
| 修复策略 | 加门控、加密度检查、加合成段落 | 合并 REQ+DES+TASK 为 FEATURE，用阶段演进替代类型拆分 |
| STR 角色 | 保持独立类型，加强合成要求 | 移除独立类型，替换为 CLI 自动聚合 |
| LOG 角色 | 保持独立类型 | 合并到 FEATURE 的 History 段落 |
| 卡片类型数 | 10 种（不变） | 5 种 |
| 核心思路 | 让现有模型运行得更正确 | 换一个更正确的模型 |

---

## 14. 设计决策记录

| 决策 | 理由 |
|------|------|
| REQ+DES+TASK → FEATURE | 三者信息高度重叠，拆分只产生维护成本不产生新洞察 |
| 保留 CONV/MOD/DEC/FIND | 这些是横切关注点，天然跨功能生效，不宜塞入单一 FEATURE |
| FEATURE 用阶段而非子类型表达成熟度 | 阶段演进让内容在卡片内部累积，而非在卡片之间拆分 |
| STR 改为自动生成 | STR 手动维护的合成价值不足以覆盖维护成本；自动聚合更准确 |
| 门控在 CLI 层面强制执行 | 不依赖 SKILL 自律，减少"流程被跳过"的风险 |
| Implementation Plan 步骤禁止跨卡引用 | 解决 TASK 变成"指向器"的根本问题——强制信息自足 |
