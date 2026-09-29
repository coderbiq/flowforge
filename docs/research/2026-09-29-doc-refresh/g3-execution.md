# G3 执行/执行者契约组 — 面向用户的 feature 事实提取

- 范围：`docs/proposals/` 下本组 5 个 proposal：`executor-loop-hardening`、`executor-value-measurement`、`fast-executor-reliability`、`lightweight-execution-contract`、`ticket-refinement-contract`。
- 方法：逐文件精读 requirements/spec/design 与 `issues/`（Status、Changes、Completion evidence），再对照当前 `README.md`、`docs/architecture.md`、`docs/skill-system.md`、`docs/cli-design.md` 的小节定位过时/缺失点。
- 引用约定：`path:line` 为本次读取时的行号；命令输出逐字引用。所有"已交付"断言均给出源文件行号或实测命令输出。
- 边界：本文件只陈述事实，不做跨组综合、不提改进建议。

---

## 1. executor-loop-hardening（执行者循环硬化）

### 1.1 Feature 一句话

配置 `agents.max_steps` 后，flowforge-implementer 的 OpenCode 部署产物携带宿主 `steps` 迭代预算（默认 200，`-1` 关闭）；同时实现者固化 prompt 常驻携带 fail-fast 防循环摘要，`flowforge check` 会对"同一命令非零退出出现 ≥3 次"报 `evidence-repeat-failure`。

### 1.2 表面变更清单

CLI / 诊断
- 新诊断码 `evidence-repeat-failure`（warning；`--strict` 判失败）— `docs/proposals/executor-loop-hardening/issues/03-repeat-failure-diagnostic.md:20`（Delivery）；常量落地 `internal/tracker/catalog.go:58`。
- `flowforge check --help` 诊断清单新增第 6 行"Repeated failure diagnostics (same command reporting non-zero exit 3+ times in one ticket)" — `internal/command/check.go:29`（实测 `./bin/flowforge check --help` 输出逐字一致）。
- 阈值 3 为不可配置常量；复用既有豁免通道 `evidence.exempt_proposals`，不新增配置面 — `docs/proposals/executor-loop-hardening/issues/03-repeat-failure-diagnostic.md:60-61`。

config 键
- `agents.max_steps`（`yaml:"max_steps,omitempty"`）：`0`/未配置 → 默认 200；`-1` → 不产出字段（关闭）；`>0` → 原值；其余负数报配置错误 — `docs/proposals/executor-loop-hardening/design.md:23`；字段落地 `internal/config/config.go:43`。
- 仅作用于 `flowforge-implementer`；其余 subagent 产物零变化 — `docs/proposals/executor-loop-hardening/design.md:26`。

部署产物
- OpenCode implementer 产物 frontmatter 增 `steps: <N>`（仅 OpenCode；Claude Code/Codex 无对应能力，不动）— `docs/proposals/executor-loop-hardening/design.md:25`、`docs/proposals/executor-loop-hardening/issues/01-steps-budget-channel.md:20`。
- 未配置时 implementer 产物差异仅一行 `steps: 200`；五个非 implementer 产物 `grep -c "steps:"` 均为 0 — `docs/proposals/executor-loop-hardening/issues/01-steps-budget-channel.md:116`。

agent / skill 资产
- `assets/subagents/flowforge-implementer.md` body 在 Identity 与 Boundaries 之间新增英文 `## Non-negotiables` 小节：fail-fast（`at most 2 times`、第二次失败即 `STATUS: BLOCKED`、绝不做第三次相同尝试）、修复上限（`5 failed repair rounds`）、预算收束 — `docs/proposals/executor-loop-hardening/issues/02-pinned-prompt-contracts.md:20`；落地 `assets/subagents/flowforge-implementer.md:18-29`。
- 结构断言锁定 `assets/skills/flowforge-implement/SKILL.md`、定义源、部署产物三处关键词一致 — `docs/proposals/executor-loop-hardening/issues/02-pinned-prompt-contracts.md:47`；测试 `internal/command/assets_deploy_test.go:368`。

wiki / 追踪器结构：无变更（`docs/proposals/executor-loop-hardening/requirements.md:34` 明示不修改 Issue Schema 头规范）。

### 1.3 实施状态

**issues 全部 closed。**
- `docs/proposals/executor-loop-hardening/issues/01-steps-budget-channel.md:15` `**Status:** closed`
- `docs/proposals/executor-loop-hardening/issues/02-pinned-prompt-contracts.md:15` `**Status:** closed`
- `docs/proposals/executor-loop-hardening/issues/03-repeat-failure-diagnostic.md:16` `**Status:** closed`

完成证据（本次实测复验）：
- `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command/ ./internal/tracker/ -run '^TestOpenCodeStepsBudget$'` → `ok flowforge/internal/command 0.461s`
- `... -run '^TestEvidenceRepeatFailureThreshold$'` → `ok flowforge/internal/tracker 0.151s`
- `... -run '^TestImplementerPromptPinsLoopContracts$'` → `ok flowforge/internal/command 0.329s`
- `./bin/flowforge check --dir docs/proposals/executor-loop-hardening --strict` → exit 0（healthy）

### 1.4 文档影响（对照现有小节）

| 文档小节（文件:行） | 事实 |
|---|---|
| `docs/cli-design.md:23` `## flowforge check` 节 | 仅写"检查循环依赖、悬空依赖、自依赖和 Catalog 诊断"与 `--strict` 语义；未列 `evidence-repeat-failure`（`docs/cli-design.md:29` 全行无该词） |
| `docs/cli-design.md:13` `## Artifact Catalog` 的"确定性诊断覆盖"列表（`docs/cli-design.md:15-21`） | 六条枚举（role/位置、identity/authority、语义链接、open item/waiver、完成证据）中无任何 evidence 诊断族，也无 `evidence-repeat-failure` |
| `docs/cli-design.md:47` `## 其他命令` / `README.md:101` `## 核心文档` | 全文无 `agents.max_steps`、`steps:`、模型/预算配置的任何说明（`grep -n "max_steps" docs/cli-design.md docs/architecture.md docs/skill-system.md README.md` 无匹配） |
| `docs/architecture.md:62` `assets/subagents` 一行 | 只描述资产目录用途，未提 implementer 的预算/常驻契约固化层 |
| `docs/skill-system.md:53` `## Subagent 委派与协作` 表（`docs/skill-system.md:57-64`） | `flowforge-implementer` 行（`docs/skill-system.md:62`）只写"TDD 实现与轻量自检 / 读写受限于 Write set"，未提 OpenCode `steps` 预算、Non-negotiables、模型钉定或测试文件保护 |
| `README.md:20` `## 一个需求如何实际推进` 的 5 步叙述（`README.md:20-70`） | 无执行预算/循环熔断一节；`README.md:55` 只写 `check` 校验 DAG/角色/revision/链接/waiver/完成证据 |

### 1.5 候选用户场景

批量派发弱模型执行时，一张票反复重试同一条命令近三百次、烧掉整晚额度——配置 `agents.max_steps` 后预算到限即收束，`check` 还会在 ≥3 次重复失败时直接点名那条命令。

---

## 2. executor-value-measurement（执行者价值度量）

### 2.1 Feature 一句话

提供一条零依赖 Python 命令，从本机 OpenCode 会话库只读提取 `flowforge-implementer` 会话的成本/时长/循环指标并按部署纪元与票难度分档聚合，再按预先登记的 G1–G4 规则给出"保持 flash / 回退旗舰"的判定。

### 2.2 表面变更清单

- 新脚本 `scripts/executor_metrics.py`（python3 标准库、零第三方依赖、`file:...?mode=ro` 只读打开 DB）— `docs/proposals/executor-value-measurement/requirements.md:32`；文件在位（`scripts/executor_metrics.py`，52469 bytes）。
- 新测试 `scripts/executor_metrics_test.py`（unittest）— `docs/proposals/executor-value-measurement/issues/01-metrics-extractor.md:31`；实测 `python3 -m unittest discover -s scripts -p 'executor_metrics_test.py'` → `Ran 108 tests ... OK`。
- 两个子命令与参数：`extract --db <path> --project <dir> --out <obs.md>`（幂等键 session id、`--since`、`--epochs`、`--price-override`）— `docs/proposals/executor-value-measurement/issues/01-metrics-extractor.md:20,84`；`report --obs <obs.md> --epochs "name:<T,name:T1-T2,name:>T2" --project-root <dir>` — `docs/proposals/executor-value-measurement/issues/02-epoch-strata-report.md:84`。
- 新产物文件（proposal 目录内数据契约）：`observations.md`（追加式，17 列表行）、`report.md`（当日替换制）、`DECISION.md`（收口模板）— `docs/proposals/executor-value-measurement/design.md:112-120`；模板落地 `docs/proposals/executor-value-measurement/DECISION.md:16-21`。
- 预登记决策门 G1 循环安全 / G2 经济性（<0.7× 且 n≥5）/ G3 时长（≤1.5×，仅记录）/ G4 复发熔断（`rep_max≥20 或 steps≥400` 立即回退）— `docs/proposals/executor-value-measurement/design.md:91-96`；观察期 2026-09-14 → 2026-09-20 — `docs/proposals/executor-value-measurement/design.md:89`。
- 价目表为脚本内档位代理常量（flash 0.30/2.50/0.075；旗舰 0.60/2.20/0.113），`--price-override` 可覆盖 — `docs/proposals/executor-value-measurement/requirements.md:34`。
- CLI 无变更：明确"不新增 flowforge CLI 子命令"、"不修改 flowforge 产品代码与 CLI 接口" — `docs/proposals/executor-value-measurement/requirements.md:33,52`。
- 运行手册：每日一次 extract + report（约 2 分钟），无守护进程/定时器 — `docs/proposals/executor-value-measurement/design.md:100-110`、`requirements.md:36`。

### 2.3 实施状态

**issues 全部 closed。**
- `docs/proposals/executor-value-measurement/issues/01-metrics-extractor.md:15` `**Status:** closed`
- `docs/proposals/executor-value-measurement/issues/02-epoch-strata-report.md:15` `**Status:** closed`
- `docs/proposals/executor-value-measurement/issues/03-decision-gates-playbook.md:15` `**Status:** closed`

完成证据（proposal 内自成文件，非推理）：
- `docs/proposals/executor-value-measurement/observations.md:8` 表头；`:10-14` 逐会话指标行（含事故会话 `ses_f66f341e5ffe5i4UyyAW1jrIoo`：`dur_min 86.4 / steps 895 / tools 893 / bash_n 509 / rep_max 184 / status COMPLETED`）。
- 观察期已完成收口：`docs/proposals/executor-value-measurement/report.md:8` `generated: 2026-09-18T17:35:41+08:00`，`:9` `observations: ... (72 sessions)`，且 `:12` 窗口行显示 `2 day(s) remaining`，`:14` 有 `> ALARM: G1 loop safety FAILED (hardened-deepseek-open)`。
- 收敛记录：`docs/proposals/executor-value-measurement/DECISION.md` 已含 2026-09-18 的收敛结论文本 — 收敛结论同时记录于一份既有 research 笔记（内部项目调度模式分析，已脱敏引用）。
- 本次实测复验：`python3 -m unittest discover -s scripts -p 'executor_metrics_test.py'` → `Ran 108 tests in 0.066s / OK`。
- `./bin/flowforge check --dir docs/proposals/executor-value-measurement --strict` → exit 0。

（仅陈述状态事实：本 proposal 不属于 FlowForge 产品 CLI 面，属分析工具与数据产物。）

### 2.4 文档影响（对照现有小节）

| 文档小节（文件:行） | 事实 |
|---|---|
| `README.md:101` `## 核心文档`（`README.md:103-106` 四条链接） | 只链接 architecture / skill-system / cli-design / documentation-contract-refinement；无"如何度量执行者收益/成本"的入口，也无 `scripts/` 工具链说明 |
| `docs/skill-system.md` 全文 | 无度量/纪要类 skill 或脚本的条目（文件内 94 行，无 `executor_metrics`/`metric` 匹配） |
| `docs/cli-design.md` 全文 | 与本 proposal 无对应小节（该 proposal 显式不新增 CLI 命令，`requirements.md:33`），因此 cli-design 无过时点、也无缺失项需要承载 |
| `docs/architecture.md:55` `## 实现边界`（`docs/architecture.md:57-62`） | 逐包列出 `internal/*` 与 `assets/*`；未列 `scripts/` 下新增的分析工具，也无"度量产物不进 CLI"的边界声明 |

### 2.5 候选用户场景

换用 flash 执行一个月后，团队说不清到底省了钱还是只把问题推迟——`extract` 每日追加会话指标、`report` 按纪元与难度分档出中位数成本与失控计数，用预登记的四道门给出保留或回退的结论。

---

## 3. fast-executor-reliability（快/弱执行者可靠交付）

### 3.1 Feature 一句话

勾选 `- [x]` 的 Change 必须携带命令锚定的证据四元组（`cmd`/`exit`/`output`/`artifact`），`flowforge check` 机械核对并在 `--strict` 下落成失败；同时弱执行者获得收窄行为契约、Review 增加廉价 Round 0 差距审计、部署物可钉定执行者模型并默认保护测试文件。

### 3.2 表面变更清单

CLI / 诊断
- 四个新诊断码（全 warning，`--strict` 判失败）：`evidence-missing`（勾选但无四元组）、`evidence-incomplete`（缺任一键）、`evidence-exit-nonzero`（`exit` ≠ 0）、`evidence-artifact-missing`（artifact 路径相对仓库根不存在）— `docs/proposals/fast-executor-reliability/design.md:55-60`。
- 解析范围限 `role: ticket`；纯文档票（Constraints 无 `Write set:`）跳过校验 — `docs/proposals/fast-executor-reliability/design.md:62`。
- CLI 签名不变（复用既有 `--strict` 与诊断管线）— `docs/proposals/fast-executor-reliability/design.md:65`。
- 实测 `./bin/flowforge check --help` 第 5 行：`5. Checked-change evidence quadruple diagnostics (missing/incomplete/non-zero exit/absent artifact)`（`internal/command/check.go` help 文案）。

config 键
- `evidence.exempt_proposals: [<dir>]`：命中则该 proposal 下全部 evidence 诊断不产生 — `docs/proposals/fast-executor-reliability/design.md:64`；字段落地 `internal/config/config.go:29,37`。
- `agents.models.tool-capable` / `agents.models.tool-capable-read-only`：配置后 OpenCode 产物携带显式 `model` 字段，未配置回退继承主会话 — `docs/proposals/fast-executor-reliability/issues/05-agent-model-pinning-testfile-guard.md:19`；字段落地 `internal/config/config.go:44`。
- `agents.test_file_globs`（覆盖默认测试 glob 集）与 `agents.disable_test_guard`（关闭保护）— `docs/proposals/fast-executor-reliability/issues/05-agent-model-pinning-testfile-guard.md:38,48`；字段落地 `internal/config/config.go:48-49`。
- 未知 models 键报错，不静默忽略 — `docs/proposals/fast-executor-reliability/issues/05-agent-model-pinning-testfile-guard.md:38`。

部署产物
- OpenCode implementer 产物默认含 `permission.edit` deny 规则，默认 glob 集 `**/*_test.go`、`**/src/test/**`、`**/src/integrationTest/**`、`**/__tests__/**`、`**/*.test.ts`、`**/*.test.tsx`、`**/*.spec.ts` — `docs/proposals/fast-executor-reliability/issues/05-agent-model-pinning-testfile-guard.md:48`。
- 宿主 enforcement 本轮仅 OpenCode；Claude Code / Codex 只出文档示例 — `docs/proposals/fast-executor-reliability/issues/05-agent-model-pinning-testfile-guard.md:58`。

agent / skill 角色
- `assets/AGENTS.md` 新增 `## Execution unit policy` 小节（一票一上下文、入口必读工件清单、测试分离含 opencode `permission` 与 Claude Code `disallowedTools` 示例）— `docs/proposals/fast-executor-reliability/issues/04-agents-md-execution-unit-policy.md:19,34,39`；小节在 `assets/AGENTS.md:108-127`。
- `flowforge-implement` SKILL.md lightweight mode 六项弱执行者行为契约（editor 收窄、两段式执行、失败即停 2 次/5 轮、错误逐字回喂、勾选与 `exit: 0` 同次编辑、Constraint + 全部 Execution scenario 预检遍历）— `docs/proposals/fast-executor-reliability/design.md:77-82`、`docs/proposals/fast-executor-reliability/issues/02-implement-weak-executor-contract.md:19`；落地 `assets/skills/flowforge-implement/SKILL.md:35-76`（含 `Phase 0`/`Phase 0b`/`3a`/`3b`/`3c`）。
- 模式选择从执行者自评改为派发方/票面声明 `**Mode:** lightweight` — `docs/proposals/fast-executor-reliability/design.md:84`；读端 `assets/skills/flowforge-implement/SKILL.md:29-31`，生产端 `assets/skills/flowforge-plan/SKILL.md:74,135`。
- `flowforge-review` SKILL.md 新增串行前置 `Round 0: converge audit`（单代理 tool-capable-read-only；意图源 = ticket + linked authorities；gap 分类 `missing`/`partial`/`contradicts`/`unrequested` 且每条带 `file:line`；零 gap 才升双轴）— `docs/proposals/fast-executor-reliability/issues/03-review-round0-audit.md:19,34-49`；落地 `assets/skills/flowforge-review/SKILL.md:66,76,80`。

### 3.3 实施状态

**issues 全部 closed（5/5）。**
- `docs/proposals/fast-executor-reliability/issues/01-evidence-gate-check-diagnostics.md:15` `**Status:** closed`
- `docs/proposals/fast-executor-reliability/issues/02-implement-weak-executor-contract.md:15` `**Status:** closed`
- `docs/proposals/fast-executor-reliability/issues/03-review-round0-audit.md:15` `**Status:** closed`
- `docs/proposals/fast-executor-reliability/issues/04-agents-md-execution-unit-policy.md:15` `**Status:** closed`
- `docs/proposals/fast-executor-reliability/issues/05-agent-model-pinning-testfile-guard.md:15` `**Status:** closed`

完成证据（本次实测复验）：
- `... -run '^TestPackagedSkillPointersResolve$'` → `ok flowforge/internal/command 0.267s`
- `... -run '^TestAgentsBlockContainsExecutionUnitPolicy$'` → `ok flowforge/internal/command`（`internal/command/assets_deploy_test.go:534`）
- `... -run '^TestReviewSkillCarriesRound0Audit$'` → `ok flowforge/internal/command`（`internal/command/assets_deploy_test.go:568`）
- `... -run '^TestPlanSkillDeclaresModeLine$'` → `ok flowforge/internal/command`（`internal/command/assets_deploy_test.go:554`）
- `./bin/flowforge check --dir docs/proposals/fast-executor-reliability --strict` → exit 0
- 本仓库实证副本：`.flowforge/config.yaml` 已使用 `agents.models` / `models_by_name`（`agents.models.tool-capable-read-only: cpa/deepseek-v4.1-flash`、`flowforge-implementer: cpa/deepseek-v4.1-flash`），与 feature 5 的实现面一致。

### 3.4 文档影响（对照现有小节）

| 文档小节（文件:行） | 事实 |
|---|---|
| `docs/cli-design.md:11` `## Artifact Catalog` / `docs/cli-design.md:15-21` 确定性诊断列表 | 列表只覆盖 role/位置、identity/authority、语义链接、open item/waiver、closed ticket 完成证据；**完全没有 evidence 四元组诊断族**（`docs/cli-design.md` 全文无 `evidence-missing`/`四元组` 匹配） |
| `docs/cli-design.md:23` `## flowforge check`（`docs/cli-design.md:29`） | 只写"默认 warning 与 gap 保持可见但不让 check 失败"；未描述勾选证据四元组的四条新诊断与其 warning/`--strict` 语义 |
| `docs/cli-design.md:47` `## 其他命令` / `docs/cli-design.md:65` `## 稳定边界` | 无 `evidence.exempt_proposals`、`agents.test_file_globs`、`agents.disable_test_guard`、`agents.models(.tool-capable)` 等新配置键的任何条目（`grep -n "models_by_name\|max_steps" docs/cli-design.md` 无匹配；`docs/cli-design.md:60` 仅提到 `agents.test_file_globs` 在 PI 扩展段落内被顺带引用） |
| `docs/skill-system.md:15` `flowforge-implement` 行 | 写"轻量模式：执行 unchecked Changes、规范 pre-flight 检查、机检自检、写 Implementation note 后停止"——未提六项行为契约（fail-fast 2 次/5 轮、错误逐字回喂、勾选与 `exit: 0` 同次编辑、`**Mode:**` 由派发方声明） |
| `docs/skill-system.md:16` `flowforge-review` 行 | 只写双轴与 fix 翻译；**没有 Round 0 审计层**（零 gap 才升双轴、gap 四分类）的表述 |
| `docs/skill-system.md:57-64` Subagent 委派表 | `flowforge-implementer` 行未提模型钉定；无"测试文件默认只读"；无 `permission.edit` deny 的宿主表达 |
| `docs/skill-system.md:81` `## 继续、跳过与返回` | 无"批量执行用一票一上下文、跨 ticket 状态只经工件"的执行单元策略条目 |
| `docs/architecture.md:64`（Subagent 名册段） | 未提 `agents.models` 钉定优先级链或测试文件保护 |

### 3.5 候选用户场景

让 flash 模型跑批量 ticket，最怕它"勾了 `[x]` 却什么都没做"——现在勾选必须同时写 `cmd/exit/output/artifact`，`check --strict` 当场拦下；测试文件默认只读，执行者也不再偷偷改断言混过验收。

---

## 4. lightweight-execution-contract（轻量执行契约）

### 4.1 Feature 一句话

Ticket 变成三层信息结构（人类优先 / 共享执行契约 / agent 细节），并引入"轻量模式"：flash 执行者按顺序机械执行未勾选 Changes、自检后写 Implementation note 即停；Review 在双轴之后把发现翻译成 `Fix:` Change 追加回票，零发现才写 Completion evidence 并关闭。

### 4.2 表面变更清单

skill 资产
- `assets/skills/_shared/ARTIFACT-CONTRACT.md` 新增三个 ticket 小节定义（`Execution detail`、`Implementation note`、`Review rounds`）、Constraints 内的 `Write set:` 行、Changes 的 `- [ ]` 勾选格式 — `docs/proposals/lightweight-execution-contract/issues/01-artifact-contract-execution-density.md:15,29-34`；小节顺序定义 `:57`。
- `flowforge-plan`：三层信息模型 + `---` 分隔的 `## Execution detail` 小节 + 完整 ticket 模板；Touch points 必须给文件路径与符号，Done and verify 必须给预期结果 — `docs/proposals/lightweight-execution-contract/issues/02-plan-three-tier-tickets.md:15,29-35`；落地 `assets/skills/flowforge-plan/SKILL.md:76-135`。
- `flowforge-implement`：新增轻量模式（Determine execution mode → 执行未勾选 Changes → 机检自检 → 勾选 → 写 Implementation note → 停止；不 review、不关票、不做设计决策），保留 full mode — `docs/proposals/lightweight-execution-contract/issues/03-implement-lightweight-mode.md:15,29-35`；落地 `assets/skills/flowforge-implement/SKILL.md:27-88`。
- `flowforge-review`：新增 fix-planning 步骤（可修复 finding 写 `Fix:` Change 续编，架构/seam/接口类 finding 记 design return），零发现时写 Completion evidence 并 `closed` — `docs/proposals/lightweight-execution-contract/issues/04-review-fix-planning.md:15,29-35`；落地 `assets/skills/flowforge-review/SKILL.md`。

wiki / tracker 结构
- `## Execution detail`（`---` 分隔，含 Settled decisions / Expected tests / Conventions）、`## Implementation note`（非 evidence，是给 review 的结构化状态报告）、`## Review rounds`（累积每轮固定点、发现、fix Change、处置）— `docs/proposals/lightweight-execution-contract/spec.md:23,31,33`。
- `assets/agents/issue-tracker.md` 增 "Execution and review loop" 子节与 closed-ticket 定义强化（最近一轮零发现）— `docs/proposals/lightweight-execution-contract/issues/05-document-loop-in-tracker-and-skills.md:32-33`；落地 `assets/agents/issue-tracker.md:24`（`## Execution and review loop`）。
- `docs/skill-system.md` 主交付链表与"继续、跳过与返回"节增补轻量模式/fix planning/review loop — `docs/proposals/lightweight-execution-contract/issues/05-document-loop-in-tracker-and-skills.md:34-35`。
- Fix: Change 协议、lightweight/full 模式边界、三层信息模型保留不动 — `docs/proposals/fast-executor-reliability/requirements.md:36`（该需求声明只修订本 spec 的两条 out-of-scope 决定）。

CLI / config：无新增（`docs/proposals/lightweight-execution-contract/spec.md:47-49`：不做 CLI 解析新段落、不做 write-set 自动强制、不改 `internal/`）。

### 4.3 实施状态

**issues 全部 closed（5/5）。**
- `docs/proposals/lightweight-execution-contract/issues/01-artifact-contract-execution-density.md:11` `**Status:** closed`
- `docs/proposals/lightweight-execution-contract/issues/02-plan-three-tier-tickets.md:11` `**Status:** closed`
- `docs/proposals/lightweight-execution-contract/issues/03-implement-lightweight-mode.md:11` `**Status:** closed`
- `docs/proposals/lightweight-execution-contract/issues/04-review-fix-planning.md:11` `**Status:** closed`
- `docs/proposals/lightweight-execution-contract/issues/05-document-loop-in-tracker-and-skills.md:11` `**Status:** closed`

完成证据：
- 票面 Done and verify 为结构断言（`docs/proposals/lightweight-execution-contract/issues/03-implement-lightweight-mode.md:45-50`）。
- 复验事实：`grep -n "3a\|3b\|3c\|Execution detail\|Write set" assets/skills/flowforge-implement/SKILL.md` 命中 `assets/skills/flowforge-implement/SKILL.md:35-85`；`assets/skills/_shared/ARTIFACT-CONTRACT.md:87`（`## Execution and review loop`）；`assets/agents/issue-tracker.md:24`（`## Execution and review loop`）；`docs/skill-system.md:15` 已含"轻量模式"措辞。
- **过时证据缺口（事实陈述）**：5 张票的勾选 Change **均无证据四元组**，`./bin/flowforge check --dir docs/proposals/lightweight-execution-contract --strict` 报 20+ 条 `warning: evidence-missing ...` 并以 `Error: proposal diagnostics violate selected policy` 退出（该 proposal 关闭于 evidence 门禁交付之前，属定时顺序事实，非本组判断对象）。

### 4.4 文档影响（对照现有小节）

| 文档小节（文件:行） | 事实 |
|---|---|
| `README.md:44` `### 3. 发布可执行 ticket`（`README.md:46`） | 只写"必须包含一项可观察交付、设计上下文、稳定触点、有序变更、容易违反的约束、成对的完成条件和验证命令"；**无三层信息模型、无 `Execution detail`/`Implementation note`/`Review rounds`、无 `Write set:` 行、无 `- [ ]` 机械步骤要求** |
| `README.md:57` `### 4. 实现、审查并留下证据`（`README.md:59-66`） | 只描述单会话模式（`flowforge-implement` 用 TDD 交付、`flowforge-review` 双轴、Implement 处理后关闭）；**无轻量模式、无 `Fix:` Change 回填、无零发现收敛条件** |
| `docs/skill-system.md:15` `flowforge-implement` 行 | 已含"轻量模式：…写 Implementation note 后停止"，但未含 fix 回填后的再执行闭环 |
| `docs/skill-system.md:81-90` `## 继续、跳过与返回` | `docs/skill-system.md:89` 只写"翻译为 `Fix:` Change 追加到 ticket，轻量 implementer 重新执行"，**无零发现收敛与 closed 判定的成对表述**（`docs/skill-system.md:90` 提到"零 findings 收敛"） |
| `docs/cli-design.md` 全文 | 本 proposal 不改 CLI（`spec.md:47-49`），cli-design 无对应过时点 |
| `docs/architecture.md:43` `## 自适应工件，而非固定阶段` | 未提三层 ticket 信息模型或轻量/完整执行模式的区分 |

### 4.5 候选用户场景

弱模型执行时总爱"自由发挥"——现在 Plan 把票写成机械步骤加预期结果的契约，flash 只按顺序执行、自检后写下 Implementation note 就停，review 再把发现翻译成 `Fix:` Change 让它重新跑，直到零发现才关票。

---

## 5. ticket-refinement-contract（Ticket 细化契约与 Repair 票）

### 5.1 Feature 一句话

新增 `flowforge-refine-ticket` 与实施前 preflight：缺五段机器执行契约的 ticket 报 `execution-contract-incomplete` gap 并被 frontier 默认排除、Implement 无条件拒绝；Review 的实质发现改为创建可追溯的 repair ticket，原票进入 `needs-repair`（非终态、不可执行），下游改依赖 repair。

### 5.2 表面变更清单

CLI / 诊断
- 新诊断码 `execution-contract-incomplete`（scoped gap；默认从 frontier 排除，`--include-gaps` 保留可见）— `docs/proposals/ticket-refinement-contract/issues/01-derive-execution-contract-gap.md:22,37`；常量落地 `internal/tracker/catalog.go:52`。
- catalog 对 open/ready-for-agent ticket 结构化检查五段（`Verified contracts`/`Execution scenarios`/`Expected tests`/`Generated artifacts`/`Conventions`），不得只靠 `Execution detail` 标题通过；空任务项（如 `- [ ]`）视为占位内容 — `docs/proposals/ticket-refinement-contract/design.md:33`、`issues/01-derive-execution-contract-gap.md:40`。
- 新状态与字段：`**Status:** needs-repair`（非终态、非可执行；`ComputeFrontier` 只派发 executable 状态）与 `**Repair of:** <id>`（provenance，不是 DAG edge）— `docs/proposals/ticket-refinement-contract/issues/03-add-repair-ticket-lifecycle.md:22,37-38`；落地 `internal/tracker/catalog.go:163`（reciprocal 校验）与 tracker 测试 `internal/tracker/catalog_test.go:329-332`。
- 不引入持久化 `ready`/`execution-ready` 状态 — `docs/proposals/ticket-refinement-contract/issues/01-derive-execution-contract-gap.md:45`。
- 不新增配置键、不改 CLI 签名；`frontier --include-gaps` 语义沿用 — `docs/proposals/ticket-refinement-contract/design.md:48`、`issues/02-add-refine-ticket-and-implement-preflight.md:48`。

skill / agent 资产
- 新 skill `flowforge-refine-ticket`（`assets/skills/flowforge-refine-ticket/SKILL.md`）：只细化一张无 DAG blocker 的候选票，读取仓库证据填五段，未知/冲突事实返回 Explore / Solution Design，不猜测 — `docs/proposals/ticket-refinement-contract/requirements.md:25,41`、`docs/proposals/ticket-refinement-contract/design.md:19`；SKILL 六步流程 `assets/skills/flowforge-refine-ticket/SKILL.md:13-56`。
- `flowforge-plan` 只发布带标题的 `## Execution detail` 骨架，事实内容交给 refine-ticket — `docs/proposals/ticket-refinement-contract/issues/02-add-refine-ticket-and-implement-preflight.md:39`；落地 `assets/skills/flowforge-plan/SKILL.md:55`。
- `flowforge-implement` 在任一模式前执行 readiness preflight，发现 execution-contract gap 必须停止，`--include-gaps` 不构成绕过 — `docs/proposals/ticket-refinement-contract/design.md:21`、`issues/02-add-refine-ticket-and-implement-preflight.md:40`。
- `flowforge-review` 处置规则改为三条：局部机械 finding 仍走 `Fix:` 循环；High/Critical、契约/验收/跨 write set/设计返回类 finding 创建 repair ticket（原票 `needs-repair`、下游重连、原票 Review rounds 记录引用）— `docs/proposals/ticket-refinement-contract/requirements.md:29`、`issues/04-route-review-findings-to-repair-tickets.md:38`；共享契约文本 `assets/skills/_shared/ARTIFACT-CONTRACT.md:97`、`docs/agents/issue-tracker.md:34`。
- 人读文档：`README.md`、`docs/agents/issue-tracker.md` 与受管资产同述修复规则 — `docs/proposals/ticket-refinement-contract/design.md:45`、`issues/04-route-review-findings-to-repair-tickets.md:40`（`README.md` 在票面 Write set，`issues/04-route-review-findings-to-repair-tickets.md:47`）。

wiki 结构：`Write set` 保持只在 `Constraints` 中作为唯一权威，Refine Ticket 只核验其存在与足够窄 — `docs/proposals/ticket-refinement-contract/design.md:32`。

### 5.3 实施状态

**issues 全部 closed（4/4）。**
- `docs/proposals/ticket-refinement-contract/issues/01-derive-execution-contract-gap.md:18` `**Status:** closed`
- `docs/proposals/ticket-refinement-contract/issues/02-add-refine-ticket-and-implement-preflight.md:18` `**Status:** closed`
- `docs/proposals/ticket-refinement-contract/issues/03-add-repair-ticket-lifecycle.md:18` `**Status:** closed`
- `docs/proposals/ticket-refinement-contract/issues/04-route-review-findings-to-repair-tickets.md:18` `**Status:** closed`

完成证据（本次实测复验）：
- 新 skill 已打包进资产并被解包测试覆盖：`... -run '^TestPackagedSkillPointersResolve$'` → `ok flowforge/internal/command 0.267s`（`internal/command/assets_deploy_test.go:24` 经 `assertRequiredArtifactContractPointers` 遍历 `assets/skills`，`internal/command/assets_deploy_test.go:617`）。
- 状态、parser、frontier 行为在测试中固定：`internal/tracker/catalog_test.go:329`（needs-repair 不可执行）、`:370`（原票状态）、`:466`（非可执行）、`:479`（不进 ready）。
- **域内证据缺口（事实陈述）**：4 张票的勾选 Change 均无证据四元组，`./bin/flowforge check --dir docs/proposals/ticket-refinement-contract --strict` 报 `warning: evidence-missing ...` 并以 `Error: proposal diagnostics violate selected policy` 退出。

### 5.4 文档影响（对照现有小节）

| 文档小节（文件:行） | 事实 |
|---|---|
| `README.md:44` `### 3. 发布可执行 ticket`（`README.md:46`） | 只写"必须包含一项可观察交付、…以及成对的完成条件和验证命令"；**无五段机器执行契约、无 `execution-contract-incomplete` gap、无 `--include-gaps` 与 preflight 的拒绝关系** |
| `README.md:57` `### 4. 实现、审查并留下证据`（`README.md:66`） | 只写"Implement 处理完两轴发现…最后才把 ticket 设为 `closed`"；**无 repair ticket 分流、无 `needs-repair`、无 `Repair of:` 与下游重连** |
| `README.md:16`（"`--include-gaps`"段） | 只有一句"可显式使用 `--include-gaps` 或精确 waiver"，未区分"能看"与"Implement 不可绕过" |
| `docs/skill-system.md:5` `## 主交付链` 表（`docs/skill-system.md:7-16`） | **表内没有 `flowforge-refine-ticket` 行**；`flowforge-plan`（`:14`）与 `flowforge-implement`（`:15`）行也未提 refine-ticket 交接与 preflight 拒绝 |
| `docs/skill-system.md:40` `## 支持与特殊路径`（`docs/skill-system.md:42-51`） | 列表未列 `flowforge-refine-ticket` |
| `docs/skill-system.md:81-90` `## 继续、跳过与返回` | 无 repair ticket 路径（`docs/skill-system.md:89` 仅 design return → Solution Design） |
| `docs/cli-design.md:11` `## Artifact Catalog`（`docs/cli-design.md:13-21`） | 只列"只有 `issues/*.md` 中的 ticket 能进入 DAG"与五类确定性诊断；**无 `execution-contract-incomplete` gap，也无 `needs-repair` 状态的 frontier 派发门控** |
| `docs/architecture.md:45` `## 自适应工件，而非固定阶段`（`docs/architecture.md:45-53`） | "状态字段只表示 ticket 执行生命周期"一段未包含 `needs-repair` 这一非终态、非可执行状态 |
| `docs/architecture.md:68` `## 完成不变量` | 未提 repair 闭环（repair 通过自身 review 后回写原票证据并关闭两票） |

### 5.5 候选用户场景

弱模型票常因契约缺失反复返工——Plan 只发标题骨架、refine-ticket 先核验仓库事实再填五段，Implement 对残缺票直接拒收；遇到实质发现时 Review 不再往原票堆 Fix，而是开一张可追溯的 repair 票并让下游等它修完。

---

## 附：本次复验命令与逐字结果（证据）

| 命令 | 结果 |
|---|---|
| `GOPROXY=https://goproxy.cn,direct go test -count=1 ./internal/command/ ./internal/tracker/ -run '^TestPackagedSkillPointersResolve$'` | `ok flowforge/internal/command 0.267s` |
| 同上 `-run '^TestImplementerPromptPinsLoopContracts$'` | `ok flowforge/internal/command 0.329s` |
| 同上 `-run '^TestOpenCodeStepsBudget$'` | `ok flowforge/internal/command 0.461s` |
| 同上 `-run '^TestReviewSkillCarriesRound0Audit$'` | `ok flowforge/internal/command 0.218s` |
| 同上 `-run '^TestEvidenceRepeatFailureThreshold$'` | `ok flowforge/internal/tracker 0.151s` |
| `python3 -m unittest discover -s scripts -p 'executor_metrics_test.py'` | `Ran 108 tests in 0.066s` / `OK` |
| `./bin/flowforge check --dir docs/proposals/{executor-loop-hardening,executor-value-measurement,fast-executor-reliability} --strict` | 三组均 exit 0 |
| `./bin/flowforge check --dir docs/proposals/{lightweight-execution-contract,ticket-refinement-contract} --strict` | 均以 `Error: proposal diagnostics violate selected policy` 退出，诊断全为 `evidence-missing` |
| `./bin/flowforge check --help` | 输出第 5/6 行分别为 evidence 四元组诊断族与重复失败诊断 |
