# flowforge-design SKILL 草案

> 版本：draft
>
> 目标：把 [Design SKILL 工作流](./design-skill-workflow.md) 压缩成可部署 SKILL 的草案，明确 SKILL 本体应写什么、reference 应拆什么、哪些 CLI 能力是实现前置。
>
> 关键 CLI 输出契约参考：[Design SKILL CLI 契约设计](./design-skill-cli-contracts.md)。

## 1. 核心设计判断

`flowforge-design` 的难点不是“让 Agent 写设计”，而是让 Agent **不要直接写长设计**。

它必须强制 Agent 先建立 proposal 工作面，再用卡片和关系推进：

```text
proposal inspect
  -> 需求索引树
  -> requirement
  -> analysis task
  -> library discovery
  -> design card
  -> implementation task
  -> 单轮输出
```

SKILL 本体只保留流程和强约束；卡片模板、walkthrough、判断细则放 reference。

## 2. SKILL Description 草案

```yaml
name: flowforge-design
description: Use when the user wants to analyze, clarify, design, or decompose a requirement or proposal before implementation. Guides the agent to build and update FlowForge proposal cards through the CLI: requirement index tree, requirement cards, analysis tasks, design cards, logs, and implementation tasks. Do not use for executing an existing task, handling test feedback or bugs, archiving knowledge, or simple card lookup.
```

### 2.1 触发命中目标

应该命中的用户表达：

- “分析这个需求”
- “设计这个功能”
- “把这个 proposal 拆成任务”
- “继续完善当前 proposal”
- “先梳理需求边界”
- “这个实现前需要先设计一下”

不应该命中的表达：

- “执行 TASK-xxx”
- “这个测试失败了，处理反馈”
- “归档这个 proposal”
- “查一下某张卡”

## 3. SKILL.md 主体草案

下面是未来 `assets/skills/flowforge-design/SKILL.md` 的主体草案。它应尽量短，避免把模板和 walkthrough 全塞进去。

```markdown
# flowforge-design

You guide requirement analysis and design for FlowForge proposals. Do not write long proposal documents. Grow a card network through `flowforge` CLI commands.

## Start

1. Resolve the current project and proposal:
   - `flowforge project current`
   - `flowforge proposal current`
   - `flowforge proposal inspect <proposal-id>`
   - `flowforge context proposal --proposal <proposal-id>`
2. If no project exists, ask the user to create or select one.
3. If no proposal exists and the user is starting new work, suggest `flowforge proposal create "<title>"`.
4. Never edit wiki files directly. Use FlowForge CLI only.

## Workflow

For each design turn:

1. Update the requirement index tree.
   - Keep STR cards to 7-15 direct entries.
   - Use `flowforge structure add/remove` for STR changes.
   - Do not put logs, findings, or long design text into STR cards.

2. Create or update atomic requirement cards.
   - One user-visible behavior, constraint, or acceptance point per requirement.
   - Put unresolved issues in `Open Questions`.

3. Create analysis tasks for uncertainty.
   - Use `flowforge task create --type a`.
   - Create analysis tasks when code impact, library rules, cross-project boundaries, or acceptance criteria are unclear.

4. Discover library context through CLI only.
   - Use `flowforge library suggest --for <card-id>`.
   - Use `flowforge card search <query> --scope library`.
   - Read only selected cards with `flowforge card read <id> --summary` or `--section`.
   - Never grep or directly read `02-library/`.

5. Create design cards when conclusions stabilize.
   - Link requirements, findings, module cards, conventions, and decisions.
   - Do not create a design card that summarizes the whole proposal.

6. Create implementation tasks only when executable.
   - Use `flowforge task create --type i`.
   - A ready task must link requirement, design, and constraints.
   - If based on assumptions, create it as `not_ready` or blocked and link the open question.

7. Record the turn.
   - Use `flowforge log create --kind <kind>`.
   - Logs point to the relevant proposal, task, requirement, or design.
   - Do not backfill every log into root, task, or requirement cards.

## Output

End each turn with:

- Cards added or updated: IDs and purpose.
- Relations added: short summary.
- Current gaps: open questions, not_ready tasks, blocked tasks.
- Next step: continue design, execute one ready implementation task, or ask the user.

## Hard Rules

- CLI is the only write path.
- Do not load the whole proposal or whole library.
- Do not create tasks with only a title.
- Do not attach every evidence card to root, requirement, or task cards; rely on backlinks.
- Do not execute implementation work inside this skill.
```

## 4. Reference 拆分设计

未来部署时建议目录：

```text
assets/skills/flowforge-design/
  SKILL.md
  references/
    card-templates.md
    library-discovery.md
    walkthrough-flowforge-v2.md
```

### 4.1 `card-templates.md`

包含：

- Requirement 卡模板
- Analysis Task 卡模板
- Design 卡模板
- Implementation Task 卡模板
- Log 卡模板
- 单轮输出格式

触发读取时机：

- Agent 要创建或审查卡片正文时读取。
- Agent 不确定某类卡最小段落时读取。

不放进 SKILL 本体的原因：

- 模板较长。
- 后续会随真实使用调整。
- 每次触发 design SKILL 不一定都需要全部模板。

### 4.2 `library-discovery.md`

包含：

- library 查询三层模型
- `library suggest` 输出字段
- 候选筛选规则
- 定点读取规则
- 关联写入规则
- 未命中处理

触发读取时机：

- Agent 需要查规范、模块知识、历史设计、finding。
- Agent 准备把 library 卡链接到 analysis/design/task。

### 4.3 `walkthrough-flowforge-v2.md`

包含：

- 用 FlowForge v2 设计自身作为验证用例
- 从 proposal inspect 到 implementation task 的端到端示例
- 验证点和失败信号

触发读取时机：

- 迭代 SKILL 本身。
- 验证 design workflow 是否跑通。
- 新 agent 不理解卡片生长方式时参考。

## 5. SKILL 本体覆盖检查

| workflow 要求 | SKILL 本体覆盖方式 |
|---------------|-------------------|
| 不写长 proposal | Hard Rules + Workflow 第 1/5 条 |
| 先建立工作面 | Start 第 1 条 |
| proposal / project 缺失处理 | Start 第 2/3 条 |
| 需求索引树 | Workflow 第 1 条 |
| requirement 原子化 | Workflow 第 2 条 |
| analysis task 驱动不确定点 | Workflow 第 3 条 |
| library 只能通过 CLI 查 | Workflow 第 4 条 + Hard Rules |
| design card 只表达稳定焦点 | Workflow 第 5 条 |
| implementation task 不空心 | Workflow 第 6 条 + Hard Rules |
| log 卡记录过程 | Workflow 第 7 条 |
| 中心卡不回写证据 | Hard Rules |
| 单轮用户输出 | Output |
| 不执行实现 | Hard Rules |

覆盖结论：SKILL 本体可以保持短小，核心约束已经覆盖；模板和判断细节应由 reference 提供。

## 6. 每轮运行手册

这一节定义 SKILL 真实运行时的顺序。它比 `SKILL.md` 主体更详细，但仍然是业务设计，不是代码实现。

### 6.1 固定命令顺序

每轮 design 开始时，Agent 按顺序执行：

```bash
flowforge project current
flowforge proposal current
flowforge proposal inspect <proposal-id>
flowforge context proposal --proposal <proposal-id>
```

分工：

- `proposal inspect`：判断 proposal 状态、缺口、open question、not_ready / blocked 任务。
- `context proposal`：拿本轮可消费的上下文包和 deep read 建议。

禁止把 `context proposal` 当成全量 proposal 加载器。它可以包含 inspect 摘要，但不能替代 inspect 的体检语义。

### 6.2 选择本轮工作模式

读取 inspect/context 后，Agent 只能选择一种主要模式推进：

| 模式 | 触发条件 | 主要产物 |
|------|----------|----------|
| 索引整理 | 新需求尚未进入需求索引树，或 STR 超过 15 条 | STR 更新、requirement 卡 |
| 需求澄清 | requirement 缺 Acceptance / Scope / Open Questions | requirement 更新、analysis task |
| 分析推进 | 存在未闭合 analysis task 或影响范围未知 | finding、design、log |
| library 发现 | 当前设计缺规范、模块知识、历史设计 | library 候选、链接、log |
| 设计定稿 | 分析结论稳定，能形成一个设计焦点 | design card |
| 任务拆分 | design card 足够明确，可交给 implement | ready / not_ready implementation task |

一轮可以产生多个卡片，但用户汇报时只突出一个主线，避免看起来像散乱批处理。

### 6.3 Reference 读取时机

SKILL 本体默认不加载 references。只有在需要时读取：

| reference | 读取条件 |
|-----------|----------|
| `card-templates.md` | 要创建或审查 requirement、analysis task、design、implementation task、log |
| `library-discovery.md` | 要查询或链接 library 中的 convention、module、design、finding |
| `walkthrough-flowforge-v2.md` | 正在迭代 SKILL 本身，或需要端到端示例验证 |

这样可以保持 SKILL 触发时上下文轻量。

### 6.4 停止条件

一轮 design 应在以下任一条件满足时停止并汇报：

- 已创建或更新本轮需要的 requirement / design / task / log。
- 发现关键问题需要用户确认。
- 生成了一个可执行的 ready implementation task。
- 生成了 not_ready task，并明确阻塞它的 open question。
- 需要切换到 implement / feedback / archive。

不要在同一轮里继续执行 implementation task。

## 7. 最复杂问题的实现方案

这里的“实现”指业务方案，不是代码。

### 7.1 上下文入口

设计 SKILL 的第一条命令应该是 `proposal inspect`，而不是直接 `card list` 或读取文件。

`proposal inspect` 应返回：

- proposal root 摘要
- 顶层需求索引入口
- 直接子索引摘要
- ready / not_ready / blocked / in_progress task 数量
- open question 摘要
- 最近关键 log 摘要

这能让 Agent 建立工作面，同时避免加载整个 proposal。

### 7.2 Library Discovery

library 查询不能依赖 Agent 的文件能力，必须作为 CLI 原语。

第一版实现方案：

1. `library suggest --for <card-id>` 从 sqlite 的 `card_index`、`card_search`、`card_link`、`card_backlink` 查候选。
2. 默认只返回摘要、命中原因、建议关系。
3. Agent 根据候选摘要选择少量卡。
4. Agent 用 `card read --summary/--section` 定点读取。
5. Agent 只把确认相关的卡链接到当前 analysis/design/task。

不做：

- 不让 Agent grep `02-library/`。
- 不一次性读取所有 convention / module。
- 不在第一版引入 embedding。

### 7.3 索引树维护

需求索引树必须有 `structure add/remove`，不能让 Agent 手写 STR markdown。

第一版语义：

- `structure add STR-X REQ-Y` 同时更新 STR 条目和 `indexes` 关系。
- 添加后超过 15 条必须提示拆分。
- 命令只维护索引条目，不编辑任意正文段。

这样可以避免 STR 卡变成长文档，也避免关系索引和 markdown 内容漂移。

### 7.4 基于假设的任务

design SKILL 可以提前拆 implementation task，但要保护 implement 阶段不误执行。

规则：

- 信息完整：创建 `ready` implementation task。
- 有假设但任务边界有价值：创建 `not_ready` implementation task。
- 关键问题未闭合：优先创建 analysis task。

`not_ready` task 必须链接 open question、finding 或 analysis task，`task ready` 不返回它。

## 8. CLI 前置清单

写入真实 SKILL 前，CLI 设计至少需要承认这些命令：

- `flowforge proposal inspect`
- `flowforge context proposal`
- `flowforge library suggest`
- `flowforge card search --scope library`
- `flowforge card read --summary/--section`
- `flowforge structure add/remove`
- `flowforge log create --kind`
- `flowforge task create --type a/i`
- `flowforge task ready --type a`
- `flowforge card link`
- `flowforge card related --direction backlinks`

其中代码实现可交给其它 agent；本草案只定义业务语义和 SKILL 依赖。

这些命令的输出分区、上下文裁剪和 ready / not_ready 判定详见 [Design SKILL CLI 契约设计](./design-skill-cli-contracts.md)。

## 9. 暂不进入 SKILL 本体

以下内容不要写入 `SKILL.md` 主体：

- 完整卡片模板。
- FlowForge v2 walkthrough 全文。
- library candidate 排序算法细节。
- 具体 markdown section 更新策略。
- archive / feedback / implement 的完整流程。

这些内容会扩大上下文，破坏薄适配器原则。

## 10. 可部署资产规格

这一节是给后续实现 agent 的落地规格。当前只设计，不在本轮创建 `assets/` 文件。

### 10.1 文件职责

未来部署资产建议拆成四个文件：

```text
assets/skills/flowforge-design/
  SKILL.md
  references/
    card-templates.md
    library-discovery.md
    walkthrough-flowforge-v2.md
```

职责边界：

| 文件 | 职责 | 不应包含 |
|------|------|----------|
| `SKILL.md` | 激活边界、固定启动命令、主流程、硬约束、单轮汇报格式 | 完整模板、示例 walkthrough、检索排序算法 |
| `card-templates.md` | requirement / analysis task / design / implementation task / log 的最小正文结构和审查标准 | CLI 命令语义、library 检索策略 |
| `library-discovery.md` | 如何通过 CLI 发现、筛选、读取、链接 library 卡 | 卡片正文模板、实现任务模板 |
| `walkthrough-flowforge-v2.md` | 用 FlowForge v2 自身验证 workflow 的端到端样例 | 真实项目专属规则、会过期的实现细节 |

### 10.2 `SKILL.md` 必须保持的约束

`SKILL.md` 应控制在“Agent 每次触发都值得加载”的长度内。它必须包含：

- 触发后先解析 project / proposal / inspect / context。
- CLI 是唯一写入路径。
- 不直接读 wiki 文件，不 grep `02-library/`。
- 需求先进入 STR 索引树，再拆 requirement。
- 不确定点优先变成 analysis task。
- library 上下文必须通过 `library suggest` / `card search --scope library` / `card read` 获取。
- design card 只表达一个稳定设计焦点。
- implementation task 只有在可执行时才进入 ready。
- 每轮记录 log，log 主动链接上下文卡。
- 中心卡不累计所有证据，依赖 backlink / sqlite 视图。
- 每轮结束必须汇报卡片、关系、缺口和下一步。

如果某条规则需要超过两三句话解释，应下沉到 reference。

### 10.3 Reference 读取协议

SKILL 本体只告诉 Agent “什么时候读 reference”，不默认加载 reference。

读取规则：

| 场景 | 读取文件 | 读取目的 |
|------|----------|----------|
| 要创建或审查卡片正文 | `card-templates.md` | 确认最小 section、ready 条件、禁止空心任务 |
| 要查询规范、模块、历史设计或 finding | `library-discovery.md` | 确认查询、筛选、定点读取、链接规则 |
| 要验证 SKILL 自身是否跑通 | `walkthrough-flowforge-v2.md` | 用固定用例检查 workflow 是否闭环 |

reference 读取后也不能突破硬约束：仍然不能直接读取目标项目 wiki 文件，仍然不能把 library 全文批量塞入上下文。

### 10.4 实现验收清单

后续实现 agent 创建部署资产后，应按以下清单验收：

- `SKILL.md` 没有完整卡片模板。
- `SKILL.md` 没有让 Agent 直接读写 `ff-wiki/` 或 `02-library/`。
- `SKILL.md` 明确 design skill 不执行 implementation task。
- `card-templates.md` 中的 implementation task 模板要求 requirement、design、constraints 和 acceptance。
- `library-discovery.md` 明确候选发现只返回摘要，全文读取必须定点。
- `walkthrough-flowforge-v2.md` 能覆盖 requirement index、analysis task、library discovery、design card、implementation task。
- 任一示例任务都不是只有标题。

## 11. 上下文选择核心算法

design SKILL 最容易失败的地方是上下文选择：拿少了会胡猜，拿多了会退化成 v1 长上下文。第一版采用固定的三层预算。

### 11.1 第一层：工作面摘要

每轮必须加载，但只能是摘要：

- 当前 project 和 proposal 指针。
- proposal root 摘要。
- 顶层 STR 和直接子 STR 摘要。
- ready / not_ready / blocked / in_progress 任务计数。
- 当前 open question 摘要。
- 最近关键 log 摘要。

来源命令：

```bash
flowforge proposal inspect <proposal-id>
flowforge context proposal --proposal <proposal-id>
```

这一层回答“现在 proposal 处于什么状态”，不回答所有历史细节。

### 11.2 第二层：本轮焦点卡

Agent 根据用户最新输入和 inspect/context 结果，选择 1 个主焦点：

- 一个 STR 条目
- 一个 requirement
- 一个 analysis task
- 一个 design card
- 一个 not_ready implementation task

本轮只围绕这个主焦点扩展上下文。可以读取它的摘要、必要 section、直接关系和 backlink 摘要，但不递归展开整张图。

选择优先级：

1. 用户明确点名的卡或需求。
2. 阻塞 ready task 的 open question。
3. 当前 proposal 中最影响拆任务的 analysis task。
4. 新需求尚未进入 STR 的索引整理。
5. 已稳定设计后的 implementation task 拆分。

### 11.3 第三层：定点深读

只有当本轮要做出设计判断或创建 ready task 时，才进入深读。

允许深读：

- 被焦点卡直接链接的 requirement / design / convention / module。
- `library suggest` 高相关候选中被 Agent 确认需要引用的少量卡。
- 与 open question 直接相关的 finding 或历史 decision。

不允许深读：

- proposal 下所有 requirement。
- proposal 下所有 log。
- library 下所有 convention。
- 与当前焦点只有弱关键词匹配的候选全文。

### 11.4 上下文不足时的动作

当上下文不足以继续设计时，Agent 不应扩大读取范围到整库，而应选择以下动作之一：

- 创建 analysis task，明确要查什么。
- 用更窄的 query 调用 `card search --scope library`。
- 读取候选卡的 `--summary` 或指定 section。
- 创建 open question 并询问用户。
- 创建 not_ready implementation task 暴露缺口，但不能把它交给 implement。

这条规则保证“不确定”会被转化为卡片或问题，而不是靠长上下文和 Agent 记忆硬撑。

## 12. 下一步

建议下一步不是直接写 CLI 代码，而是准备可交给实现 agent 的 SKILL 资产任务包：

1. `assets/skills/flowforge-design/SKILL.md`：只放第 3 节压缩版。
2. `assets/skills/flowforge-design/references/card-templates.md`：来自 workflow 第 14 章。
3. `assets/skills/flowforge-design/references/library-discovery.md`：来自 workflow 第 9 章。
4. `assets/skills/flowforge-design/references/walkthrough-flowforge-v2.md`：来自 workflow 第 16 章。

在真正创建 `assets/` 文件前，需要确认这些文件会部署到目标项目，并且符合“assets 是部署边界”的约束。
