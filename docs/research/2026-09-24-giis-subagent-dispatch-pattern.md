# GIIS 项目 pi 宿主的 subagent 派发模式调研（调度协议与通用角色的实证来源）

日期：2026-09-24
数据源：`~/.pi/agent/sessions/--Users-qiangbi-develop-projects-Bytesforce-giis--/`（09-22 ~ 09-23 主会话 18 个 + `subagent-artifacts/` worker transcript 42 个，合计约 11MB 会话记录）
研究问题：用户在没有角色资产与调度契约的项目里，如何用提示词驱动 subagent 完成大规模分析/迁移工作？哪些能力反复出现、值得固化为通用角色？

## 结论摘要

1. **编排模式（旗舰/flash 分工 + 批次迭代）已实证有效**：旗舰（GLM-5.3）负责拆任务链、定目标、Review 收敛；flash（mimo-v2.6-flash）负责机械批量执行。单日 30+ worker 并行，编排会话按"规划一批 → 派发 → Review → 再规划下一批"迭代。
2. **角色全部是提示词临场定义**（任务生灭，无资产沉淀），共出现 10+ 种角色名；其中约 5 种能力反复出现（见角色谱系表），是固化为通用角色的候选。
3. **worker 统一结构**：pi fork 派发自带包装（继承父会话为只读参考）+ 内联角色定义 + 结构化任务（【目标】【输入】【输出】，带精确文件路径）。这个结构本身就是可固化的**任务模板**。
4. **模型分层在派发指令中指定**（"交给 mimo-v2.6-flash"），不走配置文件——固化时应改走 flowforge 既有模型通道（agents.models*）。
5. 调用上下文靠 fork 继承而非 skill 文件传递：方法论写在任务提示里。固化后通用角色仍可保留 fork 继承（pi 特性），但角色能力契约本身应进资产。

## 一、编排侧证据（主会话用户提示词，逐字）

> "完成规划实施步骤，简单的执行任务创建 subagent 交给 mimo-v2.6-flash 模型来做，你只需要做 review 和复杂的工作，交给 subagent 做的工作要明确必须让 flash 模型可以机械完成"

> "把剩下的所有可以简单直接分析的字段分组然后明确分析执行步骤创建 subagent 使用 mimo-v2.6-flash 模型批量快速执行，你只需要提供明确的工作目标以及对最终结果进行 Review"

> "继续，再规划下一批子任务并创建 subagent 派发执行"

> "根据这个目标整体来规则任务，简单的派发给 Subagent 做，复杂的和决策性的你自己做"

派发边界约定（用户反复强调）：**flash 只接"机械可完成"的任务；决策、复杂工作、Review 留编排会话。**

## 二、worker 侧统一结构

pi fork 派发的固定包装（pi-subagents 自动注入）：

> "You are a delegated subagent running from a fork of the parent session. Treat the inherited conversation as reference-only context, not a live thread to continue... Your sole job is to execute the task below and return a focused result."

之后是编排会话写的任务体，结构稳定：

```
你是<能力角色>。项目根：<绝对路径>。
【目标】<一句话可验收目标>
【输入】<精确文件/目录路径，含参照物>
【输出】<工作台文档路径>（全部要求带引用：文件路径+行号或命令输出）
```

工具面：read / bash / write / edit（fork 继承编排会话工具）；worker 体量 50KB~920KB transcript（工具调用 2~46 次/任务）。

## 三、角色谱系（42 worker 统计）

| 临场角色名 | 频次 | 通用能力抽象 | 固化候选 |
|---|---|---|---|
| 码表映射分析员 | 13 | 批量"提取+对照+汇总"分析，按字段组并行切分，产出带引用工作台文档 | **batch-analyst**（flash 档） |
| 数据结构侦察员 / 配置探查员 | 2 | 定向探查一个有界问题（源结构/FK 指向/环境配置），闭合事实缺口 | 侦察探查角色——与 flowforge-investigator 能力面重叠，泛化它或新造（设计裁决） |
| 迁移规格撰写员 / 回填员 | 3+3 | 按模板产出/修订结构化文档（只落格式不改语义） | **spec-writer / doc-editor**（flash 档） |
| 迁移工具开发员 / 执行员 / 镜像同步开发员 | 3+1+1 | 开发生成器、批量运行机械转换 | builder/executor（tool-capable 档） |
| 映射方案分析员 / 数据分析师 / 产品配置分析员等 | 各 1 | **为决策准备素材表，明确"不做决策"** | **brief-preparer**（flash 档）——与 flash-worksurface-expansion P3"设计事实简报"同源 |
| workflow 结构分析员 / 迁移依赖分析员等 | 各 1 | 提取既有结构作为模板参照 / 完成度核验 | batch-analyst 变体 |

## 四、对 FlowForge 固化的含义

1. **通用角色资产**：上表五类能力固化为流程无关的持久 subagent（能力契约 + 模型档位 + 输出契约=带引用工作台文档），取代临场造角色。
2. **调度协议**（写进 AGENTS.md 通用调度段）：任务链规划 → 结构化任务（目标/输入/输出模板）→ 批量派发 → 旗舰 Review 收敛 → 规划下一批；"机械可完成才下放 flash"是派发边界条款。
3. **pi 宿主强化面**：fork 继承参考上下文、异步批量执行与完成唤醒（pi-subagents bg_wait/named runs）、用户级 agentOverrides 模型覆盖、项目级 pi extension 原生工具——都是现成机制，固化时可选作增强（宿主无关基线仍以 AGENTS.md 为准）。
4. **产物落点**：调研/分析产物全部是带引用 Markdown 工作台文档——落 `<docs_dir>/research/`（wiki 根相对，跟随项目 `docs_dir` 配置，不强制新顶层目录；本仓库自举 wiki 根为 `docs/`，故本笔记落在 `docs/research/`），可直通 flowforge-import。
5. **模型通道**：GIIS 用指令指定模型；固化后走 agents.models* 配置通道（与 flash-worksurface-expansion 的钉扎机制复用）。

## 证据范围与可复现方式

- worker 任务与工具统计：解析 `subagent-artifacts/*_worker_transcript.jsonl`（recordType: message/tool_start），提取首个含 "delegated subagent" 的 user 消息与 tool_start 计数。
- 编排提示词：主会话 jsonl 中 role=user 消息按关键词（派发/subagent/并行/批量/规划）过滤。
- 会话记录含隐私提示词脱敏标记（"[prompt redacted]"），本笔记仅引用与研究问题相关段落。
