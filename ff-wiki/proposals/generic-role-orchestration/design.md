---
flowforge:
  schema: 1
  role: design
  id: generic-role-orchestration-design
  revision: 1
  consumes:
    requirements:
      generic-role-orchestration-requirements: 1
  areas:
    roster:
      revision: 1
      anchor: d-roster
    dispatch:
      revision: 1
      anchor: d-dispatch
    pi-boost:
      revision: 1
      anchor: d-pi-boost
    research-intake:
      revision: 1
      anchor: d-research-intake
---

<a id="generic-role-orchestration-design"></a># 通用角色与任务链调度方案

依据：[通用角色与任务链调度需求](requirements.md#generic-role-orchestration-requirements)。实证基础：[GIIS 派发模式调研](../../research/2026-09-24-giis-subagent-dispatch-pattern.md)（42 worker 的编排模式、任务模板与角色谱系）。

## <a id="d-roster"></a>d-roster：角色清单与契约（混合模型落点）

**新增 3 个通用角色资产**（`assets/subagents/`，flowforge-\* 命名，部署机制与既有角色同一）：

| 角色 | 能力契约（description 骨架） | model_profile | default_skill（方法材料） | permission |
|---|---|---|---|---|
| `flowforge-batch-analyst` | 批量"提取+对照+汇总"：对可分组并行的一组分析单元，产出带引用工作台文档；不做决策、不综合跨组结论 | `tool-capable` | `flowforge-research` | `workspace-write`（语义标签） |
| `flowforge-scribe` | 按给定模板与材料撰写/回填结构化文档；只落格式与既定内容，不引入新语义 | `tool-capable` | `flowforge-writing-for-agents` | `workspace-write` |
| `flowforge-executor` | 机械执行既有命令/生成器的批次并逐字报告输出；不开发新工具、不修代码 | `tool-capable` | `flowforge-implement`（仅取其 fail-fast 与证据纪律；Boundaries 声明不进入票务流程） | `workspace-write` |

每个角色 Body 三段：Identity（能力）/ Boundaries（MUST NOT 决策、MUST NOT 改代码〔executor 同〕、每条产出带可验证引用）/ Default Skill（"或直接 read `.agents/skills/...`"双通道，与既有角色同构）。

**复用而非新造**：
- **定向探查 = `flowforge-investigator` 原样复用**（能力行入通用调度段，资产零改动——其 description 本就是能力型）。替代（否决）：新造 `flowforge-scout`——与 investigator 能力重复，造成两个 read-only 探查角色的钉扎与调度歧义。
- **决策素材简报 = investigator × flowforge-research 简报契约**：简报是 investigator 的一种输出形态，契约文本由 flash-worksurface-expansion 06 票承载（流程内 SKILL 修订），本设计只在通用调度段建能力行指向它。替代（否决）：新造 `flowforge-brief-writer`——与 investigator+research 重复；把 06 吸进本 proposal——否决，06 是流程契约修订，归属 flash 提案，两边共享契约文本即可。
- **工具开发不设通用角色**：开发新工具=代码变更=implementer 票道（执行单元政策不动）。替代（否决）：`flowforge-builder`——越权进代码面。

**模型档位与钉扎**：资产只声明 profile；flash 钉扎是项目侧 config（`agents.models_by_name` / 后续 `models_by_host`），资产不写死模型——未钉扎项目产物零变化（继承宿主默认）。`workspace-write` 是新语义标签：编译无特殊化（与 requirement-authority 等既有标签同待遇），写约束由 Boundaries 文本承载；def 级 edit-deny 硬化记 open item。

## <a id="d-dispatch"></a>d-dispatch：AGENTS.md 通用调度段（双节结构）

`assets/AGENTS.md`（60 行）在 `## Agent skills` 与 `## Subagent delegation` 之间插入 `## Generic capability dispatch`——通用在前、流程在后，声明顺序即适用优先级（任何工作先按能力找角色，flowforge 流程态再查流程委派表）：

1. **能力键调度表**：

   | 当前要做的事（能力键） | 角色 | 档位建议 |
   |---|---|---|
   | 定向探查一个有界问题（源结构/环境/既有行为），闭合事实缺口 | `flowforge-investigator` | flash 可钉 |
   | 批量提取+对照+汇总，按组并行切分 | `flowforge-batch-analyst` | flash 可钉 |
   | 按模板撰写/回填结构化文档 | `flowforge-scribe` | flash 可钉 |
   | 机械执行既有命令/工具批次 | `flowforge-executor` | tool-capable |
   | 决策素材简报（研究输出契约） | `flowforge-investigator`（简报形态，契约见 flowforge-research） | flash 可钉 |

2. **任务链协议**（GIIS 提炼）：编排会话（旗舰）规划任务链 → 逐任务用结构化模板派发 → 可并行批量派发 → Review 收敛产物 → 规划下一批。结构化任务模板钉进模板文本：

   ```
   你是<角色>。项目根：<绝对路径>。
   【目标】<一句话可验收目标>
   【输入】<精确文件/目录路径与参照物>
   【输出】<输出文档路径；每条结论带引用（文件路径+行号或命令输出）>
   ```

3. **派发边界条款**："机械可完成才下放 flash 档；任务链规划、跨组综合、决策、Review 收敛留在编排会话。"

同步面：`.agents/skills/` 侧 AGENTS 模板随 upgrade 收敛；flowforge 仓自身 AGENTS.md（自举）同步该段。

## <a id="d-pi-boost"></a>d-pi-boost：pi 宿主强化（第一版 = 指引文本）

通用调度段末尾附宿主提示小节（仅 pi 生效的表述）：pi 下派发经 fork 继承编排会话上下文为只读参考（worker 无需重读背景）；批量任务可异步并行派发并等待完成唤醒收敛；模型/thinking 可用用户级 `agentOverrides` 按名覆盖。**不改 `compile_pi.go`、不改 extension 代码**——基线体验宿主无关（纯 AGENTS.md 文本），pi 提示是增强。extension 原生 dispatch 助手工具记 open item（oi-pi-dispatch-tool）。

## <a id="d-research-intake"></a>d-research-intake：调研产物落点与导入衔接

- **落点约定**：讨论期/分析期产物落 `<docs_dir>/research/YYYY-MM-DD-<slug>.md`，正文带引用；不引入新顶层目录（需求红线）。
- **导入衔接**：`flowforge-import` SKILL 的 Inputs 节增补一句——"`<docs_dir>/research/` 是讨论期产出的标准源位置，research 笔记按 Source fact / Requirement candidate / … 既有分类流转"；其余分类与 hand-off 逻辑零改动（import 本就消费本地文档，只补标准源的可知性）。
- **proposal 创建衔接**：创建 proposal 的会话（align 起点或 plan 起点前）检查 `<docs_dir>/research/` 未消费笔记，提示经 import 导入——落为 AGENTS 调度段的一句指引，不新增命令。

## Standards clauses

- must 纯本地确定性文件操作，无网络、无 LLM 调用（源：[AGENTS.md](../../../AGENTS.md) 核心设计原则 2，[Constraints]）。
- must 变更后运行 `go test ./internal/...`（源：[AGENTS.md](../../../AGENTS.md) boundaries，[Conventions]）。
- must 不改 CLI 命令签名与 Issue Schema 头规范；不改编译器行为（源：[AGENTS.md](../../../AGENTS.md) Ask first 边界 + 本设计 d-roster/d-pi-boost，[Constraints]）。
- must `assets/` 只放部署内容（角色资产与 AGENTS 模板属部署内容）（源：[AGENTS.md](../../../AGENTS.md) 🚫 Never，[Constraints]）。
- must flash 档角色产出必须带可验证引用（源：[flash-worksurface-expansion 需求](../flash-worksurface-expansion/requirements.md#flash-worksurface-expansion-requirements) 约束，[Conventions]）。
- must 通用角色 description 按能力书写、不含流程术语（源：[本需求](requirements.md#generic-role-orchestration-requirements) 验收 1，[Conventions]）。

## 兼容与迁移

- 未钉扎模型的项目：新角色部署产物无 `model` 字段（继承宿主默认），既有六角色产物零变化。
- 已部署项目 `flowforge upgrade` / `agents deploy` 后获得 3 个新角色与新版 AGENTS.md；`agents disabled` 可按名关闭不需要的角色（既有机制）。
- AGENTS.md 是 managed asset：升级覆盖用户自改会触发既有 preserved/漂移语义（assets 部署先例），无新迁移路径。
- 名册测试（`subagent_source_test`、`assets_deploy_test` 的 builtin 名单）随新角色扩展——属 preset 测试，Plan 票内显式授权执行者修改。

## 验证策略

- 单元/结构：新资产经 `Parse` 通过（name=filename、default_skill 解析到真实 skill 目录）；四宿主编译产物回归（未钉扎零变化）；AGENTS.md 模板断言含 `Generic capability dispatch` 段锚点、能力表关键词、任务链模板与边界条款；import SKILL 断言提及 research 标准源。
- 名册：preset 名单测试扩展（授权在票内）。
- e2e（人工、记 Completion evidence）：pi 宿主按协议派发 ≥3 并行 batch-analyst 子任务（复用本仓真实分析需求），观察 fork worker 产物带引用、旗舰 Review 收敛、（选用时）完成唤醒链路。

## Open items

- **oi-pi-dispatch-tool**（gap，area: pi-boost）：pi extension 原生 dispatch 助手（任务链模板渲染/批量派发记账），第一版不做；等 dogfood 暴露文本协议的摩擦点再立项。
- **oi-def-edit-deny**（gap，area: roster）：def 级 edit-deny（把 workspace-write 约束从 Boundaries 文本硬化为编译期 edit-deny），依赖 per-def deny 通道设计，暂以文本约束。
