# 卡片架构不变量与修正方案

> 日期：2026-06-14  
> 状态：已固化为实现约束

本文固化本轮自举实践暴露出的卡片架构问题，以及 FlowForge CLI 必须保证的不变量。目标是让 proposal 可以从 `ROOT-<proposalId>` 稳定进入，沿受控导航理解核心需求、设计、任务，同时避免日志、证据、发现卡造成中心卡链接爆炸。

## 1. 问题结论

当前自举过程中暴露出三类框架级问题：

1. 正文链接使用 Obsidian wikilink，但文件名是 `ID_slug.md`，导致 `[[CARD-ID]]` 在 Obsidian 外不可用，在多数场景下也无法定位真实文件。
2. ROOT、STR、REQ、TASK 的关系边界不清，导致 requirement index 混入错误类型，子任务缺少真实父子关系。
3. CLI 允许写入断链、孤儿卡片和错误索引，破坏“从 ROOT 串起 proposal”的基本前提。

修正方向：

- frontmatter `links` 是结构化事实关系。
- 正文内部卡片导航由 FlowForge 生成，渲染为标准 Markdown 链接，作为人类可读导航层。
- CLI 写入前必须保证显式链接目标存在，validate 必须发现断链、wikilink、错误索引和孤儿卡片。

## 2. 三层模型

### 2.1 事实层

事实层写在 YAML frontmatter 的 `links` 中，是查询、上下文组装、索引重建的来源。

事实层遵循“子卡主动指向上游”的原则，避免中心卡反复回写：

- log 主动 `records -> TASK/REQ/DES/ROOT`
- finding 主动 `discovers -> TASK/REQ/DES`
- task 主动 `implements/designs/satisfies/requires/constrains -> 上游卡`
- proposal 内普通卡主动 `belongs_to -> ROOT-<proposalId>`

### 2.2 导航层

导航层写在正文中，面向人类阅读。内部卡片导航只能由 FlowForge CLI 根据 frontmatter 关系格式化并插入，渲染结果使用标准 Markdown 链接。

Agent 不手写内部卡片链接。Agent 在卡片正文中手写 Markdown 链接时，只能用于外部资料引用。

导航层只保留主线，不收纳所有证据：

- `ROOT -> STR-<proposal>-REQ`
- `STR -> requirement 或 child STR`
- `REQ -> 关键 analysis/design/main task`
- `DES -> main implementation task`

log、finding、临时 evidence 不直接进入 ROOT/STR 主线导航。需要看证据链时通过 `card related --direction backlinks` 或 sqlite 索引视图查询。

### 2.3 约束层

约束层由 CLI 和 `validate` 实现：

- 创建/修改显式链接前必须检查目标卡存在。
- 正文中 FlowForge 生成的本地 Markdown 链接必须能解析到真实文件。
- wikilink 不是 FlowForge 标准链接格式，validate 应报错。
- 除 proposal root 与 `STR-HOME` 外，所有卡片必须至少有一个 outbound frontmatter link。
- requirement index 只能 `indexes` requirement 或 structure。
- 文件名创建后不因 title 更新而重命名；validate 只检查文件名以 card ID 开头。

## 3. 卡片类型边界

| 类型 | 前缀 | 职责 |
|------|------|------|
| `proposal` | `ROOT` | proposal 稳定入口，只由 `proposal create` 创建 |
| `structure` | `STR` | 索引卡，组织 7-15 张同主题卡片或子索引 |
| `requirement` | `REQ` | 原子需求点 |
| `design` | `DES` | 针对需求或任务的设计方案 |
| `task` | `TASK` | 可执行任务或分析任务 |
| `log` | `LOG` | proposal 生命周期内的过程记录 |
| `finding` | `FIND` | 探索发现 |
| `decision` | `DEC` | 技术或产品决策 |
| `convention` | `CONV` | 可执行规范或约定 |
| `module` | `MOD` | 系统/模块认知 |

`ROOT` 不再伪装成 `structure`。它可以指向顶层 STR，但自身不是索引树节点。

## 4. 关系模型

| relation | 语义 | 典型方向 |
|----------|------|----------|
| `belongs_to` | 归属某 proposal/root | `REQ/DES/TASK/LOG -> ROOT` |
| `indexes` | STR 或 ROOT 纳入导航入口 | `ROOT -> STR`，`STR -> REQ/STR` |
| `decomposes` | 子任务指向父任务，表达父任务被拆解 | `TASK-sub -> TASK-parent` |
| `analyzes` | 分析任务分析某需求/索引 | `TASK-a -> REQ/STR` |
| `designs` | 设计卡或任务设计某需求/任务 | `DES/TASK -> REQ/TASK` |
| `implements` | 实施任务实现某设计 | `TASK -> DES` |
| `satisfies` | 任务或设计满足某需求 | `TASK/DES -> REQ` |
| `requires` | 执行需要某输入、需求或知识 | `TASK/DES -> REQ/DEC/FIND` |
| `constrains` | 规范约束当前工作 | `TASK/DES -> CONV` |
| `records` | 日志记录某对象过程 | `LOG -> TASK/REQ/DES/ROOT` |
| `discovers` | 发现来自某对象或探索 | `FIND -> TASK/REQ/DES` |
| `references` | 一般引用 | 任意卡指向参考卡 |
| `blocks` | 任务阻塞关系 | `TASK -> TASK` |
| `supports/questions/refines/extends/contradicts/supersedes/related/produced` | 兼容性关系 | 按具体语义使用 |

`related` 只能作为弱关联兜底，不用于父子任务、归属或索引。

## 5. Requirement Index 规则

顶层需求索引是 `STR-<proposalId>-REQ`。

它的 `indexes` 目标只允许：

- `requirement`
- `structure`

不允许进入 requirement index 的类型：

- `design`
- `task`
- `log`
- `finding`
- `decision`
- `convention`
- `module`

索引卡直接条目建议上限为 15 条。超过时应裂变为子 STR，形成索引树。

## 6. CLI 写入规则

### 6.1 proposal create

`flowforge proposal create` 必须创建：

- `ROOT-<proposalId>.md`，type 为 `proposal`
- `STR-<proposalId>-REQ.md`，type 为 `structure`

ROOT：

- `indexes -> STR-<proposalId>-REQ`
- 正文 Entries 由 FlowForge 生成，使用 Markdown 链接指向 `STR-<proposalId>-REQ.md`

REQ index：

- `belongs_to -> ROOT-<proposalId>`
- 初始 Entries 为 `- None`

### 6.2 card create

proposal 作用域内创建的普通卡片必须自动补：

- `belongs_to -> ROOT-<proposalId>`

非 proposal 作用域创建卡片时，如果没有任何 `--links`，CLI 必须拒绝写入。

### 6.3 task sub

`flowforge task sub <task-id>` 创建的子任务必须：

- 使用 `{parentTaskId}-a/b/c` ID
- `decomposes -> <parentTaskId>`
- 如果父任务属于 proposal，则同时 `belongs_to -> ROOT-<proposalId>`

### 6.4 structure add/refresh

`structure add/refresh` 必须同步维护：

- frontmatter `indexes` 关系
- 正文 `## Entries` 中由 FlowForge 生成的标准 Markdown 相对路径链接

CLI 不再生成 `[[wikilink]]`。

## 7. Library 导入统一原则

library 导入有两类来源：

- 外部长文或参考资料，经独立 skill 拆分为候选原子卡片
- proposal 归档时沉淀出的稳定知识候选

两者进入 library 前的来源处理不同，但进入 library 的流程相同：

1. 形成候选卡片。
2. 检查粒度、类型、标题、正文自足性。
3. 建立必要 outbound link。
4. 写入 library。
5. 重建 sqlite 派生索引。

FlowForge 框架只固化有限 card type 和 relation，不把项目词汇表、facet、tag 当作强语义匹配基础。tag 可用于粗分类和召回提示，但最终链接由 Agent 读取摘要后确认。

## 8. 实施计划

本轮优先完成可阻断坏数据继续产生的能力：

1. 固化 `proposal` 根卡类型和关系白名单。
2. 将 CLI 生成的正文导航从 wikilink 改为 Markdown 链接。
3. `validate` 检查 frontmatter 断链、正文 Markdown 断链、wikilink、孤儿卡片、错误 requirement index。
4. `card create/task create/task sub/log create/structure add` 在写入前检查目标存在和必要链接。
5. 更新测试覆盖上述不变量。
6. 同步知识系统与 CLI 设计文档。
7. 提供 `card refresh` 刷新 REQ/DES 的 CLI 生成导航。
8. 提供 `library import/promote` 支撑结构化候选入库和 proposal 知识沉淀。

后续再推进：

- library ingestion 的 scan/plan/apply 批处理命令。
- 基于 sqlite 的 backlink/navigation 视图。
