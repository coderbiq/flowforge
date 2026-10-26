# FlowForge Wiki：问题、原理与解法

> 状态：01 篇已形成正文初稿，02–24 篇保留结构稿，后续逐篇完善。

这套 Wiki 的目标是从 Context 与 Agent 的工作原理出发，讲清楚 FlowForge 为什么存在、要解决哪些具体问题，以及各项功能通过什么机制解决这些问题。

总的出发点是：复杂需求需要人与 Agent 多轮协作，而每一步能够可靠使用的信息有限。FlowForge 尝试通过持续保存有效结论、提前准备任务所需信息、显式组织依赖和保留交付证据，帮助 Agent 逐步完成分析、设计与开发。后续各篇将分别解释这条因果链及其边界。

## 阅读路线与写作范围

全文采用“原理 → 问题 → 总体解法 → 功能分解”的总分结构。每部分回答不同问题：原理解释为什么会有这些困难，问题篇说明复杂项目中的表现，总体解法连接问题与机制，功能篇逐一展开机制。

文档文件名使用 01–24 的连续序号，便于按阅读顺序浏览。第一篇已展开正文与官方来源引用；其余各篇保留标题、概要、正文规划与例子及依据准备，作为后续逐篇讨论的入口。完善时应区分模型行为、宿主策略、工程经验和 FlowForge 的实现约定，并为关键判断补充依据。

这套 Wiki 聚焦解释问题与设计理由；安装与命令操作可链接现有[项目 README](../../README.md)、[场景导览](../scenarios.md)和 [CLI 文档](../cli-design.md)。

## 第一部分：Context 与 Agent 的原理

先建立必要概念，再讨论容量、信息使用和上下文组织的取舍。

| 文档 | 核心问题 |
|---|---|
| [01 · Context 是什么](1-principles/01-context-fundamentals.md) | 模型、宿主与文件存储分别持有哪些信息？什么内容进入一次请求？ |
| [02 · 有效上下文](1-principles/02-effective-context.md) | 内容进入窗口后，能否被可靠检索、理解与执行？ |
| [03 · Agent 如何组织 Context](1-principles/03-agent-context-processing.md) | 读取、工具调用、压缩、恢复和子代理如何改变可用信息？ |
| [04 · 精准上下文的取舍](1-principles/04-precision-vs-recall-in-context.md) | 如何兼顾局部相关、必要信息充分与上下文预算？ |

## 第二部分：复杂项目中的核心问题

把上述原理落实到多轮需求讨论、设计、开发和交接的具体困难。

| 文档 | 核心问题 |
|---|---|
| [05 · 结论传递](2-problems/05-problem-continuity.md) | 已确认结论为何会在阶段切换和长对话中丢失？ |
| [06 · 信息取用](2-problems/06-problem-retrieval.md) | 内容已经写进文档，执行时为何仍会漏用？ |
| [07 · 任务上下文](2-problems/07-problem-task-context.md) | 为什么任务清单还不足以直接指导实施？ |
| [08 · 会话与执行者交接](2-problems/08-problem-handoff.md) | 如何继承实际进展、失败尝试和有效版本？ |
| [09 · 项目规范适用性](2-problems/09-problem-standards.md) | 每项任务应遵守哪些规范，为什么容易判断错或遗漏？ |

## 第三部分：FlowForge 的总体解法

先给出完整信息流，再解释各组成部分为何分工协作。

| 文档 | 核心问题 |
|---|---|
| [10 · 总体解法](3-approach/10-approach-overview.md) | 持久结论、任务上下文、依赖图和证据反馈如何共同工作？ |
| [11 · Skill、Markdown 与 CLI 的分工](3-approach/11-skill-and-cli-split.md) | 哪些判断由人与 Agent 承担，哪些事实由程序计算？ |

## 第四部分：每项功能解决什么问题

按核心信息流展开，再解释支持其落地的能力。每篇正文围绕“具体问题 → 机制与例子 → 与其他机制的衔接 → 适用条件和边界”组织。

| 文档 | 问题与机制 |
|---|---|
| [12 · 需求与设计持久化](4-features/12-authority-artifacts.md) | 结论难以跨会话继承 → 持续维护职责明确的需求与设计来源 |
| [13 · 任务执行契约](4-features/13-ticket-execution-contract.md) | 执行信息不足 → 局部上下文、明确动作与事实细化 |
| [14 · 规范转录](4-features/14-standards-transcription.md) | 执行时反复判断规则 → 前置识别并把适用条款带入任务 |
| [15 · DAG 与诊断](4-features/15-dag-and-diagnostics.md) | 顺序和阻塞依赖记忆 → 显式关系、图检查与可推进队列 |
| [16 · 新上下文与角色协作](4-features/16-fresh-context-execution.md) | 历史负担和职责混杂 → 有界工作与持久工件交接 |
| [17 · 版本与引用追踪](4-features/17-revision-tracking.md) | 下游信息过时 → 来源定位、消费记录与变化提示 |
| [18 · 验证与双轴审查](4-features/18-evidence-and-review.md) | 状态无法证明完成 → 可观察验证、审查处置与完成证据 |
| [19 · 失败反馈](4-features/19-failure-feedback-loop.md) | 失败经验丢失 → 保存尝试、修正契约并安排修复 |
| [20 · 模型分工与执行边界](4-features/20-model-tiering.md) | 工作能力需求与成本不同 → 角色配置、有界执行与审查 |
| [21 · 外部材料接入](4-features/21-external-material-intake.md) | 旧材料混有不同性质的内容 → 分类、保留出处并交给相应责任方确认 |
| [22 · 调查与设计辅助](4-features/22-research-and-design-support.md) | 缺少事实或决定依据 → 调查、诊断、原型与决策梳理 |
| [23 · 前端视觉反馈](4-features/23-frontend-visual-feedback.md) | 文字与代码难以证明界面结果 → 截图、自审和视觉复核 |
| [24 · 多宿主部署与配置](4-features/24-host-and-model-configuration.md) | 方法在不同工具与机器上难以一致运行 → 资产部署、本机配置与漂移检查 |

## 后续逐篇完善方式

先完善 01–04，建立可靠的原理基础；随后用贯穿案例完善问题篇与总体解法；最后逐项核对功能机制。核心功能与配套能力都回到具体问题解释，篇数可随内容成熟合并或细分。

原理篇需要官方文档或研究依据；问题篇需要可理解的工程场景；功能篇需要当前 Skill、CLI 或宿主实现依据。机制存在、预期作用与实测效果分别说明。参考讨论素材为[第一轮问题与方案讨论稿](../research/2026-10-26-flowforge-problem-and-approach.md)，现有架构和操作文档作为补充来源。
