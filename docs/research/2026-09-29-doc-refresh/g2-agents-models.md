# g2 — Subagent 生命周期 / 通用角色调度 / 模型保留与模型方案 / PI 宿主集成事实提取（文档更新消费用）

> 分组：`subagent-lifecycle`、`generic-role-orchestration`、`agent-model-preservation`、
> `model-sets-switching`、`pi-host-integration`
> 方法：逐文件引用事实（文件路径 + 行号 或命令输出），不做跨组综合，不提改进建议。
> 每组内"表面变更"均已在当前仓库实物中复核（下方 `assets/…` / `internal/…` 行号为实测）。
> 本组不涉及任何代码或源工件修改；仅产出本文档。

复核用命令与观测（本组全部 proposal 的收口状态）：

```
$ ./bin/flowforge status --dir docs/proposals/subagent-lifecycle
FlowForge Local Tracker Status: 8/8 resolved (100%)
Frontier Ready: 0 | In Progress: 0 | Blocked: 0

$ ./bin/flowforge status --dir docs/proposals/generic-role-orchestration
FlowForge Local Tracker Status: 4/4 resolved (100%)
Frontier Ready: 0 | In Progress: 0 | Blocked: 0

$ ./bin/flowforge status --dir docs/proposals/agent-model-preservation
FlowForge Local Tracker Status: 2/2 resolved (100%)
Frontier Ready: 0 | In Progress: 0 | Blocked: 0

$ ./bin/flowforge status --dir docs/proposals/model-sets-switching
No features or issues found in docs/proposals/model-sets-switching

$ ./bin/flowforge status --dir docs/proposals/pi-host-integration
FlowForge Local Tracker Status: 3/3 resolved (100%)
Frontier Ready: 0 | In Progress: 0 | Blocked: 0
```

---

## 1. subagent-lifecycle

### 1.1 Feature 一句话
安装或升级 FlowForge 后项目里就有一批内置角色（analyst / architect / planner / implementer / reviewer / investigator），用户可用 `flowforge agents deploy|remove|status` 把任意角色（含自定义角色）确定性编译成宿主原生子代理文件，并可指定只启用哪些宿主。

### 1.2 表面变更清单
| 表面 | 事实 | 引用 |
|---|---|---|
| 新 CLI 命令组 | `flowforge agents`，含三个子命令 `deploy` / `remove` / `status` | 命令实测：`flowforge agents --help` → `Manage subagent definitions for Claude Code, OpenCode, Codex, and PI` + `deploy Remove a subagent… status Compare deployed…`；`internal/command/agents.go:30`（Long 描述四宿主目录） |
| CLI flag | `flowforge agents status --json`（输出格式对齐 `assets verify --json`） | `docs/proposals/subagent-lifecycle/design.md:207`；`docs/cli-design.md:52` |
| 新 agent 角色（6 个内置流程角色） | `flowforge-analyst` / `flowforge-architect` / `flowforge-planner` / `flowforge-implementer` / `flowforge-reviewer` / `flowforge-investigator` | `docs/proposals/subagent-lifecycle/design.md:27-34`（角色清单表）；实物 `assets/subagents/` 下 12 份 `.md`（含本提案 6 份） |
| 新资产目录 + 权威定义 schema | `assets/subagents/*.md`，frontmatter 键集 `flowforge_agent.{name,description,model_profile,default_skill,detour_skills,permission,after,before,returns_to}`；正文五段式（Identity / Boundaries / Workflow Position / Default Skill / Result Contract） | `docs/proposals/subagent-lifecycle/design.md:175-189`（schema 代码块）、`design.md:191`（目录来源）、`design.md:231`（实现边界-新增目录） |
| 模型档位映射 | `model_profile` 三档 `high-capability` / `tool-capable` / `tool-capable-read-only` → 宿主字段（claude `model`、opencode `model`、codex `model_reasoning_effort`） | `docs/proposals/subagent-lifecycle/design.md:193-197` |
| 部署目标目录（初版三宿主） | `.claude/agents/`、`.opencode/agent/`、`.codex/agents/` | `docs/proposals/subagent-lifecycle/design.md:203` |
| config 键（新增） | `agents.disabled`（内置角色停用列表）、`agents.hosts`（启用宿主集合，`claude`/`opencode`/`codex` 非空子集，未配置默认全三） | `docs/proposals/subagent-lifecycle/design.md:206`、`design.md:211`、`design.md:213-216`、`design.md:232`；实物 `internal/config/config.go:41`（`disabled`）、`internal/config/config.go:42`（`hosts`） |
| 新 skill | 无新增 skill；本提案是部署/生命周期机制 | 全组 diff 面见 `docs/proposals/subagent-lifecycle/design.md:231-234`（实现边界只列目录/CLI/文档/测试） |
| AGENTS.md 模板 | 新增 `## Subagent delegation` 段（含"里程碑责任人 → subagent → 绑定 skill"路由表） | `docs/proposals/subagent-lifecycle/design.md:150-167`（模板全文）、`design.md:233`；实物 `assets/AGENTS.md:64`、仓库自身 `AGENTS.md:93` |
| wiki 结构 | 无新增目录（仍是 `<docs_dir>/{CONTEXT.md,adr/,proposals/}`） | `docs/proposals/subagent-lifecycle/design.md:231` 只新增 `assets/subagents/` 与三宿主部署目录 |

### 1.3 实施状态
**全部 closed（8/8）。**

证据：
- ticket 状态行：`docs/proposals/subagent-lifecycle/issues/01-subagent-source-files.md:17`、`02-agents-md-delegation-table.md:17`、`03-parser-and-compiler-core.md:18`、`04-agents-deploy-command.md:18`、`05-agents-remove-and-disable.md:18`、`06-agents-status-command.md:18`、`07-init-upgrade-auto-deploy.md:18`、`08-host-scoped-deployment.md:16` — 全部 `**Status:** closed`。
- 命令输出：`flowforge status --dir docs/proposals/subagent-lifecycle` → `8/8 resolved (100%)`，`Frontier Ready: 0`（本文档首部命令块）。
- 交付证据：`issues/01-subagent-source-files.md:94`（"`assets/subagents/` contains exactly 6 markdown files …each with valid `flowforge_agent` frontmatter and 5-section body"）、`01-subagent-source-files.md:96`（4 个测试函数通过，commit `077a519`）。
- 宿主收窄部署交付证据：`issues/08-host-scoped-deployment.md:88`（`go test ./internal/...` 全通过，含 5 个新用例 `TestAgentsDeployHonorsHostSelection`、`TestAgentsDeployCleansDeselectedHosts`、`TestAgentsHostsValidation`、`TestAgentsStatusScopedToSelectedHosts`、`TestAgentsRemoveScopedToSelectedHosts`）。
- 未完成度：0/8；`flowforge check` 对本 proposal 只报 waived `upstream-changed`（`issues/01…07`：Consumed revision 1 is behind 3，均已填 waiver 理由），无 error 级诊断（命令输出见本文档"复核用命令"节外的 `flowforge check` 尾部 `✓ Dependency graph is healthy`）。
- 修订边界（事实，非缺陷）：本提案设计到 revision 3，实际部署宿主已扩为四个（`pi` 由 pi-host-integration 提案加入）：`internal/command/agents.go:169-180`（`allHostTargets` 含 `{"pi", filepath.Join(".pi","agents"), …}`）。

### 1.4 文档影响（对照现有小节）
| 受影响文档小节 | 现状 | 引用 |
|---|---|---|
| `README.md` §`## 一个需求如何实际推进`（L20；其中 §`### 1. 找到下一位责任人` L26） | 只写"agent 读取各 skill 的 description 自选合适的 skill"，全篇 0 次出现 `subagent` / `agents deploy`；未提内置角色与委派表 | `README.md:20`、`README.md:26-34`；`grep -c "subagent" README.md` → `0` |
| `README.md` §`## 安装与初始化`（L72；部署清单在 L92） | 称 init "部署 `.agents/skills/` 与 `ff-wiki/agents/`"，未提四个宿主的 agent 目录与 `agents.hosts` | `README.md:92` |
| `docs/architecture.md` §`## 实现边界`（L55；`internal/config` 行 L57） | L57 已列 `agents.disabled`，未列 `agents.hosts`；L60 已列 deploy/remove/status | `docs/architecture.md:57`、`docs/architecture.md:60` |
| `docs/skill-system.md` §`## Subagent 委派与协作`（L53） | 已落地 6 角色表，是本提案的文档成果，不因它过时 | `docs/skill-system.md:53-66` |
| `docs/cli-design.md` §`## 其他命令`（L47；agents 三行 L50-52） | 已列 `agents deploy/remove/status`（含 `agents.disabled` 语义），已更新 | `docs/cli-design.md:50-52` |
| `docs/cli-design.md` §`## 目录解析`（L5；init 清单 L9） | L9 仍只写 "部署 `<docs_dir>/agents/`、`.agents/skills/`"，未提 `.claude/agents/`、`.opencode/agent/`、`.codex/agents/`、`.pi/agents/` | `docs/cli-design.md:9` |

### 1.5 候选用户场景
> 装完 FlowForge 我不用再手写角色提示词：`agents deploy` 一次性把 analyst 到 reviewer 六个角色变成 Claude Code / OpenCode / Codex / PI 各自认得的原生子代理，只想用两个宿主就配 `agents.hosts`。

---

## 2. generic-role-orchestration

### 2.1 Feature 一句话
旗舰会话现在可以把"批量提取分析 / 结构化撰写回填 / 机械执行命令"这类与流程无关的活，按能力键派给三个通用角色（batch-analyst / scribe / executor，flash 档可钉），并在 AGENTS.md 里拿到任务链协议与结构化派发模板；讨论期笔记统一落 `<docs_dir>/research/` 再经 import 汇入 proposal。

### 2.2 表面变更清单
| 表面 | 事实 | 引用 |
|---|---|---|
| 新 agent 角色（3 个通用角色） | `flowforge-batch-analyst`（`default_skill: flowforge-research`、`model_profile: tool-capable`、`permission: workspace-write`）、`flowforge-scribe`（`flowforge-writing-for-agents`）、`flowforge-executor`（`flowforge-implement`，Boundaries 声明不进入票务流程） | `docs/proposals/generic-role-orchestration/design.md:31-37`（三角色表）、`design.md:39`；实物 `assets/subagents/flowforge-batch-analyst.md:1-12`（frontmatter 实测） |
| 复用而非新造 | 定向探查复用 `flowforge-investigator`（零改动）；决策素材简报 = investigator × `flowforge-research` 简报契约；不新造 `flowforge-scout` / `flowforge-brief-writer` / `flowforge-builder` | `docs/proposals/generic-role-orchestration/design.md:42-44` |
| AGENTS.md 模板新增段 | 新增 `## Generic capability dispatch`，插在 `## Agent skills` 与 `## Subagent delegation` 之间（顺序即适用优先级） | `docs/proposals/generic-role-orchestration/design.md:50`；实物 `assets/AGENTS.md:21`、`assets/AGENTS.md:64`；仓库自身 `AGENTS.md:50`、`AGENTS.md:93` |
| AGENTS.md 段内容三件套 | ① 能力键调度表（含档位建议，五基础行 + 后续扩展行）② 任务链协议 ③ 派发边界条款 | `docs/proposals/generic-role-orchestration/design.md:52-71`；实物 `assets/AGENTS.md:26-35`（能力表）、`AGENTS.md:37-40`（协议）、`AGENTS.md:42-47`（【目标】【输入】【输出】模板）、`AGENTS.md:49-51`（边界条款） |
| AGENTS.md 宿主提示（pi） | 末尾附 pi 提示小节：fork 继承编排会话上下文为只读参考、批量异步并行 + 完成唤醒收敛、用户级 `agentOverrides` 按名覆盖 | `docs/proposals/generic-role-orchestration/design.md:77`（明确"不改 `compile_pi.go`、不改 extension 代码"）；实物 `assets/AGENTS.md:53-57`、`AGENTS.md:86`（仓库自身） |
| AGENTS.md 调研落点句 | 讨论期产物落 `<docs_dir>/research/YYYY-MM-DD-<slug>.md` 带引用；创建 proposal 前检查未消费笔记并经 `/flowforge-import` 流转 | `docs/proposals/generic-role-orchestration/design.md:81-83`；实物 `assets/AGENTS.md:59-62` |
| 既有 skill 增补 | `flowforge-import` Inputs 节增补 `<docs_dir>/research/` 为标准源位置（"research 笔记按既有分类流转，无新分类"），其余零改动 | `docs/proposals/generic-role-orchestration/design.md:82`；实物 `assets/skills/flowforge-import/SKILL.md:14` |
| wiki 结构 | 约定 `<docs_dir>/research/`（不引入新顶层目录）；工作台中间产物另用 `docs/research/workbench/` 子目录 | `docs/proposals/generic-role-orchestration/design.md:81`；`issues/04-dogfood-batch-drill.md:30`（产物路径） |
| config 键 | 不新增任何 config 键；flash 钉扎复用既有 `agents.models` / `agents.models_by_name`（及后续 `models_by_host`）通道 | `docs/proposals/generic-role-orchestration/requirements.md:25`（目标 5）、`design.md:46`（"资产只声明 profile，不写死模型"） |
| 新 skill | 无新增 skill | `docs/proposals/generic-role-orchestration/design.md:31-44` 只列角色与复用关系 |

### 2.3 实施状态
**全部 closed（4/4）。**

证据：
- ticket 状态行：`docs/proposals/generic-role-orchestration/issues/01-generic-role-assets.md:15`、`02-agents-generic-dispatch.md:15`、`03-import-research-intake.md:15`、`04-dogfood-batch-drill.md:15` — 全部 `**Status:** closed`。
- 命令输出：`flowforge status --dir docs/proposals/generic-role-orchestration` → `4/4 resolved (100%)`、`Frontier Ready: 0`。
- 交付证据（角色资产）：`issues/01-generic-role-assets.md:88`（`go test ./internal/...` 全 ok，0 failures，含 `TestAgentsDeploy/Status/Remove`、`TestInitDeploys`、`TestUpgradeSync` 系）。
- 交付证据（协议真实运行）：`issues/04-dogfood-batch-drill.md:36-39`（Change 1：`ls .pi/agents/flowforge-batch-analyst.md` exit 0，"存在（部署后 pi 热加载零摩擦）"）；`04-dogfood-batch-drill.md:99-105`（Implementation note：3 × `flowforge-batch-analyst` 并行 run ID `3aabd040` / `2276cd41` / `2cd35f70`；三份 workbench 引用密度 103/67/40+，语义抽检 22/22、21/21 通过）；收敛产物 `docs/research/2026-09-25-dual-track-wiki-config.md`（`04-dogfood-batch-drill.md:105`）。
- 未完成项：0/4；设计 open item `oi-pi-dispatch-tool`、`oi-def-edit-deny` 为明确记录的非目标（`design.md:109-110`），不构成票务未完成。

### 2.4 文档影响（对照现有小节）
| 受影响文档小节 | 现状 | 引用 |
|---|---|---|
| `docs/skill-system.md` §`## Subagent 委派与协作`（L53；角色表 L57-64） | 表内只有 6 个流程角色，无 `flowforge-batch-analyst` / `flowforge-scribe` / `flowforge-executor`，未叙混合模型 | `docs/skill-system.md:57-64`；`grep -rn "batch-analyst" docs/skill-system.md` → 无命中（命中仅在 `docs/architecture.md:64`） |
| `README.md` §`## 一个需求如何实际推进`（L20） | 全篇 0 次出现 `subagent` / `batch-analyst` / 能力键调度；未叙"先按能力找角色、流程态再查委派表" | `README.md:20-70`；`grep -c "subagent" README.md` → `0` |
| `README.md` §`## 安装与初始化`（L72；资产清单 L92） | 未提新增的 3 个通用角色随 init/upgrade 部署 | `README.md:92` |
| `docs/architecture.md` §`## 实现边界`（L55；L64 混合名册行） | 已叙"6 流程角色 + 3 通用能力角色"与段序优先级，已更新（非过时） | `docs/architecture.md:64` |
| `docs/cli-design.md` §`## 目录解析`（L5；`docs_dir` 默认值 L7） | L7 称 `docs_dir` "默认为 `docs`"，与实物 `ff-wiki` 冲突（`internal/config/config.go:17` `DefaultDocsDir = "ff-wiki"`），也与 README:92 的 `ff-wiki` 冲突；research 落点段未在此文档出现 | `docs/cli-design.md:7`；`internal/config/config.go:17`；`README.md:92` |
| `docs/cli-design.md` §`## Artifact Catalog`（L11；角色枚举 L13） | 已含 `research` 角色枚举，与本提案落点约定相容 | `docs/cli-design.md:13` |

### 2.5 候选用户场景
> 上周我只想把"这套调度协议"固化下来：现在我一句"照模板拆三个批次分析员"，三个 flash 子代理并行产出带行号引用的笔记，我只看收敛结果，不再每次重写提示词。

---

## 3. agent-model-preservation

### 3.1 Feature 一句话
项目在已部署的 agent 文件 frontmatter 里手写过 `model:` 时，升级或再次 `agents deploy` 不再把它静默抹掉——会保留该值并在 stderr 提示，而 config 里显式钉的模型始终优先。

### 3.2 表面变更清单
| 表面 | 事实 | 引用 |
|---|---|---|
| CLI 行为（无新命令/flag） | 变更落在既有 `agents deploy` / `init --force` / `upgrade` 的部署路径内，不改 CLI 签名 | `docs/proposals/agent-model-preservation/requirements.md:25`（"不改 Issue Schema、不改 CLI 接口签名；合并发生在 deploySubagents 写文件前"）；`docs/proposals/agent-model-preservation/design.md:20` |
| 保留合并（preserve-merge）语义 | 重编译前读取既有部署文件 frontmatter 的 `model:`，新内容缺该字段时注入旧值 | `docs/proposals/agent-model-preservation/requirements.md:37`（术语）、`design.md:19-20` |
| 优先级链 | `Model`（config）> `FallbackModel`（文件残留）> profile 默认 | `docs/proposals/agent-model-preservation/design.md:21`、`design.md:23` |
| stderr 提示文案 | 采用 fallback 时输出一行提示；设计文本为 `info: preserved local model %q for %s (set agents.models in .flowforge/config.yaml to pin explicitly)` | `docs/proposals/agent-model-preservation/design.md:24`；实物文案（后者）：`internal/command/agents.go:623` → `set agents.models_by_name/models_by_host in .flowforge/config.yaml to pin explicitly`（与设计文本的键名不同，逐字见 `internal/command/agents_test.go:1627`） |
| 适用宿主 | claude 与 opencode（有 `model:` 字段）；pi 因 revision 2 的有条件注入同样适用；codex 无 model 字段，`FallbackModel` 被忽略 | `docs/proposals/agent-model-preservation/requirements.md:23`、`design.md:25`；pi 侧 `internal/subagent/compile_pi_test.go:41-42` |
| 实现缝 | `CompileOptions` 增 `FallbackModel string`；claude 适配器补 opts 布线 | `docs/proposals/agent-model-preservation/design.md:20-21`；实物 `internal/subagent/compile_opencode.go:19`、`compile_claude.go:18` |
| status 语义同步 | `computeSubagentStatus` 期望内容纳入 preserve-merge，preserved 文件报 `current` 不再误报 `drifted`；config 钉死仍 drifted、body 手改 drifted、文件缺失 missing | `docs/proposals/agent-model-preservation/issues/02-status-preserve-semantics.md:112`；实物 `internal/command/agents_status.go:90-92` |
| config 键 | 无新增键（沿用 `agents.models` / `models_by_name` / `models_by_host`） | `docs/proposals/agent-model-preservation/requirements.md:22`、`design.md:20` |
| 新 skill / 新 agent 角色 / wiki 结构 | 均无 | `docs/proposals/agent-model-preservation/design.md:19-25`（变更仅限部署合并逻辑与提示） |

### 3.3 实施状态
**全部 closed（2/2）。**

证据：
- ticket 状态行：`docs/proposals/agent-model-preservation/issues/01-preserve-merge-model.md:15`、`02-status-preserve-semantics.md:15` — 均为 `**Status:** closed`。
- 命令输出：`flowforge status --dir docs/proposals/agent-model-preservation` → `2/2 resolved (100%)`、`Frontier Ready: 0`。
- 交付证据：`issues/02-status-preserve-semantics.md:113`（`go test ./internal/...` 全绿：新增 `TestAgentStatusTreatsPreservedModelAsCurrent` 3 子测试 + `TestDeployPreservesLocalModel` 8 子测试 + 既有 5 个 status 用例；CLI e2e：手编 model → deploy → status current（exit 0）→ 手改 body → status drifted（exit 1））。
- 一致性事实：ticket 01 声称三态表（文件有+config 无→保留+提示；文件有+config 有→config；文件无→无）——`docs/proposals/agent-model-preservation/design.md:42`；实物测试同时覆盖 codex 无 model 与首次部署回归（`design.md:42`）。

### 3.4 文档影响（对照现有小节）
| 受影响文档小节 | 现状 | 引用 |
|---|---|---|
| `docs/cli-design.md` §`## 其他命令`（L47；`agents deploy` 行 L50、`agents status` 行 L52） | 两行只写"编译并部署到宿主原生目录"与"报告 current/missing/drifted/project-owned"，未提 preserve-merge 与 stderr 保留提示、config 优先规则 | `docs/cli-design.md:50-52` |
| `README.md` §`## 安装与初始化`（L72；`init --force` 说明 L99） | L99 只说 "`init --force` 同步受管资产，但保留已有项目配置"，未提部署产物 frontmatter 中的本地 `model:` 会被保留 | `README.md:99` |
| `README.md` §`## 一个需求如何实际推进`（L20） | 未涉及模型钉扎/保留 | `README.md:20-70` |
| `docs/architecture.md` §`## 实现边界`（L55；`internal/config` 行 L57） | 已列 `agents.disabled`，未提及 preserve-merge 属部署语义 | `docs/architecture.md:57` |

### 3.5 候选用户场景
> 我在 `.opencode/agent/flowforge-reviewer.md` 手写的 `model: cpa/glm-5.3` 曾经被一次升级悄悄抹掉；现在升级后它还在，stderr 还告诉我"已保留本地模型、要钉死请写 config"。

---

## 4. model-sets-switching

### 4.1 Feature 一句话
一台机器可以在多套命名模型方案之间一键切换（如旗舰从 glm-5.3 换成限额时段的 mimo-2.6-pro），切换会自动重部署 agent 文件，失败则回滚，不会留下"指针切了产物没切"的中间态。

### 4.2 表面变更清单
| 表面 | 事实 | 引用 |
|---|---|---|
| 新 CLI 命令组 | `flowforge model-set`，含 `list` / `use <name>` / `show [name]` | `docs/proposals/model-sets-switching/requirements.md:29-33`；`docs/proposals/model-sets-switching/design.md:62-66`；实物 `internal/command/modelset.go:25`（`Use: "model-set"`）、`modelset.go:52`、`modelset.go:88`、`modelset.go:152`；命令实测 `flowforge model-set --help` → `list List model sets and mark the active one / show Show the effective per-agent model table for a model set / use Switch the active model set and redeploy agents` |
| 新 config 键 | `agents.model_sets`：命名 overlay 集，set 内复用现有三层 `models` / `models_by_name` / `models_by_host`（增量覆盖，未声明的角色继续用基础层） | `docs/proposals/model-sets-switching/requirements.md:22-25`；`docs/proposals/model-sets-switching/design.md:15-21`（yaml 样例）、`design.md:27-34`（Go 结构）；实物 `internal/config/config.go:47`（`yaml:"model_sets,omitempty"`）、`internal/config/modelset.go:15`（`ModelSetConfig`） |
| 新状态文件 | `.flowforge/model-set.active`（一行 set 名；不存在 = 基础层），保持 gitignored，不写回 config.yaml | `docs/proposals/model-sets-switching/design.md:41-45`；实物 `internal/config/modelset.go:25-33`（`ActiveModelSetPath` / `ReadActiveModelSet`）、`modelset.go:48`（`WriteActiveModelSet`）、`modelset.go:58`（`ClearActiveModelSet`） |
| 切换原子性 | 顺序：校验 set → 写指针 → 走既有 deploy 路径 → deploy 报错则回滚指针并以非 0 退出 | `docs/proposals/model-sets-switching/requirements.md:36-39`；`docs/proposals/model-sets-switching/design.md:68-75` |
| list 输出标注 | 基础层显示为 `default` 并标 `ACTIVE` | `docs/proposals/model-sets-switching/requirements.md:31`；命令实测 `flowforge model-set list` → `ACTIVE  default (base agent model layers)` 下一行 `        offpeak` |
| show 输出 | 展示 set 与基础层**合并后**的有效模型表，并标注来源 `base` / `set` | `docs/proposals/model-sets-switching/design.md:65`；命令实测 `flowforge model-set show` → `model set: default` + 12 行 `flowforge-*  <model>  base` |
| 非目标（明确排除） | 不自动检测限额/惩罚时段、不做单角色临时 override 命令、不做跨机器同步 | `docs/proposals/model-sets-switching/requirements.md:43-45` |
| 新 skill / 新 agent 角色 / wiki 结构 | 均无（纯 config + CLI + 状态文件） | `docs/proposals/model-sets-switching/design.md:7-45`（"结构"节只含 config 层、Go 结构、指针文件） |

### 4.3 实施状态
**实现完成（小型提案，无 issues 目录）。**

证据：
- 需求文件自述状态：`docs/proposals/model-sets-switching/requirements.md:4` → `> Type: proposal(小型)` / `> Status: implemented (2026-09-29,全部测试绿+本仓双向往返验收通过)`；同处 L5 记录三项决策（"A 增量覆盖 / B 切换自动重部署+失败回滚 / C set 内复用现有三层结构"）。
- 无 ticket 可判：`docs/proposals/model-sets-switching/` 下只有 `design.md` 与 `requirements.md`（无 `issues/`）；命令输出 `flowforge status --dir docs/proposals/model-sets-switching` → `No features or issues found in docs/proposals/model-sets-switching`。
- 载体提交：`git log --oneline` → `6df209a feat(model-set): 命名模型方案与快速切换`（当前 HEAD）。
- 实物回归：测试文件 `internal/config/modelset_test.go`（overlay / 指针读写四例）与 `internal/command/modelset_cmd_test.go` 均存在。

### 4.4 文档影响（对照现有小节）
| 受影响文档小节 | 现状 | 引用 |
|---|---|---|
| `docs/cli-design.md` §`## 其他命令`（L47；命令清单 L49-55） | 清单中无 `flowforge model-set` 任何一行 | `docs/cli-design.md:49-55`；`grep -rn "model-set" docs/cli-design.md` → 无命中 |
| `README.md`（全文 115 行） | `grep -rn "model-set" README.md` → 无命中；`grep -c "model" README.md` → `0`；`## 安装与初始化` 与 §`## 一个需求如何实际推进` 均未提模型方案切换 | `README.md:72-99`、`README.md:20-70` |
| `docs/architecture.md` §`## 实现边界`（L55；`internal/config` 行 L57） | L57 的 config 键清单未含 `model_sets`；无 model-set 状态文件的承载说明 | `docs/architecture.md:57`；`grep -rn "model_sets" docs/architecture.md` → 无命中 |
| `docs/skill-system.md`（全文 94 行） | 无模型方案切换相关内容（属 CLI + config 面，非 skill 面） | `grep -rn "model-set" docs/skill-system.md` → 无命中 |
| `docs/cli-design.md` §`## 稳定边界`（L65） | L67-70 的四条稳定边界未涉及部署态与指针一致性 | `docs/cli-design.md:65-70` |

### 4.5 候选用户场景
> 白天限额被限速时，我敲一句 `flowforge model-set use offpeak`，四个决策/审查角色立刻换成 mimo-2.6-pro 并自动重部署；夜里切回 `default` 也一样，切错失败还会自动回滚。

---

## 5. pi-host-integration

### 5.1 Feature 一句话
PI 成为第四个一等部署宿主：`agents deploy` 会写出 `.pi/agents/*.md` 原生子代理与项目级 `.pi/extensions/flowforge.ts`（拦截受保护 test 文件的 write/edit、把 frontier/check 注册为 LLM 原生工具），并且 `flowforge frontier --pi-workflow` 能直接输出可被 pi-subagents 执行的批次脚本。

### 5.2 表面变更清单
| 表面 | 事实 | 引用 |
|---|---|---|
| 新宿主 | `agents.hosts` 支持枚举值 `pi`；部署目标 `.pi/agents/<name>.md`；`remove` / `cleanDeselectedHosts` / `status` 经既有 hostTarget 抽象自动覆盖 | `docs/proposals/pi-host-integration/requirements.md:18`；`docs/proposals/pi-host-integration/design.md:33`、`design.md:26`；实物 `internal/command/agents.go:180`（`{"pi", filepath.Join(".pi","agents"), ".md", …}`）、`agents.go:195`、`agents.go:209`（错误消息列出 `claude, opencode, codex, pi`）、`agents.go:30`（命令 Long 含 `.pi/agents/`） |
| 新编译目标 | `CompilePiWithOptions(def, opts)`；PI frontmatter 字段集 `name` / `description` / `model`（有条件）/ `thinking` / `tools` / `skills` / `inheritSkills` | `docs/proposals/pi-host-integration/design.md:21`、`design.md:23-31`（映射表）；实物 `internal/subagent/compile_pi.go:27`（`CompilePiWithOptions`）、`compile_pi.go:28-29`（字段固定顺序注释）、`compile_pi.go:34`、`compile_pi.go:37` |
| 档位映射 | `high-capability` → `thinking: high`，其余 → `medium`；只读角色写 `tools: read, grep, find, ls`，其余不写 `tools` | `docs/proposals/pi-host-integration/design.md:26-27`；验收实测 `issues/01-pi-host-compile-target.md:127`（抽查 `flowforge-investigator.md` 含 `thinking: medium`、`tools:` 白名单、`skills:`、`inheritSkills: false`） |
| skill 绑定 | `DefaultSkill` + `DetourSkills` → `skills`（列表）+ `inheritSkills: false`（不继承全局目录） | `docs/proposals/pi-host-integration/design.md:28`；实物 `internal/subagent/compile_pi.go:64-65`（`piSkills`） |
| 新 CLI flag | `flowforge frontier --pi-workflow`（优先于 `--json` / `--quiet`） | `docs/proposals/pi-host-integration/design.md:37`；实物 `internal/command/frontier.go:139`；命令实测 `docs/cli-design.md:34`（用法行）、`cli-design.md:45`（语义行） |
| flag 产物形态 | pi-subagents workflowScript：顺序 fail-fast；每 ready ticket 一次顶层 `await runs.run({agent: "flowforge-implementer", task: …, gate: {command: "flowforge check --dir <docs>/proposals"}})`；顶层语句、显式 `return`、禁嵌套函数/箭头函数；空批次输出单行提示脚本 | `docs/proposals/pi-host-integration/design.md:39-45`；实物 `internal/command/frontier.go:207-208`（任务前后缀文本常量）、`frontier.go:229`（`jsString("flowforge-implementer")` + `jsString("fresh")` + gate 命令） |
| 新资产文件（项目级 pi extension） | `assets/pi/flowforge.ts` → 部署为 `.pi/extensions/flowforge.ts`；随 `pi` 宿主启用与否收敛 | `docs/proposals/pi-host-integration/design.md:51`；实物 `assets/pi/flowforge.ts`（存在）；`internal/command/agents.go:441`（`piExtensionRelPath`） |
| extension 功能 1：test 写保护 | `pi.on("tool_call")` 拦截 `write`/`edit` 且路径匹配保护 glob 时返回 `{block: true, reason}` | `docs/proposals/pi-host-integration/design.md:53`；实物 `assets/pi/flowforge.ts:223`、`flowforge.ts:224`、`flowforge.ts:230` |
| extension 功能 2：原生工具 | 注册 `flowforge_frontier`（透传 `--json`）与 `flowforge_check`，execute 内 spawn CLI（PATH 优先，回退 `<project>/bin/flowforge`） | `docs/proposals/pi-host-integration/design.md:54`；实物 `assets/pi/flowforge.ts:255`、`flowforge.ts:258`、`flowforge.ts:178`（二进制解析） |
| extension 功能 3：逃生阀 | `agents.disable_test_guard: true` 时不注册写拦截（与 opencode 同语义） | `docs/proposals/pi-host-integration/design.md:55`；实物 `assets/pi/flowforge.ts:13`、`flowforge.ts:96` |
| config 键 | 复用 `agents.test_file_globs`（glob 单一真相，缺省回退内置默认集）；`agents.disable_test_guard`；`agents.hosts` 增 `pi` 枚举值。无新增键 | `docs/proposals/pi-host-integration/design.md:53`、`design.md:55`；实物 `internal/config/config.go:48`（`test_file_globs`）、`config.go:49`（`disable_test_guard`）；默认 glob 集实物 `internal/command/agents.go:226-228`（设计文档 `design.md:53` 引用为 `agents.go:222`） |
| 新 skill / 新 agent 角色 | 均无新增角色定义；PI 复用同一批 `assets/subagents/*.md` 权威定义 | `docs/proposals/pi-host-integration/design.md:21`（"正文不改写"）、`design.md:25`（正文 Default Skill 段双通道回退在 PI 下天然生效） |

### 5.3 实施状态
**全部 closed（3/3）。**

证据：
- ticket 状态行：`docs/proposals/pi-host-integration/issues/01-pi-host-compile-target.md:19`、`02-frontier-pi-workflow-format.md:19`、`03-pi-extension.md:19` — 均为 `**Status:** closed`。
- 命令输出：`flowforge status --dir docs/proposals/pi-host-integration` → `3/3 resolved (100%)`、`Frontier Ready: 0`。
- 交付证据（01）：`issues/01-pi-host-compile-target.md:126`（`go test ./internal/... -v` 全部 ok；新增 `TestCompilePiFields`、`TestAgentsDeployPiHost`、`TestAgentsDeployCleansPiWhenDeselected` 均 PASS）、`issues/01…:127`（实机 `agents deploy` exit 0，`.pi/agents/` 6 份文件，抽查 `flowforge-investigator.md` 含 `thinking: medium`/`tools:`/`skills:`/`inheritSkills: false`）。
- 交付证据（02）：`issues/02-frontier-pi-workflow-format.md:134`（`go test -count=1 ./internal/...` 全部 ok、`go vet` 无告警）。
- 交付证据（03）：`issues/03-pi-extension.md:40`（Node 冒烟 `ALL SMOKE CHECKS PASSED`，含 `binary resolved: /home/biqiang/.local/bin/flowforge`、`frontier output bytes: 34804`）、`issues/03…:144`（glob 9 断言 / 配置解析 4 断言 / 拦截器 5 断言 / 工具真实执行 1 断言）。
- 残留未做项（票内明示，非未完成票）：`issues/03-pi-extension.md:130`（交互式 PI 会话手工冒烟留交操作者："需真实 PI 会话，本执行环境无交互 TTY"）；`issues/02…:134`（交互式 workflowScript 冒烟留给操作者）。
- 诊断状态：三个 ticket 均带 `upstream-changed` waiver（Consumed revision 1 is behind 2，理由指向 `deploy-artifact-localization`）——`issues/01…:8-11`、`issues/02…`、`issues/03…:8-11`；`flowforge check` 报为 waived warning（见本文档首部说明）。

### 5.4 文档影响（对照现有小节）
| 受影响文档小节 | 现状 | 引用 |
|---|---|---|
| `docs/cli-design.md` §`## 其他命令`（L47；PI 宿主说明 L57-63） | 已完整覆盖 `.pi/agents/`、扩展两处部署、作用范围、逃生阀与手工冒烟，已更新（非过时）；L60 已含 `agents.hosts` 与 `agents.test_file_globs` | `docs/cli-design.md:57-63` |
| `docs/cli-design.md` §``flowforge frontier``（L31；用法 L34、flag 语义 L45） | 已列 `--pi-workflow` 及其语义与优先级，已更新 | `docs/cli-design.md:34`、`docs/cli-design.md:45` |
| `docs/architecture.md` §`## 实现边界`（L55；L58 四宿主行、L62 assets 行） | 已叙"四宿主（Claude Code、OpenCode、Codex、PI）"与 `assets/pi/flowforge.ts`，已更新 | `docs/architecture.md:58`、`docs/architecture.md:62` |
| `README.md`（全文 115 行） | 0 次出现 `pi` / `PI`（`grep -c "PI" README.md` → `0`）；§`## 安装与初始化`（L72）只列 `.agents/skills/` 与 `ff-wiki/agents/`，§`## 一个需求如何实际推进`（L20）未提 PI 批次委派；README 无 PI 前置依赖（pi-subagents 扩展）说明 | `README.md:72-99`、`README.md:20-70` |
| `docs/skill-system.md` §`## Subagent 委派与协作`（L53；宿主枚举 L55） | L55 只列 "Claude Code Subagents、OpenCode Agent Tool / `@mention`、Codex 子会话"，未含 PI | `docs/skill-system.md:55` |

### 5.5 候选用户场景
> 我在 PI 会话里直接问"下一批该干什么"，`flowforge_frontier` 工具就返回 JSON；让 PI 自己跑批次时用 `frontier --pi-workflow` 生成脚本，每张票都由一个干净上下文的 implementer 干完再自动 `flowforge check`。

---

## 附：本组事实的可复核命令清单

| 用途 | 命令 |
|---|---|
| 收口状态（逐 proposal） | `./bin/flowforge status --dir docs/proposals/<id>` |
| 全仓诊断（无 error） | `./bin/flowforge check`（尾部 `✓ Dependency graph is healthy`） |
| CLI 面实测 | `flowforge agents --help`；`flowforge model-set --help`；`flowforge model-set list`；`flowforge model-set show` |
| 资产实物 | `ls assets/subagents/`（12 份）；`ls assets/pi/`（`flowforge.ts`） |
| 文档缺失面 | `grep -rn "model-set\|model_sets\|agents.hosts\|batch-analyst\|PI" README.md docs/architecture.md docs/skill-system.md docs/cli-design.md` |
