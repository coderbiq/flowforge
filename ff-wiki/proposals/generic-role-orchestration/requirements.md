---
flowforge:
  schema: 1
  role: requirement
  id: generic-role-orchestration-requirements
  revision: 1
---

<a id="generic-role-orchestration-requirements"></a>
# 通用角色与任务链调度：需求

## 问题

1. **调度契约只覆盖 flowforge 流程阶段**：AGENTS.md 委派表以流程状态为键（"需求未定 → flowforge-analyst"），proposal 未立的讨论期与流程外的批量分析工作没有调度入口——编排会话不知道该按什么协议拆分与派发。
2. **反复出现的通用能力没有角色资产**：GIIS 项目一周内 42 个 fork worker 临场造了 10+ 种角色（批量对照分析、定向探查、结构化撰写/回填、决策素材简报、工具开发/执行——见[GIIS 派发模式调研](../../research/2026-09-24-giis-subagent-dispatch-pattern.md)），全部提示词现造、任务生灭、无法沉淀复用。这是能力缺失下的替代品，不是目标形态。
3. **旗舰/flash 分工已实证但未协议化**：旗舰拆任务链 + 定目标 + Review 收敛、flash 机械批量执行的批次制模式（单日 30+ worker 并行）没有写进任何引导，每次靠用户手写提示词撑着。
4. **讨论期产物与流程衔接缺失**：flowforge-import skill 已存在，但“调研产物落 `<docs_dir>/research/` → proposal 创建时分类导入”的源头约定没有，导入靠人肉搬运。（`docs_dir` 是 wiki 根配置；新 init 项目默认已改为 `ff-wiki`，不再占用项目惯用的 `docs/`）

## 目标

1. **通用角色资产**：固化 4-5 个能力化角色（候选见调研笔记角色谱系表：批量分析员、定向探查员〔与 flowforge-investigator 能力重叠，泛化还是新造归设计〕、撰写/回填员、决策素材简报员〔与 flash-worksurface-expansion P3 简报契约同源，汇流载体归设计〕、工具开发/执行员）；每个角色含能力契约（description 按能力写，不含流程术语）、模型档位建议、输出契约（带引用工作台文档：文件路径+行号或命令输出）。
2. **调度协议进 AGENTS.md 通用调度段**：任务链规划 → 结构化任务（【目标】【输入】【输出】+ 精确路径）→ 批量派发 → 旗舰 Review 收敛 → 规划下一批；含派发边界条款（"机械可完成才下放 flash 档"）；与既有 flowforge 流程委派表并存（混合模型：流程原生角色 flowforge-\* 保持绑定不动）。
3. **PI 宿主强化**（用户补充）：结合 pi 特性强化调度体验——fork 继承参考上下文、异步批量执行与完成唤醒、用户级 `agentOverrides` 模型覆盖、项目级 pi extension 原生工具；**基线体验宿主无关**（AGENTS.md 即可用），pi 强化是增强不是依赖；具体强化形态归 solution-design。
4. **旅程衔接**：讨论期调研产物落 `<docs_dir>/research/`（wiki 根相对，跟随项目配置，不引入新顶层目录，带引用），proposal 创建时经 flowforge-import 分类进 requirements 背景事实/约束与待裁决项。
5. **模型通道复用**：通用角色的 flash 钉扎走既有 `agents.models` / `agents_by_name`（及后续 `models_by_host`）通道，不新增配置键。

## 范围与约束

- 仅讨论**已 init 项目**（无部署场景出范围）；临时角色运行时出范围——固化的目标正是消灭临场造角色。
- flowforge 流程原生角色（flowforge-analyst/architect/planner/implementer/reviewer/investigator）保留流程绑定，**不追求完全解耦**（用户裁决：混合模型）。
- 不改 CLI 命令签名与 Issue Schema 头规范；pi 强化不得使其他宿主（opencode/claude/codex）体验降级。
- 与 flash-worksurface-expansion 的关系：05（reviewer-lite）继续走其快速通道，与本 proposal 互不阻塞；06（设计事实简报）与本 proposal 简报角色同源，设计阶段裁决合并载体（一张票吃掉还是两处并存）。

## 可观察验收

1. 通用角色资产存在于 `assets/subagents/`（或等价源），`flowforge agents deploy` 部署到启用宿主，description 按能力书写、不含流程术语。
2. AGENTS.md 含通用调度段（能力键调度表 + 任务链协议 + 派发边界条款），与流程委派表并存且互不遮蔽。
3. pi 宿主按协议派发一批 ≥3 个并行子任务到通用角色：批量执行可观察（worker 产物/记录）、旗舰 Review 收敛点在协议中有明确步骤；若设计选用完成唤醒，派发-唤醒-收敛链路可观察。
4. 讨论期产物落 `<docs_dir>/research/` 且带引用；此后新建的 proposal 其 requirements 引用该产物并把事实导入。
5. flash 档通用角色经 config 钉扎后，部署产物 frontmatter `model` 为钉扎值（复用既有通道，回归零）。

## 术语

- **通用角色**：按能力定义、流程无关的持久 subagent 资产（与"流程原生角色"相对；两类并存构成混合模型）。
- **任务链调度协议**：编排会话规划结构化任务批次、按角色能力派发、Review 收敛、迭代下一批的契约。
- **工作台文档**：调研/分析产出的带引用 Markdown（文件路径+行号或命令输出），与 `<docs_dir>/research/` 约定同构。

## Standards 识别清单（转 solution-design 转换为 must 条款）

- [AGENTS.md] 核心设计原则 2：纯本地确定性文件操作，无网络、无 LLM 调用。
- [AGENTS.md] boundaries：变更后运行 `go test ./internal/...`。
- [AGENTS.md] Ask first：修改 Issue Schema 头规范、变更 CLI 接口签名。
- [AGENTS.md] 🚫 Never：`assets/` 只放部署内容（角色资产与 AGENTS 模板属部署内容，合规）。
- [flash-worksurface-expansion 需求] flash 化产物必须带可验证引用——通用角色输出契约与之同构。
