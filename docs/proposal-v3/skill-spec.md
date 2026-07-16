# FlowForge SKILL v3：Agent 工作流方法论

> 日期：2026-07-09
>
> 定义 flowforge-design 和 flowforge-implement 两个 SKILL 的完整方法论：
> 需求分解、信息探索、设计推理、任务拆分、约束构建、Token 感知执行。
> 卡片模型见 `proposal-card-model-v3.md`，CLI 命令规格见 `proposal-cli-spec-v3.md`。

---

## 1. Token 消耗控制设计

### 1.1 问题

FEATURE 卡片完整内容可能达到 200-400 行（~3000-6000 tokens）。如果 Agent 每次执行一个步骤都读取整张卡，
在 5-10 个步骤的执行周期中 Token 消耗将呈线性增长。需要在三个层面协同控制。

### 1.2 卡片结构层面

FEATURE 模板设计时已考虑分段可提取性：

- 每个 `### Step N:` 是自足的独立单元，包含执行该步骤所需的全部信息（Goal, Files, Approach, Edge Cases, Dependencies, Verification）
- 步骤状态用 HTML 注释标记（`<!-- step-status: done -->`），CLI 可解析
- `## Constraints` 独立段落，可单独提取
- 层级标题明确（H2 为大段落，H3 为子结构），支持 `--section "Implementation Plan.Step 3"` 级别的精确提取

### 1.3 CLI 层面

三个命令协同控制 Token：

| 命令 | 用途 | Token 量级 |
|------|------|-----------|
| `card read --summary` | 快速了解 FEATURE 元信息和阶段 | ~100 tokens |
| `card read --section "Implementation Plan.Step 3"` | 仅读取当前步骤 | ~300 tokens |
| `context feature --feature <id> --step 3` | 获取步骤上下文 + 约束 + 库引用 + 依赖状态 | ~400 tokens |
| `card read --section "Constraints"` | 仅读取约束 | ~100 tokens |
| `card read --section "Design.Key Decisions"` | 仅读取设计决策 | ~150 tokens |

Agent 执行 Step 3 时只需一条命令：
```bash
flowforge context feature --feature FEAT-001 --step 3
```

返回的上下文集（步骤详情 + 约束 + 相关库卡片 + 依赖状态）完全自足，
不需要再读取整张卡或其他卡片。预期 Token 消耗 ~400 tokens，远低于完整读取的 ~3000-6000 tokens。

### 1.4 SKILL 层面

SKILL 必须明确引导 Agent 使用正确的读取策略：

```markdown
## Token-Aware Reading Rules

- **Start execution**: Run `context feature --feature <id> --step <n>`. This gives you everything.
- **Need design rationale**: `card read --section "Design.Key Decisions"` — not the whole Design section.
- **Verify overview**: `card read --summary` — not the whole card.
- **Check proposal health**: `proposal inspect <id>` — not loading every FEATURE card.
- **All constraints in context**: `context feature --step` already includes Constraints section and library cards.
- **Never read the whole FEATURE card during step execution.** If context feature --step is not enough, read the specific missing section.
```

### 1.5 场景化读取矩阵

| 场景 | 命令 | 读什么 |
|------|------|--------|
| 设计阶段：了解功能全貌 | 直接读 .md 文件 | 全量（此时需要全貌做设计决策） |
| 设计阶段：查库约束 | `library suggest --for FEAT-xxx` | 相关 CONV/DEC/MOD/FIND 摘要 |
| 实施阶段：执行 Step N | `context feature --feature <id> --step N` | 步骤上下文 + 约束 + 依赖 |
| 实施阶段：确认某条约束细节 | `card read --section "Constraints"` | Constraints 段落 |
| 实施阶段：了解被依赖 FEATURE 接口 | `card read --section "Design.Architecture" <dep-id>` | 被依赖 FEATURE 的架构段 |
| 审查阶段：看整体进展 | `proposal inspect <id>` | 聚合视图 |
| 审查阶段：看某 FEATURE 的 Design | `card read --section "Design" <id>` | Design 段落 |

---

## 2. flowforge-design SKILL

### 2.1 工作流总览

```
seed → clarify → enrich → plan
  创建      追问      填充     填充
 FEATURE   用户      Design  Implementation
 (draft)          段落    Plan 段落
```

每轮不是"创建新卡片"，而是**丰富现有 FEATURE 卡片的更多段落**。唯一创建新卡片的时机：
1. 识别出新的独立功能 → 创建新的 FEATURE
2. 发现需要拆分的过大 FEATURE → `card split` 创建子 FEATURE
3. 发现横切知识 → 创建 CONV/DEC/MOD/FIND

### 2.2 需求分解为 FEATURE 的方法论

**Step 1 — 理解需求边界**

- 阅读用户输入，提取核心问题陈述
- 追问：**"这个需求完成后，用户能做什么以前不能做的事？"**
- 输出：一段 Summary（1-3 句话），写入 FEATURE 卡片

**Step 2 — 识别能力点**

从需求中提取用户可感知的独立能力。判断标准：

> 一个能力 = 一个可独立演示/验证的产出（一个 API endpoint、一个完整交互流程、一个数据处理链路）

示例："FileProcessor Clone" →
- 能力 1：通过 API 触发 clone（后端能接收请求并返回结果）
- 能力 2：clone 过程正确复制所有子对象（数据完整性）
- 能力 3：前端页面支持 clone 操作（用户可交互）

**Step 3 — 检验粒度（→ 映射为 FEATURE）**

对每个候选能力，追问：

1. **能否在一个迭代内完成？** 能 → 通过，创建 FEATURE（draft）。不能 → 继续分解为子能力
2. **能否独立验证？** 能独立运行测试/手动验证 → 通过。需要等其他能力完成才能验证 → 合并到依赖的能力中
3. **是否只是技术步骤（如"写单元测试"）？** 是 → 不独立成卡，写入 Implementation Plan 的步骤中

**Step 4 — 建立依赖**

分析能力之间的时序关系：
- A 的产出是 B 的输入？→ B `depends_on` A
- A 和 B 可并行？→ 无 depends_on 关系
- A 的接口契约是 B 的设计前提？→ B 标注依赖并写明接口契约，不阻塞（可用 mock/stub）

**Step 5 — 识别横切关注点**

- 是否存在跨 FEATURE 的约束？→ 创建 CONV
- 是否存在影响系统架构的决策？→ 创建 DEC
- 是否存在可复用的模块知识？→ 创建 MOD

**硬约束：**
- 单轮最多创建 5 张 draft FEATURE（防止机械导入 PRD 条目）
- 每张 draft FEATURE 的 Summary + Motivation 有效行数 ≥ 5
- 如果拆不出足够大的独立能力，先在一张 FEATURE 内用 Open Questions 标注待澄清点

**PROP 卡片更新触发点：**

PROP 的 `## Feature Map` 和 `## Architecture Overview` 应在以下时机更新：

| 触发点 | 更新内容 |
|--------|---------|
| 首个 FEATURE 创建后 | 在 Feature Map 中追加一行 |
| 任一 FEATURE 完成 `card evolve --stage designed` | 更新该 FEATURE 的职责描述（从占位符改为实际总结） |
| `card split` 执行后 | 将父 FEATURE 拆分信息反映到 Feature Map + Architecture Overview |
| 任一 FEATURE 完成 `card evolve --stage done` | 更新阶段标记 |
| 发现所有 FEATURE 的依赖关系稳定后 | 填充 Architecture Overview（协作关系、共享技术决策） |

Agent 不应在每个小变更后都更新 PROP——仅在上述触发点更新。
`proposal inspect` 会检测 `prop_feature_map_stale` 并在 Feature Map 与实际状态不一致时告警。

### 2.3 信息探索与现状调查

在 enrich FEATURE 时，Agent 必须按优先级探索三个信息源：

```
1. FlowForge Library (最高优先)
   → library suggest --for <feature-id>
   → 获取相关 CONV/DEC/MOD/FIND
   → 提取约束（写入 Constraints）和参考设计（写入 Design.Alternatives Considered）

2. 项目源码 (理解现状)
   → 读取相关模块的现有代码
   → 理解现有接口、数据结构、模式
   → 注意：源码是事实，不是规范——可以改变

3. 外部知识源 (补充)
   → 配置的 knowledge_sources
   → PRD、设计文档、API 文档
```

**探索的终止条件：**
- 当前 FEATURE 的 Open Questions 中没有阻塞设计决策的问题
- 连续两轮探索没有产生新信息和新搜索方向
- Library、源码、外部源都已检查且均无 actionable 结果

**记录探索：** 使用 `card log --kind progress` 记录探索过程（搜索了什么、匹配了什么、遗漏了什么）。

### 2.4 设计时的思维逻辑

填充 FEATURE 的 Design 段落时，Agent 遵循以下推理链：

```
1. 明确问题
   → "我们要解决什么？"（已有 Summary）
   → "为什么需要解决？"（已有 Motivation）

2. 分析约束
   → CONV 的 Rule（从 library 来的强制规则）
   → MOD 的模块边界（不能跨模块职责）
   → DEC 的架构决策（已有技术选择）
   → 业务规则（PRD 中的硬性要求）
   → 写入 Constraints 段落，标注来源

3. 生成备选方案
   → 每个方案记述：核心思路、涉及改动、优缺点
   → 至少考虑 2 个方案；如果只有一个方案，追问自己"除了这个，还有其他可能性吗？"

4. 评估与选择
   → 对照 Constraints 排除不合适的方案
   → 选择理由写入 Key Decisions（每条决策 ≤3 行，包含"为什么"而不仅是"选了什么"）

5. 记录放弃的方案
   → 写入 Alternatives Considered
   → 如果只有一个方案且无可替代，写 None 并说明原因

6. 检查横切影响
   → 这个决策是否影响多个 FEATURE？→ 是：提取为 DEC 卡，FEATURE 中 references DEC
   → 是否定义了可复用的模式？→ 是：提取为 CONV 卡
```

**关键思维检查点：**
- 每个 Key Decision 的"理由"是否 ≥ 1 句具体分析而不仅是"这样更简单"？
- 每个 Constraint 是否标注了来源（CONV ID / MOD ID / 业务规则）？
- Design 是否与 Open Questions 有对应关系（每个问题都有回答或标注为假设）？

### 2.5 实施计划拆解标准

填充 Implementation Plan 时，步骤的拆解标准：

> **一个步骤 = 一个可独立验证的代码变更单元。**

| 维度 | 标准 |
|------|------|
| **边界** | 每个步骤有明确的文件范围和产出物 |
| **可验证** | 每个步骤完成后可以独立测试/验证（不需要等其他步骤） |
| **有序** | 步骤按依赖关系排序，先基础后上层 |
| **自足** | 每个步骤包含足够信息，不需要跨卡跳转 |

**每个步骤的必填字段（门控强制）：**

| 字段 | 要求 |
|------|------|
| Goal | 1 句话：本步骤交付的可验证结果 |
| Files | 文件路径列表（相对项目根目录），可标注新建/修改 |
| Approach | 关键方法签名或伪代码、算法选择、数据结构、状态流转逻辑 |
| Edge Cases | 至少 1 个边界条件及其处理方式 |
| Dependencies | 依赖的 FEATURE 及原因；如果依赖未就绪，**必须写明等待策略**（如 "可用 mock FileProcessorCloneService 解耦等待"） |
| Verification | 测试场景或关键断言 |

**拆分数量标准：**
- 理想：3-7 个步骤（太少说明不够细，太多说明 FEATURE 该拆分）
- Implementation Plan 超过 10 个步骤 → 触发拆分建议（`card split`）

**反模式：**
- "Step 1: 写实现 → Step 2: 写测试" — 测试应嵌入每个步骤的 Verification，不独立成步
- "参考 DES-xxx 的设计" — 禁止跨卡引用
- 步骤只有 Goal 没有 Approach — 太薄，无法执行

**并行步骤标记：**

Implementation Plan 的步骤默认串行，但两个修改不同文件、无共享状态的步骤可以标记为并行：

```markdown
### Step 2: Clone API 路由
- **Parallel**: yes
- **Dependencies**: Step 1 必须完成（但 Step 2 和 Step 3 可并行）

### Step 3: Clone 请求参数校验
- **Parallel**: yes
- **Dependencies**: Step 1 必须完成
```

并行标记规则：
- `Parallel: yes` 必须同时满足：(1) 与同一 FEATURE 的其他 `Parallel: yes` 步骤修改的文件无交集，(2) 无共享运行时状态
- 如果两个并行步骤依赖同一个前置步骤（如都依赖 Step 1），这不是问题
- 如果 Agent 不确定是否可以并行，默认 `Parallel: no`
- 实施 SKILL 可以为 `Parallel: yes` 的步骤启动多个并行 Agent

### 2.6 约束构建

约束从多个来源汇总到 FEATURE 的 Constraints 段落，每条标注来源：

```markdown
## Constraints

- [CONV-001] 所有 clone 操作必须使用独立 Cmd + Service 模式（来源：library）
- [DEC-002] clone 过程保持单事务（来源：library）
- [MOD-003] 不得在 adapter 层写业务编排逻辑（来源：library）
- fileProcessorCode 由前端提供且必须唯一（来源：业务规则）
- 不复制关联的定时任务配置（Out of Scope）
```

**约束构建流程：**
1. `library suggest --for <feature-id>` → 获取 CONV/DEC/MOD
2. 询问用户确认业务约束
3. 对照 Constraints 验证每个 Design.Key Decisions —— 确保没有决策违反约束
4. 如果 Design 中的某个决策引入了新的跨功能约束 → 评估是否需要创建 CONV

### 2.7 硬规则

```markdown
## Hard Rules

- Use `card init --type feature` to create cards; then edit the file directly.
- Use `card link`/`card unlink` for all link operations.
- Use `card evolve` for stage transitions — never hand-edit status in frontmatter.
- Use `card log` for progress recording — never hand-edit ## History.
- Run `flowforge validate all` after any .md file changes.
- Never edit auto-generated navigation sections (if present).
- Never create a FEATURE with <5 lines of effective business content in Summary + Motivation.
- Never skip stages: draft → designed → planned must be sequential via `card evolve`.
- All Open Questions must be cleared or marked as assumptions before `card evolve --stage designed`.
- Implementation Plan steps must include Files, Approach, and Edge Cases — no "参考其他卡片" style pointer steps.
- Never create >5 draft FEATURE cards in a single round.
- Before enriching Design, always run `library suggest --for <feature-id>`.
- Each Key Decision must include a "why" (≥1 sentence), not just "what".
- Each Constraint must have a source annotation (CONV ID / MOD ID / business rule).
- After completing enrichment, run `proposal inspect <id>` and address all health issues.
```

---

## 3. flowforge-implement SKILL

### 3.1 Token 感知的执行流程

```
1. 确定目标步骤
   → card read --summary FEAT-xxx   (确认当前阶段和步骤状态，~100 tokens)

2. 获取执行上下文
   → context feature --feature FEAT-xxx --step N   (~400 tokens)
   → 返回：步骤详情 + Constraints + 库引用 + 依赖状态
   → 不再读取整张卡或其他卡片

3. 实现步骤
   → 按 Approach 描述实施
   → 对照 Edge Cases 处理边界条件
   → 对照 Verification 定义编写测试

4. 记录进展
   → card steps FEAT-xxx --status done N
   → card log FEAT-xxx --event "..." --kind progress

5. 验证
   → 运行测试
   → flowforge validate all

6. 进入下一步或完成
   → 如果所有步骤完成：card evolve FEAT-xxx --stage done
   → 否则：继续 Step N+1（重复步骤 2-5）
```

### 3.2 执行中的额外信息获取

如果 `context feature --step <n>` 返回的信息不足以完成步骤，按需获取补充信息——**不读取整张卡**：

| 缺什么 | 怎么获取 |
|--------|---------|
| 需要更详细的设计背景 | `card read --section "Design.Key Decisions" <id>` |
| 需要看被依赖 FEATURE 的接口定义 | `card read --section "Design.Architecture" <dep-id>` |
| 需要了解某个 CONV 的具体要求 | `card read --section "Rule" <conv-id>` |
| 实施中发现 plan 有遗漏 | 直接编辑 FEATURE card，补充缺失的 Approach/Edge Cases/Dependencies |

**实现中发现设计问题的处理流程：**
1. 暂停实施
2. 在 FEATURE 的 Open Questions 中追加发现的问题
3. `card log --kind blocked --event "..."` 记录阻塞原因
4. `card steps --status blocked N --reason "..."` 标记步骤阻塞
5. 切换回 flowforge-design 模式，更新 Design/Constraints
6. 更新 Implementation Plan 中受影响的步骤
7. `card steps --status in_progress N` 解除阻塞后继续

### 3.3 硬规则

```markdown
## Hard Rules

- Start each step with `context feature --feature <id> --step <n>`. This is your primary context.
- Never read the whole FEATURE card during step execution. Use section-level reading for supplemental info.
- Execute steps in order; skip blocked steps (marked `<!-- step-status: blocked -->`).
- After each step, use `card log` to record progress and `card steps` to update status.
- Use `card evolve --stage done` when all steps are complete.
- If implementation reveals a missing detail, edit the FEATURE card directly to add it, then continue.
- If a design issue is found, stop and record it via `card log --kind blocked` + `card steps --status blocked`.
- Run tests and `flowforge validate all` when card state changes.
- CLI for structured ops only (link, evolve, log, steps); direct file editing for body content.
```

---

## 4. Feedback 修正协议

当 Agent 在实施过程中发现设计问题或执行受阻时，需要遵循明确的反馈闭环协议。
本协议替代 v2 中"创建 design-flaw requirement 卡片 → 路由到 design SKILL"的流程。

### 4.1 状态机

```
实施中发现问题
  ├── 小问题（可在当前步骤内修正）
  │     → 暂停实施 → 直接编辑 FEATURE 补充设计细节 → 继续实施
  │
  ├── 中等（影响当前 FEATURE 的设计）
  │     → card steps --status blocked N --reason "..."
  │     → card log --kind blocked --event "..."
  │     → 更新 Open Questions 追加发现的问题
  │     → [切换思维到设计模式]
  │     → 修改 Design / Constraints / Implementation Plan（直接编辑 .md）
  │     → 解除阻塞：card steps --status in_progress N
  │     → 继续实施
  │
  └── 严重（设计有根本缺陷，需要回退阶段）
        → card steps --status blocked N --reason "..."
        → card log --kind blocked --event "..."
        → card evolve <id> --stage designed --regress
        → 重新分析、修改 Design
        → 更新受影响的 Implementation Plan 步骤
        → card evolve <id> --stage planned
        → card steps --status in_progress 1（从头开始）
```

### 4.2 各场景详述

#### 场景一：小修正（实施细节补充）

**触发条件：** Plan 缺少某个边界条件或方法参数细节，但不影响整体设计。

**协议：**
1. 暂停当前步骤实施
2. 直接编辑 FEATURE .md 文件，在对应步骤的 Edge Cases 或 Approach 字段补充
3. `card log --kind progress --event "Step N: 补充了 Edge Case X"`
4. 继续实施

**不需要标记为 blocked。** 这是 Plan 的自然完善。

#### 场景二：中等修正（设计修改但不回退阶段）

**触发条件：** 发现某个设计决策在当前实现中不可行，需要修改 Design 段落，
但修改不影响已完成的步骤和其他 FEATURE。

**协议：**
1. 标记当前步骤阻塞：
   ```bash
   flowforge card steps FEAT-xxx --status blocked N --reason "API 响应格式与前端约定不兼容"
   flowforge card log FEAT-xxx --kind blocked --event "Step N blocked: API response format mismatch"
   ```
2. 在 Open Questions 追加：
   ```markdown
   ## Open Questions
   - [阻塞 Step N] API 响应格式是否需要包含嵌套文件树结构？
   ```
3. 切换到设计思维：修改 Design.Key Decisions 或 Architecture
4. 检查受影响的其他步骤——如果 Step N+1 也依赖相同设计，需要同步更新
5. 解决 Open Questions
6. 解除阻塞：
   ```bash
   flowforge card steps FEAT-xxx --status in_progress N
   flowforge card log FEAT-xxx --kind progress --event "Step N resumed: API response format redesigned"
   ```

#### 场景三：严重修正（阶段回退）

**触发条件：** 设计有根本缺陷——比如架构方案完全不适用、需要推翻重来。
已完成的步骤也可能受影响。

**协议：**
1. 记录阻塞：
   ```bash
   flowforge card steps FEAT-xxx --status blocked N --reason "单事务方案不适用于异步 clone 场景"
   flowforge card log FEAT-xxx --kind blocked --event "Design flaw: single-transaction approach incompatible with async clone"
   ```
2. 阶段回退：
   ```bash
   flowforge card evolve FEAT-xxx --stage designed --regress
   ```
   CLI 输出确认："FEAT-xxx: planned → designed。2 个步骤状态已重置。History 保留。"
3. 重新设计：修改 Design、Constraints、更新 Open Questions
4. 重建 Implementation Plan：删除无效步骤，编写新的步骤列表
5. 重新通过门控：
   ```bash
   flowforge card evolve FEAT-xxx --stage planned
   ```
6. 恢复实施：`card steps FEAT-xxx --start 1`

**回退时已完成步骤的代码变更不会自动回滚。** Agent 需要评估是否需要 revert 已有的代码变更。

### 4.3 何时使用 feedback SKILL 而非内置修正

FEATURE 模型下的修正协议覆盖了大部分场景。但以下情况应触发独立的 `flowforge-feedback` SKILL：

| 场景 | 处理方式 |
|------|---------|
| 发现的问题影响**多个** FEATURE | 触发 feedback SKILL，创建 FIND 或 CONV 卡 |
| 发现的模式应该成为库知识 | 触发 feedback SKILL → `library import` |
| 用户报告的 bug（非 Agent 自己发现） | 触发 feedback SKILL，创建追踪 FEATURE |
| 设计讨论需要人的决策 | 暂停实施，等待用户输入 |

**简单规则：** 只影响当前 FEATURE → 内置修正协议；影响跨 FEATURE 或跨 proposal → feedback SKILL。

---

## 5. 现有 SKILL 的 v3 适配

### 5.1 flowforge-feedback

**当前行为：** 接收发现 → 分类（bug/finding/knowledge/missing-requirement/design-flaw）→ 创建对应卡片 → 记录 log。

**v3 变更：**

| 维度 | v2 行为 | v3 行为 |
|------|--------|--------|
| bug 追踪 | 创建 `task` 卡（not_ready）→ 链接到 source | 创建 FEATURE（draft）或标注已有 FEATURE 的 Open Questions |
| missing-requirement | 创建 `requirement` 卡 + `structure add` | 创建 FEATURE（draft）；不再有 `structure add` |
| design-flaw | 创建 `requirement` 卡（design change request）→ 路由到 design SKILL | 直接在受影响 FEATURE 上执行修正协议（参见 §4）；跨 FEATURE 的创建 FIND |
| knowledge 路由 | `library import` / `library promote` | 保持不变 |
| 进展记录 | `log create --kind feedback` | `card log --kind feedback` 追加到受影响 FEATURE 的 History |
| CLI-only 约束 | "CLI is the only read/write path" | 放宽为职责分离（参见 `cli-spec.md` §1.2） |
| STR 索引 | `structure add` 将 requirement 加入索引 | 废弃——PROP 的 Feature Map 由 Agent 在触发点更新 |

**classification-rules.md 修改：**

- `bug → task` 改为 `bug → feature (draft) 或标注已有 feature`
- `missing-requirement → requirement + structure add` 改为 `missing-requirement → feature (draft)`
- `design-flaw → requirement (design change request)` 改为 `design-flaw → 评估影响范围：单 FEATURE 修正协议 / 跨 FEATURE 创建 FIND`
- 所有 `log create` 引用改为 `card log`
- 所有 `card create --type task/requirement` 引用改为 `card init --type feature`

### 5.2 flowforge-curate

**当前行为：** Mode A（外部导入）+ Mode B（提案归档）。两者都创建 STR 索引 + 原子卡片。

**v3 变更：**

| 维度 | v2 行为 | v3 行为 |
|--------|--------|--------|
| 知识组织 | 按概念聚类 → 创建 STR 索引卡 | 按概念聚类 → 创建 PROP（library 中的主题 proposal）或在 library 中直接组织 CONV/DEC/MOD/FIND |
| STR 创建 | `card batch` 中创建 STR 类型索引卡 | 废弃——library 中的组织方式用 library facets + 标签替代 |
| Mode B 过滤 | 跳过 `log`、`requirement`、`task`、`ROOT`、`STR` | 跳过 `feature`（除非需要从 FEATURE 的 Design 提取 DEC）、`log` 已不存在 |
| Mode B 触发 | 扫描 `03-completed/` 目录 | 用 `proposal list --status completed` 定位已完成 proposal（不再依赖目录位置） |
| 提取类型映射 | `design → library design card` | 从 FEATURE 的 Design.Key Decisions 提取→创建 library DEC 卡 |
| CLI-only 约束 | "CLI is the only read/write path" | 放宽为职责分离 |
| 创建方式 | `card batch --manifest` + `card create` | 保留 `card batch` 和 `card init`；废弃 `--type requirement/design/structure` |
| 归档后文件位置 | 物理移动到 `03-completed/` | 不移动——proposal 始终在 `01-workspace/` 下，仅 PROP status 变更 |

**extraction-guide.md 修改：**
- 移除 `design` 和 `requirement` 作为库卡片类型（库中只有 CONV/DEC/MOD/FIND）
- 新增：从提案 FEATURE 提取 Design 决策为 DEC 的指南

**workflow-rules.md 修改：**
- "Cluster and Plan" 步骤：不再创建 STR 索引卡，改为创建 library topic PROP 或直接在卡片 tags 中标记聚类
- `structure add` / `@ref:indexes` 引用 → 废弃

### 5.3 flowforge-implement

已在本文档 §3 完整定义。核心变化：
- 入口从 `task ready` 改为 `context feature --step <n>`
- 进展从 `log create` 改为 `card log` + `card steps`
- 遇到设计问题从"创建 design-flaw requirement → 路由到 design SKILL" 改为内置修正协议
- 卡片编辑从 CLI-only 放宽为直接文件编辑

---

## 6. 与 remedation 方案的 SKILL 对比

| 维度 | 旧 SKILL（remediation 后） | 新 SKILL |
|------|--------------------------|---------|
| 设计入口 | proposal create → STR 索引 | `card init --type feature` 创建单张 FEATURE |
| 设计流程 | index → clarify → analyze → discover → design → split tasks | seed → clarify → enrich → plan（同一张卡渐进填充） |
| 实现入口 | `task ready --type i` 或 task ID | `context feature --feature <id> --step <n>` |
| 实现上下文 | `context task --task <id>` → 跨卡跳转读 REQ/DES | 一次命令获取全部上下文（步骤+约束+库+依赖） |
| 进展记录 | `log create --kind progress`（创建新卡） | `card log --event "..."`（追加到 FEATURE History） |
| 步骤状态 | `task done` / `task block`（独立命令） | `card steps --status done/blocked <n>` |
| 卡片编辑 | 全部通过 `card update --body` / `--section` | 直接编辑 .md 文件（content），CLI 管不变式 |
| Token 控制 | 无专门设计（读全量） | 三层设计：结构+CLI+SKILL，执行单步骤 ~400 tokens |
