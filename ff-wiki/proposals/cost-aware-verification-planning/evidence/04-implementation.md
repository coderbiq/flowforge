---
flowforge:
  schema: 1
  role: evidence
  id: cost-aware-verification-planning-04-implementation
  revision: 1
---

# 04 阶段证明门禁范围与执行记录修复

本记录覆盖 [04 修复票](../issues/04-stage-proof-scope-and-execution-record.md) 的 Changes 1–3。它证明源规则收窄、记录可溯源与源→嵌入→隔离部署一致，不证明模型遵从新语义；三组语义反例由独立 Round 0 核源，02 另做模型实测。

## 固定点与实际差分

- 起点 `b118c98`（04 合同发布提交）；本轮实施固定点 `32722d7327bda96975c9a06a4ee202d7a6fe09af`，内容快照如下表。后续A06仅更正记录，不改变该产品快照。
- 实际差分：2 个源规则句、`evidence/01-implementation.md` 的 Fix 轮次记录、01 票 Fix 11–13 记录段、04 票、本证据；`internal/command/assets/**` 与 `bin/flowforge` 仅由 `make dev` 生成（被 `.gitignore` 忽略）。

| 文件 | SHA-256 |
| --- | --- |
| `assets/skills/_shared/ARTIFACT-CONTRACT.md` | `fcdc87d778bd9aff6f49726dafe2630a9cbcbd136e3b623184c3fd1fb1ca82fd` |
| `assets/skills/flowforge-plan/SKILL.md` | `94e0ba8d274bad4affbe12a8d3ecd677aaeb6467379f0f2a2fdb88946039fa80` |
| `ff-wiki/proposals/cost-aware-verification-planning/evidence/01-implementation.md` | `ff317ab085bac70a4ec2fd8d5430724459793282bab5a3eaf9ce0d171f74b277` |
| `ff-wiki/proposals/cost-aware-verification-planning/issues/01-rule-package.md` | `f8ed8d734517e5b4673de187d0525732460d6d0311b21187da33c2a6666837fc` |

`bin/flowforge` SHA-256 `4d916a94439e0973ca97ee223fc0cc71ed354b489dba27cd80f0e793b7e677cc`。

## Change 1：门禁范围恢复（A04）

- shared `### Stage proof before dispatch` 首句改为条件化：只有需求/设计/票面点名某既有阶段正向证明为入场条件时才受门禁；单纯有上游交付不加要求。三条正负例逐条可核对：**G** 消费局部 A/B 并交付组合证明，不以自身将产生的同一证明为入场条件（另有的更早不同证明仍单独检查）、**P** 需完整证明的发布/实际使用在 G 正向记录存在后才派发、**普通 A→B** 用自身局部证明，不凭空造阶段证明。
- 派发侧拒绝段与收票侧 `STATUS: BLOCKED` 段原文保留：从未执行、无记录的 `wontfix`/`resolved`/`closed`、过期或异版本记录仍不派发，范围重定仍归需求/设计权威。
- plan 消费句改为同一条件化表述（`A ticket the requirement, design, or ticket text gates on ...`），并写明三种情形的派发后果；`flowforge-implement/SKILL.md:30` 与 `issue-tracker.md:72` 本已只针对“需要阶段证明的票”，核对后未改，不扩约束。

## Change 2：执行记录溯源（A05）

对照审计日志 `rule-package-delivery.audit-01-1.history-evidence.log`（原始 `61b27f54-…_flowforge-implementer_transcript.jsonl` 行摘录）逐条核对：

| 轨迹 | 事实 |
| --- | --- |
| `:188–202` | 首批 S1：`make dev` exit 0、源=嵌入；全内包 240 PASS/0 FAIL/3 SKIP（command 9.850s、tracker 1.108s）；lint `0 issues.` |
| `:204–210` | S1 两处隔离部署 `/tmp/ff-fixA-7XHacA`、`/tmp/ff-fixB-ZwyS8a`：52/52 `current`、七面 `same:` |
| `:255–265` | S1 之后真实修改 `flowforge-implement/SKILL.md` 标题与 Phase 0 句 → 第二批必要 |
| `:275–289` | 第二批 S2：`make dev`（`bin` `95f7898…5710`）、全内包 240 PASS（9.902s/0.765s）、lint `0 issues.`、`diff --check` exit 0、两处新目录 52/52 `current` 且七面无 `DIFF` |
| `:299–313` | 证据用 `python3` 截断 + `cat >>` heredoc 写入；01 票四元组用 `python3` 勾选插入；标成“执行前”的 restatement/preflight 于执行后插入 |

已原位改写（不是追加一句让原说法继续有效）：

- `evidence/01-implementation.md`：标题“（本轮各一次）”→“S1、S2 两批，每批内各计一次”；表拆两批并引入 S1/S2 与真实目录；固定点段写明 S1 hash 缺失、不作当前验收、S2 对应 `c990efb`；“唯一一次全内包运行”改为“该批”；“复用与不重复”的计数改为两批各一次；新增“记录更正与溯源”段列出原句、事实行号、现记法与未做事项。
- 01 票：`Phase 0 restatement`/`Phase 0b preflight` 标题改为“事后回填”并声明非执行前；“全部 edit/write，无 CLI 生成”改为如实记录 `python3`/heredoc 流程偏离（04 起改用 `edit`/`write`）；执行结果表拆 S1/S2；收口覆盖确认只把 S2 当当前验收。

未倒填执行顺序、未用本轮补跑冒充 01 历史、未抹掉 S2 的真实 PASS、未删除任何早批次记录，未改 01 状态与验收。

## Change 3：当轮验证

| 命令 | 退出码 | 观察输出 |
| --- | --- | --- |
| `GOPROXY=https://goproxy.cn,direct make dev` | 0 | `MAKE_EXIT=0`；随后 `git diff --no-index --stat assets internal/command/assets` 无输出（`DIFF_EXIT=0`） |
| `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` | 0 | `ok flowforge/internal/command 10.540s`、`config (cached)`、`subagent (cached)`、`ok flowforge/internal/tracker 1.284s`、`update (cached)`、`? flowforge/internal/version [no test files]`；`PASS=240 FAIL=0 SKIP=3`；定向集合 14 项命中、`TestCostAwareRulePackageSourceEmbeddedDeployed` 7 个子项 |
| 新二进制 `init` + `assets verify --json` + 七面 `cmp` + 指针/标题解析（`/tmp/ff-04A-XDwbei`，默认 `ff-wiki`） | 0 | `✓ Managed assets verified current.`；`entries: 52 {'current': 52}`；`seven-surface cmp done` 无 `DIFF`；项目根 `.agents/skills/_shared/ARTIFACT-CONTRACT.md` 存在，4 个标题命中，新门禁句与三条正负例在位，旧越界句计数 0 |
| 同上（`/tmp/ff-04B-CwRsIo`，预写 `docs_dir: nested/wiki`） | 0 | 同上一行（`✓ Configured …/nested/wiki/agents/ rules`），`52 {'current': 52}`，无 `DIFF`，pointer/标题/新句检查全为 True |
| `git diff --check` | 0 | 无输出 |
| `./bin/flowforge check --dir ff-wiki/proposals` | 0 | `Checked 107 issues`；`gap:execution-contract-incomplete: 1`（范围外 `02-controlled-acceptance-and-deploy.md`，非本轮新增）、`warning:evidence-artifact-missing: 19`（与基线同）、`evidence-missing: 202`、`legacy-metadata: 10`、`missing-required-metadata: 5`、`untracked-upstream: 4`；`Dependency graph is healthy` |

3 项 SKIP 为既有环境性跳过，不称其行为已测：`TestRepositoryAndTangramProposalLayouts`（`catalog_test.go:531 Tangram compatibility checkout is not available`）、`TestE2E_SelfupdateApply_ReplacesBinaryCorrectly` 与 `TestFullUpgradePipeline`（`upgrade_test.go:219/347` Darwin 无 `/proc/self/exe`）。

- 收尾前另做一次完整内部 Go 运行（同一最终工作树，`go test ./internal/...`）：`ok flowforge/internal/command 9.281s`、`ok flowforge/internal/tracker 0.483s`、其余 `(cached)`，exit 0。这是超出预定一次安排的重复运行，未发现新的源变更或失败需要它，不能事后赋予其必要性；不替代上面的逐用例日志，两次均保留。

## 复用与不重复

- lint 复用 `c990efb` 的整文件结果 `0 issues.`：`git diff --stat c990efb..HEAD -- '*.go' 'go.mod' 'go.sum' '.golangci*'` 无输出，工具仍为 `golangci-lint 2.15.0`（Go 源码/测试/依赖/lint 配置/工具版本均未变），条件成立，本轮**未重跑**，也不把旧结果写成新执行；全仓既有失败按 [lint 基线](01-lint-baseline.md) 保留。
- 三个 Go 测试（`assets_deploy_test.go`、`cost_aware_assets_test.go`、`cost_aware_evidence_test.go`）本轮只读，未改任何断言；`assets/pi/**` guard、保护与模型配置未动。
- 04实施会话实际计数（截至 `32722d7`）：`make dev` 1次、全内包Go 2次（逐用例1次、额外非verbose1次）、lint运行0次（复用）、隔离部署2个目录各一次init/verify、`git diff --check` 2次与 `git diff --cached --check` 1次、`flowforge check` 7次。定向集合由逐用例全内包日志覆盖，没有独立启动定向套件。后续独立审查和主会话记录检查另计，不混入实施次数。

### A06 记录更正与流程偏离

主会话依据独立 `rule-package-repair.audit-04.md` 及原始 `8ebbd7dd-0949-43f4-939d-a4b08d73f9ea_flowforge-implementer_transcript.jsonl` 正常文件编辑更正：Go调用在140/196行；check调用在152/156/220/224（两次）/232/259行；diff空白检查在152/156/259行。原“check一次”及票内“Go一次”不成立，现按真实次数统一。部分check命令经过head/tail/grep管道，末尾shell退出码不能当作每个上游check的独立退出码；保留实际调用与输出，最终图健康另由独立Round0验证。220行frontier输出解析脚本失败、后续成功的轨迹也保留，不将其改成首次即成功。

额外Go运行和多次check属于实际开销，未获验证安排授权的重复运行如实记录为流程偏离；本次只修复记录，不补跑冒充历史合规。比较基线应为 `c990efb` 至当前：Go/依赖/lint配置无变化；`16cce8e` 至当前仍包含G0的三个测试文件，不能称该范围Go差分为空。原审查报告和原始轨迹保留，不改写其历史结果。

## 边界（如实保留）

- 结构/一致性测试只证明源→嵌入→隔离部署与既有门禁未回退，不证明模型在真实受控情景遵从 A04 新语义；`G`/`P`/普通局部三组语义反例由独立 Round 0 核源，模型实测归 02。
- 本仓自身 `.agents/skills/**`（每机部署）未刷新，`ff-wiki/agents/issue-tracker.md` 未改；项目部署归 02。
- 未关票、未改 01 状态/验收、未改 Schema/CLI/DAG/Catalog/四元组/默认 Mode、未 push、未委派、未自审。
