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

- [x] 1. 合法Plan/Refine作者保留原安全断言，在两项command测试接缝准备新Plan默认集中检查、关键前置仍保留、单次执行确认、reviewer只读且不再委派/写票的断言；扩展源/嵌入/隔离部署覆盖到本票角色和AGENTS面。仅替换被设计明确改变的旧文案断言，记录对应关系与实际基线，冻结测试后撤销授权。此项不是轻量实施者工作。
  - cmd: `/tmp/flowforge-verification-tools/golangci-lint-2.15.0/golangci-lint run --timeout 2m --new-from-rev 16cce8e --whole-files --max-issues-per-linter 0 --max-same-issues 0 ./cmd/... ./internal/...`
  - exit: 0
  - output: `0 issues.`
  - artifact: `ff-wiki/proposals/cost-aware-verification-planning/evidence/05-implementation.md`
  - 作者验证归属：以上为最终冻结测试的整文件lint实跑，并非新行为验收通过。最终定向作者命令实际exit 1：三项新语义测试族预期red；同次安全及部署检查pass，无SKIP。已记录旧断言对应、全部失败及冻结哈希，主会话明确撤权后只读查询为missing。此勾选仅完成G0作者交付；Changes2–5及整票仍未完成。
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

- 消费requirements rev3、design rev4，基线`08e87cb19001aef9b70a0bb047dee7b08b76003e`；04已closed，初始`check`/`frontier`确认05只缺本五节，02仍依赖05。没有待消费的Blocked evidence。主会话固定的本票G0/G1覆盖当前部署Skill的旧逐Change/全文复述规则，候选assets不成为本会话入口。
- `assets/skills/_shared/ARTIFACT-CONTRACT.md:Execution and review loop/Verification schedule/Effective verification entry/Stage proof before dispatch`承载共同规则。现有Fix verification→ticket schedule→旧逐Change的优先序、需要阶段证明才门禁、产证票不自等、普通局部依赖不额外门禁均保留；按design#allocation/#schedule改新Plan默认，不改tracker算法。
- `assets/skills/flowforge-plan/SKILL.md:Draft tracer increments/Publish execution contracts`当前仍称Optional schedule；`flowforge-solution-design/SKILL.md:Verification strategy check`当前有“add no early verification stage”。按design#allocation补完整逻辑、拆分理由及具体风险前置理由；不可把后一旧句保留为无条件禁止前置。
- `assets/skills/flowforge-implement/SKILL.md:Phase 0/Phase 0b/3a–3e`当前分别要求逐Change复述和全约束/场景清单；按design#evidence合为一次本组/检查引用与关键失败/写集/测试权限确认，保留四元组、失败原文、2次/5轮上限。`flowforge-refine-ticket/SKILL.md:Gather verified repository evidence/Fill the five execution-contract sections`保留真实五节，未变且版本/范围相同的事实用精确引用。
- `assets/skills/flowforge-review/SKILL.md:Pin the fixed point and scope/Round 0/Spawn both sub-agents/Fix planning`当前只读分支仍会走派发；`assets/subagents/flowforge-reviewer.md:Default Skill/Workflow Position`仍要求派发双轴及Fix写票，reviewer-lite仍称flagship唯一closeout。按design#owners/#fix-contract由主会话owner编排/裁决/写回，叶子只做指定审计/轴并返回，默认reviewer做Spec；同步四角色、AGENTS和tracker，权限profile/default Skill/角色集合不变。
- `internal/command/assets_deploy.go:locateAssetsDir/deployManagedAssets/applyAgentsBlock`优先embedded，skills复制到`.agents/skills`，tracker到配置docsRoot，AGENTS源内容trim后包管理标记；`internal/command/agents.go:discoverSubagentSources/deploySubagents`另行编译角色。隔离测试使用`initializeTestProject(t.TempDir())`及默认config；不消费本项目模型配置。
- 已核实`Makefile:dev`先复制assets至internal/command/assets再构建bin/flowforge；Go为1.26.0 darwin/arm64，票列lint绝对路径存在且版本2.15.0（Go1.27.0构建）。最终作者整文件lint exit0，两个Go文件本来已跟踪。完整G1尚未运行，历史全仓lint债只引用[原记录](../evidence/01-lint-baseline.md)，不改写为绿。
- G0测试已冻结且主会话已撤权；实际session/instance、哈希、原始命令/exit/输出与旧断言对应均在[05作者证据](../evidence/05-implementation.md#g0-author)。G1仅消费冻结测试，连续完成Changes2–5后执行上方集中检查；本票Write set已核实与此范围一致，不额外新增写集。

### Execution scenarios

- 成功：普通接口/调用方关联逻辑在新Plan中组成默认一个连续组，紧凑schedule引用已有命令/预期和结果位置；Implement一次确认关键约束，组内不插入未安排验收，组末检查通过后才勾选。`TestCostAwareContinuousDeliveryContracts`钉住必要条款，独立语义审查判断所有入口无相反指令；真实轨迹由02承担。
- 必要前置：已确定的关键假设/不可逆使用风险保留有具体理由的提前检查；design与shared不能一律禁止提前。`TestCostAwareContinuousDeliveryContracts`检查相关条款；02 B2才观察失败后停止且后续动作不发生，本票不冒称已运行。
- 兼容/失败：旧票无schedule仍旧规则；声明残缺/矛盾/命令不存在先返回owner；本轮Fix优先消费Fix verification。阶段证明缺失/过期/版本不明时，不因结构通过而放行；未改的tracker兼容测试在G1完整Go中执行，语义审查对照design#schedule/#allocation。
- 角色：owner明确调度独立Round0→零gap→双轴并写回；只读reviewer默认Spec、reviewer-lite做Standards，禁止再派发/写文件/关票；必要输入缺失返回缺口。四角色文案由`TestCostAwareRoleContracts`覆盖，PI白名单/Codex只读及四宿主正文传递由部署测试覆盖。
- 生成失败：任一源与embedded不同、管理块内容不同、任一角色宿主正文丢节或只读权限漂移均使部署测试失败；不得编辑测试或手改embedded获得绿灯，按G1顺序生成，真实失败返回主会话。

### Expected tests

- 作者命令采用上方Verification schedule中给出的定向`go test ./internal/command -run ... -count=1`；冻结前最终加`-v`观察各case。`TestImplementSkillCarriesWeakExecutorContract`、`TestCostAwareContinuousDeliveryContracts`、`TestCostAwareRoleContracts`在未实现的源上实际red/exit1，匹配design rev4替换点；未来G1应全部pass。这不是作者失败，也不是产品通过。
- 旧`Restate before executing`及全场景复述文字pin已替换为一次execution confirmation、关键约束/Write set/受保护测试；其余失败停止、同次四元组、Round0、repair、安全原文断言原样保留。条款匹配忽略空白/大小写并接受近义词，精确拒绝仅限已确定相反的旧指令；不能用匹配通过代替语义审查。
- `TestCostAwareRulePackageSourceEmbeddedDeployed`现覆盖6份skill/shared、tracker及AGENTS管理块；`TestCostAwareRolesSourceEmbeddedDeployed`现覆盖4角色×Claude/OpenCode/Codex/PI，Codex允许既有Default Skill调用改写，其余段落保持。reviewer/reviewer-lite实际PI工具必须为read/grep/find/ls、Codex sandbox必须read-only。作者实跑全部pass、无SKIP。
- G1在所有源面完成后一次执行Done and verify所列`make dev`、`diff -qr assets internal/command/assets`、完整内部Go、2.15.0整文件lint、`git diff --check`及proposal check，预期exit0；完整Go已覆盖作者子集，不追加重复子集。结构检查保留02缺五节及历史waiver；任何新增真实失败返回，不自行改命令/跳过失败。

### Generated artifacts

- `assets/**` → `make dev`复制 → `internal/command/assets/**` → `internal/command/embed.go:embeddedAssets`嵌入 → `bin/flowforge`；源/嵌入用`diff -qr`及包测试字节比较核对。此G0不生成这些产物。
- embedded → `deployManagedAssets` → 临时目录的`.agents/skills/**`、配置docsRoot/agents/issue-tracker.md和`AGENTS.md`管理块 → Skill/项目入口；新隔离目录对齐为本票验收，定制文档保留行为不改。
- embedded/subagents → `discoverSubagentSources`/`deploySubagents`及`internal/subagent/compile_{claude,opencode,codex,pi}.go` → 临时四宿主角色文件 → 宿主注入；新测试比较源/embedded及正文传播/只读权限。不把源md假设为直接部署到`.agents/subagents`。
- G0/G1实际记录共用`evidence/05-implementation.md`。本项目`.agents`/宿主文件部署与真实角色B1/B2链属于02，05只用临时测试目录，不触及旧未跟踪02-scenarios。

### Conventions

- must 将共同语义保留在单一权威处，其他入口精确引用，删除无新信息的复述 — [信息价值](../../../../assets/skills/_shared/ARTIFACT-CONTRACT.md#information-value)
- 共同语义保留shared已有锚点，Plan/Refine/Implement/Review和tracker精确链接；`TestPackagedSkillPointersResolve`会在G1完整Go中验证相对链接。角色frontmatter的permission/model_profile/default_skill不因正文职责调整而修改。
- 源测试从`internal/command`包工作目录用`../../assets`读取；隔离部署使用`t.TempDir()`。Markdown资产测试仅钉住条款/同步，完整解释由独立审查确认；模板/否定句即使含关键词也不证明遵从。
- Changes2–5按一组连续实施，不按文件分组；Implementation note只补本轮结果/例外并精确引用既有核实。保持本票open，主会话独立裁决/写回/关闭；不让轻量实施者自行选测或裁决证据复用。

## Implementation note

G0/Change1与Refine完成；仅两测试、本票和[05证据](../evidence/05-implementation.md)改变。最终作者定向检查exit1为预期新语义red，安全/部署case通过，整文件lint exit0；冻结及撤权已证实。Changes2–5/G1、独立审查及02均未执行。全部修改在本次限定写集内。

## Review rounds
