---
flowforge:
  schema: 1
  role: ticket
  id: cost-aware-verification-planning-01
  revision: 1
  consumes:
    requirements:
      cost-aware-verification-planning-requirements: 2
    design:
      cost-aware-verification-planning-design: 3
---

# 01: 规则实现及必要保护性测试

**Blocked by:** 03
**Status:** open
**Mode:** lightweight

## Delivery

交付一致的成本感知验证规则包及必要保护性测试：有权角色分配验证、执行者按已批准安排机械操作、Review 独立判定证据复用；旧票继续原执行规则。本票证明规则包完整性、交付接缝和既有门禁兼容；模型受控遵从及本项目部署由后继验收交付。

## Design context

采用现有正文表达检查分组与 Fix verification，共享契约统一语义，五个技能分别落实职责。新安排必须完整明确；缺失走旧规则，残缺先阻塞。一次真实执行可明确覆盖多个 Change，新 Fix 必须有本轮实跑证据，旧结论由实际差分上的独立审查确认。阶段发布要求适用版本的正向证明，CLI 终态不替代它。

消费 [需求 revision 2](../requirements.md#cost-aware-verification-planning-requirements) 与 [设计 revision 3](../design.md#cost-aware-verification-planning-design)，执行依据为 [检查安排](../design.md#schedule)、[Fix 合同](../design.md#fix-contract)、[证据生命周期](../design.md#evidence)、[落点](../design.md#surfaces)。

## Touch points

- `assets/skills/_shared/ARTIFACT-CONTRACT.md` — Hand-offs / Execution and review loop / Diagnostics。
- `assets/skills/flowforge-solution-design/SKILL.md` — Maintain authority incrementally / Derive planning coverage。
- `assets/skills/flowforge-plan/SKILL.md` — Resolve effective authority / Draft tracer increments / Publish execution contracts。
- `assets/skills/flowforge-refine-ticket/SKILL.md` — Gather verified repository evidence / Fill the five execution-contract sections / Verify readiness。
- `assets/skills/flowforge-implement/SKILL.md` — Phase 0/0b、3a–3e、Full mode、Publish the next frontier。
- `assets/skills/flowforge-review/SKILL.md` — Pin fixed point / Round 0 / Aggregate / Fix planning / Review rounds / closeout。
- `assets/agents/issue-tracker.md` — Publishing / Execution and review loop / Deterministic checks。
- `internal/command/assets_deploy_test.go:TestIssueTrackerDescribesRepairWorkflow`；`internal/command/cost_aware_assets_test.go:TestCostAwareRulePackageSourceEmbeddedDeployed`、`internal/tracker/cost_aware_evidence_test.go:TestCostAwareSharedRunPreservesEvidenceGate`（[已应用测试补丁](../evidence/test-authoring.patch)，G0实跑见 [测试作者证据](../evidence/01-test-authoring.md)）。
- `Makefile:dev`、`internal/command/embed.go:embeddedAssets`、`internal/command/assets_deploy.go:deployManagedAssets` — 源复制、编译嵌入、部署接缝。

## Changes

- [ ] 1. 在 `assets/skills/_shared/ARTIFACT-CONTRACT.md` 的 hand-offs / execution loop 统一定义验证所有者、局部/阶段正向证明、完整 schedule 的覆盖/时机/命令/预期/记录、旧票兼容、共享实跑与历史复用区别；保存四元组，明确实际差分后的独立确认、失效后 repair 与发布门禁。其他技能只引用共同语义。
- [ ] 2. 在 `assets/skills/flowforge-solution-design/SKILL.md` 的验证策略检查纳入行为/风险→局部证明→阶段证明/触发表；真实环境若为本票承诺必需不能移走。增加重型验证必要性/共享/成本未知的说明，不增前置阶段。
- [ ] 3. 在 `assets/skills/flowforge-plan/SKILL.md` 的切票与发布流程纳入唯一阶段归属、真实开始依赖、发布需正向证明、纯验证票的非空记录及关联实现范围；在 Done and verify 给可选 Verification schedule 的完整要素及重型理由。保留 title/Delivery/依赖批准、Mode、五节空骨架及 standards 转抄。
- [ ] 4. 在 `assets/skills/flowforge-refine-ticket/SKILL.md` 五节核实步骤加入组范围、观察时机、真实命令/环境、成功失败反例与逐验收项覆盖；检查残缺安排及失效引用返回原规划方，不能猜值、删验收或修改保护测试。
- [ ] 5. 在 `assets/skills/flowforge-implement/SKILL.md` Phase 0/0b、3a–3e、Full mode 及 Publish frontier 协调三条路径：完整授权安排按组末验证后勾选；无安排保留原每 Change build 和全部检查；声明采用但不完整先 Blocked evidence 并返回。一次实跑可共享但不能记成多次；收口只确认完整有效覆盖，失败即停/重试限制不变；不得自行分组、选测、拟定复用或改测试。
- [ ] 6. 在 `assets/skills/flowforge-review/SKILL.md` Fix planning / Review rounds 发布针对固定点的 Fix verification（发现/编号/写集/关键反例、本轮实跑、拟复用原版本环境与覆盖、扩大条件）；Round 0 对实际差分检查新证据和旧适用性，零 gap 后独立双轴。保留关闭所有者及 design return/High-Critical/跨写集等 repair 边界；汇齐同固定点发现为机械批次，分别记录 review 调用/Fix批次/partial复发，禁止凭轮数豁免或扩张验收。
- [ ] 7. 在 `assets/agents/issue-tracker.md` Publishing、Execution and review loop、Deterministic checks 对齐新解释与旧票兼容，纠正旧 Tier 3 描述为现有五节，引用共享权威；不增加状态/Schema/CLI 能力，不将 frontier Ready 或阶段终态解释为正向验收。
- [x] 8. 由合法Plan/refine作者在G0通过已批准的实例授权，用正常文件工具将 `internal/command/assets_deploy_test.go:TestIssueTrackerDescribesRepairWorkflow` 的读取路径 `docs` 改为 `ff-wiki`；保留 `strings.Contains(string(data), "repair")` 与原失败消息。同时落地预设 `TestCostAwareRulePackageSourceEmbeddedDeployed` 和 `TestCostAwareSharedRunPreservesEvidenceGate`，记录实际运行与基线。路径维护有真实旧失败→修后通过；两个兼容回归可基线green，不能伪造red。后续轻量实现者对全部测试只读。G0仅执行本项；精确补丁、实例及撤权记录见 [测试作者证据](../evidence/01-test-authoring.md)。
    - cmd: GOPROXY=https://goproxy.cn,direct go test -v ./internal/...
    - exit: 0
    - output: "两个新增测试及旧路径修正通过；全部内包执行用例通过，0 failures，3项既有环境SKIP"
    - artifact: ff-wiki/proposals/cost-aware-verification-planning/evidence/01-test-authoring.md
- [ ] 9. 按 `Makefile:dev` 执行 `make dev`，从 `assets/**` 生成 `internal/command/assets/**` 并构建新二进制；将最终 diff 对照 [实施前漂移基线](../evidence/planning-baseline.md#源与嵌入漂移)。允许同步既有 PI feature 参数漂移，不修改 PI 源或引入无关差异。
- [ ] 10. 在 `evidence/01-implementation.md` 记录固定点（含未提交文件时完整快照）、实际命令/退出/输出、源→嵌入差分归因、每个 Change 的覆盖及未执行项；在本票 Implementation note 写真实结果后停止，由主会话收敛独立 Review，不自行关票。

## Constraints

- must 直接用文件工具维护 proposal Markdown，CLI 仅承担现有图与诊断计算 — [项目核心原则](../../../../AGENTS.md#核心设计原则)
- must not 改变 Issue Schema 头或 CLI 签名；本方案用现有正文及依赖表达 — [项目边界](../../../../AGENTS.md#boundaries)
- must 保持测试作者/轻量实现者分离、每 ticket 的执行上下文策略和独立审查 — [执行单元政策](../../../../AGENTS.md#execution-unit-policy)、[R1](../requirements.md#reliability)
- must not 把旧结果当新执行、把结构诊断通过当行为证明，证据仍须独立核实 — [证据契约](../../../../assets/skills/_shared/ARTIFACT-CONTRACT.md#hand-offs)、[审查交付核对](../../../../assets/skills/flowforge-review/references/review-protocol.md#test-pass-does-not-equal-requirement-met)
- must 只让有具体失败依据的发现形成修复义务，保留真实失败及已约定标准 — [审查反猜测规则](../../../../assets/skills/flowforge-review/references/review-protocol.md#no-speculation)
- must 不在 assets 放不部署的情景材料，按源→嵌入→部署接缝验证一致性 — [项目边界](../../../../AGENTS.md#boundaries)、[构建入口](../../../../Makefile)
- 轻量执行者不得更改任何预设测试；Change 8路径维护及两项新增测试全部由合法Plan/refine作者在G0先行完成，只有主会话向该实际实例签发这三个具体路径。宿主保护不得绕过。
- 不修改 DAG/Catalog 算法、evidence 四字段、默认 Mode 或旧票；不迁移 GIIS。保持 Round 0、独立双轴、零未处置发现关票与 substantive finding → repair。
- Write set: assets/skills/_shared/ARTIFACT-CONTRACT.md; assets/skills/flowforge-solution-design/SKILL.md; assets/skills/flowforge-plan/SKILL.md; assets/skills/flowforge-refine-ticket/SKILL.md; assets/skills/flowforge-implement/SKILL.md; assets/skills/flowforge-review/SKILL.md; assets/agents/issue-tracker.md; internal/command/assets/**（仅 make dev 生成）; bin/flowforge（仅构建）; internal/command/assets_deploy_test.go（仅合法测试作者执行Change 8；轻量实现者只读）; internal/command/cost_aware_assets_test.go 与 internal/tracker/cost_aware_evidence_test.go（仅合法测试作者）；ff-wiki/proposals/cost-aware-verification-planning/issues/01-rule-package.md; ff-wiki/proposals/cost-aware-verification-planning/evidence/01-implementation.md。

## Done and verify

- 规则包源/嵌入/隔离部署逐字节一致且原骨架、Mode、链接、执行保护仍在：`go test ./internal/command -run 'Test(CostAwareRulePackageSourceEmbeddedDeployed|PackagedSkillPointersResolve|PlanPublishesExecutionDetailSkeleton|PlanSkillDeclaresModeLine|ArtifactContractDocumentsExecutionContractSections|AgentRulesDescribeGenericCapabilityDispatch|ImplementSkillCarriesWeakExecutorContract|ReviewSkillCarriesRound0Audit|ReviewSkillDescribesRepairEscalation|IssueTrackerDescribesRepairWorkflow)$' -count=1` — 全部通过，0 failures；结构测试不充当模型行为证明。
- 分组共享四元组兼容、缺新 Fix 证据仍报错，旧票与 Catalog/依赖/journal 门禁无回退：`go test ./internal/tracker -run 'Test(CostAwareSharedRunPreservesEvidenceGate|ExecutionContractCompletenessAppliesOnlyToManagedExecutableTickets|EvidenceQuadrupleParsing|NeedsRepairIsNonExecutableAndRepairProvenanceIsParsed)$' -count=1` — 全部通过；完整内包检查包含该定向集合，不另重复执行同一固定点的相同命令。
- 完整回归：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` — 全部通过。
- 修改文件lint门禁（设计rev3明确）：`/tmp/flowforge-verification-tools/golangci-lint-2.15.0/golangci-lint run --timeout 2m --new-from-rev 16cce8e --whole-files --max-issues-per-linter 0 --max-same-issues 0 ./cmd/... ./internal/...` — exit 0。全部新增/修改Go文件必须已跟踪并在该包范围，整文件纳入，不禁用linter、不豁免新失败。工具缺失或新失败立即阻塞。全仓/全产品包的既有失败如实保留于 [实施前lint基线](../evidence/01-lint-baseline.md)，不宣称全仓lint通过。
- 规则逐项语义覆盖：独立 Review 对照 Changes 1–7、设计责任表、三条兼容路径、Fix/repair 分界和四元组诚实性；Round 0 零 gap 后独立双轴零未处置发现。自动测试通过不替代此项。

### Verification schedule

本安排由主会话设计/规划所有者在实施前明确授权。G0由合法Plan/refine作者执行Change 8及两项预设测试，记录实际结果并撤销实例授权，再完成Refine；该项可用自身四元组勾选。G1覆盖Changes 1–7、9，消费G0证据：先完成全部互相引用的源规则，再执行Change 9的 `make dev`，一次运行共享记录；组内不对半完成规则包构建或勾选。随后运行本节定向规则包检查和 tracker 定向检查。G2 覆盖全部 Changes：规则包定向通过后运行一次完整内包检查与 lint，保存于 `evidence/01-implementation.md`；Change 10 随记录完成。独立审查前确认每项已有适用证据，不以收口名义重跑相同命令。没有合法作者落地和实跑记录时 G1 不开始。

全内包检查复用既有工具，保护跨部署/图计算兼容；构建与嵌入同步一次共享。真实模型情景及项目部署归 02，因为它们必须消费完整包；本票不承诺其结果。成本测量未知，不估算时长或提速。

---

## Execution detail

### Verified contracts

- `Makefile:dev` 删除并复制整个 `internal/command/assets` 后构建 `bin/flowforge`；`internal/command/embed.go:embeddedAssets` 编译嵌入。原4文件漂移及哈希已固定在 [规划基线](../evidence/planning-baseline.md#源与嵌入漂移)，00/03已消除该漂移；G0资产兼容回归七子项通过，但Changes 1–7/9的源修改与生成同步尚未执行。
- `internal/command/assets_deploy.go:locateAssetsDir` 优先提取嵌入资产；`deployManagedAssets` 覆盖 skills、保留既有 agent docs。必须分别证明源→嵌入和 binary→deployed。
- 已解决的旧路径失败：规划时 `TestIssueTrackerDescribesRepairWorkflow` 实跑失败为 `open ../../docs/agents/issue-tracker.md: no such file or directory`。G0已验证正确读取路径为 `ff-wiki/agents/issue-tracker.md`，保留 `strings.Contains(string(data), "repair")` 和 `t.Fatal("issue-tracker.md must describe the repair workflow")`；修后精确测试 exit 0。
- G0合法入口已由主会话按真实 session/instance 签发2h `role=plan` 三路径授权恢复，执行后由主会话撤销；正常 `flowforge_test_authority` 最终返回 `grant.status=missing`。宿主保护未关闭，未以bash或其他工具写测试。精确补丁经 `git apply --check --reverse` exit 0确认完整应用，`gofmt -d` 无输出；hash、授权及撤销详见 [测试作者证据](../evidence/01-test-authoring.md)。
- 两项新增测试现已落地并实际运行，均 baseline green，不存在它们落地前的red。`TestCostAwareRulePackageSourceEmbeddedDeployed` 仅证明七面源/嵌入/隔离部署逐字节一致；`TestCostAwareSharedRunPreservesEvidenceGate` 证明共享四元组格式兼容、有/无schedule均不豁免新checked Fix的evidence-missing。它们不证明未写源规则的语义合规或真实模型行为。
- 既有 `TestExecutionContractCompletenessAppliesOnlyToManagedExecutableTickets`、`TestEvidenceQuadrupleParsing`、`TestNeedsRepairIsNonExecutableAndRepairProvenanceIsParsed` 与新增tracker用例的合并定向命令通过；完整内包命令 exit 0、0 failures，但有3项既有环境SKIP，具体原因见G0证据。

原 Blocked evidence 作为已解决历史保留如下；其恢复条件已满足，不再代表当前阻塞：

> - 实际操作：`functions.write` 创建 `internal/command/cost_aware_assets_test.go`。工具拒绝，无进程 exit code：`protected test file (matched **/*_test.go, source: FlowForge defaults): preset acceptance tests are authored by Plan; if the task requires editing it, report STATUS: BLOCKED instead`。
> - 后续事实：主会话核实保护没有合法作者例外，禁止关闭配置或通过其他通道写入；仅保存 [未应用测试补丁](../evidence/test-authoring.patch)。`git apply --check ff-wiki/proposals/cost-aware-verification-planning/evidence/test-authoring.patch` exit=0，仅校验适用性，没有应用。
> - 尚缺：合法作者落地新测试、实际 red/green 记录以及旧路径修正。下一步由主会话裁决合法作者执行入口；问题解决后才消费本段、完成 Refine 并重新计算 frontier。

### Execution scenarios

- Success（G0已实跑）：合法作者取得本实例受限授权，正常write/edit落地三路径，新测试实际匹配执行且green，旧路径修正后green，随后主会话撤销授权。后续G1/G2仍须按Verification schedule完成源规则、一次make dev、适用定向检查与独立审查，不能把G0扩大解释为规则语义完成。
- Failure（历史已解决）：合法入口出现前，正常write曾被保护拒绝；当时停止测试变更并保存可审阅补丁，未关闭保护或换工具绕过。现由00/03的合法作者入口和实例授权解决，历史原文保留在Verified contracts。
- Failure（后续执行边界）：源规则、嵌入或部署出现真实差异，或新checked Fix缺四元组时，相应新增保护测试应失败；实际失败仍按既有即停/重试限制处理，不修改预设断言凑绿。新测试baseline green不制造这种失败，也不模拟02真实模型情景。

### Expected tests

G0已实际执行且exit 0：两个新增测试的精确regex、旧路径维护测试、Done and verify中的command/tracker定向命令，以及 `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`。精确命令、七个资产子项、tracker四组合、3项SKIP与实际耗时见 [G0测试作者证据](../evidence/01-test-authoring.md#实际测试)；不为Refine或收口重复运行相同命令。

指定临时工具已实际执行 `/tmp/flowforge-verification-tools/golangci-lint-2.15.0/golangci-lint run --timeout 2m --new-from-rev 16cce8e --whole-files --max-issues-per-linter 0 --max-same-issues 0 ./cmd/... ./internal/...`，三个Go文件已跟踪，exit 0、输出 `0 issues.`。这是G0修改文件门禁，不是全仓lint通过；[既有失败基线](../evidence/01-lint-baseline.md)不变。

后续轻量实现者消费G0有效证据，不重复相同固定点检查；在Changes 1–7/9实际改变源/嵌入固定点后，按Verification schedule运行本票Done and verify的适用检查。G1/G2证据写入 `evidence/01-implementation.md`。01结构/兼容性通过不替代逐项语义独立Review；真实模型行为仍归02。

### Generated artifacts

G0产生两个新Go测试及 [01-test-authoring.md](../evidence/01-test-authoring.md)，并修改旧路径测试；没有生成嵌入资产或二进制。后续 `assets/**` → `make dev` → `internal/command/assets/**` → `embeddedAssets` / `bin/flowforge` → `deployManagedAssets` → 隔离 `.agents/skills` 和文档目录仍归Changes 9/10。01只同步源包与嵌入并以现有部署seam验证，真实项目部署归02；不得手写嵌入规则，最终diff需与原四文件漂移基线区分。

### Conventions

- must 只写有独立信息的成本理由与复用说明，复用共同权威链接而非复制长模板 — [信息价值规则](../../../../assets/skills/_shared/ARTIFACT-CONTRACT.md#information-value)
- G0只授权精确三测试路径；授权撤销后测试只读，Refine不得借文档维护再改测试。兼容回归允许baseline green，不制造red；全仓lint旧失败与G0指定lint通过分开表述。[Conventions]

## Implementation note

## Review rounds
