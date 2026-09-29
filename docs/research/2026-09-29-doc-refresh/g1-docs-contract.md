# g1 — 文档契约 / 技能路由 / 规范注入 / wiki 单轨 / 外部材料：面向用户的 feature 事实提取

> 分组：`documentation-contract-refinement`、`skill-routing-simplification`、`standards-injection`、
> `wiki-config-single-track`、`external-material-intake`
> 方法：只陈述文件与命令输出里可复核的事实（`文件:行号` 或命令原文），不做跨组综合、不提改进建议、不做决策。
> 事实来源：各 proposal 的 `requirements.md` / `design.md` / `spec.md` / `issues/*.md`（Status 与 Completion evidence），
> 外加当前仓库实物（`README.md`、`docs/architecture.md`、`docs/skill-system.md`、`docs/cli-design.md`、`assets/…`、`internal/…`）。
> 复核时点：2026-09-29，二进制 `bin/flowforge`（12,428,418 bytes，mtime Sep 29 19:31）。
> 全部"表面变更"均已用命令在实物中复核，命令与输出见 §7。
> 简写约定：各小节中的 `…/` 指 `docs/proposals/<该小节 proposal id>/`；`…/issues/NN-….md` 即该 proposal 的 ticket 文件。

---

## 1. documentation-contract-refinement

### 1.1 Feature 一句话
把「需求 → 设计 → ticket → 证据」四类工件分工固化成可确定性检查的契约：`check` / `frontier` 只从当前文件事实算出可执行工作，warning 默认放行、gap 默认不进前沿、blocker 不可绕过，ticket 关闭前必须留下真实验证证据。

### 1.2 表面变更清单

| 表面 | 事实 | 引用 |
|---|---|---|
| CLI 命令/flag | `flowforge check [--dir <path>] [--json] [--strict]` | `docs/cli-design.md:26`；实物 `internal/command/check.go:115-116`（`--json`、`--strict` 定义） |
| CLI 命令/flag | `flowforge frontier [--dir <path>] [--json] [--quiet] [--strict] [--include-gaps] [--pi-workflow]` | `docs/cli-design.md:34`；实物 `internal/command/frontier.go:136-139` |
| CLI 退出策略 | 默认 warning/gap 可见但不失败，blocker 始终失败；`--strict` 让未豁免 warning/gap 也失败 | `docs/cli-design.md:29`；proposal 决策 `docs/proposals/documentation-contract-refinement/spec.md:239` |
| 工件角色 | 四个逻辑角色 Requirement / Solution design / Ticket / Evidence 各自的唯一含义 | `docs/proposals/documentation-contract-refinement/spec.md:91-96`（spec 内正文自 `spec.md:91` 起）；落地文档 `README.md:9-12`、`docs/architecture.md:32-37` |
| wiki/目录结构 | `<docs_dir>/proposals/<feature>/{requirements,design,spec}.md` 按信息价值选用；`issues/*.md` 唯一可执行；`evidence/*.md` 需独立生命周期时才建 | `docs/architecture.md:18-22`；proposal 决策 `spec.md:236-238`（`<docs_dir>/proposals/<feature>/issues/*.md` 是唯一可执行位置） |
| schema（机器层） | 按需 Schema v1：id / revision / area / `consumes` / open item / waiver | `docs/architecture.md:39`、`docs/cli-design.md:17-21`（确定性诊断清单）；proposal 决策 `spec.md:240` |
| 新增诊断 | closed ticket 缺非空 `Completion evidence` → `missing-completion-evidence` warning，`check --strict` 视为失败 | `docs/cli-design.md:21`、`docs/architecture.md:70`；ticket `docs/proposals/documentation-contract-refinement/issues/08-implementation-review-evidence.md:38` |
| 新增 skill | `flowforge-solution-design`（模块责任/接口/seam/迁移/验证策略的 owner） | ticket `…/issues/05-solution-design-skill.md:34`；落地文档 `docs/skill-system.md:12` |
| 新增 skill | `flowforge-to-spec` 变成可选导航（compact work 可跳过，非权威） | ticket `…/issues/09-optional-spec-navigation.md:33-34`；落地文档 `docs/skill-system.md:13` |
| 共享契约资产 | 所有生产 skill 共同读取 `assets/skills/_shared/ARTIFACT-CONTRACT.md` | ticket `…/issues/04-deploy-artifact-contract.md:35`；落地文档 `docs/skill-system.md:70` |
| skill 增补 | Align 持久化 requirement authority；Plan 输出 Delivery / Design context / Blocked by / Touch points / Changes / Constraints / Done and verify | ticket `…/issues/06-align-and-route-authority.md:35`、`…/issues/07-plan-high-information-tickets.md:35`；落地文档 `docs/skill-system.md:11,14` |
| agent 角色 | 本 proposal 未新增 subagent 角色（skill 侧的 `flowforge-solution-design` 之外无角色变更证据） | `…/issues/05-solution-design-skill.md:34`（只记 Skill）；`docs/skill-system.md:57-64` 角色表 |
| config 键 | 本 proposal 未新增 config 键（`docs_dir` / `--dir` 语义沿用） | `docs/cli-design.md:5-9`；`docs/proposals/documentation-contract-refinement/spec.md:237`（仅以 `<docs_dir>` 指路径，未定义配置键） |

### 1.3 实施状态

**issues 全部 closed（10/10），并已产出独立 end-to-end evidence。**

- 状态行逐票：`issues/01-artifact-catalog-discovery.md:5`、`02:5`、`03:5`、`04:5`、`05:5`、`06:5`、`07:5`、`08:5`、`09:5`、`10:5` 均为 `**Status:** closed`（`**Completion evidence` 节分别位于 34/33/35/33/32/33/33/33/31/35 行）。
- 需求文档自身状态：`docs/proposals/documentation-contract-refinement/spec.md:14` → `**Status:** Implemented by tickets 01–10; verified by [end-to-end evidence](evidence/end-to-end.md)`。
- 证据文件存在且有 schema envelope：`docs/proposals/documentation-contract-refinement/evidence/end-to-end.md:1-12`（`role: evidence`、`consumes.requirements.documentation-contract-refinement: 2`）。
- 关闭声明引用：`…/issues/10-end-to-end-migration-proof.md:37`（一个 compact run + 本 proposal 十票 complex run、12 个 coordination case 映射证据）。
- 当前二进制复核（见 §7）：`flowforge check --dir docs/proposals/documentation-contract-refinement --strict` → **exit=1，15 条 warning（10 `legacy-metadata` + 5 `missing-required-metadata`）**，全部落在 `scenario-fixtures/` 下无 schema 的示例文件上，未涉及 `issues/`。

### 1.4 文档影响（对照四份现行文档）

| 文档小节 | 现状事实 | 引用 |
|---|---|---|
| `README.md` §"核心文档"（`README.md:101-106`） | 该节第 106 行把 `spec.md` 标为「本次文档契约重构规格」；被指向文件自述已完成（tickets 01–10） | `README.md:106` + `docs/proposals/documentation-contract-refinement/spec.md:14` → 该条目表述与文件状态不符 |
| `README.md` §"4. 实现、审查并留下证据"（`README.md:57-66`） | 该节只讲 ticket 内记录与 `closed`，未提 evidence 作为独立工件角色/`evidence/*.md` 提升条件 | `README.md:57-66`；对照 `docs/architecture.md:22`（`evidence/*.md` 布局）与 `…/issues/08-implementation-review-evidence.md:19`("promote evidence only when …") → README 该节缺 promoted evidence 分支 |
| `README.md` §"整体设计理论"（`README.md:5-18`） | 四角色与 `--include-gaps` 已在该节出现 | `README.md:9-12`、`README.md:16` → 与 proposal 一致，无过时 |
| `README.md` §"3. 发布可执行 ticket"（`README.md:44-55`） | 已写 `issues/` 路径、`check`/`frontier` 校验范围 | `README.md:46,55` → 与 `docs/cli-design.md:26,29` 一致 |
| `docs/architecture.md` §"自适应工件，而非固定阶段"（`docs/architecture.md:41-53`） | DAG edge / warning / gap / blocker / waiver 五类事实与默认策略齐备 | `docs/architecture.md:47-51` → 与 `spec.md:239` 一致，无过时 |
| `docs/architecture.md` §"完成不变量"（`docs/architecture.md:68-70`） | 已写 `missing-completion-evidence` 与严格检查失败语义 | `docs/architecture.md:70` → 无过时 |
| `docs/skill-system.md` §"主交付链"（`docs/skill-system.md:5-18`） | 已列 `flowforge-solution-design`（12 行）与可选 `flowforge-to-spec`（13 行） | `docs/skill-system.md:12-13` → 无过时 |
| `docs/skill-system.md` §"工件协作规则"（`docs/skill-system.md:68-79`） | 已写共享 `ARTIFACT-CONTRACT.md` 六条与 Schema v1 兼容层 | `docs/skill-system.md:70-79` → 无过时 |
| `docs/cli-design.md` §"`flowforge check`"（`docs/cli-design.md:23-29`）/§"`flowforge frontier`"（`docs/cli-design.md:31-45`） | flag 清单、退出策略、`--include-gaps` 优先级、`--quiet` 流分离均已写全 | `docs/cli-design.md:26,29,34,43,45` → 无过时 |
| `docs/cli-design.md` §"Artifact Catalog"（`docs/cli-design.md:11-21`） | 七类 artifact role、`issues/*.md` 唯一可执行、五类确定性诊断齐备 | `docs/cli-design.md:13,17-21` → 无过时 |

### 1.5 候选用户场景
「票里写清复用 seam 和验证命令，收工没填 Completion evidence 会被 check 拦下。」

---

## 2. skill-routing-simplification

### 2.1 Feature 一句话
删掉中心路由 skill 与 front-matter 调用闸门后，agent 直接读各 skill 的扁平 `description` 自选责任人；每个 `description` 必须自带「触发短语 + 下游所有权 + 负边界」，用户说「检查 review 发现」时不再误入设计或拆票。

### 2.2 表面变更清单

| 表面 | 事实 | 引用 |
|---|---|---|
| skill 删除 | 删除 `flowforge-route` skill 目录（源与部署侧），"选择 next owner" 责任溶解到 description 自选 | `docs/proposals/skill-routing-simplification/requirements.md:17`、`design.md:41,45`；ticket `…/issues/01-remove-dispatch-mechanism.md:43-44,48` |
| front-matter | 从 14 个 SKILL.md 移除 `disable-model-invocation: true`；`user-invocable` 全库从未存在 | `…/issues/01-remove-dispatch-mechanism.md:19,45,53`；需求 `requirements.md:16`；闭环记录 `…/issues/01-remove-dispatch-mechanism.md:127,135`（`user-invocable` field never existed in codebase (grep 0)） |
| 新增 skill | 无（本 proposal 只删与改 description） | `requirements.md:24-33`（In/Out of scope 均无新增 skill） |
| skill 数量 | 删 route 后现为 28 个 `assets/skills/flowforge-*` 目录；`.agents/skills/` 同样 28 个 | §7 命令输出 `ls -d assets/skills/flowforge-* \| wc -l → 28` / `.agents/skills/flowforge-* → 28`（含后续其他 proposal 新增的 skill，route 已不在列） |
| description 内容约定 | 每 skill 的 `description` 必含三段：触发短语 / 下游所有权 / 负边界（NOT for X）；触发短语禁用正文术语 | 需求 `requirements.md:18,64`；设计 `design.md:51,55`；落地文档 `docs/skill-system.md:26-38` |
| P0 三角重写 | `flowforge-review`（owns fix planning）、`flowforge-solution-design`（NOT for review-fix design）、`flowforge-plan`（NOT for fix changes） | 设计 `design.md:73,77,81`；ticket `…/issues/03-rewrite-p0-triangle-descriptions.md:15`(closed) |
| 主交付链与支持路径 description 批量重写 | 主轴 / 支持路径两批各自的 description 重写 | ticket `…/issues/04-rewrite-main-chain-descriptions.md:15`、`…/issues/05-rewrite-support-path-descriptions.md:15`（均 closed） |
| deploy 行为（用户可见） | `deployManagedAssets` 在 skills 拷贝后调用 `cleanupRemovedSkillDirs`：清理 source 已不存在的 `flowforge-*` 目标目录，保留 `_shared` 与用户自建目录，幂等 | ticket `…/issues/06-deploy-cleanup-removed-skills.md:15,132`（closed）；实物 `internal/command/assets_deploy.go:38` |
| 文档权威 | `docs/skill-system.md` 新增「Description-driven dispatch」节 + 两个子节；`assets/skills/flowforge-writing-for-agents/SKILL-MECHANICS.md` 新增「Description content convention」 | `…/issues/02-document-dispatch-convention.md:36-37,90`（closed）；落地 `docs/skill-system.md:20-38` |
| agent 角色 | 无新增/删除 subagent 角色 | `requirements.md:24-33`（作用面仅 skill front-matter、AGENTS 路由表、文档、init 提示、测试） |
| config 键 | 无新增 config 键 | `design.md:19-33`（未列 config 变更） |
| wiki 结构 | 无变更 | `requirements.md:24-33` |

### 2.3 实施状态

**issues 全部 closed（6/6）；DELETE 与描述重写已落地，deploy 清理逻辑已在实物中。**

- 状态行逐票：`issues/01-remove-dispatch-mechanism.md:15`、`02-document-dispatch-convention.md:15`、`03-rewrite-p0-triangle-descriptions.md:15`、`04-rewrite-main-chain-descriptions.md:15`、`05-rewrite-support-path-descriptions.md:15`、`06-deploy-cleanup-removed-skills.md:15` 均为 `**Status:** closed`。
- 实物复核：`assets/skills/flowforge-route` 与 `.agents/skills/flowforge-route` 均 `No such file or directory`；`AGENTS.md` / `assets/AGENTS.md` grep `flowforge-route\|Route & Guide` → 0 命中；`README.md:28` 现为「agent 读取各 skill 的 description 自选合适的 skill」（对比 `…/issues/01:34` 记的旧状态「README.md — line 28 route mention」）。
- 关闭证据：`…/issues/01-remove-dispatch-mechanism.md` completion evidence（`grep -rn "flowforge-route" … → 0 hits, PASS`、`flowforge assets verify → all 49 entries "current"`、`check --strict → Dependency graph is healthy`）。
- 当前二进制复核（见 §7）：`flowforge check --dir docs/proposals/skill-routing-simplification --strict` → **exit=1，48 条 `evidence-missing` warning**（针对已 checked 的 Change 缺 evidence 四元组）；票面当时的 `--strict` 通过记录见 `…/issues/02-document-dispatch-convention.md:78,100`。

### 2.4 文档影响（对照四份现行文档）

| 文档小节 | 现状事实 | 引用 |
|---|---|---|
| `README.md` §"1. 找到下一位责任人"（`README.md:26-34`） | 已改为 description 自选，无 route 提及 | `README.md:28`；对照 `…/issues/01:48` → 已同步 |
| `docs/skill-system.md` §"Description-driven dispatch"（`docs/skill-system.md:20-38`） | 三段约束与「触发短语禁用正文术语」已写入 | `docs/skill-system.md:26-38`；对照 `…/issues/02:36` → 已同步 |
| `docs/skill-system.md` §"Description-driven dispatch"（第 22 行陈述） | 该行写「没有任何 front-matter 闸门」；实物 `assets/skills/flowforge-refine-ticket/SKILL.md:4` 现存 `disable-model-invocation: true`（该 skill 由 `docs/proposals/ticket-refinement-contract/` 引入，属本组外来源，此处仅列事实） | `docs/skill-system.md:22` vs `assets/skills/flowforge-refine-ticket/SKILL.md:4`、`git log --oneline -3 -- assets/skills/flowforge-refine-ticket/` → `f50dfcc feat: add ticket refinement contract, repair lifecycle, and refine-ticket skill` |
| `docs/skill-system.md` §"Subagent 委派与协作"（`docs/skill-system.md:53-66`） | 第 55 行写「主会话依据 `flowforge frontier` 的就绪状态与 AGENTS.md 路由表，将工作分配给对应的 Subagent」；与本 proposal 记录「AGENTS.md 路由表保留为人类参考、不作为 agent dispatch 路径」及同文档第 24 行「agent 不自动读取它们做兜底」并列存在 | `docs/skill-system.md:55`、`docs/skill-system.md:24`；`…/requirements.md:57`（"AGENTS.md 路由表作为人类参考"）、`…/design.md:45`（"不作为 agent dispatch 路径、也无 agent 自动入口"）→ 该节与本次决策表述冲突 |
| `docs/skill-system.md` §"Subagent 委派与协作"（第 55 行宿主清单） | 该行列 3 个宿主（Claude Code Subagents、OpenCode Agent Tool/`@mention`、Codex 子会话）；`docs/architecture.md:58` 记四宿主（含 PI） | `docs/skill-system.md:55` vs `docs/architecture.md:58` → 宿主清单过时（本组内可归属：`…/design.md:103` 明确要求"双宿主均走 description 自选"的说明落在 skill-system.md） |
| `docs/cli-design.md` §"目录解析"（`docs/cli-design.md:5-9`） | 该节只讲 `init` 创建与 `--force` 刷新，未记载 upgrade/init 会清理 source 已删除的 skill 目录 | `docs/cli-design.md:9`；对照实物 `internal/command/assets_deploy.go:38`（`cleanupRemovedSkillDirs(...)`）与 `…/issues/06-deploy-cleanup-removed-skills.md:132`（闭环证据）→ 该节缺该用户可见行为 |
| `docs/cli-design.md` §"其他命令"（`docs/cli-design.md:47-55`） | 未列 skill 目录清理语义 | 同上 |

### 2.5 候选用户场景
「用户说『检查 review 发现』，agent 不再误入设计或拆票，直接进 review 定 Fix 项。」

---

## 3. standards-injection

### 3.1 Feature 一句话
项目规范在分析阶段就进场：Align 按项目自带的「规范提取说明」挑出适用规范交给 Design，Design 把规范写成 `must`/`must not` 并决定 Constraints/Conventions 归属，Plan 只做机械转写，implementer 开工前若票里没有规范就退回。

### 3.2 表面变更清单

| 表面 | 事实 | 引用 |
|---|---|---|
| config 键（新增） | `standards.guide`（默认 `agents/standards.md`，相对 `docs_dir`），支持 `config get/set` | 设计 `docs/proposals/standards-injection/design.md:51`；实物 `internal/config/service.go:41,71,100` |
| 受管资产（新增） | `<docs_dir>/agents/standards.md`，与 `agents/domain.md`、`agents/issue-tracker.md` 同级，由 `init`/`upgrade` 部署内置默认版 | 设计 `design.md:47,49`；实物 `assets/agents/standards.md`、`docs/agents/standards.md`（均存在） |
| skill 职责（Align） | 读提取说明 → 识别适用规范 → 传给 Design；只识别、不转换 | 设计 `design.md:59-65`；落地 `docs/skill-system.md:11`；ticket `…/issues/09-align-extraction.md:15`(closed) |
| skill 职责（Solution Design） | 接收规范清单；合规性作为设计比较维度；把规范转成 `must`/`must not`（含 tier 与规范源链接）写入设计 authority（唯一权威来源） | 设计 `design.md:73-92`；落地 `docs/skill-system.md:12`；ticket `…/issues/10-design-conversion.md:15` |
| skill 职责（Plan） | 从设计 authority 机械转写 clauses 到 `## Constraints` 或 `### Conventions`；不读提取说明、不判 tier；缺 clauses 时标 `standards: pending` 退回 Design | 设计 `design.md:102-112`；落地 `docs/skill-system.md:14`；ticket `…/issues/11-plan-transcription.md:15` |
| skill 职责（Implement） | step 1 pre-flight：票内有 must/must not → 通过；`standards: pending` → 退回 Plan；`standards: none found` → 通过；有 Changes+Write set 却无标记 → 退回 Plan | 设计 `design.md:118-128`；落地 `docs/skill-system.md:15`；ticket `…/issues/05-implement-preflight.md:15` |
| skill 职责（Review 收窄） | Standards 轴只查票内已注入规范，不重跑提取、不查遗漏；仓库通用 coding-standards 查找保留 | 设计 `design.md:134-140`；落地 `docs/skill-system.md:16`；ticket `…/issues/06-review-standards-narrowed.md:15` |
| skill 职责（Setup） | `flowforge-setup` 新增 Section D「Standards」：检查 `agents/standards.md`、展示默认版摘要、引导特化（不强制） | 设计 `design.md:146-153`；ticket `…/issues/07-setup-section-d.md:15`(closed) |
| 共享契约 | `_shared/ARTIFACT-CONTRACT.md` 的 Standards clauses 格式改为「设计 authority 产物、Plan 转写」，新增设计 authority 中该段落格式 | 设计 `design.md:166`；落地 `docs/skill-system.md:77`；ticket `…/issues/12-artifact-contract-update.md:15` |
| ticket 契约 | 规范陈述书写格式（`must` / `must not` + 规范源语义链接 + tier 标注），不新增 tier、不记规范 revision | 需求 `docs/proposals/standards-injection/requirements.md:40,66-69` |
| CLI 命令 | 无新增/变更命令与 flag（提取逻辑明确不入 CLI） | 需求 `requirements.md:47`（"CLI 解析提取逻辑…在 Out of scope"）、`design.md:181`（`CLI standards match 子命令` 已否决） |
| agent 角色 | 无新增/删除 subagent 角色 | 设计 `design.md:155-168`（实现边界只列 Config / Managed asset / Skills / Tracker / 文档） |
| wiki 结构 | 新增 `<docs_dir>/agents/standards.md` 一个受管文件（目录结构不变） | 设计 `design.md:47` |

### 3.3 实施状态

**issues 全部 closed（13/13）；Config / 受管资产 / Setup / Implement / Review 部分在设计中被记为 v5.4.0 已实现。**

- 状态行逐票：`issues/01:15`、`02:15`、`03:15`、`04:15`、`05:15`、`06:15`、`07:15`、`08:16`、`09:15`、`10:15`、`11:15`、`12:15`、`13:17` 均为 `**Status:** closed`（`## Completion evidence` 节位于 98/99/100/102/101/102/103/87/66/70/68/67/67 行）。
- 设计侧已实现标记：`design.md:157`（Config `standards.guide` — v5.4.0 已实现）、`design.md:158`（`assets/agents/standards.md` 已存在）、`design.md:163-165`（Implement pre-flight / Review step 3 / Setup Section D 已实现）。
- 文档类票的闭环：`…/issues/08-doc-update.md:73-83`（Changes 1-5 completed + grep 计数 2/1/3）、`…/issues/13-doc-update-revised.md:69-71`（`assets verify` 全 current、`check --strict` green）。
- 当前二进制复核（见 §7）：`flowforge check --dir docs/proposals/standards-injection --strict` → **exit=1，59 条 `evidence-missing` warning**（已 checked 的 Change 缺 evidence 四元组，主要集中在票 08/13 及 skills 类票）。

### 3.4 文档影响（对照四份现行文档）

| 文档小节 | 现状事实 | 引用 |
|---|---|---|
| `docs/architecture.md` §"实现边界"（`docs/architecture.md:55-66`） | 第 57 行已列 `standards.guide` 配置；第 62 行已写 `assets/agents`（含 `standards.md`） | `docs/architecture.md:57,62`；对照 `…/issues/08-doc-update.md:36`（该节 Change 1）→ 已同步 |
| `docs/cli-design.md` §"其他命令"（`docs/cli-design.md:47-55`） | 第 53 行已列 `standards.guide`，但未给默认值 `agents/standards.md` | `docs/cli-design.md:53`；对照设计 `design.md:51`（默认 `agents/standards.md`，相对 `docs_dir`）→ 该行缺默认值事实 |
| `docs/cli-design.md` §"目录解析"（`docs/cli-design.md:5-9`） | 第 9 行 `init` 部署清单只写 `<docs_dir>/agents/`（目录粒度），未点名 `standards.md` | `docs/cli-design.md:9` vs `docs/architecture.md:62` |
| `docs/skill-system.md` §"主交付链"（`docs/skill-system.md:5-18`） | Align/Solution Design/Plan/Implement/Review 五行已写入规范识别/转换/转写/pre-flight/收窄职责 | `docs/skill-system.md:11,12,14,15,16`；对照 `…/issues/08-doc-update.md:36-40` 与 `…/issues/13-doc-update-revised.md:35-37` → 已同步 |
| `docs/skill-system.md` §"支持与特殊路径"（`docs/skill-system.md:40-51`） | 该节逐条列举支持类 skill，`grep -n "setup" docs/skill-system.md` → 0 命中；`flowforge-setup`（Section D 的载体）未在此列出现 | `docs/skill-system.md:40-51`；对照设计 `design.md:146`（Section D 落在 `flowforge-setup`）与 ticket `…/issues/07-setup-section-d.md:20` → 该节缺 setup 条目 |
| `README.md` §"安装与初始化"（`README.md:72-99`） | 第 92 行写 init 部署 `.agents/skills/` 与 `ff-wiki/agents/`，未点名 `standards.md` | `README.md:92` vs `docs/architecture.md:62` |

### 3.5 候选用户场景
「改分层依赖时，规范以前靠人肉抄进每张卡；现在 Align 挑规范、Design 写 must，缺规范的票被退回。」

---

## 4. wiki-config-single-track

### 4.1 Feature 一句话
wiki 根目录只由 `docs_dir` 一轨决定：旧的 `wiki.root` / `projects[].wikiRoot` 键在 `load` 时于 stderr 提示一行然后被忽略，不会再出现「配置里改了 wiki root 但目录纹丝不动」。

### 4.2 表面变更清单

| 表面 | 事实 | 引用 |
|---|---|---|
| config 键（删除） | 删除 `Config.Wiki` / `WikiConfig` / `WikiRoot()` / `WikiRootForProject()` / `projectWikiRoot()` / `primaryProject()` / `ProjectConfig.WikiRoot` 与 `v.SetDefault("wiki.root", …)` | 设计 `docs/proposals/wiki-config-single-track/design.md:32`；ticket `…/issues/01-remove-wiki-track.md:35-37` |
| config 键（删除） | `config list` 不再输出 `project.<id>.wikiRoot`；`config set project.<id>.wikiRoot` 从「可写」变为报错 `unknown project config field: wikiRoot`（裁决内破坏面） | 设计 `design.md:33,41`；实物 `internal/config/config_test.go:425`（`TestConfigSetProjectWikiRootRejected`） |
| config 键（用户可见告警） | `Load` 在 `v.ReadInConfig()` 成功后、unmarshal 前逐键向 stderr 输出：`warning: config key "wiki.root" is deprecated and ignored; wiki root is decided by "docs_dir" only` / `warning: config key "projects[%s].wikiRoot" is deprecated and ignored; use "docs_dir"` | 设计 `design.md:37-39`；ticket `…/issues/01-remove-wiki-track.md:38-40` |
| config 默认值 | `docs_dir` 默认 `ff-wiki` 收敛为单点常量 `DefaultDocsDir` | 设计 `design.md:45`；实物 `internal/config/config.go:17`（`DefaultDocsDir = "ff-wiki"`） |
| CLI 命令/flag | 不改 CLI 命令签名（需求红线）；用户可见变化仅限上述 config 键面与 load 告警 | 需求 `docs/proposals/wiki-config-single-track/requirements.md:32`、`design.md:59` |
| skill / agent 角色 | 无 skill、无 subagent 角色变更 | 需求 `requirements.md:32`（"只动配置解析、服务键空间、默认值常量与文档"）；设计 `design.md:26-53` 无 skills/agents 条目 |
| wiki 结构 | 结构不变；仅本仓自举 `.flowforge/config.yaml` 删除 `projects[].wikiRoot` 残留行（untracked 本地文件） | 设计 `design.md:47`；ticket `…/issues/02-docs-accuracy.md:52` |
| 文档修正 | 3 处 d04e955 前过时表述清零：`README.md:92`、`docs/proposals/pi-host-integration/issues/03-pi-extension.md:136`、`docs/proposals/generic-role-orchestration/issues/04-dogfood-batch-drill.md:64` | 需求 `requirements.md:19,28`；ticket `…/issues/02-docs-accuracy.md:20,37-51` |

### 4.3 实施状态

**issues 全部 closed（2/2）。**

- 状态行：`issues/01-remove-wiki-track.md:15`、`issues/02-docs-accuracy.md:15` 均为 `**Status:** closed`（Completion evidence 位于 `01:142`、`02:129`）。
- 代码面证据：`…/issues/01:48-56` 五条 Change 的 cmd/exit/output 四元组（如 `grep -rn '"ff-wiki"' internal/ --include='*.go' | grep -v _test.go` → 仅 `internal/config/config.go:17` 单点；`grep -rn 'WikiRoot' internal/… | grep -v _test.go | wc -l` → `0`）。
- 文档面证据：`…/issues/02:37-56` 四条 Change 的四元组（README `docs/` → `ff-wiki/`；`grep -c 'ff-wiki' README.md` → `2`）+ `…/issues/02:105-115` 修正前后文逐条对照。
- 当前二进制复核（见 §7）：`flowforge check --dir docs/proposals/wiki-config-single-track --strict` → **exit=0，0 warning**（"Checked 2 issues … Dependency graph is healthy"）。
- 需求验收 3「`grep -rn '"ff-wiki"'` 仅常量定义一处」在当前仓库成立：`internal/config/config.go:17`。

### 4.4 文档影响（对照四份现行文档）

| 文档小节 | 现状事实 | 引用 |
|---|---|---|
| `docs/cli-design.md` §"目录解析"（`docs/cli-design.md:5-9`） | 第 7 行仍写「`docs_dir` 指定文档根目录，**默认为 `docs`**」；现行默认值为 `ff-wiki` | `docs/cli-design.md:7` vs `docs/architecture.md:66`（"`docs_dir` 默认为 `ff-wiki`（d04e955 起）"）与实物 `internal/config/config.go:17` → **该节过时** |
| `docs/cli-design.md` §"目录解析"（同上） | 未记载 wiki 轨键的 load 弃用告警与忽略语义 | `docs/cli-design.md:5-9`；对照设计 `design.md:37-39` 与 ticket `…/issues/01:38-40` → 该节缺告警事实 |
| `docs/cli-design.md` §"其他命令"（`docs/cli-design.md:53`） | 该行未提示 wiki 轨键已失效/会告警（仅列 `docs_dir`、`standards.guide`、version check） | `docs/cli-design.md:53` vs `design.md:33,41` |
| `README.md` §"安装与初始化"（`README.md:72-99`） | 第 92 行已为 `ff-wiki/CONTEXT.md`、`ff-wiki/adr/`、`ff-wiki/proposals/`、`ff-wiki/agents/`；第 95 行示例 `config set docs_dir ff-wiki-v5` | `README.md:92,95`；对照 `…/issues/02:37-41,105-108` → 已同步 |
| `docs/architecture.md` §"实现边界"（`docs/architecture.md:66`） | 默认值与单轨语义（含无 config 时回落 `docs/proposals`）已写在该行 | `docs/architecture.md:66`；`…/requirements.md:16`（本仓活体分歧记录）→ 已同步 |
| `docs/skill-system.md` | 该 proposal 不涉及 skill 语义，本文件无对应小节改写记录 | `…/requirements.md:32`（范围排除 CLI 签名与 Issue Schema，未含 skill 文档） |

### 4.5 候选用户场景
「配置里同时有 docs_dir 和 wikiRoot，改后者目录没变；现在旧 wiki 键提示一行后被忽略。」

---

## 5. external-material-intake

### 5.1 Feature 一句话
可以拿本地 PRD、旧 proposal 或会议记录当起点：`flowforge-import` 把来源分成事实/需求候选/设计决定/证据/未知五类再交接给 Align 或 Solution Design（不做格式转换），同时可用 `flowforge assets verify` 只读核对项目受管 Skill 与二进制嵌入资产是否一致。

### 5.2 表面变更清单

| 表面 | 事实 | 引用 |
|---|---|---|
| 新增 skill | `flowforge-import`：输入为来源路径 + 目标 feature + 可选目标语言；把保留内容分五类并交给 Align / Solution Design | 设计 `docs/proposals/external-material-intake/design.md:19,21`；ticket `…/issues/01-import-skill-and-writing-contract.md:20,32`；落地文档 `docs/skill-system.md:10`、`README.md:32` |
| 新增 skill 产物角色 | `role: research` 的 `source-notes.md` 只在多来源被独立评审/存在实质冲突/需复查判读理由时创建 | 设计 `design.md:25`；需求 `requirements.md:32` |
| CLI 命令（新增） | `flowforge assets verify [project]`（含 `--json`）：只读比较，输出 `current` / `missing` / `drifted` / `project-owned`；missing/drifted 返回非零，project-owned 仅信息展示且不失败 | 设计 `design.md:49-58`；ticket `…/issues/05-assets-verify-command-and-sync-reporting.md:20,49`；实物 `internal/command/assets.go:17`（`Use: "verify [project]"`） |
| CLI 行为变更（init/upgrade 报告） | `init`/`upgrade` 先显式同步，再用同一比较结果决定是否报「已同步」；不一致时给出具体路径并指向 `flowforge assets verify` | 设计 `design.md:58`；ticket `…/issues/05:33,50`；落地文档 `docs/cli-design.md:75` |
| authority 发布自检 | 新建/修订带 schema metadata 的 requirement/design/research/spec 后运行 `flowforge check --dir <feature-dir> --strict`；只读诊断，不写 readiness、不阻塞其他 feature | 设计 `design.md:37-45`；ticket `…/issues/02-strict-authority-publication-check.md:20,48`；落地文档 `README.md:40`、`docs/skill-system.md:94` |
| skill 集成 | Import → Align（需求候选）→ 必要时 Solution Design → 每次发布后 strict check → Plan 展示 title/Delivery/真实 DAG 边并等用户接受才写 issue | 设计 `design.md:23`；需求 `requirements.md:29-30`；ticket `…/issues/04-integrate-authoring-flow.md:20,49-51`；落地文档 `README.md:32-34` |
| 写作契约 | Import/Align/Solution Design 共享「目标语言语义重写、不逐句翻译、代码标识与术语稳定」规则；信息价值规则写入共享契约 | 设计 `design.md:29-33`；ticket `…/issues/01:33,51`（共享 `ARTIFACT-CONTRACT.md#source-intake-and-semantic-rewrite`）；落地文档 `docs/skill-system.md:76`（信息价值条目） |
| agent 角色 | 无新增/删除 subagent 角色 | 设计 `design.md:62-67`（实现边界只列 Skill/共享 reference、CLI、测试、文档） |
| config 键 | 无新增 config 键 | 设计 `design.md:62-67`；`docs/cli-design.md:53` 的键面未含新键 |

### 5.3 实施状态

**issues 全部 closed（5/5）。**

- 状态行：`issues/01-import-skill-and-writing-contract.md:18`、`02-strict-authority-publication-check.md:18`、`03-managed-assets-comparison.md:18`、`04-integrate-authoring-flow.md:18`、`05-assets-verify-command-and-sync-reporting.md:18` 均为 `**Status:** closed`（Completion evidence 位于 48/46/46/47/47 行）。
- 实物复核：`assets/skills/flowforge-import/` 存在；`internal/command/assets.go:17` 有 `verify [project]` 子命令；`docs/cli-design.md:71-75` 有该命令的用户可见说明。
- 关闭证据要点：`…/issues/01:50-51`（Import skill + 共享 source-intake 契约 + 中文混合来源 walk-through）、`…/issues/05:49-51`（human/JSON 双投影、非零策略、init/upgrade 一致才报成功、CLI 文档；`go test ./internal/...` 与 `check --strict` 通过）。
- 当前二进制复核（见 §7）：`flowforge check --dir docs/proposals/external-material-intake --strict` → **exit=0，0 warning**。

### 5.4 文档影响（对照四份现行文档）

| 文档小节 | 现状事实 | 引用 |
|---|---|---|
| `README.md` §"1. 找到下一位责任人"（`README.md:26-34`） | 外部材料分支已写入（Import 分类五类、需求候选回 Align、模块/seam 交 Solution Design） | `README.md:32`；对照 ticket `…/issues/04:35,51` → 已同步 |
| `README.md` §"安装与初始化"（`README.md:72-99`） | 该节讲安装、`init`、`config set docs_dir`、`init --force`；`grep -n "assets verify" README.md` → 0 命中 | `README.md:86-99`；对照 `docs/cli-design.md:73,75`（该命令的完整说明）→ README 缺 `assets verify` 用户入口 |
| `README.md`（第 99 行同步语义） | 只写「`init --force` 同步受管资产，但保留已有项目配置」，未含「仅全部 current 才报同步成功」的新报告语义 | `README.md:99` vs `docs/cli-design.md:75` 与 `…/issues/05:33` |
| `docs/cli-design.md` §"Managed asset verification"（`docs/cli-design.md:71-75`） | 命令、四类结果、`--json`、非零语义、init/upgrade 复用同一比较均已写明 | `docs/cli-design.md:73,75`；对照 `…/issues/05:51`（Change 3 的文档交付）→ 已同步 |
| `docs/architecture.md` §"实现边界"（`docs/architecture.md:55-66`） | 第 60 行 `internal/command` 职责列举「Cobra 命令、策略投影、初始化、Subagent 生命周期管理（deploy/remove/status）和受管资产部署」，未列 `assets verify` 只读比较 | `docs/architecture.md:60` vs `internal/command/assets.go:17`、`docs/cli-design.md:73` → 该行缺 verify |
| `docs/skill-system.md` §"主交付链"（`docs/skill-system.md:5-18`） | 第 10 行已写 `flowforge-import` 的拥有场景与「不转换 authority」边界 | `docs/skill-system.md:10` → 已同步 |
| `docs/skill-system.md` §"继续、跳过与返回"（`docs/skill-system.md:81-94`） | 第 94 行已写 `check --dir <feature-dir> --strict` 发布自检及其非状态机语义 | `docs/skill-system.md:94`；对照设计 `design.md:43` → 已同步 |

### 5.5 候选用户场景
「旧 PRD 不知哪段还算数时，先 import 分类交接，再用 assets verify 查受管资产漂移。」

---

## 6. 组内横向事实（仅陈述，不综合）

- 五个 proposal 的 `issues/*.md` 状态行全部为 `**Status:** closed`：DCR 10、skill-routing 6、standards-injection 13、wiki-config 2、external-material-intake 5，合计 36 张。
- 当前二进制 `check --strict` 对五个目录的结果不一致（§7 原文）：`wiki-config-single-track` 与 `external-material-intake` exit=0 / 0 warning；`documentation-contract-refinement` exit=1 / 15 warning（10 `legacy-metadata` + 5 `missing-required-metadata`，全在 `scenario-fixtures/`）；`skill-routing-simplification` exit=1 / 48 `evidence-missing`；`standards-injection` exit=1 / 59 `evidence-missing`。
- 与本组 group 直接相关的现行文档小节被标为过时或缺项共 9 项（分属 `README.md`、`docs/cli-design.md`、`docs/skill-system.md`、`docs/architecture.md` 四份文档；引用见各节 1.4/2.4/3.4/4.4/5.4）：`docs/cli-design.md:7`（默认 `docs`）、`docs/cli-design.md:9`（缺 skill 目录清理）、`docs/cli-design.md:53`（缺 `standards.guide` 默认值）、`docs/skill-system.md:40-51`（缺 `flowforge-setup`）、`docs/skill-system.md:55`（按 AGENTS.md 路由表派发 + 3 宿主）、`docs/architecture.md:60`（缺 `assets verify`）、`README.md:72-99`（缺 `assets verify` 入口）与 `README.md:99`（缺 current-only 同步语义）、`README.md:106`（"本次"框架）、`README.md:57-66`（缺 promoted evidence）。

---

## 7. 复核命令与原始输出

```
$ ls -d assets/skills/flowforge-* | wc -l
28
$ ls -d .agents/skills/flowforge-* | wc -l
28
$ ls -d assets/skills/flowforge-route .agents/skills/flowforge-route
ls: .agents/skills/flowforge-route: No such file or directory
ls: assets/skills/flowforge-route: No such file or directory
$ grep -rn 'flowforge-route\|Route & Guide' AGENTS.md assets/AGENTS.md
(no output, exit=1)
$ grep -rn "assets verify" README.md docs/architecture.md docs/skill-system.md docs/cli-design.md
docs/cli-design.md:73:`flowforge assets verify [project]` compares the running binary's embedded Skills and agent rules with the selected project without changing files. ...
docs/cli-design.md:75:`flowforge init` and `flowforge upgrade` explicitly synchronize managed assets, then use this same comparison before reporting synchronization success. ...
$ grep -n "assets verify" README.md docs/architecture.md docs/skill-system.md
(no output, exit=1)
$ grep -n "setup" docs/skill-system.md
(no output, exit=1)
$ grep -rn "standards.guide" internal/config/service.go
internal/config/service.go:41:	case key == "standards.guide":
internal/config/service.go:71:	case key == "standards.guide":
internal/config/service.go:100:	result["standards.guide"] = guide
$ grep -rn "WikiRoot" internal/ --include='*.go' | grep -v _test.go
(no output, exit=1)
$ grep -rn '"ff-wiki"' internal/ --include='*.go' | grep -v _test.go
internal/config/config.go:17:	DefaultDocsDir        = "ff-wiki"
$ grep -rn "cleanupRemovedSkillDirs" internal/command/assets_deploy.go
internal/command/assets_deploy.go:38:	if err := cleanupRemovedSkillDirs(filepath.Join(assetsDir, "skills"), filepath.Join(targetDir, ".agents", "skills")); err != nil {
$ sed -n '4p' assets/skills/flowforge-refine-ticket/SKILL.md
disable-model-invocation: true
```

```
$ for d in documentation-contract-refinement skill-routing-simplification standards-injection wiki-config-single-track external-material-intake; do
    ./bin/flowforge check --dir docs/proposals/$d --strict ; done
documentation-contract-refinement: exit=1 warnings=15
skill-routing-simplification: exit=1 warnings=48
standards-injection: exit=1 warnings=59
wiki-config-single-track: exit=0 warnings=0
external-material-intake: exit=0 warnings=0
```

warning 分类（`sed -E 's/^warning: ([a-z-]+).*/\1/' | sort | uniq -c`）：

```
documentation-contract-refinement: 10 legacy-metadata / 5 missing-required-metadata
skill-routing-simplification:      48 evidence-missing
standards-injection:               59 evidence-missing
```
