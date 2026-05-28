---
doc_type: "note"
title: "提案归档生成文档结构草案"
status: "draft"
workspace: "default"
module_scope: []
system_scope: []
convention_scope: []
ownership:
  - type: "system"
    target: "architecture/proposal-archive-document-structure.md"
    role: "primary"
  - type: "module"
    target: "modules/workflow-core"
    role: "secondary"
information_class: "exploration"
topics: []
related_docs: []
archive_target: "default:architecture/proposal-archive-document-structure.md"
created: "2026-05-22T08:17:52.067Z"
updated: "2026-05-22T08:17:52.067Z"
---

# 提案归档生成文档结构草案

## Ownership summary

- Primary module: modules/workflow-core
- System / architecture targets: architecture/proposal-archive-document-structure.md
- Convention targets: none
- Canonical reading path: proposal-archive-document-structure/artifacts/archive-document-structure-draft.md

## 1. Purpose

这份草案定义的是“提案归档后，长期文档应该长什么样”的分层稳定结构。目标不是把 proposal 原文搬过去，而是把 proposal 转译成可长期维护、可追加、可校验、也能继续扩展深度的项目文档。

## 2. Design principles

1. 类型专属，正文分离。
2. 只统一追踪信息，不统一所有正文。
3. 归档写入采用追加式幂等更新，不覆盖人工维护内容。
4. 主目标承载核心叙述，次目标承载配套沉淀。
5. 模块归档是目录级工件，architecture 和 decision 归档是单文档工件。
6. 小项目保持低摩擦，大项目通过索引页和子文档扩展深度，不强迫所有信息塞进一个文件。

## 3. Common traceability header

所有归档目标都应该保留一小段共用追踪信息，用于回答“这份文档从哪个 proposal 来、当前状态是什么、最后一次归档是什么时候”。

建议共用字段：

- Proposal ID
- Proposal title
- Archive date
- Source proposal link
- Primary or secondary role
- Target status

建议不要把所有正文都抽成共享模板。社区经验更倾向于把元信息统一，把正文保持在各自的文档类型里。
对复杂项目来说，共用头部的作用是让所有文档先“可定位”，不是让所有文档都“长得一样”。

## 4. Module archive structure

模块归档目标采用目录结构，建议固定为：

- `README.md`
- `design.md`
- `api.md`
- `history.md`

### 4.1 README.md

用于模块入口和边界说明，保持最短路径阅读。

建议内容：

- Status
- Primary proposal
- Purpose
- Key behaviors
- Important links

### 4.1.1 Large-project extension

对于复杂模块，`README.md` 应该承担“模块地图”的职责，补充以下内容：

- 子能力或子职责分区
- 关键入口和调用路径
- 相关设计/接口/运行文档索引
- 模块内最重要的不变量或约束摘要

当单个模块内部已经有明显的子域时，不要强行把所有内容塞进 `design.md`，而是允许增加子文档，例如：

- `subsystems/`
- `operations.md`
- `migrations.md`
- `scenarios.md`
- `decision-log.md`

### 4.2 design.md

用于描述当前模块形态和内部约束。

建议内容：

- Current shape
- Dependencies
- Invariants

### 4.2.1 Large-project extension

对于复杂模块，`design.md` 只保留“总览级设计”，更细的内容应该拆到子文档：

- 架构分层和边界
- 关键数据流
- 扩展点和失败模式
- 重要约束如何被代码、测试或运行规则保证

这样可以避免设计文档变成一段简短摘要，同时保留从总览下钻到细节的路径。

### 4.3 api.md

用于描述模块暴露的外部接口，仅在模块确实对外提供接口时填写。

建议内容：

- Public commands, APIs, or surface area
- Compatibility notes
- Non-goals if needed

### 4.3.1 Large-project extension

复杂模块可以继续保留 `api.md` 作为外部契约入口，但如果接口面过大，应该再拆出：

- 命令参考
- API/CLI 细节
- 配置说明
- 兼容性与迁移说明

### 4.4 history.md

用于追加式记录模块演进。

建议内容：

- Date
- Proposal ID
- Summary of what changed
- Source proposal reference

## 5. Architecture archive structure

architecture 目标采用单文件主文档结构，适合系统视角和跨模块关系说明。

建议章节：

- Status header
- Scope
- Components
- Relationships
- Views to maintain

### 5.0 Large-system stance

对于复杂系统，`system.md` 不应该承载全部知识，而应该成为系统文档的导航中心。

建议把它当成：

- 总览入口
- 视图索引
- 关键约束摘要
- 指向更细分文档的目录页

复杂系统应当允许拆分出更多架构子文档，例如：

- `context.md`
- `container.md`
- `domain/*.md`
- `flows/*.md`
- `constraints.md`
- `adrs/`

### 5.1 Recommended content shape

- Scope: 这份架构文档覆盖什么范围
- Components: 关键构件有哪些
- Relationships: 它们之间如何协作
- Views to maintain: 哪些视图需要长期维护，例如 system context、container、必要时的 dynamic view

### 5.1.1 Depth expectation

如果系统很大，`Views to maintain` 不应该只是一个标题列表，而应该明确每个视图的目的和读者：

- 哪些视图是架构稳定入口
- 哪些视图是分析某个子系统时必须看的
- 哪些视图只在复杂交互或问题排查时使用

### 5.2 Update rule

归档时不改写正文主结构，而是在文档末尾追加新的 proposal block，记录：

- proposal id
- summary
- source
- required follow-through

如果系统已经很复杂，追加块只是历史记录，不应替代专门的子系统架构页。

## 6. Decision archive structure

decision 目标采用单文件 ADR 结构，适合稳定决策记录。

建议章节：

- Status header
- Context
- Decision
- Alternatives
- Consequences

### 6.1 Recommended content shape

- Context: 发生了什么，为什么需要决策
- Decision: 最终选择了什么
- Alternatives: 还考虑过什么
- Consequences: 正负后果是什么

### 6.2 Update rule

和 architecture 一样，归档时采用追加式更新块，保留原始决策正文，避免重复覆盖历史判断。

## 7. Write behavior

归档实现建议遵循下面的行为：

1. 检查是否已经存在当前 proposal 的 marker。
2. 如果存在，跳过重复写入。
3. 如果不存在，追加一个新的归档 block。
4. 保留人工编辑内容，不做整文件覆盖。
5. 主目标先更新，次目标随后更新。

这个模型和当前实现的 marker-based append-only 行为一致，也符合社区里常见的“决策记录不可轻易改写、只能持续追加”的做法。

## 8. Knowledge landing and merge rules

这一节回答两个最关键的问题：

- 探索阶段识别出来的知识，最终应该落到哪里
- 提案中真正发生变更的知识，应该如何与已有内容融合

### 8.1 Knowledge source types

归档输入通常有两类：

1. 探索阶段识别到的知识
2. 提案中确认要变更、补充或替换的知识

两者的处理方式不同：

- 探索知识更像“证据和背景”，优先进入能够长期承载结构和约束的地方
- 变更知识更像“已决内容”，优先进入最终阅读路径和稳定正文

### 8.2 Where knowledge should land

#### 8.2.1 Module knowledge

适合落到模块文档的内容：

- 模块边界和职责
- 子能力或子职责分区
- 关键行为和使用者视角
- 依赖关系
- 不变量和约束
- 外部接口
- 运行、迁移、故障、场景等专题知识

建议分配如下：

- `README.md`: 模块入口、边界、能力地图、总览索引
- `design.md`: 总体结构、关键流转、边界条件、约束
- `api.md`: 外部契约、命令、接口、兼容性
- `history.md`: 本次提案真正改变了什么、为什么变、何时变

如果模块很复杂，探索里识别到的专题知识不应全部挤进 `design.md`，而应拆成子文档或专题页，例如：

- `operations.md` 记录运行和维护知识
- `migrations.md` 记录演进和迁移知识
- `scenarios.md` 记录典型场景和边界场景
- `constraints.md` 记录关键约束和验证方式

#### 8.2.2 Architecture knowledge

适合落到 architecture 文档的内容：

- 系统范围
- 关键构件
- 构件之间的关系
- 运行视图
- 交互视图
- 约束和决策之间的连接

建议分配如下：

- `system.md` 或总览页: 系统边界、主要视图、关键关系、阅读导航
- `context.md`: 外部环境、上下游系统、系统存在的原因
- `container.md`: 核心容器、职责划分、部署或进程形态
- `flows/*.md`: 关键流程、调用链、状态流转
- `constraints.md`: 不变量、架构约束、必须保持的规则
- `adrs/`: 已经稳定下来的关键决策

如果提案识别到的是“系统层的知识”，不要优先写进某个模块文档，而应该优先写进 architecture 文档，让整个系统的知识可见。

#### 8.2.3 Decision knowledge

适合落到 decision/ADR 的内容：

- 为什么必须做这个选择
- 备选方案有哪些
- 为什么放弃其他方案
- 这个决定带来的正负后果
- 哪些约束因此被固定下来

探索阶段发现的技术取舍，只要最终形成稳定选择，就应该进入 ADR，而不是只停留在 proposal 或 notes 里。

### 8.3 Merge rules with existing content

如果最终文档已经存在，归档不是“整篇重写”，而是“按知识类型融合”。

#### 8.3.1 When to append

以下情况优先追加，不改原结构：

- 新知识是对已有内容的补充，而不是替换
- 新知识属于新案例、新约束、新视图
- 新知识是提案后的历史记录
- 新知识只影响局部专题页

追加的典型位置：

- `history.md`
- architecture/decision 文档末尾的 proposal block
- 专题子文档末尾

#### 8.3.2 When to edit existing sections

以下情况应当直接融合到现有正文，而不是新增重复段落：

- 提案改变了模块边界定义
- 提案改变了系统范围定义
- 提案改变了不变量或约束的正式表述
- 提案纠正了之前写错的事实

编辑原则：

- 保留旧内容所表达的语义，不保留已失效的叙述
- 更新正文中的正式表述
- 如有必要，在 `history.md` 或追加块里记录“该段已被本次提案更新”

#### 8.3.3 When to split into a new page

如果新增知识满足以下任一条件，就不应继续塞进原页：

- 内容已经跨越多个主题
- 读者需要不同视角才能理解
- 单页开始同时承载总览和细节
- 后续还会持续增长

这时应该新建专题页，而不是把原文撑成巨石。

### 8.4 Conflict handling

如果探索知识、提案内容和现有文档冲突，优先级建议是：

1. 已确认并经过归档的最终文档
2. 当前提案中已经批准的变更
3. 最新探索中的证据和发现
4. 历史 notes、临时记录和过时表述

冲突处理方式：

- 如果旧内容只是过时，直接改正文
- 如果旧内容仍有历史意义但不再是当前事实，在 `history.md` 或 `archived notes` 中保留痕迹
- 如果冲突说明原有结构不够承载，就拆专题页

## 9. Knowledge base maintenance model

大型知识库的长期维护，通常不是依赖“单个大文档写得更全”，而是依赖一组稳定的治理习惯。我们可以把这些习惯转成归档后的默认运维模型。

### 9.1 Layered knowledge model

建议把最终知识库分成四层：

1. **Navigation layer**：入口页、索引页、目录页，让读者知道去哪找东西
2. **Canonical layer**：模块、architecture、decision 的正式正文，承载稳定知识
3. **Operational layer**：history、migrations、scenarios、constraints 之类的专题页，承载演进和细节
4. **Evidence layer**：explorations、findings、proposal notes，承载证据和推理过程

这和行业里常见的层次化文档做法一致：总览文档负责导航，专题文档负责深度，历史记录负责追溯，探索材料负责证据留存。

### 9.2 Single source of truth with derivative views

知识库里必须明确哪些内容是“正式事实”，哪些只是“解释或派生视图”。

- 正式事实放在 canonical layer
- 派生视图只做索引、摘要或视角化表达
- 任何派生视图都要能回链到正式事实

这借鉴了 ADR、C4 和 arc42 的共同做法：结构可以多视图，但正式定义要稳定，不能每个视图都各说各话。

### 9.3 Review and freshness discipline

大型知识库会失真，不是因为写得少，而是因为写完后没人持续检查。

建议增加三种维护动作：

- **Staleness check**：检查总览页是否还指向正确的专题页
- **Consistency check**：检查正文与 ADR、模块设计、架构视图是否互相冲突
- **Coverage check**：检查关键系统/模块是否至少有一个正式入口和一个可追溯历史

### 9.4 Knowledge capture rule

当提案完成归档后，它不应该只“结束在 proposal 里”，而应该沉淀为最终知识库的一部分。

建议规则：

- 新提案默认先检查已有 canonical docs
- 如果已有知识足够，直接复用并仅记录 delta
- 如果已有知识不够，探索先补证据，再更新 canonical docs
- 如果提案导致旧知识失效，直接更新正式正文，同时在 history 或 ADR 中保留演进痕迹

### 9.5 Default exploration target

后续任何新的探索，都应该默认把已归档的知识库当作首要目标，而不是从空白开始。

推荐探索顺序：

1. 当前问题对应的模块文档
2. 相关 architecture 文档
3. 相关 ADR
4. 相关 proposal 和 exploration 历史
5. 最后才是新增探索

这能把探索从“重复造轮子”变成“在现有知识上找缺口和冲突”。

### 9.6 Maintenance outcome

如果这套规则落实，最终产物就不只是归档目录，而是一个可持续扩展的知识库：

- 新知识有明确落点
- 旧知识有明确去向
- 探索默认围绕已有事实展开
- 归档目标不会沦为一次性产物

## 10. Community patterns referenced

### ADR / MADR

ADR 社区的共识是：单个 ADR 应该记录单个决策及其 rationale，核心信息通常围绕 context、drivers、options、outcome、consequences 展开。MADR 进一步强调模板应当轻量、可读、可写，并允许按场景使用短版或展开版。

这对我们的启发是：

- decision 目标适合保持单文档、稳定章节、轻量追加
- 不要把正文做成过度泛化的“万能文档”

### C4

C4 强调架构文档应该按抽象层次组织，不同视图对应不同读者和变化频率。系统上下文、容器、组件等视图并不是都必须存在，但应该按价值选择。

这对我们的启发是：

- architecture 文档应该偏“系统视角”，而不是记录提案过程
- 章节可以稳定，但视图内容可以按需选择，不必强制堆满所有图

### Diataxis

Diataxis 关注文档按用户需求组织，强调让内容结构服务于维护和使用，而不是把所有信息压成同一模式。

这对我们的启发是：

- module、architecture、decision 不应共享同一正文模板
- 共享的部分应该仅限追踪和导航信息

## 11. Operational stances

这些问题已经收口为操作立场：

- 共用元信息头部不需要独立模板文件，保持在 schema、模板占位符和渲染逻辑中统一即可。
- architecture 的 static / dynamic view 作为专题内容按需拆分，不强制成为主模板的一部分。
- `api.md` 只有在模块确实暴露接口时才需要存在，不必为了完整性留空。
- 追加块 marker 统一由归档实现生成，不单独再抽一个文件模板。
- `edit`、`append` 和 `split` 的判断按第 8 节的知识落点规则执行，不再引入额外自动判定。

## 12. Draft conclusion

如果当前只需要一个可执行的起点，那么建议先采用以下最小规范：

- module: `README.md` + `design.md` + `api.md` + `history.md`
- architecture: `system.md` 单文件，正文固定为 scope/components/relationships/views
- decision: `ADR-template.md` 单文件，正文固定为 context/decision/alternatives/consequences
- 共用：保留 proposal id、标题、时间、来源、主次角色等追踪信息
- 写入：只追加，不覆盖；靠 marker 保证幂等

但对于复杂系统，这个最小规范只是底层约束，不是完整答案。完整答案还应包含：

- 模块索引页和子文档
- 架构总览页和子视图页
- 决策记录与架构文档之间的交叉引用
- 约束、运行、迁移、场景等专题文档

换句话说，最终输出不是“一个更大的模板”，而是“一个能从总览一路下钻到细节的文档体系”。

## 13. Sources

- [ADR homepage](https://adr.github.io/)
- [MADR template](https://adr.github.io/madr/decisions/adr-template.html)
- [MADR examples](https://adr.github.io/madr/examples.html)
- [arc42](https://arc42.org/)
- [C4 model homepage](https://c4model.com/)
- [C4 diagrams overview](https://c4model.com/diagrams)
- [Diataxis](https://diataxis.fr/)

## 14. Candidate workflow decision

This draft has been formalized as [ADR-002](/Users/qiangbi/develop/projects/Syl/tangram-v2/tg-workflow/docs/decisions/ADR-002-archived-knowledge-base-as-default-exploration-baseline.md).
