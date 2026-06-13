# Design SKILL 工作流设计

> 版本：draft
>
> 目标：定义 `flowforge-design` 如何引导 Agent 从一个模糊需求开始，逐步建立 proposal 的需求索引树、需求卡、分析任务、设计卡和实现任务，让 proposal 的分析设计过程可追踪、可迭代、可被后续 implement / feedback / archive 消费。
>
> CLI 输出契约参考：[Design SKILL CLI 契约设计](./design-skill-cli-contracts.md)。

## 1. 设计定位

`flowforge-design` 不是“生成设计文档”的 SKILL，而是“驱动需求分析与设计卡片生长”的 SKILL。

它要解决 v1 的几个核心问题：

- 不再把 proposal 写成一篇长文档。
- 不再一进入 SKILL 就加载完整 proposal context。
- 不再让任务只有标题。
- 不再让 Agent 直接跳到方案结论。
- 不再把过程记录塞进单个 `notes.md`。

`flowforge-design` 的职责是把用户需求转化为一组可执行、可追溯、粒度受控的卡片网络，并在设计足够稳定时生成实现任务。

## 2. 触发边界

### 2.1 应该触发

当用户表达以下意图时，应触发 `flowforge-design`：

- 创建、分析、澄清、拆解一个新需求。
- 对已有 proposal 继续补充需求或设计。
- 把模糊想法整理成可执行任务。
- 讨论某个功能应该怎么做、影响哪些模块、需要哪些任务。
- 发现实现前需要先补设计、补需求边界、补约束。

典型用户表达：

- “帮我分析这个需求”
- “设计一下这个功能”
- “把这个 proposal 拆成任务”
- “这个需求还不清楚，先梳理一下”
- “继续完善当前 proposal 的设计”

### 2.2 不应该触发

以下场景不应由 `flowforge-design` 主导：

- 用户明确要求执行某个已存在任务：交给 `flowforge-implement`。
- 用户反馈测试失败、bug、遗漏需求：优先交给 `flowforge-feedback`。
- 用户要求沉淀知识或关闭 proposal：交给 `flowforge-archive`。
- 用户只想查询已有卡片或上下文：使用 CLI 查询，不进入设计流程。

## 3. 核心产物

`flowforge-design` 的输出不是单个文档，而是一组卡片和关系。

| 产物 | 作用 |
|------|------|
| proposal root card | proposal 的稳定入口，只承载摘要、状态和导航 |
| 需求索引树 | 管理需求地图，单张索引卡直接条目上限 15 |
| requirement 卡 | 原子需求点，一个用户可感知的功能点或约束 |
| analysis task 卡 | 驱动调研、澄清、影响分析、边界分析 |
| design 卡 | 一个接口、行为、模块协作或技术方案的设计结论 |
| finding 卡 | 分析过程中发现的事实、限制、风险或复用知识 |
| log 卡 | design 过程中的事件记录和证据链 |
| implementation task 卡 | 设计稳定后可交给 implement SKILL 执行的任务 |

## 4. 总体流程

```text
用户提出需求
  -> 解析当前 project / proposal
  -> 确保 proposal root card 和需求索引入口存在
  -> 建立或更新需求索引树
  -> 拆出原子需求卡
  -> 识别不确定点和分析任务
  -> 执行或排队分析任务
  -> 形成 finding / design / log 卡
  -> 判断设计是否足够拆任务
  -> 生成 implementation task 卡
  -> 输出下一步建议
```

这个流程不是一次性线性阶段。任意一步都可以回退：

- 新反馈可以补需求卡。
- 新发现可以补 finding 或 design 卡。
- 设计变化可以追加新的 analysis task。
- 实现前发现上下文不足，可以回到 design 补卡。

## 5. 进入流程时的第一轮动作

Agent 激活 `flowforge-design` 后，第一轮动作必须先建立工作面，而不是直接写方案。

### 5.1 解析当前工作上下文

Agent 应先确认：

- 当前 project 是哪个。
- 当前 proposal 是哪个。
- 是否已有 proposal root card。
- 是否已有顶层需求索引卡。
- 当前 proposal 下已有多少需求、任务、设计、log。

建议 CLI：

```bash
flowforge project current
flowforge proposal current
flowforge proposal inspect <proposal-id>
flowforge context proposal --proposal <id>
```

`proposal inspect` 用于体检 proposal 状态和缺口；`context proposal` 用于拼装本轮 Agent 可消费的上下文。两者可以共享底层数据，但 CLI 语义不能合并，避免 `context proposal` 重新膨胀成 v1 式长上下文。

如果没有当前 project，提示先运行 `flowforge project use` 或 `flowforge project create`。

如果没有当前 proposal，提示先运行 `flowforge proposal create "<title>"`。当用户已经明确要创建新 proposal 时，design SKILL 可以建议创建，但不应绕过 proposal 命令直接写目录。

### 5.2 读取最小上下文

第一轮通过 `context proposal` 只读取：

- proposal root card 摘要
- 顶层需求索引卡摘要
- 当前活跃 analysis task 摘要
- 与用户最新需求直接相关的候选 library 卡摘要

不读取完整 proposal 卡片集合，不读取所有日志。

### 5.3 记录 design 启动 log

当本轮确实开始分析新需求或大幅更新设计时，应创建一张 log 卡记录入口事件：

- 用户原始意图摘要
- 当前 proposal
- 本轮目标
- 关联 root card 或需求索引卡

log 卡主动 `records -> ROOT-<proposal>` 或 `records -> TASK-...`，不回写 root card。

## 6. 需求索引树工作流

### 6.1 初始需求索引

`proposal create` 默认创建空的顶层需求索引入口。只有当创建命令显式提供初始摘要或初始需求条目时，才把 3-5 个粗粒度条目写入索引卡。

如果历史 proposal 只有 root card，没有需求索引入口，`flowforge-design` 应创建顶层需求索引卡。

顶层需求索引卡的内容应保持短小：

- 3-5 个初始需求主题
- 明确哪些是已知、待澄清、待分析
- 指向已存在的 requirement 卡或子索引卡

它不是需求文档，不写完整方案。

### 6.2 拆原子需求卡

当索引条目已经表达一个独立用户目标、行为、约束或验收点时，应拆成 requirement 卡。

拆卡标准：

- 能被单独验证。
- 能被一个或多个任务满足。
- 不依赖整篇 proposal 才能理解。
- 可以关联设计、任务或反馈。

不应该把“整个功能模块”“整个页面族”“完整后端实现”写成一张 requirement 卡。

### 6.3 索引裂变

单张需求索引卡的健康范围是 7-15 个直接条目。超过 15 条时必须拆分；接近 15 条且已经出现清晰子主题时，应提前拆分。

拆分规则：

- 父索引卡保留主题级入口。
- 一组同主题需求替换成子索引卡链接。
- 子索引卡只索引同一模块、用户目标或业务主题下的需求。
- 已有任务或需求卡不批量迁移历史链接。
- 后续新卡指向新的子索引卡即可。

### 6.4 索引更新边界

需求索引卡只维护导航关系，不承载所有证据。

应该写入：

- 子索引卡链接
- 需求卡链接
- 少量状态摘要

不应该写入：

- 所有 log
- 所有 finding
- 所有实现任务执行记录
- 大段设计说明

这些内容通过反向链接和 context 查询获取。

## 7. 分析任务工作流

design SKILL 的关键不是 Agent 自己臆测所有答案，而是把不确定点变成 analysis task。

### 7.1 什么时候创建 analysis task

出现以下情况时应创建 analysis task：

- 需求边界不清楚。
- 需要查看现有代码才能判断影响范围。
- 需要查询 library 中既有规范或历史设计。
- 多个方案需要比较。
- 涉及跨项目、前后端、数据模型或兼容性风险。
- 用户问题不足以直接拆实现任务。

### 7.2 analysis task 必须包含

analysis task 不能只是标题，应至少包含：

- 分析目标
- 输入来源：需求索引卡、requirement 卡、用户原始描述
- 需要检查的代码目录或模块
- 需要查询的 library 主题
- 预期输出：finding / design / requirement update / implementation task
- 验收方式：什么信息补齐后可结束分析

### 7.3 analysis task 的链接

analysis task 主动链接稳定上下文：

- `analyzes -> REQ-*` 或 `analyzes -> STR-*`
- `references -> MOD-* / CONV-* / DEC-*`
- `constrains -> CONV-*`
- 子任务使用 `decomposes -> TASK-*`

分析过程中产生的 log/finding 主动链接 analysis task。analysis task 不持续回写所有证据卡。

## 8. 用户引导规则

`flowforge-design` 必须主动引导用户，而不是只输出设计结论。

### 8.1 何时追问用户

当缺失信息会影响需求方向、验收标准、数据边界或任务拆分时，应追问用户。

典型追问点：

- 目标用户是谁。
- 哪些场景必须支持，哪些可以暂缓。
- 输入输出和验收标准是什么。
- 是否存在兼容旧数据、权限、性能、迁移要求。
- 前后端或多项目之间的边界怎么划分。
- 哪些行为是用户明确要求，哪些只是 Agent 推断。

### 8.2 追问方式

一次只问关键问题，避免把用户变成填表。

推荐输出：

- 先说明当前已识别出的需求地图。
- 再列出阻塞设计推进的 1-3 个问题。
- 同时说明如果用户暂时不回答，Agent 可以基于哪些假设继续推进。

### 8.3 推断必须落卡

如果 Agent 基于假设继续推进，必须把假设写入卡片：

- 不确定需求写入 requirement 卡或索引卡正文中的 `Open Questions` 段落。
- 设计假设写入 design 卡。
- 风险和限制写入 finding 卡。
- 推断过程写入 log 卡。

不能只在对话里说“我假设”。

## 9. Library Discovery 工作流

`flowforge-design` 需要能找到 library 中已有的规范、模块知识、历史设计和 finding，但 Agent 不应直接遍历 wiki 文件、grep markdown 或自行解析卡片目录。library 的查找、筛选、摘要和读取都必须通过 CLI。

### 9.1 查询原则

library 查询遵循三层模型：

1. **候选发现**：CLI 根据当前需求、任务或关键词返回候选卡片摘要。
2. **结构化筛选**：Agent 根据类型、标签、领域、关系和匹配理由缩小候选范围。
3. **定点读取**：只对少量高相关卡片调用 `flowforge card read` 读取全文。

禁止行为：

- Agent 直接读取 `02-library/` 文件。
- Agent 用 shell grep 遍历卡片库。
- Agent 在没有候选筛选的情况下批量读取 library 全文。
- Agent 把不确定相关的 library 卡全部链接到当前任务。

### 9.2 候选发现

当进入 design 或执行 analysis task 时，Agent 应从当前上下文提取查询线索：

- requirement / analysis task 的 title、summary、tags。
- 用户原始描述中的领域词、模块词、动作词。
- 当前 project 的 `srcDirs`、模块目录、技术栈线索。
- 已有关联卡的 tags、domain、type。
- 需求索引树中的父主题。

推荐 CLI：

```bash
flowforge library suggest --for REQ-xxx --types convention,module,design,finding --limit 10
flowforge library suggest --for TASK-xxx --relation constrains --limit 10
flowforge card search "分页 查询 条件" --scope library --type convention,module,design --limit 10
```

这一层只返回摘要和匹配理由，不返回全文。

推荐输出字段：

| 字段 | 说明 |
|------|------|
| id | 卡片 ID |
| type | 卡片类型 |
| title | 标题 |
| summary | 短摘要 |
| tags / domain | 筛选依据 |
| matchedBy | 命中原因，如 keyword / tag / relation / structure |
| suggestedRelation | 建议关系，如 `constrains` / `references` |
| score | 排序分数，仅作参考 |

### 9.3 结构化筛选

Agent 对候选结果继续筛选时，应优先考虑：

- `convention`：是否直接约束当前设计或实现任务。
- `module`：是否描述当前要修改的模块边界。
- `design`：是否是同类需求的历史设计。
- `decision`：是否是仍然有效的架构决策。
- `finding`：是否揭示风险、限制或历史坑点。

筛选规则：

- `must` / `should` 重要性高的卡优先。
- `active` / `accepted` 状态优先。
- 同 project / 同 domain / 同 module 的卡优先。
- 被当前需求索引树或相关模块索引卡引用的卡优先。
- deprecated / superseded 卡默认不进入上下文，除非是为了理解历史迁移。

如果候选过多，Agent 应继续使用 CLI 缩小条件，而不是读取更多全文。

### 9.4 定点读取

只有当候选卡满足以下任一条件时，才读取全文：

- 将被链接到 analysis task、design card 或 implementation task。
- 需要确认某条规范是否适用。
- 需要引用历史设计中的关键约束。
- 需要判断某个 finding 是否影响当前方案。

推荐 CLI：

```bash
flowforge card read CONV-012
flowforge card read MOD-004 --summary
```

如果 CLI 支持字段裁剪，design SKILL 应优先读取摘要、规则段、约束段，而不是全文。

### 9.5 关联写入

真正确认相关后，Agent 才写入链接。

推荐关系：

- analysis task `references -> MOD-* / DEC-* / FIND-*`
- analysis task `constrains -> CONV-*`
- design card `references -> DEC-* / FIND-*`
- design card `constrains -> CONV-* / MOD-*`
- implementation task `constrains -> CONV-*`
- implementation task `references -> MOD-* / DES-*`

不要把所有候选 library 卡都链接到当前卡。未确认但可能相关的候选可以写入 log，作为“查询过但未采用”的过程证据。

### 9.6 未命中时的处理

如果 library 没有命中可用规范或历史设计：

- 不要编造已有规范。
- 在 analysis task 或 design card 中标记“未找到现有约束”。
- 创建 log 记录查询条件和结果。
- 如果分析中形成可复用知识，创建 finding 卡。
- 等 archive 时再将稳定知识合成为 library 卡。

## 10. 设计卡工作流

### 10.1 什么时候创建设计卡

当分析已经形成一个稳定设计结论时，创建 design 卡。

设计卡适合表达：

- 一个接口或命令的行为设计。
- 一个数据结构或配置结构。
- 一个模块协作方式。
- 一个错误处理或状态流转规则。
- 一个跨项目协作边界。

不适合表达：

- 整个 proposal 的完整方案。
- 所有后端设计。
- 所有前端页面设计。
- 一串没有决策依据的实现步骤。

### 10.2 设计卡必须包含

- 设计目标
- 关联需求
- 关键决策
- 约束和不做什么
- 影响范围
- 验收或验证方式
- 后续可拆任务

### 10.3 设计卡链接

设计卡主动链接：

- `designs -> REQ-*`
- `references -> DEC-* / FIND-*`
- `constrains -> CONV-* / MOD-*`
- `satisfies -> REQ-*`，当设计直接满足某需求时使用

设计讨论过程的 log 主动 `records -> DES-*`。

## 11. 实现任务生成工作流

design SKILL 负责生成 implementation task，但不负责执行任务。

### 11.1 什么时候可以生成实现任务

满足以下条件时，可以生成 ready 状态的 implementation task：

- 需求点已经原子化。
- 关键设计卡已经存在。
- 任务目标和验收标准明确。
- 相关规范卡或模块边界已经关联。
- 不再存在阻塞实现的 open question。

如果还缺信息，应先创建 analysis task 或追问用户。

如果 Agent 必须基于假设先拆出 implementation task，该任务不得进入 ready，应标记为 blocked 或 `not_ready`，并链接对应假设、风险或 open question。默认策略是先创建 analysis task，而不是把假设任务直接交给 implement。

### 11.2 implementation task 必须包含

- 任务目标
- 输入来源
- 预期输出
- 关联需求
- 关联设计
- 关联规范
- 验收方式
- 不做事项

### 11.3 任务拆分原则

任务应按可验证交付物拆分，而不是按文件名机械拆分。

推荐拆分维度：

- 配置/schema 变更
- 命令或 API 行为
- 核心业务逻辑
- 测试验证
- 文档更新
- 迁移或兼容处理

任务过大时，先创建父任务，再拆子任务。子任务最多两层。

## 12. Log 与 finding 工作流

### 12.1 design 过程中的 log

以下事件应记录 log：

- 开始一轮需求分析。
- 用户补充了关键约束。
- Agent 做出重要设计假设。
- 分析任务完成。
- 方案发生重要调整。
- 发现需要 feedback 或 implement 跟进的问题。

每张 log 只记录一个事件，并主动链接上下文卡。

### 12.2 finding 的边界

finding 是可复用或可追溯的认知，不是普通流水账。

适合创建 finding：

- 代码现状与预期不一致。
- 发现历史约束、兼容性问题或隐含规则。
- 某种实现方式不可行。
- 某个规范或模块边界需要被后续任务遵守。

如果只是“我创建了某张卡”“我读了某个文件”，用 log，不创建 finding。

## 13. 完成一次 design 迭代的标准

一次 design 迭代不要求 proposal 完成，只要求本轮目标闭合。

本轮可以结束的条件：

- 用户当前需求已经进入需求索引树。
- 关键需求点已拆成 requirement 卡，或明确保留为待澄清索引项。
- 不确定点已变成 analysis task 或 open question。
- 已形成的结论写入 design/finding/log 卡。
- 可执行工作已变成 implementation task。
- 下一步应该 design、implement、feedback 还是等待用户输入已经明确。

输出给用户时，应简要说明：

- 本轮新增或更新了哪些卡片。
- 当前还缺什么。
- 下一步建议执行哪个任务或继续分析哪个问题。

## 14. 最小卡片模板

这一节定义 design SKILL 第一版必须稳定写出的卡片正文结构。模板只规定最小段落，不把卡片写成长文档。

### 14.1 Requirement 卡

Requirement 卡表达一个原子需求点。

最小结构：

```markdown
# <需求标题>

## Summary

一句话说明用户可感知的目标、行为或约束。

## Source

- Proposal: <proposal-id>
- Input: 用户原始描述或来源卡

## Acceptance

- 可验证条件 1
- 可验证条件 2

## Scope

- In:
- Out:

## Open Questions

- 未决问题；没有则写 None
```

写入要求：

- `Summary` 必须能脱离 proposal 长文独立理解。
- `Acceptance` 至少有一条可验证条件；如果暂时无法验证，应保留 open question，不要伪造验收标准。
- `Scope` 用来控制边界，避免需求卡膨胀成模块级长文档。

### 14.2 Analysis Task 卡

Analysis task 驱动调研、澄清和影响分析。

最小结构：

```markdown
# <分析任务标题>

## Goal

本次分析要回答的问题。

## Inputs

- Requirement / STR:
- User input:
- Known constraints:

## Investigation Plan

- 要查看的代码目录、模块或 library 主题
- 要比较的方案或要确认的边界

## Expected Outputs

- requirement update / finding / design card / implementation task

## Done When

- 什么信息补齐后可结束分析
```

写入要求：

- 必须链接被分析的 requirement 或 STR。
- 如果需要查询 library，要写明查询主题，不让 Agent 靠记忆猜规范。
- Analysis task 的完成产物必须落到 requirement、finding、design 或 task 卡，不能只留在对话中。

### 14.3 Design 卡

Design 卡表达一个稳定设计结论。

最小结构：

```markdown
# <设计标题>

## Goal

本设计解决哪个需求或分析结论。

## Decision

采用的设计结论。

## Rationale

为什么这样设计，引用哪些 requirement / finding / convention / module。

## Constraints

- 必须遵守的规范或模块边界
- 不做事项

## Impact

- 影响的命令、接口、配置、模块或项目

## Verification

- 如何验证设计被正确实现

## Follow-up Tasks

- 可拆出的 implementation task 草案
```

写入要求：

- `Decision` 只写一个设计焦点，不汇总整个 proposal。
- `Constraints` 必须来自已确认上下文或 library 候选读取结果；不能凭空引用不存在的规范。
- 设计仍有未决问题时，不应生成 ready 状态 implementation task。

### 14.4 Implementation Task 卡

Implementation task 是交给 implement SKILL 的执行单元。

最小结构：

```markdown
# <任务标题>

## Goal

本任务要交付的可验证结果。

## Inputs

- Requirement:
- Design:
- Constraints:

## Deliverables

- 需要修改或新增的行为、命令、配置、测试、文档

## Acceptance

- 验收条件

## Out of Scope

- 明确不做什么

## Read Before Work

- Agent 执行前必须读取的卡片 ID
```

写入要求：

- `Inputs` 至少关联 requirement 和 design；如果还缺 design，只能创建 `not_ready` task。
- `Read Before Work` 只列稳定上下文卡，不列执行过程 log。
- 任务必须能被单独认领和完成。

### 14.5 Log 卡

Log 卡记录一个过程事件。

最小结构：

```markdown
# <log 标题>

## Kind

progress | bug | finding | knowledge | blocked

## Event

发生了什么。

## Context

- Proposal:
- Related card / task:

## Result

本事件带来的结论、后续动作或无结果说明。
```

写入要求：

- 每张 log 只记录一个事件。
- log 主动 `records -> <context-card>`，不回写中心卡。
- log 不是最终结构化产物；可复用内容应进一步形成 finding / design / convention / module。

## 15. 单轮输出格式

design SKILL 每轮结束时，应向用户输出稳定格式，避免“写了很多但不知道系统状态如何”。

推荐格式：

```markdown
本轮完成：
- 新增/更新卡片：REQ-..., DES-..., TASK-..., LOG-...
- 建立关系：TASK-... analyzes REQ-...，DES-... constrains CONV-...

当前缺口：
- Open question:
- not_ready / blocked task:

下一步：
- 继续 design：执行 TASK-... 分析 ...
- 或进入 implement：执行 TASK-...
- 或等待用户确认：...
```

输出规则：

- 不贴完整卡片内容，只列 ID、标题和目的。
- 明确哪些任务是 `ready`，哪些是 `not_ready`。
- 如果需要用户回答，最多列 1-3 个关键问题。
- 如果可以进入 implement，只推荐一个最优先任务，不一次性展开所有任务。

## 16. FlowForge v2 Walkthrough

FlowForge v2 自身就是 design SKILL 的第一条验证用例。这个 walkthrough 用“设计 flowforge-design SKILL”作为需求，检验 workflow 是否能真实驱动卡片生长。

### 16.1 用户输入

示例输入：

```text
开始设计 flowforge-design skill，让 Agent 能从模糊需求开始建立需求索引树、分析任务、设计卡和实现任务。
```

### 16.2 进入工作面

Agent 先执行：

```bash
flowforge project current
flowforge proposal current
flowforge proposal inspect <proposal-id>
flowforge context proposal --proposal <proposal-id>
```

如果 proposal 不存在，先建议：

```bash
flowforge proposal create "FlowForge v2 design skill 工作流"
```

`proposal create` 默认生成：

- `ROOT-<proposal>.md`
- `STR-<proposal>-REQ.md`
- `90-cards/`

### 16.3 初始需求索引树

顶层需求索引卡先写 3-5 个主题入口：

- Design SKILL 触发边界
- 需求索引树和 requirement 卡生成
- Library Discovery
- Analysis task 与 design card
- Implementation task 生成与移交

如果某个主题继续增长，例如 Library Discovery 下面出现规范检索、模块检索、历史设计检索、候选筛选、定点读取、关联写入、未命中处理等 7 个以上条目，可以裂变为子索引：

```text
STR-<proposal>-REQ
  -> STR-<proposal>-REQ-LIBRARY-DISCOVERY
```

### 16.4 原子需求卡示例

从索引中拆出 requirement：

- `REQ-...`：Design SKILL 必须通过 CLI 查询 library 候选，而不是直接读取 `02-library/` 文件。
- `REQ-...`：Design SKILL 每轮结束必须输出新增卡、缺口和下一步。
- `REQ-...`：基于假设生成的 implementation task 必须是 `not_ready`。

每张 requirement 都写 `Acceptance` 和 `Open Questions`。例如：

```markdown
## Acceptance

- Agent 使用 `library suggest` 或 `card search --scope library` 获取候选摘要。
- Agent 只对确认相关卡片使用 `card read` 定点读取。
- Agent 不直接 grep 或遍历 `02-library/`。

## Open Questions

- 第一版是否需要 embedding？结论：不需要，预留接口即可。
```

### 16.5 Analysis task 示例

创建 analysis task：

```text
TASK-...-a-...：分析 design SKILL 需要哪些 MVP CLI 原语
```

链接：

- `analyzes -> STR-<proposal>-REQ`
- `references -> docs/cli-design 对应卡片`，如果已有文档卡
- `records <- LOG-*` 由后续 log 反向查询

任务目标：

- 判断哪些命令是 MVP 原语。
- 判断哪些命令应延后。
- 输出 CLI 调整建议。

完成后产生 finding：

- `FIND-...`：`card append-section` 过早设计会绑定未稳定的正文模板，应延后。
- `FIND-...`：`context requirement` 应等真实 design 轮次稳定后再封装。

### 16.6 Library Discovery 示例

Design SKILL 查询已有规范：

```bash
flowforge library suggest --for REQ-xxx --types convention,module,design,finding --limit 10
flowforge card search "SKILL design thin adapter" --scope library --type convention,design --limit 10
```

候选只返回摘要：

```text
CONV-001  convention  SKILL 薄适配器原则  matchedBy=tag:skill suggestedRelation=constrains
DES-004   design      v1 design context 问题 matchedBy=keyword:context suggestedRelation=references
```

Agent 只读取确认相关卡：

```bash
flowforge card read CONV-001 --section Rules
```

然后链接：

- design card `constrains -> CONV-001`
- analysis task `references -> DES-004`

### 16.7 Design card 示例

形成设计卡：

```text
DES-...：Design SKILL 采用“需求索引 -> analysis task -> design card -> implementation task”的循环
```

关键内容：

- Goal：让 Agent 不直接写长文档，而是驱动卡片网络生长。
- Decision：每轮先 inspect proposal，再围绕需求索引树推进。
- Constraints：SKILL 本体保持短，详细模板留在 docs/reference。
- Verification：用 FlowForge v2 自身需求跑一轮，能产出 requirement、analysis task、design、implementation task。

### 16.8 Implementation task 示例

当设计稳定后，生成 ready 任务：

```text
TASK-...-i-...：实现 flowforge-design SKILL 草案
```

Inputs：

- Requirement：Design SKILL 触发边界、Library Discovery、单轮输出格式
- Design：Design SKILL 工作流设计卡
- Constraints：SKILL 薄适配器原则、SKILL 文件短小、CLI 唯一入口

Acceptance：

- `assets/skills/flowforge-design/SKILL.md` 只包含流程和强约束。
- 详细模板不塞进 SKILL 本体。
- SKILL 明确禁止 Agent 直接读 library 文件。

如果还有未确认问题，例如卡片模板仍未稳定，则任务创建为 `not_ready`，并链接 open question。

### 16.9 Walkthrough 验证点

这条用例验证：

- proposal root card 是否足够作为入口。
- 需求索引树是否能避免长文档。
- analysis task 是否能驱动不确定点闭合。
- library 查询是否能通过 CLI 拿到规范和历史设计。
- design card 是否能表达稳定结论。
- implementation task 是否带足上下文。
- 单轮输出是否能让用户知道下一步该做什么。

如果 walkthrough 中某一步必须靠 Agent 直接读文件、手写复杂 markdown、或加载大量卡片全文，说明 CLI 原语或卡片模板还缺设计。

## 17. 反推 CLI 能力

根据 design SKILL workflow，CLI 能力分为两类：

- **MVP 原语**：没有这些命令，design SKILL 很难稳定执行。
- **延后封装**：有价值，但依赖卡片模板和真实上下文组成稳定后再设计。

原则是先补 SKILL 必须依赖的原语，不提前设计便捷封装。

### 17.1 MVP 必需能力

| 命令能力 | 用途 |
|----------|------|
| `flowforge context proposal` | 获取 root、需求索引入口、活跃任务摘要 |
| `flowforge proposal inspect` | 汇总 proposal root、索引树、任务状态和缺口 |
| `flowforge library suggest --for <card-id>` | 为需求、任务或设计查找 library 候选摘要 |
| `flowforge card search <query> --scope library` | 在 library 中做关键词和类型筛选 |
| `flowforge card read <id> --summary/--section <name>` | 定点读取摘要或指定段落，避免全文加载 |
| `flowforge card create --type requirement/structure/design/finding/log` | 创建核心设计卡片 |
| `flowforge task create --type a/i --status ready/not_ready` | 创建分析和实现任务 |
| `flowforge task ready --type a` | 找出可执行分析任务 |
| `flowforge structure add/remove` | 维护 STR 索引条目并处理 7-15 条上限提示 |
| `flowforge log create --kind <kind>` | 高频创建过程记录 |
| `flowforge card link <from> <to> --relation <rel>` | 写入类型化关系 |
| `flowforge card related <id> --direction backlinks` | 查询反向证据链 |
| `flowforge index rebuild` | 重建 sqlite 查询索引 |

### 17.2 MVP 命令语义补充

- `proposal inspect` 不返回所有卡片全文，只返回 proposal root、顶层需求索引、活跃任务、open question、not_ready / blocked 任务摘要。
- `structure add/remove` 只维护 STR 索引条目和关系，不承担任意 markdown 编辑；当直接条目超过 15 时必须提示拆分。
- `log create --kind` 负责统一 log frontmatter、note kind 和上下文链接，避免 design SKILL 手写 log 模板。
- `card read --summary/--section` 是 library 定点读取的基础能力，默认仍可读取全文，但 design SKILL 应优先使用裁剪读取。
- `task ready --type a` 是现有 `task ready` 的过滤能力；analysis task 还必须具备 Goal、Inputs、Investigation Plan、Expected Outputs、Done When。

### 17.3 延后封装能力

这些命令暂不进入 MVP。它们解决的是“更省步骤”，不是“能否稳定执行”。

| 命令能力 | 延后原因 | 触发条件 |
|----------|----------|----------|
| `flowforge card append-section` | 依赖 requirement / design / task / log 正文模板稳定；过早设计会变成脆弱的通用文本编辑器 | 至少 3 类卡片都反复出现“安全更新固定段落”的真实需求 |
| `flowforge context requirement --requirement <id>` | requirement 级上下文组成还不稳定，过早封装可能重新变成长 context | design SKILL 真实跑通几轮后，确认 requirement 级上下文的固定 must / should 组成 |

### 17.4 暂不优先

- 自动生成完整设计文档。
- 自动移动 proposal 卡片到 library。
- 大规模 library 整理和去重。
- 复杂图谱可视化。

这些能力可以等 design / implement / feedback 主链路稳定后再加。

## 18. SKILL 本体编写建议

真正的 `assets/skills/flowforge-design/SKILL.md` 应保持短小，作为薄适配器。

建议只包含：

- 触发后先解析 project/proposal。
- 先执行 `flowforge proposal inspect` 判断工作面，再执行 `flowforge context proposal` 获取本轮上下文。
- 按“需求索引 -> 分析任务 -> 设计卡 -> 实现任务”的顺序推进。
- 通过 `library suggest` / `card search --scope library` 查找规范和历史设计，不直接读 library 文件。
- 只读取少量确认相关的 library 卡全文，并把真正相关的卡链接到任务或设计卡。
- 每轮必须写 log。
- 中心卡不回写证据，证据卡主动链接上下文。
- 设计不足时追问用户或创建 analysis task。
- 实现任务生成后交给 implement SKILL，不在 design SKILL 中执行代码。

详细卡片模板、CLI 参数和判断规则应放在 docs 或 reference 中，不要塞进 SKILL 本体。

## 19. 已确认设计决策

以下决策作为后续 SKILL 编写和 CLI 调整的依据。

### 19.1 proposal 初始索引

`proposal create` 默认创建空的顶层需求索引入口。只有当命令显式提供初始摘要或初始需求条目时，才写入 3-5 个粗粒度条目。

这样可以保证 proposal 结构完整，同时避免 `proposal create` 在没有 design 语境时过早生成低质量需求内容。

### 19.2 Open Questions 表达方式

第一版不为 requirement schema 增加显式 `openQuestions` 字段。未决问题先放在 requirement 卡或 STR 索引卡正文的 `Open Questions` 段落中。

原因是 open question 的粒度和生命周期还需要通过真实使用校准，过早进入 frontmatter 会增加 schema 复杂度。

### 19.3 STR 条目上限

STR 索引卡采用 7-15 个直接条目的健康范围。超过 15 条必须拆分；接近 15 条且已经出现清晰子主题时，应提前拆分。

这个规则适用于需求索引树，也适用于 library 中的主题索引。

### 19.4 任务类型

保留 `analysis task`，暂缓独立 `design task`。

设计结论由 design card 表达；只有当存在明确的“分析/调研/澄清”待办时才创建 analysis task。实现工作由 implementation task 承接。

### 19.5 基于假设的实现任务

design SKILL 可以基于假设拆出 implementation task，但任务不得进入 ready。它必须标记为 blocked 或 `not_ready`，并链接对应假设、风险或 open question。

默认策略仍然是优先创建 analysis task 或追问用户，只有在任务拆分本身有助于暴露设计缺口时才创建 `not_ready` implementation task。

### 19.6 Library 检索实现

第一版 library 检索使用 sqlite FTS/BM25 + 元数据/关系排序，不引入 embedding。

索引设计可以预留 embedding 接口，但 workflow 和 CLI 不依赖它。这样可以先验证查询语义和卡片结构，再决定是否引入向量检索。

### 19.7 Library 命令边界

`library suggest` 作为独立命令存在。

它表达的是业务语义：为当前需求、任务或设计推荐规范、模块知识、历史设计和 finding。`card search --scope library` 保留为通用检索入口，`library suggest` 负责 workflow 级推荐。
