---
flowforge:
  schema: 1
  role: evidence
  id: cost-aware-verification-planning-01-implementation
  revision: 1
---

# 01 规则包实施记录（G1/G2）

本记录覆盖 Changes 1–7、9、10 的实际执行。结构/一致性测试只证明源、嵌入、部署接缝与既有门禁未回退，不证明模型受控遵从，也不替代 02 的真实情景验收与本项目部署。G0（Change 8 与两项预设测试）见 [测试作者证据](01-test-authoring.md)。

## 固定点

- 父提交（执行起点）：`ba037e78eb602d1342dbb9925da17dcd5d044392`，开始时 `git status --short` 无输出。
- 实施固定点：提交标题 `docs: deliver cost-aware verification rule package (01 G1/G2)`，直接父为上述提交；下表逐文件 SHA-256 是该提交的内容快照，`git status --porcelain --untracked-files=all` 在提交后为空。本文件随该提交写入，因此不内嵌自身所在提交的 SHA。
- 实施提交仅含下列 8 个文件与新增的本文件；`internal/command/assets/**`、`bin/flowforge` 被 `.gitignore` 忽略，由 `make dev` 产生，不作为源改动。

| 文件 | SHA-256 |
| --- | --- |
| `assets/skills/_shared/ARTIFACT-CONTRACT.md` | `c4b47b410bb34f8062af92bd68b77fad295c53af78b72bb100101d528677b22f` |
| `assets/skills/flowforge-solution-design/SKILL.md` | `e7b5b20d911d0d6304d991ff2bb2fd3698de83d8e81c2e70edc11640954147b4` |
| `assets/skills/flowforge-plan/SKILL.md` | `f9cf6bc814d2747915c1c772ee6170ca72b06574a77506e77dc6dbff2aa40db2` |
| `assets/skills/flowforge-refine-ticket/SKILL.md` | `9bc87d038ab84118455505b7b0de64397890fec2d1166ed4b6e8246a8f20c9dd` |
| `assets/skills/flowforge-implement/SKILL.md` | `9cc0053b0e760bb7bfca438a3b8bd30b5e5d404f0d8d8af950d54414bd747825` |
| `assets/skills/flowforge-review/SKILL.md` | `a07b8b452b8e5974d92daf6c1f1c1c08ef15f7e163d558f3886dffbb9cecaa9b` |
| `assets/agents/issue-tracker.md` | `5ae7f8963dab9736da47101285c3457257e3659ec2b100252042ab54bac8d82f` |
| `ff-wiki/proposals/cost-aware-verification-planning/issues/01-rule-package.md` | 随提交内容（本文件与 Implementation note 同批写入） |

## 逐 Change 覆盖

| Change | 落点（真实改动） | 覆盖的可观察行为 |
| --- | --- | --- |
| 1 | `ARTIFACT-CONTRACT.md` 新增 `Verification ownership and reuse`（本地/阶段正向证明、五方所有者）、`Verification schedule`（组覆盖/时机/命令/预期/记录、重型理由、三条执行路径、旧票兼容、已跑命令不重复）、`Shared runs and reused evidence`（四元组归属、共享实跑不重复计数、复用判据与拒绝条件、Fix 前复验安排、Round 0 对实际差分确认、CLI 格式通过不算内容诚实） | 共享语义唯一来源；其余六面只链接不复制模板 |
| 2 | `flowforge-solution-design/SKILL.md` §4 新增 `Verification strategy check.` | 行为/风险→局部证明→阶段证明/触发表、真实环境不后移、重型检查必要性与共享/成本未知、不增前置阶段 |
| 3 | `flowforge-plan/SKILL.md` §2 新增唯一阶段归属/真实依赖/纯验证票非空记录段；§3 Tier 2 新增 `Optional Verification schedule` 条目 | 阶段归属、消费方依赖正向证明、Done and verify 可选安排的完整要素与重型理由；title/Delivery/依赖批准、Mode、五节骨架、standards 转抄未动 |
| 4 | `flowforge-refine-ticket/SKILL.md` §4 三条现有小节扩写 + 新增采用安排核实段 | 组范围、观察时机、真实命令/环境、成功与失败反例、逐验收项覆盖；残缺/失效引用返回规划方，不猜值、不删验收、不改保护测试 |
| 5 | `flowforge-implement/SKILL.md` 新增 `Verification schedule pre-flight`；Phase 0/0b、3a、3b、3d、3e、4、4b、§5 协调三路 | 完整安排按组末验证后勾选、无安排保留原 build 与全部检查、残缺安排 Blocked evidence 返回；共享实跑不记多次；收口只确认完整覆盖；失败即停与 2 次/5 轮上限、测试只读不变 |
| 6 | `flowforge-review/SKILL.md` §3b 新增 `Evidence coverage`；§6 新增 `Publish the Fix verification`；Review rounds 记录新增 `Fix verification` 与 `Counts` 行 | 针对固定点的复验安排（发现/编号/写集/反例、本轮实跑、拟复用版本环境覆盖、扩大条件）、实际差分后的独立确认、review 调用/Fix 批次/部分复发分别计数；关闭所有者、design return/High-Critical/跨写集 repair 边界保留 |
| 7 | `assets/agents/issue-tracker.md` Publishing（Tier 3 纠正为五节、新增可选安排条目）、Execution and review loop（第 2、3 步）、Deterministic checks（新增结构事实≠验收一条） | 与共享解释对齐、旧票兼容，未新增状态/Schema/CLI；Ready/终态明确不是正向验收 |
| 9 | `GOPROXY=https://goproxy.cn,direct make dev` | 源→`internal/command/assets/**`→`bin/flowforge` 一次共享生成 |
| 10 | 本文件 + 票内 Implementation note | 固定点、实跑命令、差分归因、逐 Change 覆盖、未执行项 |

未执行项：无本票内的改动；02 的真实模型情景、本项目 `.agents/**` 与 `ff-wiki/agents/**` 部署、GIIS 迁移、DAG/Catalog/CLI 变更均不在本票范围，未执行也未声称。

## 实际命令与结果

均在本实例真实运行，工作目录 `/Users/qiangbi/develop/projects/Syl/tangram/flowforge`。

| 命令 | 退出码 | 观察输出 |
| --- | --- | --- |
| `GOPROXY=https://goproxy.cn,direct make dev` | 0 | `go build -ldflags="-s -w -X flowforge/internal/version.injected=v5.13.1-43-gba037e7-dirty" -trimpath -o bin/flowforge ./cmd/flowforge` |
| `git diff --no-index --stat assets internal/command/assets` | 0 | 无输出：源与嵌入逐字节一致（规划基线的 4 文件漂移已由 00/03 与本次同步消除） |
| `go test ./internal/command -run 'Test(CostAwareRulePackageSourceEmbeddedDeployed\|PackagedSkillPointersResolve\|PlanPublishesExecutionDetailSkeleton\|PlanSkillDeclaresModeLine\|ArtifactContractDocumentsExecutionContractSections\|AgentRulesDescribeGenericCapabilityDispatch\|ImplementSkillCarriesWeakExecutorContract\|ReviewSkillCarriesRound0Audit\|ReviewSkillDescribesRepairEscalation\|IssueTrackerDescribesRepairWorkflow)$' -count=1` | 0 | `ok flowforge/internal/command 0.671s`；`CostAwareRulePackageSourceEmbeddedDeployed` 的七个源/嵌入/隔离部署子项全部通过 |
| `go test ./internal/tracker -run 'Test(CostAwareSharedRunPreservesEvidenceGate\|ExecutionContractCompletenessAppliesOnlyToManagedExecutableTickets\|EvidenceQuadrupleParsing\|NeedsRepairIsNonExecutableAndRepairProvenanceIsParsed)$' -count=1` | 0 | `ok flowforge/internal/tracker 0.418s` |
| `/Users/qiangbi/develop/projects/Syl/tangram/flowforge/bin/flowforge init`（隔离目录 `/tmp/ff-iso`） | 0 | `✓ Deployed flowforge skills to .agents/skills/`、`✓ Configured /tmp/ff-iso/ff-wiki/agents/ rules`、`✓ Managed assets verified current.` |
| `/Users/qiangbi/develop/projects/Syl/tangram/flowforge/bin/flowforge assets verify --json`（同上隔离目录） | 0 | 全部条目 `"state": "current"` |
| `cmp` 七面源 ↔ 隔离目录部署副本 | 0 | 六个技能文件与 `ff-wiki/agents/issue-tracker.md` 全部一致（新二进制→部署接缝） |
| `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` | 0 | `ok` 于 command 8.982s / config / subagent / tracker 0.308s / update；`?   flowforge/internal/version [no test files]`；240 个顶层用例 PASS、0 FAIL、3 SKIP |
| `/tmp/flowforge-verification-tools/golangci-lint-2.15.0/golangci-lint run --timeout 2m --new-from-rev 16cce8e --whole-files --max-issues-per-linter 0 --max-same-issues 0 ./cmd/... ./internal/...` | 0 | `0 issues.` |

3 项 SKIP 为既有环境性跳过，未据此声称对应行为已实测：`TestRepositoryAndTangramProposalLayouts`（`catalog_test.go:531 Tangram compatibility checkout is not available`）、`TestE2E_SelfupdateApply_ReplacesBinaryCorrectly` 与 `TestFullUpgradePipeline`（`cannot read /proc/self/exe`，Darwin 无该路径）。全仓/全产品包的既有 lint 失败按 [lint 基线](01-lint-baseline.md) 保留；本次 `0 issues.` 只覆盖相对 `16cce8e` 的 cmd/internal 改动文件整文件门禁，不代表全仓 lint 通过。

## 源→嵌入差分归因

- 规划基线记录的 4 文件漂移（`assets/pi/flowforge.ts`、`ARTIFACT-CONTRACT.md`、`flowforge-implement/SKILL.md`、`flowforge-review/SKILL.md`）在 `make dev` 前已由 00/03 消除；执行前 `git diff --no-index assets internal/command/assets` 已为空。
- 本次 `make dev` 只把 7 个源规则改动带入嵌入树，无其他差异：事后 `git diff --no-index --stat assets internal/command/assets` 仍为空，`internal/command/assets/**` 与 `bin/flowforge` 均为生成物且未改动任何 PI 源。`bin/flowforge` SHA-256 为 `b1d9f0634b0bdad047cfa7a548e9e2eabec4208e8b4464cf4867bc44cc1d57af`。

## 部署边界（如实保留）

- 本仓库自身的 `.agents/skills/**`（gitignore 的每机部署）早于 `ba037e7` 的源内容，未被本票刷新；`ff-wiki/agents/issue-tracker.md`（已跟踪的部署副本）目前等于 `HEAD:assets/agents/issue-tracker.md`，即只缺本票新增内容。两者都属项目部署，按票面归 02；本票只证明源→嵌入→隔离部署一致，不声称本仓库已按新规则部署。
- 02 仍 `blocked_by: 01`，其 `execution-contract-incomplete` 不在本票范围；本票未关票、未改状态、未迁移任何旧票。

## 边界声明

- 自动结构/一致性测试通过不是模型行为证明，也不是阶段验收；Changes 1–7 的语义合规由独立 Review 逐条核对设计责任表、三条兼容路径、Fix/repair 分界与四元组诚实性。
- 一条真实 `make dev` 与各组定向/全内包/lint 运行各计一次执行；共享关系见票内四元组与 Implementation note，未把一次运行记成多次。

## Fix 轮次：Changes 11–13（消费 Round 1 `Fix verification`）

### 固定点与实际差分

- 起点：`348f0f449c5408e6008dd6195e357bd7223048e6`（Round 0 记录提交），执行前工作树干净。
- 本轮固定点为紧随其后的提交（标题 `docs: fix verification entry order and project-root pointers (01 Fix 11-13)`）；下表逐文件 SHA-256 是该提交的内容快照。工作树改动只有 6 个源 Markdown 与「本票 + 本文件」；`internal/command/assets/**`、`bin/flowforge` 由 `make dev` 生成、被 `.gitignore` 忽略，未出现在 `git status --porcelain` 中。
- 生效安排：本票 `Review rounds / Round 1 / Fix verification`（A01/A02/A03 → 11/12/13 同批）。本票同时有完整 `Verification schedule`，但按 A01 的入口顺序，本轮由该 Fix verification 支配：不重复同固定点的定向套件，改由一次全内包运行覆盖。

| 文件 | SHA-256 |
| --- | --- |
| `assets/skills/_shared/ARTIFACT-CONTRACT.md` | `71f65b5a2dc421ea698f881b588a1bf4c47df031760aaccaf24630d3dff0fcfe` |
| `assets/skills/flowforge-plan/SKILL.md` | `88f867a4cbcdd423b8ef0559d31e3f631ac092bb21ec7e4ea58233c6928f842e` |
| `assets/skills/flowforge-refine-ticket/SKILL.md` | `c96073a8cd5a712ff62af9a5c501e6fdfa9be07e389375a40221374e71cfc9ff` |
| `assets/skills/flowforge-implement/SKILL.md` | `96f6c0df2e2666699090a9fa9fd5b8e0cbc226ece7ad8b5325dc16a01fa8111c` |
| `assets/skills/flowforge-review/SKILL.md` | `8b44ff313dc91c890bf3a877847c6364c8a3c613d26a9a55142860d01ef51cdf` |
| `assets/agents/issue-tracker.md` | `9d9692ccf275b6dbbe8d744e8085720a963ded32ccec406b73d9a8b26b928480` |

`assets/skills/flowforge-solution-design/SKILL.md` 本轮未改，仍为 `e7b5b20d911d0d6304d991ff2bb2fd3698de83d8e81c2e70edc11640954147b4`。源→嵌入差分只有上述 6 面：`make dev` 后 `git diff --no-index --stat assets internal/command/assets` exit 0 且无输出；`bin/flowforge` SHA-256 为 `95f789832d0dd68bfad60cfe0783a687480cdefc36775372c01bdc376a5e5710`。

### 逐 Change 覆盖

| Change | 落点（本次提交的真实行号） | 交付的可观察行为 |
| --- | --- | --- |
| 11 | `ARTIFACT-CONTRACT.md:141` `#### Effective verification entry`（顺序、残缺处置、执行者不得选测/判复用）、`:129` schedule 收尾的优先级指针；`flowforge-implement/SKILL.md:28` `Verification schedule pre-flight`、`:44` Phase 0、`:48` Phase 0b、`:60-61` 3a、`:71` 3b、`:86` 3d、`:90` 3e、`:96` Full mode、`:130` §5；`flowforge-review/SKILL.md:83` Round 0 证据覆盖基线、`:125` 发布即本轮入口；`flowforge-refine-ticket/SKILL.md:45`；`issue-tracker.md:11,19,31,38` | 有效验证入口按 Fix verification → 完整 schedule → 旧逐 Change/全部 Done 解析；最新 Fix 安排只覆盖其列明 Fix，保留原验收项、观察方法、写集与未覆盖项，不授权删 Done；声明但残缺、引用失效或失配固定点先 Blocked 并回发布者（schedule 缺陷仍回规划方）；Full mode 与收口消费同一解析 |
| 12 | `ARTIFACT-CONTRACT.md:111` `### Stage proof before dispatch`；`flowforge-plan/SKILL.md:28`；`flowforge-implement/SKILL.md:30` `Stage proof pre-flight`；`issue-tracker.md:72` | 派发负责人在 launch 前确认存在匹配待消费/发布版本与环境、且经独立审查的实际正向验收记录；从未通过、无证明的 `wontfix`/`resolved`/`closed`、过期或不匹配均不派发；`Ready`/终态只管图；收票侧缺失/失配先 Blocked 回派发负责人；范围重定仍归需求/设计权威。未新增状态或 CLI 能力 |
| 13 | `issue-tracker.md:11`（项目根指针一次定义）、`:19`、`:38`（两处源码相对链接改为同形指针）、`:72`（新增引用同形） | shared authority 指针改为从消费项目根读取的 `.agents/skills/_shared/ARTIFACT-CONTRACT.md`，删除 `../../assets/...` 与只在默认一层 docs 目录有效的路径；未加模板替换器/部署 API/Go 改动，源与部署保持同字节 |

### A01 / A02 / A03 反例的观察位置

- **A01（Fix 安排未接入机械执行路径）**：入口顺序的唯一出处是 `ARTIFACT-CONTRACT.md:141`；`flowforge-implement/SKILL.md:28` 先解析它，`:44`、`:48`、`:60-61`、`:71`、`:90`、`:96`、`:130` 全部消费同一解析；3b 不再写「运行 Done and verify 的全部命令」，3e 不再把历史结果限定为「票内 schedule 或 Done 明文允许」。本轮自身是该反例的运行实例：本票同时有完整 `Verification schedule`，而 11–13 依 Round 1 `Fix verification` 执行，Done and verify 的定向规则包与 tracker 套件没有被单独重跑，由唯一一次全内包运行覆盖；旧 3b 文本会要求它们再跑一遍。
- **A02（初次发布缺正向证明仍被派发）**：规则落点 `ARTIFACT-CONTRACT.md:111`（拒绝条件与责任人）、`flowforge-plan/SKILL.md:28`（消费票写明 proof 引用/版本/确认角色）、`flowforge-implement/SKILL.md:30`（收票侧强制 Blocked 回派发负责人）、`issue-tracker.md:72`（同一规则）。本轮不派发需要阶段证明的票，故只给出源规则观察位置；真实派发拒绝由 02（`blocked_by: 01`）在集成验收时实测，01 不声称已行为验证。
- **A03（部署后指针失效）**：规则落点 `issue-tracker.md:11,19,38`。两个全新临时目录实测：默认 `ff-wiki`（`/tmp/ff-fixA2-RqyaJm`）与预写 `docs_dir: nested/wiki`（`/tmp/ff-fixB2-7Mzauq`）两种深度都从项目根解析到 `.agents/skills/_shared/ARTIFACT-CONTRACT.md`，`Verification schedule`、`Effective verification entry`、`Execution and review loop`、`Stage proof before dispatch` 四个标题全部命中；部署副本中 `](../../assets/...` 旧链接计数为 0；同一副本里该旧形式解析到不存在的 `<project>/assets/...`（`nested/wiki` 深度解析到 `nested/assets/...`），即被替换的正是失效路径。

### 实际命令与结果（本轮各一次）

| 命令 | 退出码 | 观察输出 |
| --- | --- | --- |
| `GOPROXY=https://goproxy.cn,direct make dev` | 0 | `go build -ldflags="-s -w -X flowforge/internal/version.injected=v5.13.1-45-g348f0f4-dirty" -trimpath -o bin/flowforge ./cmd/flowforge` |
| `git diff --no-index --stat assets internal/command/assets` | 0 | 无输出：源与嵌入逐字节一致 |
| `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` | 0 | `ok flowforge/internal/command 9.902s`、`config (cached)`、`subagent (cached)`、`ok flowforge/internal/tracker 0.765s`、`update (cached)`；`? flowforge/internal/version [no test files]`；240 个顶层用例 PASS、0 FAIL、3 项既有环境 SKIP |
| 同一份日志内的定向集合（不另起一次运行） | 0 | `PASS: TestPackagedSkillPointersResolve`、`TestAgentRulesDescribeGenericCapabilityDispatch`、`TestPlanPublishesExecutionDetailSkeleton`、`TestPlanSkillDeclaresModeLine`、`TestArtifactContractDocumentsExecutionContractSections`、`TestImplementSkillCarriesWeakExecutorContract`、`TestReviewSkillCarriesRound0Audit`、`TestReviewSkillDescribesRepairEscalation`、`TestIssueTrackerDescribesRepairWorkflow`、`TestCostAwareRulePackageSourceEmbeddedDeployed`（七个子项全 PASS）、`TestEvidenceQuadrupleParsing`、`TestCostAwareSharedRunPreservesEvidenceGate`、`TestExecutionContractCompletenessAppliesOnlyToManagedExecutableTickets`、`TestNeedsRepairIsNonExecutableAndRepairProvenanceIsParsed` |
| `golangci-lint 2.15.0 run --timeout 2m --new-from-rev 16cce8e --whole-files --max-issues-per-linter 0 --max-same-issues 0 ./cmd/... ./internal/...` | 0 | `0 issues.` |
| 新二进制 `init`（`/tmp/ff-fixA2-RqyaJm`，默认 `ff-wiki`） | 0 | `✓ Deployed flowforge skills to .agents/skills/`、`✓ Configured /tmp/ff-fixA2-RqyaJm/ff-wiki/agents/ rules`、`✓ Managed assets verified current.` |
| 新二进制 `init`（`/tmp/ff-fixB2-7Mzauq`，预写 `docs_dir: nested/wiki`） | 0 | `✓ Configured /tmp/ff-fixB2-7Mzauq/nested/wiki/agents/ rules`、`✓ Managed assets verified current.` |
| `assets verify --json`（A、B 各一次） | 0 | 两处均 `entries: 52 {'current': 52}`，无 missing/drifted |
| `cmp` 七面源 ↔ 部署副本（A、B 各一次） | 0 | 六个技能文件与 `<docs>/agents/issue-tracker.md` 全部一致（新二进制 → 部署接缝），无 `DIFF` 输出 |
| 项目根指针与标题解析（A、B 各一次） | 0 | 见 A03 段：两处 4 个标题全部命中，旧相对形式目标不存在 |
| `git diff --check` | 0 | 无输出 |
| `flowforge check --dir ff-wiki/proposals` | 0 | `Checked 106 issues`；诊断计数与实施前基线相同，唯一 gap 仍是范围外的 `02-controlled-acceptance-and-deploy.md`；Dependency graph healthy |

3 项 SKIP 为既有环境性跳过，不能称为对应行为已实测：`TestRepositoryAndTangramProposalLayouts`（`catalog_test.go:531 Tangram compatibility checkout is not available`）、`TestE2E_SelfupdateApply_ReplacesBinaryCorrectly` 与 `TestFullUpgradePipeline`（`upgrade_test.go:219/347 cannot read /proc/self/exe: open /proc/self/exe: no such file or directory`，Darwin 无该路径）。整文件 lint 的 `0 issues.` 只覆盖相对 `16cce8e` 的 cmd/internal 改动文件，不代表全仓 lint 通过；全仓既有失败按 [lint 基线](01-lint-baseline.md) 保留。

### 复用与不重复

- 三项 Go 测试来源与内容未变：`internal/command/assets_deploy_test.go` `503b574e…`、`internal/command/cost_aware_assets_test.go` `7c3fcd13…`、`internal/tracker/cost_aware_evidence_test.go` `101e4219…`，与 [测试作者证据](01-test-authoring.md#精确应用与内容对照) 逐字节一致；`git diff 16cce8e..HEAD -- '*.go'` 只有这三路径且本轮未触碰；模块依赖与 lint 配置未变，工具仍为 `golangci-lint 2.15.0`。Round 1 的复用条件成立，本轮未新增或修改任何测试。
- 整文件 lint 的复用条件同样成立（实际差分只有 6 面源规则 + 本票 + 本证据 + 生成嵌入，无 Go/测试/依赖/配置变化）。为不把旧结果写成新执行，仍在新固定点实跑一次，`0 issues.` 属本轮真实输出；全仓历史失败未重跑也未改写。
- 定向规则包与 tracker 定向命令没有重复执行：它们由唯一一次全内包运行的日志覆盖（上表第四行是同一份日志）。一次 `make dev`、一次全内包、一次 lint、一次部署验证、一次 `flowforge check` 各计一次，未把共享运行记成多次。

### 边界（如实保留）

- 结构/一致性测试只证明源→嵌入→隔离部署与既有门禁未回退，不证明模型受控遵从 A01/A02/A03 的新语义；这些仍需 02 的真实受控情景与本项目部署验收，01 不以本记录替代。
- 本仓自身 `.agents/skills/**`（每机部署、早于本轮）未刷新；`ff-wiki/agents/issue-tracker.md`（已跟踪部署副本）仍等于 `348f0f4` 时的源内容，带旧的 `../../assets/...` 指针 —— 项目部署归 02，本票只证明源与隔离部署一致。
- 未关票、未改 `**Status:**`、未 push、未委派、未改任何测试或既有断言、未新增 Schema/CLI/状态。
