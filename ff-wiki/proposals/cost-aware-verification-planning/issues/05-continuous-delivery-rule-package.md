---
flowforge:
  schema: 1
  role: ticket
  id: cost-aware-verification-planning-05
  revision: 1
  consumes:
    requirements:
      cost-aware-verification-planning-requirements: 3
    design:
      cost-aware-verification-planning-design: 4
---

# 05: 连续实现、集中验证及角色职责一致的规则包

**Blocked by:** 04
**Status:** open
**Mode:** lightweight

## Delivery

让新Plan默认组织完整逻辑与集中验证，实施/审查/角色入口一致消费这项安排，减少全文复述和无必要的记录交接；保留旧票兼容、失败停止、测试保护及独立审查。整包一次交付，不按Skill拆成独立验收票。

## Design context

消费[需求rev3](../requirements.md#cost-aware-verification-planning-requirements)和[设计rev4](../design.md#cost-aware-verification-planning-design)。[保留/替换裁决](../design.md#migration)明确不整段回滚；本票是新目标的增量，不把01/04原关闭改写成新规则已通过。依赖04是因为其正确阶段门禁和已修复基线必须保留；04已经closed，但本票仍须Refine和合法测试作者准备后才能执行。

主会话固定当前执行Skill入口与本票显式覆盖：下述G1取代旧部署Skill的逐Change检查/双重全文复述要求；其他安全边界继续生效。待修改的assets不自动成为本轮执行权威。

## Touch points

- `assets/skills/flowforge-plan/SKILL.md:Draft tracer increments/Publish execution contracts`、`flowforge-solution-design/SKILL.md:Verification strategy check` — 完整逻辑/关口理由。
- `assets/skills/_shared/ARTIFACT-CONTRACT.md:Execution and review loop/Verification schedule/Shared runs and reused evidence` — 共同语义。
- `assets/skills/flowforge-refine-ticket/SKILL.md:Gather evidence/Fill five sections`、`flowforge-implement/SKILL.md:Phase 0/0b/3a–3e` — 核实与连续执行。
- `assets/skills/flowforge-review/SKILL.md:Invocation ownership/Round 0/Fix planning`；`assets/subagents/flowforge-{planner,implementer,reviewer,reviewer-lite}.md:Identity/Default Skill/Workflow Position` — owner与只读范围。
- `assets/AGENTS.md:Subagent delegation/Execution unit policy`、`assets/agents/issue-tracker.md:Publishing/Execution and review loop` — 调度与写回。
- `internal/command/{cost_aware_assets_test.go,assets_deploy_test.go}` — 合法作者调整的预设断言；`Makefile:dev` — 嵌入生成。

## Changes

- [ ] 1. 合法Plan/Refine作者保留原安全断言，在两项command测试接缝准备新Plan默认集中检查、关键前置仍保留、单次执行确认、reviewer只读且不再委派/写票的断言；扩展源/嵌入/隔离部署覆盖到本票角色和AGENTS面。仅替换被设计明确改变的旧文案断言，记录对应关系与实际基线，冻结测试后撤销授权。此项不是轻量实施者工作。
- [ ] 2. 按[Plan与检查安排](../design.md#allocation)修改shared、solution-design和plan：新票默认一组完整逻辑和紧凑schedule；切分/提前关口需具体理由；旧票无安排仍旧规则，声明残缺仍返回owner。保持Fix入口优先级和04已修正门禁。
- [ ] 3. 按[证据与执行确认](../design.md#evidence)修改refine/implement：五节仍核实，已确认且未变部分精确引用；Phase0/0b合并成一次当前组/关键约束确认。组末统一检查后才勾选；四元组和失败上限不变，记录只补实际结果/例外。
- [ ] 4. 按[责任接口](../design.md#owners)及[修复范围](../design.md#fix-contract)统一review、四个subagent源、assets/AGENTS和tracker：主会话owner编排/裁决/写回，只读审计和双轴子代理只返回发现；保留Round0先于双轴、repair与作者分离。完整流程与叶子调用明确分支，不新增角色/权限、不为机械回填派scribe。
- [ ] 5. 全部源规则完成后，按G1集中构建/检查并记录源、嵌入、隔离部署的完整性和测试冻结状态；提交固定差分与非空结果交独立Round0/双轴。未运行02真实行为链，不部署本项目、不宣称净提速。

## Constraints

- must 直接用文件工具维护proposal，CLI只承担图与诊断计算 — [项目原则](../../../../AGENTS.md#核心设计原则)
- must not 修改Issue Schema头规范、CLI签名或现有诊断算法 — [项目边界](../../../../AGENTS.md#boundaries)、[R8](../requirements.md#compatibility)
- must 保留测试作者/实施者分离、独立Round0/双轴和真实缺陷阻塞 — [R1](../requirements.md#reliability)、[执行单元政策](../../../../AGENTS.md#execution-unit-policy)
- must 保留失败原文与证据归属；历史通过不当新实跑，结构通过不当行为通过 — [R6](../requirements.md#evidence)
- must 让只读子代理仅返回发现，由主会话承担裁决及工件写回；短机械记录不另派角色 — [R7](../requirements.md#review)
- must 不在assets放不部署的fixture；经源、嵌入、隔离部署核对后再发布 — [项目边界](../../../../AGENTS.md#boundaries)、[Makefile](../../../../Makefile)
- Write set: assets/skills/_shared/ARTIFACT-CONTRACT.md; assets/skills/flowforge-solution-design/SKILL.md; assets/skills/flowforge-plan/SKILL.md; assets/skills/flowforge-refine-ticket/SKILL.md; assets/skills/flowforge-implement/SKILL.md; assets/skills/flowforge-review/SKILL.md; assets/subagents/flowforge-planner.md; assets/subagents/flowforge-implementer.md; assets/subagents/flowforge-reviewer.md; assets/subagents/flowforge-reviewer-lite.md; assets/AGENTS.md; assets/agents/issue-tracker.md; internal/command/cost_aware_assets_test.go与internal/command/assets_deploy_test.go（仅Change1合法作者，实施者只读）；internal/command/assets/**、bin/flowforge（仅make dev生成）；本票；ff-wiki/proposals/cost-aware-verification-planning/evidence/05-implementation.md。
- 不修改权限profile、模型配置、PI保护、tracker算法或历史票结果；修改默认文案不允许删除安全断言获得绿灯。

## Done and verify

- 新默认与旧兼容在所有入口一致，reviewer工具/职责匹配；原保护语义未弱化。由独立Round0/双轴对照设计各分支判定，不以关键词或Go通过代替。
- 构建及源→嵌入：`make dev`，随后`diff -qr assets internal/command/assets` — exit 0；source/embedded/new isolated deployment在`TestCostAwareRulePackageSourceEmbeddedDeployed`及作者补充的角色部署断言中一致。
- 集中自动验证：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` — exit 0、无失败；环境SKIP原样披露，所需新增断言必须实际执行，未运行不能以SKIP通过。
- 合法作者修改Go后的整文件lint：`/tmp/flowforge-verification-tools/golangci-lint-2.15.0/golangci-lint run --timeout 2m --new-from-rev 16cce8e --whole-files --max-issues-per-linter 0 --max-same-issues 0 ./cmd/... ./internal/...` — exit 0、0 issues；工具路径和版本由Refine确认，不在实施中临时降级。新增Go须先跟踪，历史全仓lint失败继续保留。
- `git diff --check`及`bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning` — 无本票引入的结构错误。02待Refine和历史精确waiver具名保留，不为消诊断改成假完成。
- `evidence/05-implementation.md`记录新固定点、作者冻结/撤权、组实际执行、源/嵌入/隔离部署范围及偏离；独立审查通过后主会话按既定责任关票。

### Verification schedule

G0仅合法作者完成Change1并Refine，预设测试实际基线结果和写权撤销齐备后开始G1。G1覆盖Changes2–5：先连续改完所有互相关联的源文件，随后一次`make dev`、完整内部Go、整文件lint及图/差分检查；同一固定点不再另跑完整Go已覆盖的定向子集。Change1及G1共享结果按实际覆盖引用，不把同一次检查记为多次。所有结果进入同一`evidence/05-implementation.md`。真实失败停止返回；没有新差分/失败不为收口重复跑。

作者可用现有 `go test ./internal/command -run 'Test(CostAware.*|ImplementSkillCarriesWeakExecutorContract|ImplementerPromptPinsLoopContracts|BlockAndRecordContractPinned|ReviewSkillCarriesRound0Audit|ReviewSkillDescribesRepairEscalation|AgentsBlockContainsExecutionUnitPolicy|SubagentSource.*)$' -count=1` 定位预设断言；这不是G1额外验收。完整源码修改不得提前分成每Skill验收。

---

## Execution detail

### Verified contracts

### Execution scenarios

### Expected tests

### Generated artifacts

### Conventions

- must 将共同语义保留在单一权威处，其他入口精确引用，删除无新信息的复述 — [信息价值](../../../../assets/skills/_shared/ARTIFACT-CONTRACT.md#information-value)

## Implementation note

## Review rounds
