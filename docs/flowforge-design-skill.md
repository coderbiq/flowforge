# flowforge-design 正式草案

> 版本：draft-1
>
> 目标：把 `flowforge-design` 从“写设计文档”收敛成“驱动 proposal 卡片生长”的正式草案，明确它的触发边界、工作流、产物、CLI 前置能力和 reference 拆分方式。

## 1. 设计定位

`flowforge-design` 的职责不是生成长设计文档，而是把用户的模糊需求推进成一组可追踪、可迭代、可验证的卡片网络。

它的核心判断是：

- 先建立 proposal 工作面，再开始分析。
- 先拆需求索引和原子 requirement，再考虑实现任务。
- 先把不确定点转成 analysis task 或 open question，再允许结论稳定。
- 先通过 CLI 发现 library 上下文，再决定是否链接到设计或任务。

## 2. 触发边界

### 2.1 应该触发

当用户在做以下事情时，应使用 `flowforge-design`：

- 分析一个新需求。
- 澄清、拆解、补全已有 proposal。
- 讨论某个功能应该怎么设计、影响哪些模块、需要哪些任务。
- 在实现前补齐需求边界、设计约束或依赖分析。

典型表达：

- “帮我分析这个需求”
- “设计一下这个功能”
- “把这个 proposal 拆成任务”
- “先梳理一下边界”
- “这个实现前要先补设计”

### 2.2 不应该触发

以下场景不应由 `flowforge-design` 主导：

- 执行已有任务，交给 implement 流程。
- 处理测试失败或 bug 反馈，交给 feedback 流程。
- 归档、沉淀知识，交给 archive 流程。
- 只查询已有卡片或上下文，不进入设计推进。

## 3. 核心产物

`flowforge-design` 的输出不是单个文档，而是一组卡片和关系。

| 产物 | 作用 |
|------|------|
| proposal root card | 提案入口，只承载摘要、状态和导航 |
| 需求索引树 | 管理需求地图，控制单卡条目上限 |
| requirement 卡 | 原子需求点，一个用户可感知的行为、约束或验收点 |
| analysis task 卡 | 驱动调研、澄清和影响分析 |
| design 卡 | 一个稳定的设计结论 |
| finding 卡 | 分析中可复用或可追溯的事实、限制、风险 |
| log 卡 | 过程事件和证据链 |
| implementation task 卡 | 可交给 implement 执行的任务单元 |

## 4. 总体流程

```text
用户提出需求
  -> 解析当前 project / proposal
  -> 确认 proposal root card 和需求索引入口
  -> 更新需求索引树
  -> 拆原子 requirement
  -> 识别不确定点并创建 analysis task
  -> 查询 library 上下文
  -> 形成 finding / design / log
  -> 判断是否可拆 implementation task
  -> 输出下一步
```

这个流程不是一次性阶段，而是可回退的循环：

- 新反馈可以补 requirement。
- 新发现可以补 finding 或 design。
- 方案变化可以追加 analysis task。
- 上下文不足时，回到分析或追问用户。

## 5. 第一次进入流程时的动作

### 5.1 先确认工作面

第一轮必须先确认：

- 当前 project 是哪个。
- 当前 proposal 是哪个。
- 是否已有 proposal root card。
- 是否已有顶层需求索引卡。
- 当前 proposal 下有哪些活跃任务、open question 和关键 log。

建议 CLI：

```bash
flowforge project current
flowforge proposal current
flowforge proposal inspect <proposal-id>
flowforge context proposal --proposal <proposal-id>
```

`proposal inspect` 用于判断 proposal 状态和缺口，`context proposal` 用于拿到本轮可消费的最小上下文。两者不能合并成一个“大上下文”接口。

### 5.2 上下文加载原则

第一轮只加载摘要，不加载整套 proposal：

- proposal root 摘要。
- 顶层 STR 和直接子 STR 摘要。
- 活跃 analysis / task 摘要。
- 与当前需求直接相关的候选 library 摘要。

不读取完整 proposal 卡片集合，不读取所有日志。

### 5.3 记录启动 log

当本轮开始分析新需求或明显调整设计方向时，应创建一张 log 卡，记录：

- 用户原始意图摘要。
- 当前 proposal。
- 本轮目标。
- 关联的 root card 或 STR 卡。

## 6. 需求索引树规则

### 6.1 初始索引

`proposal create` 默认只创建空的顶层需求索引入口。只有在命令显式提供初始摘要或初始需求条目时，才写入 3-5 个粗粒度条目。

如果历史 proposal 只有 root card，没有索引入口，`flowforge-design` 应补建顶层需求索引卡。

### 6.2 拆原子 requirement

当一个索引条目已经表达独立用户目标、行为、约束或验收点时，应拆成 requirement 卡。

拆卡标准：

- 可单独验证。
- 能被一个或多个任务满足。
- 不依赖整篇 proposal 才能理解。
- 可以关联设计、任务或反馈。

不要把整个功能模块、整个页面族或完整后端实现写成一张 requirement 卡。

### 6.3 索引裂变

单张需求索引卡的健康范围是 7-15 个直接条目。超过 15 条必须拆分；接近 15 条且出现清晰子主题时，应提前拆分。

拆分时：

- 父索引卡保留主题入口。
- 同主题需求替换成子索引卡链接。
- 子索引卡只索引同一模块、用户目标或业务主题下的需求。
- 历史任务和需求卡不批量迁移。

### 6.4 索引边界

需求索引卡只维护导航关系，不承载所有证据。

应该写入：

- 子索引卡链接
- requirement 卡链接
- 少量状态摘要

不应该写入：

- 全量 log
- 全量 finding
- 全部实现记录
- 大段设计说明

## 7. 分析任务规则

### 7.1 什么时候创建

出现以下情况时应创建 analysis task：

- 需求边界不清楚。
- 需要看代码才能判断影响范围。
- 需要查 library 中的规范或历史设计。
- 需要比较多个方案。
- 涉及跨项目、前后端、数据模型或兼容性风险。
- 用户信息不足以直接拆实现任务。

### 7.2 必须包含什么

analysis task 至少包含：

- 分析目标
- 输入来源
- 需要检查的代码目录或模块
- 需要查询的 library 主题
- 预期输出
- 完成条件

### 7.3 链接规则

analysis task 应主动链接：

- 被分析的 requirement 或 STR
- 相关 module / convention / decision
- 过程中生成的 log 和 finding

analysis task 不应持续回写所有证据卡。

## 8. Library Discovery 规则

`flowforge-design` 需要能找到 library 中已有的规范、模块知识、历史设计和 finding，但不得直接遍历 wiki 文件或 grep markdown。

### 8.1 查询原则

library 查询分三层：

1. 候选发现。
2. 结构化筛选。
3. 定点读取。

禁止行为：

- 直接读取 `02-library/` 文件。
- 用 shell grep 遍历卡片库。
- 没有筛选就批量读全文。
- 把不确定相关的候选全部链接进来。

### 8.2 候选发现

候选发现应由以下信息驱动：

- requirement / analysis task 的 title、summary、tags
- 用户原始描述中的领域词和模块词
- 当前 project 的目录和技术栈线索
- 已有关联卡的 tags、domain、type

推荐命令：

```bash
flowforge library suggest --for REQ-xxx --types convention,module,design,finding --limit 10
flowforge library suggest --for TASK-xxx --relation constrains --limit 10
flowforge card search "分页 查询 条件" --scope library --type convention,module,design --limit 10
```

这一层只返回摘要和匹配理由，不返回全文。

### 8.3 定点读取

只有当候选卡满足以下条件之一时，才读取全文：

- 将被链接到 analysis task、design card 或 implementation task。
- 需要确认某条规范是否适用。
- 需要引用历史设计中的关键约束。

推荐命令：

```bash
flowforge card read CONV-012
flowforge card read MOD-004 --summary
```

### 8.4 关联写入

确认相关后，才写入链接：

- analysis task `references -> MOD-* / DEC-* / FIND-*`
- analysis task `constrains -> CONV-*`
- design card `references -> DEC-* / FIND-*`
- design card `constrains -> CONV-* / MOD-*`
- implementation task `constrains -> CONV-*`

不要把所有候选 library 卡都链接到当前卡。未确认但可能相关的候选可以写进 log 作为查询证据。

## 9. 设计卡规则

### 9.1 什么时候创建设计卡

当分析形成稳定设计结论时，创建 design 卡。

适合表达：

- 一个接口或命令的行为设计
- 一个数据结构或配置结构
- 一个模块协作方式
- 一个错误处理或状态流转规则
- 一个跨项目协作边界

不适合表达整个 proposal 的完整方案。

### 9.2 必须包含什么

design 卡至少包含：

- 设计目标
- 关联需求
- 关键决策
- 约束和不做什么
- 影响范围
- 验证方式
- 后续可拆任务

### 9.3 链接规则

设计卡应主动链接：

- 关联 requirement
- 相关 decision / finding
- 相关 convention / module

设计讨论过程的 log 应主动链接 design 卡。

## 10. 实现任务规则

`flowforge-design` 负责生成 implementation task，但不负责执行任务。

### 10.1 什么时候可以生成 ready 任务

满足以下条件时，才能生成 ready 的 implementation task：

- 需求点已经原子化。
- 关键 design 卡已经存在。
- 任务目标和验收标准明确。
- 相关规范或模块边界已经关联。
- 不存在阻塞实现的 open question。

如果还缺信息，应先创建 analysis task 或追问用户。

### 10.2 假设任务

如果必须基于假设先拆任务，该任务不得进入 ready，必须标记为 `blocked` 或 `not_ready`，并链接对应的假设、风险或 open question。

默认策略仍然是先创建 analysis task，而不是把假设任务直接交给 implement。

### 10.3 必须包含什么

implementation task 至少包含：

- 任务目标
- 输入来源
- 预期输出
- 关联需求
- 关联设计
- 关联规范
- 验收方式
- 不做事项

## 11. Log 与 finding 规则

### 11.1 log

以下事件应记录 log：

- 开始一轮需求分析。
- 用户补充关键约束。
- Agent 做出重要设计假设。
- 分析任务完成。
- 方案发生重要调整。
- 发现需要 feedback 或 implement 跟进的问题。

每张 log 只记录一个事件，并主动链接上下文卡。

### 11.2 finding

finding 是可复用或可追溯的认知，不是流水账。

适合创建 finding：

- 代码现状与预期不一致。
- 发现历史约束、兼容性问题或隐含规则。
- 某种实现方式不可行。
- 某个规范或模块边界需要后续任务遵守。

如果只是“我创建了某张卡”或“我读了某个文件”，应使用 log，不创建 finding。

## 12. CLI 前置能力

写入真实 SKILL 前，CLI 设计至少要承认这些命令：

- `flowforge project current`
- `flowforge proposal current`
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

代码实现可交给其它 agent；本草案只定义业务语义和 SKILL 依赖。

## 13. Reference 拆分

正式 SKILL 不应把模板和 walkthrough 全塞进去。建议拆成三个 reference：

```text
assets/skills/flowforge-design/
  SKILL.md
  references/
    card-templates.md
    library-discovery.md
    walkthrough-flowforge-v2.md
```

### 13.1 `card-templates.md`

包含：

- requirement 模板
- analysis task 模板
- design 模板
- implementation task 模板
- log 模板
- 单轮输出格式

### 13.2 `library-discovery.md`

包含：

- library 查询三层模型
- 候选筛选规则
- 定点读取规则
- 关联写入规则
- 未命中处理

### 13.3 `walkthrough-flowforge-v2.md`

包含：

- 用 FlowForge v2 设计自身的端到端示例
- 从 proposal inspect 到 implementation task 的完整流程
- 验证点和失败信号

## 14. SKILL 本体应保留什么

`SKILL.md` 只保留每次触发都必须加载的短约束：

- 先解析 project / proposal。
- 先 inspect，再 context。
- 按需求索引 -> analysis task -> design card -> implementation task 推进。
- 通过 CLI 查 library，不直接读 library 文件。
- 每轮必须写 log。
- 中心卡不回写全部证据。
- 设计不足时追问用户或创建 analysis task。
- implementation task 生成后交给 implement。

如果某条规则需要超过两三句话解释，应下沉到 reference。

## 15. 完成标准

一轮 design 迭代结束时，至少要满足以下之一：

- 当前需求已经进入需求索引树。
- 关键 requirement 已拆出，或明确保留为待澄清索引项。
- 不确定点已变成 analysis task 或 open question。
- 已形成的结论已写入 design / finding / log。
- 可执行工作已变成 implementation task。
- 下一步是继续 design、进入 implement、等待用户确认，还是回到分析，已经明确。

## 16. 当前结论

这一版正式草案已经可以作为后续实现 `assets/skills/flowforge-design/SKILL.md` 的依据。

下一步应做的是把：

1. `SKILL.md` 压缩成薄适配器。
2. `card-templates.md` 补成最小模板。
3. `library-discovery.md` 补成查询和筛选规则。
4. `walkthrough-flowforge-v2.md` 补成验证用例。
