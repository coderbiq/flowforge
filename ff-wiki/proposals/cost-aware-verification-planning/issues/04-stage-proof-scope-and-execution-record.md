---
flowforge:
  schema: 1
  role: ticket
  id: cost-aware-verification-planning-04
  revision: 1
  consumes:
    requirements:
      cost-aware-verification-planning-requirements: 2
    design:
      cost-aware-verification-planning-design: 3
---

# 04: 恢复阶段证明门禁适用范围与执行记录

**Blocked by:** None
**Status:** open
**Mode:** lightweight
**Repair of:** cost-aware-verification-planning/01

## Delivery

修复 01 修后独立只读 Round 0 的两项发现：A04 `contradicts`/Medium（阶段证明门禁被扩大到所有上游交付消费者）与 A05 `partial`/Medium（Change 10 执行记录不完整、部分过程声明与原始轨迹相反）。同时交付 A04 的正确派发边界与 A05 的可追溯记录，保持缺证/失配时的真实发布拒绝，并保留最终快照的真实 PASS。修后审查覆盖完整 01 范围 `16cce8e...HEAD`，不能只审两句 delta。不改变 requirements/design/Schema/CLI，不删除任何早批次记录。

## Design context

A04 的权威含义问题已由设计所有者明确解决（`agent_message`，主会话；转述如下）：设计 rev3 语义不变——只有需求/设计/票面明确要求某个既有阶段正向证明作为入场条件的消费/发布任务才受该门禁；单纯有上游交付不增加 stage-proof 要求。A/B 已有局部证明，G 交付首次 A+B 组合证明：G 不能以它将产生的同一证明作为入场条件，P 若要求 G 证明则缺失时不派发。普通 A→B 局部依赖无需凭空造阶段证明。若 G 本身被权威明确要求另一个更早且不同的阶段证明，仍检查该独立入场条件。此为设计返还后的既定含义恢复，无新的需求待用户确认；引用现有 [局部交付与阶段集成](../design.md#allocation)（`design.md:50–64`），其中 `:64` 的门禁前提为“派发需要阶段证明的消费/发布票前”。不把所有 producer 一概豁免。

消费 [需求 revision 2](../requirements.md#cost-aware-verification-planning-requirements) 与 [设计 revision 3](../design.md#cost-aware-verification-planning-design)。两项发现来自修后独立只读 Round 0：A04 `contradicts`（`assets/skills/_shared/ARTIFACT-CONTRACT.md:113` 与 `assets/skills/flowforge-plan/SKILL.md:28` 的实际句子）与 A05 `partial`（`evidence/01-implementation.md:112,135`、`issues/01-rule-package.md:229,235,237` 与原始 transcript）。01 因此进入 `needs-repair` 且不可执行，02 的 `**Blocked by:**` 由 01 改为本票；00/03 保持 closed。A05 属同一 repair 内的记录恢复，不另造产品范围。

## Touch points

- `assets/skills/_shared/ARTIFACT-CONTRACT.md` — `### Stage proof before dispatch`。
- `assets/skills/flowforge-plan/SKILL.md` — 消费票的 stage proof 句（`A ticket that consumes another delivery ...`）。
- `assets/skills/flowforge-implement/SKILL.md`、`assets/agents/issue-tracker.md` — 仅核对同义消费句，不扩约束。
- `ff-wiki/proposals/cost-aware-verification-planning/evidence/01-implementation.md` — Fix 11–13 记录修正段。
- `ff-wiki/proposals/cost-aware-verification-planning/issues/01-rule-package.md` — 本票对应历史记录修正段（不改状态/验收）。
- `Makefile:dev` — 生成 `internal/command/assets/**` 并构建 `bin/flowforge`。
- `ff-wiki/proposals/cost-aware-verification-planning/evidence/04-implementation.md` — 本票当轮四元组与记录。

## Changes

- [x] 1. 在 `assets/skills/_shared/ARTIFACT-CONTRACT.md` 的 `### Stage proof before dispatch` 及 `assets/skills/flowforge-plan/SKILL.md` 的相应句（`A ticket that consumes another delivery ...`）写明本票 Design context 的设计所有者裁决：门禁只适用于需求/设计/票面明确要求既有阶段正向证明作为入场条件的消费/发布任务，单纯有上游交付不增加 stage-proof 要求；给出 G/P/普通局部依赖三种正负例（G 不以其将产生的同一证明为入场条件、P 缺 G 证明不派发、普通 A→B 局部依赖不加门禁），保持缺证/失配时的真实发布拒绝。必要处只核实同义指针，不扩大约束。
  - cmd: `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`（当轮一次全内包运行，覆盖该两句的源/嵌入/部署三层检查）
  - exit: 0
  - output: `PASS=240 FAIL=0 SKIP=3`；`--- PASS: TestArtifactContractDocumentsExecutionContractSections`、`--- PASS: TestPlanPublishesExecutionDetailSkeleton`、`--- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed/skills/_shared/ARTIFACT-CONTRACT.md`、`.../skills/flowforge-plan/SKILL.md`
  - artifact: `assets/skills/_shared/ARTIFACT-CONTRACT.md`
- [x] 2. 依据本轮独立 audit 的 `history-evidence.log` 及原始 `61b27f54-eaed-4a1d-b658-3e23df77b157_flowforge-implementer_transcript.jsonl`，修正 `evidence/01-implementation.md` 与 01 票 Fix 11–13 的 Implementation note：保留两批 `make`/全内包/lint 及四个隔离 `init`/`verify` 的真实记录，说明中间源修改使第二批必要；用 S1/S2 区分快照，第一批未存完整 hash 时如实写缺失且不用它作当前验收，最终 S2 对应 `c990efb`。纠正“本轮各一次/唯一一次”“执行前”“全部 edit/write”等失实表述：restatement/preflight 是事后回填、证据与票曾用 bash/Python 写入，属已确认流程偏离；不倒填执行顺序、不补跑冒充历史、不抹掉最终真实 PASS。主会话处置为恢复准确记录 + 独立审查受影响文档/实际差分 + 后续严格文件工具，不要求无意义重新执行旧轮次。不改 01 的状态/验收，不把旧证据写成假通过。
  - cmd: `./bin/flowforge check --dir ff-wiki/proposals`（记录逐条核对 `history-evidence.log:188–313` 后的结构门禁）
  - exit: 0
  - output: `Checked 107 issues`；`gap:execution-contract-incomplete: 1`（范围外 02，非本轮新增）；`warning:evidence-artifact-missing: 19`（与基线同）
  - artifact: `ff-wiki/proposals/cost-aware-verification-planning/evidence/01-implementation.md`
- [x] 3. 全部源改完后执行一次 `make dev`；跑一次完整内部 Go、默认与 `nested/wiki` 两种隔离 `init`/`verify` + 七面 byte/pointer 检查、`git diff --check`、`flowforge check`、`flowforge frontier`；把当轮四元组与记录写入 `evidence/04-implementation.md`，不提前关票。
  - cmd: `GOPROXY=https://goproxy.cn,direct make dev`；`git diff --no-index --stat assets internal/command/assets`
  - exit: 0
  - output: `MAKE_EXIT=0`、`DIFF_EXIT=0`（无输出）；`bin/flowforge` `4d916a94439e09…e677cc`；两处隔离目录 `entries: 52 {'current': 52}` 且七面无 `DIFF`
  - artifact: `ff-wiki/proposals/cost-aware-verification-planning/evidence/04-implementation.md`

## Constraints

- must 直接用文件工具维护 proposal Markdown，CLI 仅承担现有图与诊断计算 — [项目核心原则](../../../../AGENTS.md#核心设计原则)
- must not 改变 Issue Schema 头或 CLI 签名；本方案用现有正文及依赖表达 — [项目边界](../../../../AGENTS.md#boundaries)
- must 保持测试作者/轻量实现者分离、每 ticket 的执行上下文策略和独立审查 — [执行单元政策](../../../../AGENTS.md#execution-unit-policy)、[R1](../requirements.md#reliability)
- must not 把旧结果当新执行、把结构诊断通过当行为证明，证据仍须独立核实 — [证据契约](../../../../assets/skills/_shared/ARTIFACT-CONTRACT.md#hand-offs)、[审查交付核对](../../../../assets/skills/flowforge-review/references/review-protocol.md#test-pass-does-not-equal-requirement-met)
- must 只让有具体失败依据的发现形成修复义务，保留真实失败及已约定标准 — [审查反猜测规则](../../../../assets/skills/flowforge-review/references/review-protocol.md#no-speculation)
- must 不在 assets 放不部署的情景材料，按源→嵌入→部署接缝验证一致性 — [项目边界](../../../../AGENTS.md#boundaries)、[构建入口](../../../../Makefile)
- 轻量执行者不得更改任何预设测试；01 Change 8 已由合法作者在 G0 完成，本票对全部 Go 测试只读，不修改 `assets/pi/**` 源码 guard，不修改任何保护或模型配置。
- 不修改 DAG/Catalog 算法、evidence 四字段、默认 Mode 或旧票；不迁移 GIIS。保持 Round 0、独立双轴、零未处置发现关票与 substantive finding → repair。A04 只按设计所有者裁决恢复范围，不以删除门禁代替。
- Write set: assets/skills/_shared/ARTIFACT-CONTRACT.md; assets/skills/flowforge-plan/SKILL.md; internal/command/assets/**（仅 make dev 生成）; bin/flowforge（仅构建）; ff-wiki/proposals/cost-aware-verification-planning/issues/01-rule-package.md（仅 Fix 11–13 历史记录修正段，不改状态/验收、不把旧证据写成假通过）; ff-wiki/proposals/cost-aware-verification-planning/evidence/01-implementation.md; ff-wiki/proposals/cost-aware-verification-planning/issues/04-stage-proof-scope-and-execution-record.md; ff-wiki/proposals/cost-aware-verification-planning/evidence/04-implementation.md。

## Done and verify

- 阶段证明门禁的适用条件与三种正负例在 shared `### Stage proof before dispatch` 与 plan 相应句中可逐条核对：G 不以自身将产生的同一证明为入场条件、P 缺 G 证明不派发、普通 A→B 局部依赖不加门禁，缺证/失配的真实发布拒绝仍在；结构测试通过不替代独立 Round 0 的源核对（G/P/普通局部三组语义反例由独立 Round 0 核实源，02 另做模型实测）。
- 执行记录可溯源：`evidence/01-implementation.md` 与 01 票 Fix 11–13 段保留两批 `make`/全内包/lint 与四个隔离 `init`/`verify` 的真实命令与结果，S1/S2 快照标识清楚、第一批缺失 hash 如实声明且不作当前验收，最终 S2 对应 `c990efb`；不再出现“本轮各一次/唯一一次”“执行前”“全部 edit/write”的失实表述，已写明 restatement/preflight 事后回填与 bash/Python 写入的流程偏离；未删除任何早批次、未倒填执行顺序、未以补跑冒充历史。
- 完整内部 Go 一次：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` — exit 0；240 PASS / 0 FAIL / 3 项既有环境 SKIP 据实报告（`catalog_test.go:531` 缺 Tangram compatibility checkout；`upgrade_test.go:219,347` Darwin 无 `/proc/self/exe`），不据 SKIP 声称行为已测。
- 源→嵌入→隔离部署：`make dev` 一次后 `git diff --no-index --stat assets internal/command/assets` 无输出；新二进制在默认 `ff-wiki` 与预写 `docs_dir: nested/wiki` 两个全新临时目录执行 `init` 与 `assets verify --json`，七面逐字节一致，project-root shared authority 指针与 `Verification schedule`/`Effective verification entry`/`Execution and review loop`/`Stage proof before dispatch` 标题解析实际通过。
- 修改文件 lint 门禁（设计 rev3）：`/tmp/flowforge-verification-tools/golangci-lint-2.15.0/golangci-lint run --timeout 2m --new-from-rev 16cce8e --whole-files --max-issues-per-linter 0 --max-same-issues 0 ./cmd/... ./internal/...` — 复用 `c990efb` 的整文件 lint 结果，条件为按实际差分核对 Go 源码/测试/模块依赖/lint 配置/工具版本未变；条件变化才扩大 lint，不得为“免把旧写成新”擅自再跑。全仓既有失败按 [实施前lint基线](../evidence/01-lint-baseline.md) 保留，不宣称全仓通过。
- `git diff --check`、`flowforge check --dir ff-wiki/proposals`、`flowforge frontier --json` — 记录实际结果；诊断计数不回退，本票不新增 gap。
- 逐项语义与记录独立审查：独立 Round 0 覆盖完整 01 范围 `16cce8e...HEAD`（不是两句 delta），零 gap 后才双轴；01 的 Changes/Constraints/Done 与记录溯源一并核对。自动测试通过不替代此项。

### Verification schedule

本安排随 04 合同发布。A04/A05 修复与记录恢复为同一批次：全部源修改与记录修正完成后执行一次 `make dev`，同步 `internal/command/assets/**` 与 `bin/flowforge`；随后一次完整内部 Go `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` 覆盖票内定向集合，240 PASS / 0 FAIL / 3 项既有环境 SKIP 据实报告；两种 docs_dir（默认 `ff-wiki`、预写 `nested/wiki`）各做一次全新临时目录 `init` + `assets verify --json` + 七面 `cmp` 与 pointer/标题解析；`git diff --check`、`flowforge check`、`flowforge frontier` 记录结果。A05 的记录溯源须在当轮按 `history-evidence.log` 与原始 transcript 逐项核对，两批历史与当轮 04 运行分别计数，不删除任何早批次；任何源在验证后再改必须记录新快照和必要复跑。拟复用 `c990efb` 的三个 Go 文件 G0 来源/内容与整文件 lint 结果，条件是按实际差分核对 Go/依赖/lint 配置/工具版本未变；条件变化则实跑相应检查，执行者不得豁免。三组语义反例由独立 Round 0 核实源，02 另做模型实测；文件语义检查不等于模型实跑。历史 lint 与本轮运行分开表述。

## Execution detail

### Verified contracts

- A04 的实际句子：`assets/skills/_shared/ARTIFACT-CONTRACT.md:113` 为 “A ticket that consumes another ticket's delivery, or a release or real use that needs the whole, requires positive stage proof.”；`assets/skills/flowforge-plan/SKILL.md:28` 为 “A ticket that consumes another delivery, or needs the whole before a release or real use, also writes the positive stage proof it depends on ...”。这两句把门禁扩大到所有上游交付消费者。
- 权威：`design.md:50–64` 规定局部票闭合自身契约、阶段验收票 G 消费局部交付并产生组合证明、需要完整证明的发布/实际使用票 P 才依赖 G；`:64` 的前提是“派发需要阶段证明的消费/发布票前”。01 Change 12（`issues/01-rule-package.md:99`）保留“一切需要阶段正向证明的消费/发布”的限制。共享 `:105` 把 stage proof 定义为多个独立交付的组合证明，不能重解释为单票 local proof。
- A05 的直接证据：`/Users/qiangbi/.pi/agent/sessions/--Users-qiangbi-develop-projects-Syl-tangram-flowforge--/subagent-artifacts/outputs/862b108c-74be-4c8a-b6ce-4f429752919d/rule-package-delivery.audit-01-1.history-evidence.log`；原始 `61b27f54-eaed-4a1d-b658-3e23df77b157_flowforge-implementer_transcript.jsonl` 的 `:188–202`（首批 make/全内包/lint，command 9.850s、tracker 1.108s、240 PASS/0 FAIL/3 SKIP、`0 issues.`）、`:204–210`（首批两处 `/tmp/ff-fixA-7XHacA`、`/tmp/ff-fixB-ZwyS8a` 隔离部署）、`:255–265`（随后实际修改 implement pre-flight 标题与 Phase 0 句子）、`:275–289`（第二批，command 9.902s、tracker 0.765s、lint `0 issues.`）、`:299–313`（bash/Python 截写证据段与票内勾选，之后才插入标成“执行前”的 restatement/preflight）。
- `evidence/01-implementation.md:112` 标题“实际命令与结果（本轮各一次）”、`:135`“唯一一次全内包运行”；`issues/01-rule-package.md:229,235` 把 restatement/preflight 标成“执行前”，`:237` 声称 proposal Markdown“全部 edit/write”。
- 三个 Go 测试在 `c990efb` 与 G0 逐字节相同（`503b574e…`、`7c3fcd13…`、`101e4219…`），本票对全部测试只读。

### Execution scenarios

- Success：shared 与 plan 的门禁句按设计所有者裁决收窄，G/P/普通局部三种正负例可逐条核对；记录修正后两批历史与 S1/S2 快照可追溯；`make dev` 后源=嵌入，当轮完整内部 Go、两处隔离部署与 pointer/标题解析、`diff --check`/`check`/`frontier` 均 exit 0；独立 Round 0 覆盖完整 `16cce8e...HEAD` 零 gap。
- Failure：若门禁收窄时删掉了真实发布拒绝路径（P 缺 G 证明仍派发），或把“单纯有上游交付”改成豁免真实入场条件，Independent Round 0 报 gap 并回实现者；不改验收凑绿。
- Failure：若记录修正把第一批当成当前验收、倒填执行顺序、补跑冒充历史或抹掉最终真实 PASS，属记录不诚实，回实现者重写记录。
- Failure：源/嵌入/部署出现真实差异时既有保护测试失败即停，不改断言。

### Expected tests

当轮实跑 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`（exit 0），实际 PASS 输出覆盖票内定向集合 `TestCostAwareRulePackageSourceEmbeddedDeployed`、`PackagedSkillPointersResolve`、`ImplementSkillCarriesWeakExecutorContract`、`ReviewSkillCarriesRound0Audit`、`ReviewSkillDescribesRepairEscalation`、`IssueTrackerDescribesRepairWorkflow`、`PlanPublishesExecutionDetailSkeleton`、`PlanSkillDeclaresModeLine`、`ArtifactContractDocumentsExecutionContractSections`、`AgentRulesDescribeGenericCapabilityDispatch` 及 tracker 侧 `TestCostAwareSharedRunPreservesEvidenceGate|ExecutionContractCompletenessAppliesOnlyToManagedExecutableTickets|EvidenceQuadrupleParsing|NeedsRepairIsNonExecutableAndRepairProvenanceIsParsed`；不重复同固定点的定向套件。3 项既有环境 SKIP 如实列出，不称其行为已测。lint 按复用条件核对，条件变化才实跑。

### Generated artifacts

`make dev` 生成 `internal/command/assets/**`（嵌入副本，被 `.gitignore` 忽略）与 `bin/flowforge`（仅构建）；当轮四元组、两批历史对照与命令结果写入 `evidence/04-implementation.md`。不手写嵌入规则，不修改测试或 `assets/pi/**`。

### Conventions

- must 只写有独立信息的成本理由与复用说明，复用共同权威链接而非复制长模板 — [信息价值规则](../../../../assets/skills/_shared/ARTIFACT-CONTRACT.md#information-value)
- 记录求实：早批次、事后回填与工具偏离如实保留；不以补跑或不倒填的方式伪装顺序执行；A04 恢复范围以设计所有者裁决为界，不扩大豁免。[Conventions]

---

## Implementation note

### Phase 0 restatement（执行前）

- Change 1 → 收窄 `assets/skills/_shared/ARTIFACT-CONTRACT.md` 的 `### Stage proof before dispatch` 与 `assets/skills/flowforge-plan/SKILL.md` 的 `A ticket that consumes another delivery ...` 句：门禁只在需求/设计/票面点名某既有阶段正向证明为入场条件时适用，写出 G/P/普通 A→B 三种正负例并保留缺证/失配的真实发布拒绝；验收 = `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` 内的 `TestCostAwareRulePackageSourceEmbeddedDeployed`、`TestArtifactContractDocumentsExecutionContractSections`、`TestPlanPublishesExecutionDetailSkeleton`，源语义另由独立 Round 0 逐句核对。
- Change 2 → 按 `history-evidence.log` 与原始 transcript 逐条核对两批/S1S2/事后回填/bash 写入，改写 `evidence/01-implementation.md` 与 01 票 Fix 11–13 段的失实句（`本轮各一次`/`唯一一次`/`执行前`/`全部 edit/write`）并保留更正出处；验收 = 逐条对照日志行号可核对 + `git diff --check` + `flowforge check --dir ff-wiki/proposals` 计数不回退。
- Change 3 → 全部源改完后一次 `make dev`，跑一次全内包 Go、默认 `ff-wiki` 与预写 `nested/wiki` 两处隔离 `init`/`assets verify --json`/七面 `cmp`/pointer 与标题解析、`git diff --check`、`flowforge check`、`flowforge frontier`，四元组写入 `evidence/04-implementation.md`；验收 = 上述命令 exit 0 且 `git diff --no-index --stat assets internal/command/assets` 无输出。

### Phase 0b preflight（执行前）

- Constraints → 文件工具维护 proposal Markdown：四处写入全部用 `edit`/`write`，无 bash/Python/heredoc；不改 Issue Schema 头/CLI 签名：未动 frontmatter、`internal/**` 非生成 Go、`cmd/**`、命令与 linter 参数；测试作者/实现者分离与测试只读：三个 Go 测试零修改；旧结果不充新执行：lint 复用 `c990efb` 的条件逐项核对，条件未变即不重跑；反猜测：只修有失败依据的 A04/A05；assets 可部署一致性：两处源规则 + `make dev` 生成嵌入；不改 DAG/Catalog/四元组/默认 Mode/旧票：只改两处源句 + 指定历史证据 + 本票 + 本证据。
- Write set 合规：动手前列出将改文件 = 两处源规则、`evidence/01-implementation.md`、01 票 Fix 11–13 记录段、本票、`evidence/04-implementation.md`，另由 `make dev` 生成 `internal/command/assets/**` 与 `bin/flowforge`。
- Execution scenarios → Success（门禁收窄后 G/P/普通局部三例可核对、记录可溯源、当轮全内包与两处隔离部署 exit 0）：落点 Change 1/2/3；Failure（删掉真实发布拒绝或豁免真实入场条件）：Change 1 保留 dispatcher 拒绝段与独立入场条件句；Failure（记录把首批当当前验收/倒填顺序/抹掉最终 PASS）：Change 2 只改失实句并保留两批与 S2 的 PASS；Failure（源/嵌入/部署真实差异）：Change 3 的 `TestCostAwareRulePackageSourceEmbeddedDeployed` 失败即停，不改断言。
- 复用条件核对（lint）：`git diff 16cce8e..HEAD -- '*.go' 'go.mod' 'go.sum' '.golangci*'` 无输出且工具版本仍 `golangci-lint 2.15.0` → 满足未变条件，复用 `c990efb` 的整文件 `0 issues.`，不重跑。

### Phase 0 之后的执行结果（3d）

- 完成：Changes 1、2、3；未执行项：无。
- 当轮执行计数：`make dev` 1 次、全内包 Go 1 次、lint 0 次（按条件复用 `c990efb`）、两处隔离部署、`diff --check`/`check` 各 1 次。
- 修改文件：`assets/skills/_shared/ARTIFACT-CONTRACT.md`、`assets/skills/flowforge-plan/SKILL.md`、`ff-wiki/proposals/cost-aware-verification-planning/{evidence/01-implementation.md, issues/01-rule-package.md, issues/04-stage-proof-scope-and-execution-record.md, evidence/04-implementation.md}`；`internal/command/assets/**`、`bin/flowforge` 仅 `make dev` 生成。
- Write-set 合规：All modifications within write set。三个 Go 测试只读，`assets/pi/**` guard、保护与模型配置未动；未改 01 状态、未关票、未 push、未委派、未自审。
- 独立 Round 0 与双轴审查未在本轮执行（归主会话）：结构测试通过不替代三组语义反例的源核对与 02 的模型实测。
