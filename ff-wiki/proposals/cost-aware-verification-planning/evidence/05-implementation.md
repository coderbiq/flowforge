# 05 实施证据

<a id="g0-author"></a>
## G0 / Change1：合法作者与Refine

本记录仅证明测试作者准备、基线运行与执行合同核实。源规则仍是08e87cb的旧规则；新行为预期red不是作者失败，也不是05产品通过。G1、独立Round0/双轴、02 B1/B2及本项目部署均未执行。

- 固定基线：08e87cb19001aef9b70a0bb047dee7b08b76003e；消费requirements revision 3、design revision 4。当前差分仅两测试、05票及本文件；旧未跟踪evidence/02-scenarios原样保留。
- 技能入口：已读取flowforge-plan默认入口及flowforge-refine-ticket；本任务不重规划，按主会话明确授权仅G0/Change1与五节Refine，不委派。
- 作者实际sessionId：01a1aea4-57cb-712f-b7e8-49d43512f568；instanceId：b017ad49-7304-4080-b846-d8897c32a398。
- projectRoot：/Users/qiangbi/develop/projects/Syl/tangram/flowforge；grantPath：/Users/qiangbi/develop/projects/Syl/tangram/flowforge/.flowforge/test-author-grants/01a1aea4-57cb-712f-b7e8-49d43512f568.json；eligible=true。
- 首次flowforge_test_authority返回missing。将实际身份发主会话need_decision后，主会话签发role=refine、expiresAt=2026-11-06T04:16:49.406Z、仅internal/command/cost_aware_assets_test.go及internal/command/assets_deploy_test.go；再次只读查询valid后才用正常edit/write写测试。未使用shell格式化写入、未修改grant/权限/模型。
- 写入完成后请求撤权；主会话明确回复“已核验并删除该实际session/instance的grant，文件不存在；两份测试SHA256与你报告完全匹配，冻结成立”。随后flowforge_test_authority再次返回missing / no grant record for this session。此后测试只读。

### 冻结范围与旧断言对应

| 旧检查 / 缺口 | 新检查与设计对应 | 保留边界 |
| --- | --- | --- |
| Implement的Restate before executing与Execution scenario (both Success and Failure)两个精确文字pin | TestImplementSkillCarriesWeakExecutorContract改为一次execution confirmation/check-in、关键约束、Write set、preset/protected tests，design#evidence | Mode、STATUS: BLOCKED、2次/5轮、exit:0、same edit及verbatim旧断言原样保留 |
| 新Plan仍可选schedule、双重全量复述 | TestCostAwareContinuousDeliveryContracts要求新票默认完整组、切分理由、必要风险前置、准确引用未变事实，拒绝已确定相反旧指令，design#allocation/#schedule/#evidence | 无schedule旧票、残缺停止、Fix优先级、四元组、阶段证明条款保留 |
| reviewer派双轴/写票，lite称flagship唯一closeout | TestCostAwareRoleContracts及review/AGENTS/tracker条款要求主会话owner、只读leaf不派发/写文件/关票、reviewer默认Spec；design#owners/#fix-contract | 原Round0/repair/失败停止测试不删；reviewer/lite permission与default Skill明确锁定 |
| 源/embedded/隔离部署仅覆盖skills和tracker | TestCostAwareRulePackageSourceEmbeddedDeployed增加AGENTS管理块；TestCostAwareRolesSourceEmbeddedDeployed比较四角色源/embedded、四宿主正文及PI工具白名单/Codex只读 | 只在t.TempDir部署；保留原7项字节比较；Codex仅允许现有Default Skill调用改写 |

条款检查归一化大小写/空白并接受相关近义词；没有绑定全文、版式或无关修辞。它只证明条款存在，不解析自然语言的否定/矛盾，不取代独立语义审查或02真实轨迹。不得通过添加无效关键词满足规则；G1应按已定设计真实改写全部相关入口。两份测试已经冻结，轻量实施者不得调整断言。

最终SHA256（撤权前后主会话及只读命令核实一致）：

- internal/command/cost_aware_assets_test.go：17eca4b481e495fdb57cf5ad6206da39753582b9f7cf2e17825fba03b1b90033
- internal/command/assets_deploy_test.go：e60de2b8f618405479c03e88fab29efb33c5120918ce96e6aec17b4d735c687c

### 环境与命令边界

命令均在项目根执行。Go版本实际为go1.26.0 darwin/arm64。票列绝对路径执行version返回：golangci-lint has version 2.15.0 built with go1.27.0 from 25c5369a on 2026-10-16T11:03:09Z（exit0）。没有降级工具、修改lint选项或关闭保护。历史全仓lint失败保留在[01原记录](01-lint-baseline.md)，本次整文件范围通过不表示全仓通过。

G0没有运行make dev、diff -qr或完整内部Go，也未修改assets/embedded/生产代码。包测试触发的临时隔离部署是作者接缝验证，不是本项目部署。G1仍按05 Done and verify在所有源变更完成后一次集中执行，不为已覆盖的子集重复验收。

### A0 旧测试基线（改测试前）

原始命令：

```sh
go test ./internal/command -run 'Test(CostAware.*|ImplementSkillCarriesWeakExecutorContract|ImplementerPromptPinsLoopContracts|BlockAndRecordContractPinned|ReviewSkillCarriesRound0Audit|ReviewSkillDescribesRepairEscalation|AgentsBlockContainsExecutionUnitPolicy|SubagentSource.*)$' -count=1
```

实际exit: 0

```text
ok  	flowforge/internal/command	1.152s
```

### A1 第一版新断言：预期red

原始命令：

```sh
go test ./internal/command -run 'Test(CostAware.*|ImplementSkillCarriesWeakExecutorContract|ImplementerPromptPinsLoopContracts|BlockAndRecordContractPinned|ReviewSkillCarriesRound0Audit|ReviewSkillDescribesRepairEscalation|AgentsBlockContainsExecutionUnitPolicy|SubagentSource.*)$' -count=1
```

实际exit: 1

```text
--- FAIL: TestImplementSkillCarriesWeakExecutorContract (0.00s)
    assets_deploy_test.go:347: execution confirmation: missing contract clause /(one|single|once).{0,100}(execution confirmation|execution check-in)/
    assets_deploy_test.go:347: execution confirmation: missing contract clause /(preset|protected).{0,80}test/
--- FAIL: TestCostAwareContinuousDeliveryContracts (0.02s)
    --- FAIL: TestCostAwareContinuousDeliveryContracts/skills/flowforge-plan/SKILL.md (0.00s)
        cost_aware_assets_test.go:138: skills/flowforge-plan/SKILL.md: missing contract clause /(default|normally).{0,160}(one|single).{0,80}(group|continuous)/
        cost_aware_assets_test.go:138: skills/flowforge-plan/SKILL.md: missing contract clause /(split|slicing).{0,200}(reason|justify|justification)/
        cost_aware_assets_test.go:138: skills/flowforge-plan/SKILL.md: missing contract clause /(new ticket|new plan).{0,180}(schedule|arrangement)/
    --- FAIL: TestCostAwareContinuousDeliveryContracts/skills/_shared/ARTIFACT-CONTRACT.md (0.01s)
        cost_aware_assets_test.go:138: skills/_shared/ARTIFACT-CONTRACT.md: missing contract clause /(default|normally).{0,160}(one|single).{0,80}(group|continuous)/
        cost_aware_assets_test.go:138: skills/_shared/ARTIFACT-CONTRACT.md: missing contract clause /(early|intermediate|upfront).{0,180}(risk|assumption|irreversible|rework)/
    --- FAIL: TestCostAwareContinuousDeliveryContracts/skills/flowforge-solution-design/SKILL.md (0.00s)
        cost_aware_assets_test.go:138: skills/flowforge-solution-design/SKILL.md: missing contract clause /(early|intermediate|upfront).{0,180}(risk|assumption|irreversible|rework)/
    --- FAIL: TestCostAwareContinuousDeliveryContracts/skills/flowforge-refine-ticket/SKILL.md (0.00s)
        cost_aware_assets_test.go:138: skills/flowforge-refine-ticket/SKILL.md: missing contract clause /(unchanged|still.valid).{0,180}(refer|cit|link)/
    --- FAIL: TestCostAwareContinuousDeliveryContracts/skills/flowforge-implement/SKILL.md (0.00s)
        cost_aware_assets_test.go:138: skills/flowforge-implement/SKILL.md: missing contract clause /(one|single|once).{0,100}(execution confirmation|execution check-in)/
        cost_aware_assets_test.go:138: skills/flowforge-implement/SKILL.md: missing contract clause /(unchanged|still.valid).{0,180}(refer|cit|link)/
        cost_aware_assets_test.go:138: skills/flowforge-implement/SKILL.md: missing contract clause /(test.{0,60}(permission|author|protect)|(?:preset|protected).{0,80}test)/
    --- FAIL: TestCostAwareContinuousDeliveryContracts/skills/flowforge-review/SKILL.md (0.01s)
        cost_aware_assets_test.go:138: skills/flowforge-review/SKILL.md: missing contract clause /(main|orchestrat\w*) session.{0,80}(owner|adjudicat|writ)/
        cost_aware_assets_test.go:138: skills/flowforge-review/SKILL.md: missing contract clause /(read.only|leaf).{0,160}(only|assigned|dispatched).{0,100}(audit|axis|scope|finding)/
        cost_aware_assets_test.go:138: skills/flowforge-review/SKILL.md: missing contract clause /(must not|never|do not|does not).{0,80}(spawn|delegat)/
        cost_aware_assets_test.go:138: skills/flowforge-review/SKILL.md: missing contract clause /(same|fixed).{0,100}(point|diff).{0,180}(batch|together)/
    --- FAIL: TestCostAwareContinuousDeliveryContracts/AGENTS.md (0.00s)
        cost_aware_assets_test.go:138: AGENTS.md: missing contract clause /(main|orchestrat\w*) session.{0,100}(owner|adjudicat|writ)/
        cost_aware_assets_test.go:138: AGENTS.md: missing contract clause /(mechanical|short|small).{0,100}(record|write|backfill).{0,160}(direct|scribe|delegat)/
    --- FAIL: TestCostAwareContinuousDeliveryContracts/agents/issue-tracker.md (0.00s)
        cost_aware_assets_test.go:138: agents/issue-tracker.md: missing contract clause /(default|new ticket).{0,180}(schedule|group)/
        cost_aware_assets_test.go:138: agents/issue-tracker.md: missing contract clause /(main|orchestrat\w*) session.{0,100}(owner|adjudicat|writ)/
    cost_aware_assets_test.go:150: flowforge-implement/SKILL.md: obsolete instruction "List every Constraint and every Execution scenario"
    cost_aware_assets_test.go:150: flowforge-implement/SKILL.md: obsolete instruction "for every unchecked Change, one line restating"
    cost_aware_assets_test.go:150: flowforge-plan/SKILL.md: obsolete instruction "**Optional `Verification schedule`:**"
--- FAIL: TestCostAwareRoleContracts (0.00s)
    --- FAIL: TestCostAwareRoleContracts/planner (0.00s)
        cost_aware_assets_test.go:165: planner: missing contract clause /(complete|whole).{0,80}(behavior|logic|delivery)/
        cost_aware_assets_test.go:165: planner: missing contract clause /(default|normally).{0,120}(group|concentrat|consolidat)/
    --- FAIL: TestCostAwareRoleContracts/implementer (0.00s)
        cost_aware_assets_test.go:167: implementer: missing contract clause /(continuous|consecutive|group)/
        cost_aware_assets_test.go:167: implementer: missing contract clause /(must not|never|do not|does not).{0,100}(review|close)/
    --- FAIL: TestCostAwareRoleContracts/reviewer (0.00s)
        cost_aware_assets_test.go:172: reviewer: missing contract clause /(main|orchestrat\w*) session.{0,100}(owner|adjudicat|writ)/
        cost_aware_assets_test.go:172: reviewer: missing contract clause /(must not|never|do not|does not).{0,100}(spawn|delegat)/
        cost_aware_assets_test.go:178: reviewer: missing contract clause /(default|normally).{0,100}Spec/
        cost_aware_assets_test.go:183: reviewer retains obsolete owner instruction "That Skill itself spawns"
        cost_aware_assets_test.go:183: reviewer retains obsolete owner instruction "with `Fix:` Changes appended to the ticket"
    --- FAIL: TestCostAwareRoleContracts/reviewer-lite (0.00s)
        cost_aware_assets_test.go:172: reviewer-lite: missing contract clause /(main|orchestrat\w*) session.{0,100}(owner|adjudicat|writ)/
        cost_aware_assets_test.go:172: reviewer-lite: missing contract clause /(must not|never|do not|does not).{0,100}(spawn|delegat)/
        cost_aware_assets_test.go:183: reviewer-lite retains obsolete owner instruction "sole closeout authority"
        cost_aware_assets_test.go:183: reviewer-lite retains obsolete owner instruction "those are adjudication outputs of flowforge-reviewer"
FAIL
FAIL	flowforge/internal/command	1.816s
FAIL
```

A1后单独定位安全/部署子集，确认编译、源/embedded同步、四宿主输出不是无关失败；随后补充reviewer/lite明确不写文件/不关票条款，形成最终冻结版本A2。重复运行由新增两条断言引起，不是重复尝试同一失败来制造绿。

### A1-G 安全/部署定位：green

原始命令：

```sh
go test -v ./internal/command -run 'Test(CostAwareRulePackageSourceEmbeddedDeployed|CostAwareRolesSourceEmbeddedDeployed|ImplementerPromptPinsLoopContracts|BlockAndRecordContractPinned|ReviewSkillCarriesRound0Audit|ReviewSkillDescribesRepairEscalation|AgentsBlockContainsExecutionUnitPolicy|SubagentSource.*)$' -count=1
```

实际exit: 0

```text
--- PASS: TestImplementerPromptPinsLoopContracts (0.10s)
--- PASS: TestBlockAndRecordContractPinned (0.10s)
--- PASS: TestAgentsBlockContainsExecutionUnitPolicy (0.00s)
--- PASS: TestReviewSkillCarriesRound0Audit (0.00s)
--- PASS: TestReviewSkillDescribesRepairEscalation (0.00s)
--- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed (0.04s)
--- PASS: TestCostAwareRolesSourceEmbeddedDeployed (0.07s)
--- PASS: TestSubagentSourceFilesExist (0.00s)
--- PASS: TestSubagentSourceFrontmatterValid (0.00s)
--- PASS: TestSubagentSourceDefaultSkillResolves (0.00s)
--- PASS: TestSubagentSourceHasFiveSections (0.00s)
PASS
ok  	flowforge/internal/command	0.745s
```

上段为原始关键输出摘录，所有角色×四宿主子例均实际执行，无SKIP；该次测试源hash为551ff6fef4f7e7dfc72754285ef497916939d70e10b3ee7f78b011728f470129（cost_aware）及e60de2b8f618405479c03e88fab29efb33c5120918ce96e6aec17b4d735c687c（assets_deploy）。同一版本gofmt -d无输出、git diff --check exit0；整文件lint与A2-L相同命令exit0/0 issues。A1-G是中间版本证据，最终版本以A2为准。

### A2 最终冻结版本：完整定向作者检查

原始命令：

```sh
go test -v ./internal/command -run 'Test(CostAware.*|ImplementSkillCarriesWeakExecutorContract|ImplementerPromptPinsLoopContracts|BlockAndRecordContractPinned|ReviewSkillCarriesRound0Audit|ReviewSkillDescribesRepairEscalation|AgentsBlockContainsExecutionUnitPolicy|SubagentSource.*)$' -count=1
```

实际exit: 1

```text
=== RUN   TestImplementSkillCarriesWeakExecutorContract
    assets_deploy_test.go:347: execution confirmation: missing contract clause /(one|single|once).{0,100}(execution confirmation|execution check-in)/
    assets_deploy_test.go:347: execution confirmation: missing contract clause /(preset|protected).{0,80}test/
--- FAIL: TestImplementSkillCarriesWeakExecutorContract (0.00s)
=== RUN   TestImplementerPromptPinsLoopContracts
--- PASS: TestImplementerPromptPinsLoopContracts (0.10s)
=== RUN   TestBlockAndRecordContractPinned
--- PASS: TestBlockAndRecordContractPinned (0.09s)
=== RUN   TestAgentsBlockContainsExecutionUnitPolicy
--- PASS: TestAgentsBlockContainsExecutionUnitPolicy (0.00s)
=== RUN   TestReviewSkillCarriesRound0Audit
--- PASS: TestReviewSkillCarriesRound0Audit (0.00s)
=== RUN   TestReviewSkillDescribesRepairEscalation
--- PASS: TestReviewSkillDescribesRepairEscalation (0.00s)
=== RUN   TestCostAwareRulePackageSourceEmbeddedDeployed
=== RUN   TestCostAwareRulePackageSourceEmbeddedDeployed/skills/_shared/ARTIFACT-CONTRACT.md
=== RUN   TestCostAwareRulePackageSourceEmbeddedDeployed/skills/flowforge-solution-design/SKILL.md
=== RUN   TestCostAwareRulePackageSourceEmbeddedDeployed/skills/flowforge-plan/SKILL.md
=== RUN   TestCostAwareRulePackageSourceEmbeddedDeployed/skills/flowforge-refine-ticket/SKILL.md
=== RUN   TestCostAwareRulePackageSourceEmbeddedDeployed/skills/flowforge-implement/SKILL.md
=== RUN   TestCostAwareRulePackageSourceEmbeddedDeployed/skills/flowforge-review/SKILL.md
=== RUN   TestCostAwareRulePackageSourceEmbeddedDeployed/agents/issue-tracker.md
=== RUN   TestCostAwareRulePackageSourceEmbeddedDeployed/AGENTS.md
--- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed (0.04s)
    --- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed/skills/_shared/ARTIFACT-CONTRACT.md (0.00s)
    --- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed/skills/flowforge-solution-design/SKILL.md (0.00s)
    --- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed/skills/flowforge-plan/SKILL.md (0.00s)
    --- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed/skills/flowforge-refine-ticket/SKILL.md (0.00s)
    --- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed/skills/flowforge-implement/SKILL.md (0.00s)
    --- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed/skills/flowforge-review/SKILL.md (0.00s)
    --- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed/agents/issue-tracker.md (0.00s)
    --- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed/AGENTS.md (0.00s)
=== RUN   TestCostAwareContinuousDeliveryContracts
=== RUN   TestCostAwareContinuousDeliveryContracts/skills/flowforge-plan/SKILL.md
    cost_aware_assets_test.go:138: skills/flowforge-plan/SKILL.md: missing contract clause /(default|normally).{0,160}(one|single).{0,80}(group|continuous)/
    cost_aware_assets_test.go:138: skills/flowforge-plan/SKILL.md: missing contract clause /(split|slicing).{0,200}(reason|justify|justification)/
    cost_aware_assets_test.go:138: skills/flowforge-plan/SKILL.md: missing contract clause /(new ticket|new plan).{0,180}(schedule|arrangement)/
=== RUN   TestCostAwareContinuousDeliveryContracts/skills/_shared/ARTIFACT-CONTRACT.md
    cost_aware_assets_test.go:138: skills/_shared/ARTIFACT-CONTRACT.md: missing contract clause /(default|normally).{0,160}(one|single).{0,80}(group|continuous)/
    cost_aware_assets_test.go:138: skills/_shared/ARTIFACT-CONTRACT.md: missing contract clause /(early|intermediate|upfront).{0,180}(risk|assumption|irreversible|rework)/
=== RUN   TestCostAwareContinuousDeliveryContracts/skills/flowforge-solution-design/SKILL.md
    cost_aware_assets_test.go:138: skills/flowforge-solution-design/SKILL.md: missing contract clause /(early|intermediate|upfront).{0,180}(risk|assumption|irreversible|rework)/
=== RUN   TestCostAwareContinuousDeliveryContracts/skills/flowforge-refine-ticket/SKILL.md
    cost_aware_assets_test.go:138: skills/flowforge-refine-ticket/SKILL.md: missing contract clause /(unchanged|still.valid).{0,180}(refer|cit|link)/
=== RUN   TestCostAwareContinuousDeliveryContracts/skills/flowforge-implement/SKILL.md
    cost_aware_assets_test.go:138: skills/flowforge-implement/SKILL.md: missing contract clause /(one|single|once).{0,100}(execution confirmation|execution check-in)/
    cost_aware_assets_test.go:138: skills/flowforge-implement/SKILL.md: missing contract clause /(unchanged|still.valid).{0,180}(refer|cit|link)/
    cost_aware_assets_test.go:138: skills/flowforge-implement/SKILL.md: missing contract clause /(test.{0,60}(permission|author|protect)|(?:preset|protected).{0,80}test)/
=== RUN   TestCostAwareContinuousDeliveryContracts/skills/flowforge-review/SKILL.md
    cost_aware_assets_test.go:138: skills/flowforge-review/SKILL.md: missing contract clause /(main|orchestrat\w*) session.{0,80}(owner|adjudicat|writ)/
    cost_aware_assets_test.go:138: skills/flowforge-review/SKILL.md: missing contract clause /(read.only|leaf).{0,160}(only|assigned|dispatched).{0,100}(audit|axis|scope|finding)/
    cost_aware_assets_test.go:138: skills/flowforge-review/SKILL.md: missing contract clause /(must not|never|do not|does not).{0,80}(spawn|delegat)/
    cost_aware_assets_test.go:138: skills/flowforge-review/SKILL.md: missing contract clause /(same|fixed).{0,100}(point|diff).{0,180}(batch|together)/
=== RUN   TestCostAwareContinuousDeliveryContracts/AGENTS.md
    cost_aware_assets_test.go:138: AGENTS.md: missing contract clause /(main|orchestrat\w*) session.{0,100}(owner|adjudicat|writ)/
    cost_aware_assets_test.go:138: AGENTS.md: missing contract clause /(mechanical|short|small).{0,100}(record|write|backfill).{0,160}(direct|scribe|delegat)/
=== RUN   TestCostAwareContinuousDeliveryContracts/agents/issue-tracker.md
    cost_aware_assets_test.go:138: agents/issue-tracker.md: missing contract clause /(default|new ticket).{0,180}(schedule|group)/
    cost_aware_assets_test.go:138: agents/issue-tracker.md: missing contract clause /(main|orchestrat\w*) session.{0,100}(owner|adjudicat|writ)/
=== NAME  TestCostAwareContinuousDeliveryContracts
    cost_aware_assets_test.go:150: flowforge-plan/SKILL.md: obsolete instruction "**Optional `Verification schedule`:**"
    cost_aware_assets_test.go:150: flowforge-implement/SKILL.md: obsolete instruction "List every Constraint and every Execution scenario"
    cost_aware_assets_test.go:150: flowforge-implement/SKILL.md: obsolete instruction "for every unchecked Change, one line restating"
--- FAIL: TestCostAwareContinuousDeliveryContracts (0.02s)
    --- FAIL: TestCostAwareContinuousDeliveryContracts/skills/flowforge-plan/SKILL.md (0.00s)
    --- FAIL: TestCostAwareContinuousDeliveryContracts/skills/_shared/ARTIFACT-CONTRACT.md (0.01s)
    --- FAIL: TestCostAwareContinuousDeliveryContracts/skills/flowforge-solution-design/SKILL.md (0.00s)
    --- FAIL: TestCostAwareContinuousDeliveryContracts/skills/flowforge-refine-ticket/SKILL.md (0.00s)
    --- FAIL: TestCostAwareContinuousDeliveryContracts/skills/flowforge-implement/SKILL.md (0.00s)
    --- FAIL: TestCostAwareContinuousDeliveryContracts/skills/flowforge-review/SKILL.md (0.01s)
    --- FAIL: TestCostAwareContinuousDeliveryContracts/AGENTS.md (0.00s)
    --- FAIL: TestCostAwareContinuousDeliveryContracts/agents/issue-tracker.md (0.00s)
=== RUN   TestCostAwareRoleContracts
=== RUN   TestCostAwareRoleContracts/planner
    cost_aware_assets_test.go:165: planner: missing contract clause /(complete|whole).{0,80}(behavior|logic|delivery)/
    cost_aware_assets_test.go:165: planner: missing contract clause /(default|normally).{0,120}(group|concentrat|consolidat)/
=== RUN   TestCostAwareRoleContracts/implementer
    cost_aware_assets_test.go:167: implementer: missing contract clause /(continuous|consecutive|group)/
    cost_aware_assets_test.go:167: implementer: missing contract clause /(must not|never|do not|does not).{0,100}(review|close)/
=== RUN   TestCostAwareRoleContracts/reviewer
    cost_aware_assets_test.go:172: reviewer: missing contract clause /(main|orchestrat\w*) session.{0,100}(owner|adjudicat|writ)/
    cost_aware_assets_test.go:172: reviewer: missing contract clause /(must not|never|do not|does not).{0,100}(spawn|delegat)/
    cost_aware_assets_test.go:172: reviewer: missing contract clause /(must not|never|do not|does not).{0,140}close/
    cost_aware_assets_test.go:180: reviewer: missing contract clause /(default|normally).{0,100}Spec/
    cost_aware_assets_test.go:185: reviewer retains obsolete owner instruction "That Skill itself spawns"
    cost_aware_assets_test.go:185: reviewer retains obsolete owner instruction "with `Fix:` Changes appended to the ticket"
=== RUN   TestCostAwareRoleContracts/reviewer-lite
    cost_aware_assets_test.go:172: reviewer-lite: missing contract clause /(main|orchestrat\w*) session.{0,100}(owner|adjudicat|writ)/
    cost_aware_assets_test.go:172: reviewer-lite: missing contract clause /(must not|never|do not|does not).{0,100}(spawn|delegat)/
    cost_aware_assets_test.go:172: reviewer-lite: missing contract clause /(must not|never|do not|does not).{0,140}close/
    cost_aware_assets_test.go:185: reviewer-lite retains obsolete owner instruction "sole closeout authority"
    cost_aware_assets_test.go:185: reviewer-lite retains obsolete owner instruction "those are adjudication outputs of flowforge-reviewer"
--- FAIL: TestCostAwareRoleContracts (0.00s)
    --- FAIL: TestCostAwareRoleContracts/planner (0.00s)
    --- FAIL: TestCostAwareRoleContracts/implementer (0.00s)
    --- FAIL: TestCostAwareRoleContracts/reviewer (0.00s)
    --- FAIL: TestCostAwareRoleContracts/reviewer-lite (0.00s)
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/planner
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/planner/.claude/agents
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/planner/.opencode/agent
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/planner/.codex/agents
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/planner/.pi/agents
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/implementer
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/implementer/.claude/agents
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/implementer/.opencode/agent
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/implementer/.codex/agents
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/implementer/.pi/agents
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/reviewer
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/reviewer/.claude/agents
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/reviewer/.opencode/agent
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/reviewer/.codex/agents
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/reviewer/.pi/agents
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/reviewer-lite
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/reviewer-lite/.claude/agents
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/reviewer-lite/.opencode/agent
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/reviewer-lite/.codex/agents
=== RUN   TestCostAwareRolesSourceEmbeddedDeployed/reviewer-lite/.pi/agents
--- PASS: TestCostAwareRolesSourceEmbeddedDeployed (0.05s)
    --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/planner (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/planner/.claude/agents (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/planner/.opencode/agent (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/planner/.codex/agents (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/planner/.pi/agents (0.00s)
    --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/implementer (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/implementer/.claude/agents (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/implementer/.opencode/agent (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/implementer/.codex/agents (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/implementer/.pi/agents (0.00s)
    --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/reviewer (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/reviewer/.claude/agents (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/reviewer/.opencode/agent (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/reviewer/.codex/agents (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/reviewer/.pi/agents (0.00s)
    --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/reviewer-lite (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/reviewer-lite/.claude/agents (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/reviewer-lite/.opencode/agent (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/reviewer-lite/.codex/agents (0.00s)
        --- PASS: TestCostAwareRolesSourceEmbeddedDeployed/reviewer-lite/.pi/agents (0.00s)
=== RUN   TestSubagentSourceFilesExist
--- PASS: TestSubagentSourceFilesExist (0.00s)
=== RUN   TestSubagentSourceFrontmatterValid
=== RUN   TestSubagentSourceFrontmatterValid/flowforge-analyst
=== RUN   TestSubagentSourceFrontmatterValid/flowforge-architect
=== RUN   TestSubagentSourceFrontmatterValid/flowforge-planner
=== RUN   TestSubagentSourceFrontmatterValid/flowforge-implementer
=== RUN   TestSubagentSourceFrontmatterValid/flowforge-reviewer
=== RUN   TestSubagentSourceFrontmatterValid/flowforge-investigator
=== RUN   TestSubagentSourceFrontmatterValid/flowforge-batch-analyst
=== RUN   TestSubagentSourceFrontmatterValid/flowforge-scribe
=== RUN   TestSubagentSourceFrontmatterValid/flowforge-executor
=== RUN   TestSubagentSourceFrontmatterValid/flowforge-reviewer-lite
=== RUN   TestSubagentSourceFrontmatterValid/flowforge-frontend-implementer
=== RUN   TestSubagentSourceFrontmatterValid/flowforge-frontend-reviewer
--- PASS: TestSubagentSourceFrontmatterValid (0.00s)
    --- PASS: TestSubagentSourceFrontmatterValid/flowforge-analyst (0.00s)
    --- PASS: TestSubagentSourceFrontmatterValid/flowforge-architect (0.00s)
    --- PASS: TestSubagentSourceFrontmatterValid/flowforge-planner (0.00s)
    --- PASS: TestSubagentSourceFrontmatterValid/flowforge-implementer (0.00s)
    --- PASS: TestSubagentSourceFrontmatterValid/flowforge-reviewer (0.00s)
    --- PASS: TestSubagentSourceFrontmatterValid/flowforge-investigator (0.00s)
    --- PASS: TestSubagentSourceFrontmatterValid/flowforge-batch-analyst (0.00s)
    --- PASS: TestSubagentSourceFrontmatterValid/flowforge-scribe (0.00s)
    --- PASS: TestSubagentSourceFrontmatterValid/flowforge-executor (0.00s)
    --- PASS: TestSubagentSourceFrontmatterValid/flowforge-reviewer-lite (0.00s)
    --- PASS: TestSubagentSourceFrontmatterValid/flowforge-frontend-implementer (0.00s)
    --- PASS: TestSubagentSourceFrontmatterValid/flowforge-frontend-reviewer (0.00s)
=== RUN   TestSubagentSourceDefaultSkillResolves
=== RUN   TestSubagentSourceDefaultSkillResolves/flowforge-analyst
=== RUN   TestSubagentSourceDefaultSkillResolves/flowforge-architect
=== RUN   TestSubagentSourceDefaultSkillResolves/flowforge-planner
=== RUN   TestSubagentSourceDefaultSkillResolves/flowforge-implementer
=== RUN   TestSubagentSourceDefaultSkillResolves/flowforge-reviewer
=== RUN   TestSubagentSourceDefaultSkillResolves/flowforge-investigator
=== RUN   TestSubagentSourceDefaultSkillResolves/flowforge-batch-analyst
=== RUN   TestSubagentSourceDefaultSkillResolves/flowforge-scribe
=== RUN   TestSubagentSourceDefaultSkillResolves/flowforge-executor
=== RUN   TestSubagentSourceDefaultSkillResolves/flowforge-reviewer-lite
=== RUN   TestSubagentSourceDefaultSkillResolves/flowforge-frontend-implementer
=== RUN   TestSubagentSourceDefaultSkillResolves/flowforge-frontend-reviewer
--- PASS: TestSubagentSourceDefaultSkillResolves (0.00s)
    --- PASS: TestSubagentSourceDefaultSkillResolves/flowforge-analyst (0.00s)
    --- PASS: TestSubagentSourceDefaultSkillResolves/flowforge-architect (0.00s)
    --- PASS: TestSubagentSourceDefaultSkillResolves/flowforge-planner (0.00s)
    --- PASS: TestSubagentSourceDefaultSkillResolves/flowforge-implementer (0.00s)
    --- PASS: TestSubagentSourceDefaultSkillResolves/flowforge-reviewer (0.00s)
    --- PASS: TestSubagentSourceDefaultSkillResolves/flowforge-investigator (0.00s)
    --- PASS: TestSubagentSourceDefaultSkillResolves/flowforge-batch-analyst (0.00s)
    --- PASS: TestSubagentSourceDefaultSkillResolves/flowforge-scribe (0.00s)
    --- PASS: TestSubagentSourceDefaultSkillResolves/flowforge-executor (0.00s)
    --- PASS: TestSubagentSourceDefaultSkillResolves/flowforge-reviewer-lite (0.00s)
    --- PASS: TestSubagentSourceDefaultSkillResolves/flowforge-frontend-implementer (0.00s)
    --- PASS: TestSubagentSourceDefaultSkillResolves/flowforge-frontend-reviewer (0.00s)
=== RUN   TestSubagentSourceHasFiveSections
=== RUN   TestSubagentSourceHasFiveSections/flowforge-analyst
=== RUN   TestSubagentSourceHasFiveSections/flowforge-architect
=== RUN   TestSubagentSourceHasFiveSections/flowforge-planner
=== RUN   TestSubagentSourceHasFiveSections/flowforge-implementer
=== RUN   TestSubagentSourceHasFiveSections/flowforge-reviewer
=== RUN   TestSubagentSourceHasFiveSections/flowforge-investigator
=== RUN   TestSubagentSourceHasFiveSections/flowforge-batch-analyst
=== RUN   TestSubagentSourceHasFiveSections/flowforge-scribe
=== RUN   TestSubagentSourceHasFiveSections/flowforge-executor
=== RUN   TestSubagentSourceHasFiveSections/flowforge-reviewer-lite
=== RUN   TestSubagentSourceHasFiveSections/flowforge-frontend-implementer
=== RUN   TestSubagentSourceHasFiveSections/flowforge-frontend-reviewer
--- PASS: TestSubagentSourceHasFiveSections (0.00s)
    --- PASS: TestSubagentSourceHasFiveSections/flowforge-analyst (0.00s)
    --- PASS: TestSubagentSourceHasFiveSections/flowforge-architect (0.00s)
    --- PASS: TestSubagentSourceHasFiveSections/flowforge-planner (0.00s)
    --- PASS: TestSubagentSourceHasFiveSections/flowforge-implementer (0.00s)
    --- PASS: TestSubagentSourceHasFiveSections/flowforge-reviewer (0.00s)
    --- PASS: TestSubagentSourceHasFiveSections/flowforge-investigator (0.00s)
    --- PASS: TestSubagentSourceHasFiveSections/flowforge-batch-analyst (0.00s)
    --- PASS: TestSubagentSourceHasFiveSections/flowforge-scribe (0.00s)
    --- PASS: TestSubagentSourceHasFiveSections/flowforge-executor (0.00s)
    --- PASS: TestSubagentSourceHasFiveSections/flowforge-reviewer-lite (0.00s)
    --- PASS: TestSubagentSourceHasFiveSections/flowforge-frontend-implementer (0.00s)
    --- PASS: TestSubagentSourceHasFiveSections/flowforge-frontend-reviewer (0.00s)
FAIL
FAIL	flowforge/internal/command	0.804s
FAIL
```

A2实际exit1：TestImplementSkillCarriesWeakExecutorContract、TestCostAwareContinuousDeliveryContracts及TestCostAwareRoleContracts是新语义尚未落地的预期red。其余原安全断言、source/embedded/隔离部署及SubagentSource检查实际PASS，无SKIP、无编译错误。逐条失败对应上方映射；G1修改源后这些新语义断言须pass，不授权执行者豁免或改测。

### A2-L 最终作者Go整文件lint

原始命令：

```sh
/tmp/flowforge-verification-tools/golangci-lint-2.15.0/golangci-lint run --timeout 2m --new-from-rev 16cce8e --whole-files --max-issues-per-linter 0 --max-same-issues 0 ./cmd/... ./internal/...
```

实际exit: 0

```text
0 issues.
```

### A2-F 格式只读检查、差分及冻结hash

原始命令：

```sh
gofmt -d internal/command/cost_aware_assets_test.go internal/command/assets_deploy_test.go && git diff --check && shasum -a 256 internal/command/cost_aware_assets_test.go internal/command/assets_deploy_test.go
```

实际exit: 0

```text
17eca4b481e495fdb57cf5ad6206da39753582b9f7cf2e17825fba03b1b90033  internal/command/cost_aware_assets_test.go
e60de2b8f618405479c03e88fab29efb33c5120918ce96e6aec17b4d735c687c  internal/command/assets_deploy_test.go
```

### 核实期间的失败与处置

- 一次read请求internal/command/subagents_test.go返回ENOENT；随后rg定位实际接缝internal/command/subagent_source_test.go及agents_test.go。另一次并行read中internal/command/agents_deploy.go、internal/command/subagent_test.go、internal/subagent/definition.go不存在；rg核实实际为internal/command/agents.go、agents_test.go、internal/subagent/parser.go。这些是只读定位失败，未产生代码/验收失败或据此填猜测。
- 原始诊断/搜索命令：`bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --full && rg -n 'waiv|execution-contract|rev|01|04' ff-wiki/proposals/cost-aware-verification-planning/waivers* ff-wiki/proposals/cost-aware-verification-planning/evidence/01-lint-baseline.md`，实际exit2；check本身exit0，rg输出`rg: ff-wiki/proposals/cost-aware-verification-planning/waivers*: No such file or directory (os error 2)`。未宣称该组合命令通过。历史精确waiver由check --full实际报告核实，无需猜测waiver文件名。
- 初始check/frontier实际exit0，05/02各有execution-contract-incomplete。五节填充后、证据文件写入前的临时check/frontier组合exit0：05缺节已清，但出现evidence-artifact-missing；本文件落盘后再核实消失，不隐去中间诊断。
- embed变量实际位置经`rg -n 'var embeddedAssets|go:embed' internal/command`核实为internal/command/embed.go:6；票中坐标按实修正。

### Refine完成与剩余验收

五节已按源码/命令/设计填充。04已关闭，05保持open，仅Change1勾选；其四元组引用A2-L真实exit0作者lint，并明确A2定向实际exit1，未把red写成exit0。组末G1、独立审查和02未执行。历史00/03需求rev2→3及01/04需求rev2→3、设计rev3→4的六项upstream-changed精确waiver保留；历史通过不作新行为证明。

剩余风险：自然语言条款测试有能力边界，独立审查须核对owner/leaf分支、前置理由、旧票路径和所有入口一致；02须观察真实B1/B2与发布门禁。本记录不宣称净提速或产品已完成。

### 最终结构核实

原始命令：`bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --full`；实际exit: 0

```text
Checked 6 issues in ff-wiki/proposals/cost-aware-verification-planning
gap: execution-contract-incomplete ff-wiki/proposals/cost-aware-verification-planning/issues/02-controlled-acceptance-and-deploy.md — Open ticket is missing a complete machine execution contract
warning: upstream-changed ff-wiki/proposals/cost-aware-verification-planning/issues/00-test-authority.md — Consumed revision 2 is behind 3 (waived: 需求rev3保留R10语义；本closed票仅交付已验收的受限作者入口/边界，不承担新Plan规则。)
warning: upstream-changed ff-wiki/proposals/cost-aware-verification-planning/issues/01-rule-package.md — Consumed revision 3 is behind 4 (waived: 历史closed票保留消费设计rev3；rev4替换默认和角色职责由05承担，旧证据不作新规则正向证明。)
warning: upstream-changed ff-wiki/proposals/cost-aware-verification-planning/issues/01-rule-package.md — Consumed revision 2 is behind 3 (waived: 历史closed票仍证明需求rev2下的旧范围；rev3新增连续交付及角色义务由05与修订02交付，不追写旧通过。)
warning: upstream-changed ff-wiki/proposals/cost-aware-verification-planning/issues/03-test-authority-boundary-repair.md — Consumed revision 2 is behind 3 (waived: 需求rev3保留R10语义；本closed票仅交付已验收的受限作者入口/边界，不承担新Plan规则。)
warning: upstream-changed ff-wiki/proposals/cost-aware-verification-planning/issues/04-stage-proof-scope-and-execution-record.md — Consumed revision 2 is behind 3 (waived: 历史closed票仍证明需求rev2下的旧范围；rev3新增连续交付及角色义务由05与修订02交付，不追写旧通过。)
warning: upstream-changed ff-wiki/proposals/cost-aware-verification-planning/issues/04-stage-proof-scope-and-execution-record.md — Consumed revision 3 is behind 4 (waived: 历史closed票保留消费设计rev3；rev4替换默认和角色职责由05承担，旧证据不作新规则正向证明。)
✓ Dependency graph is healthy. No cycles or dangling references found.
```

原始命令：`bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning`；实际exit: 0

```text
=== READY (Unblocked & Executable) ===
✓ [cost-aware-verification-planning] #05 连续实现、集中验证及角色职责一致的规则包 (ff-wiki/proposals/cost-aware-verification-planning/issues/05-continuous-delivery-rule-package.md)

=== BLOCKED ===
⛔ #02 集中验证连续交付行为并发布一致规则包 (Waiting on: [05])
diagnostics gap:execution-contract-incomplete: 1
gap: execution-contract-incomplete ff-wiki/proposals/cost-aware-verification-planning/issues/02-controlled-acceptance-and-deploy.md
--full for all diagnostics
```

最终git diff --check实际exit0、无输出；撤权后SHA256与A2-F一致。05的execution-contract-incomplete及临时evidence-artifact-missing均已消失，frontier由当前文件推导05可执行；没有持久化ready状态。02缺五节和六项历史waiver原样保留。

<a id="g1-implementation"></a>
## G1 / Changes2–5：集中实现与检查

本节只追加G1的真实结果，不覆盖也不改写上文G0任何记录。G0的A2 exit1仍是G0的真实结果；G1列出的exit0只适用于G1落地后的固定点。基线08e87cb19001aef9b70a0bb047dee7b08b76003e，消费requirements rev3与design rev4；固定差分见本节末`G1 固定点`。未运行02真实行为链、未在本项目部署新规则、未宣称净提速。

### 执行方式与共享归属

- G1按本票`### Verification schedule`一组执行：先连续改完所有互相关联的源面（`_shared/ARTIFACT-CONTRACT.md`、`flowforge-{solution-design,plan,refine-ticket,implement,review}/SKILL.md`、四个subagent源、`assets/AGENTS.md`、`assets/agents/issue-tracker.md`），随后一次`make dev`、完整内部Go、整文件lint、`git diff --check`及`check`/`frontier`。同一固定点没有再跑定向子集。
- 本会话按主会话固定的G1覆盖执行：旧部署Skill的逐Change build与Phase0/0b两次全文复述被本票/design rev4取代；未修改的`.agents/skills`候选文本未成为本会话执行入口。
- Changes2–5的四元组（见票）全部引用同一次完整内部Go运行（G1组#1）；Change5另引用同固定点的`make dev`/`diff -qr`、lint、`git diff --check`、`check`/`frontier`。一次运行不记成多次。
- 两份冻结测试全程只读；本轮未新增、未修改、未删除任何Go测试文件，`git diff --stat`只有6份skill/shared、4个subagent源、AGENTS.md、issue-tracker.md及本票。

### G1-C 源→嵌入→隔离部署（make dev + diff）

原始命令与关键输出（`make dev`末5行，实际exit 0）：

```text
cp -R assets internal/command/assets
mkdir -p .tmp bin
GOTMPDIR=$(pwd)/.tmp go build -ldflags="-s -w -X flowforge/internal/version.injected=v5.13.1-57-ge188b38-dirty" -trimpath -o bin/flowforge ./cmd/flowforge
chmod +x bin/flowforge
rm -rf .tmp
```

原始命令：`diff -qr assets internal/command/assets`；实际exit: 0，无输出（源与嵌入逐字节一致）。

源/嵌入/隔离部署三面对齐由同一次完整内部Go实际验证，无SKIP：

```text
--- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed (0.04s)
--- PASS: TestCostAwareRolesSourceEmbeddedDeployed (0.05s)
```

前者覆盖6份skill/shared、`agents/issue-tracker.md`及`AGENTS.md`管理块；后者覆盖4角色×Claude/OpenCode/Codex/PI，PI工具白名单仍为read/grep/find/ls，Codex仍为`sandbox_mode = "read-only"`。部署只在`t.TempDir()`临时项目内发生，未触碰本项目`.agents`/宿主文件，未读取本项目模型配置。

### G1-G 完整内部Go（唯一一次，覆盖作者子集）

原始命令：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/...`；实际exit: 0。完整未截断输出保存在`/tmp/g1-full-go.log`（1124行）。

```text
--- PASS: TestImplementSkillCarriesWeakExecutorContract (0.00s)
--- PASS: TestCostAwareContinuousDeliveryContracts (0.01s)
--- PASS: TestCostAwareRoleContracts (0.00s)
--- PASS: TestCostAwareRulePackageSourceEmbeddedDeployed (0.04s)
--- PASS: TestCostAwareRolesSourceEmbeddedDeployed (0.05s)
--- PASS: TestPackagedSkillPointersResolve (0.06s)
PASS
ok  	flowforge/internal/command	11.271s
ok  	flowforge/internal/config	(cached)
ok  	flowforge/internal/subagent	1.028s
ok  	flowforge/internal/tracker	1.197s
ok  	flowforge/internal/update	(cached)
```

统计：243个`--- PASS`，0个`--- FAIL`，3个`--- SKIP`。A2中预期red的三项新语义测试族（`TestImplementSkillCarriesWeakExecutorContract`、`TestCostAwareContinuousDeliveryContracts`、`TestCostAwareRoleContracts`）现在真实pass，未改断言、未放宽匹配。

环境SKIP原文（与本次新增断言无关，未用于通过）：

```text
    catalog_test.go:531: Tangram compatibility checkout is not available
--- SKIP: TestRepositoryAndTangramProposalLayouts (0.12s)
    upgrade_test.go:219: cannot read /proc/self/exe: open /proc/self/exe: no such file or directory
--- SKIP: TestE2E_SelfupdateApply_ReplacesBinaryCorrectly (0.00s)
    upgrade_test.go:347: cannot read /proc/self/exe: open /proc/self/exe: no such file or directory
--- SKIP: TestFullUpgradePipeline (0.00s)
```

### G1-L 整文件lint（2.15.0，同固定点一次）

原始命令：`/tmp/flowforge-verification-tools/golangci-lint-2.15.0/golangci-lint run --timeout 2m --new-from-rev 16cce8e --whole-files --max-issues-per-linter 0 --max-same-issues 0 ./cmd/... ./internal/...`；实际exit: 0。

```text
0 issues.
```

未降级工具、未改lint选项。历史全仓lint债仍在[01原记录](01-lint-baseline.md)，本次通过只针对该整文件范围。

### G1-F/G1-S 差分与图检查

- `git diff --check`：实际exit 0，无输出。
- `bin/flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --full`：实际exit 0。`05`无任何诊断（无`evidence-missing`/`evidence-incomplete`/`evidence-artifact-missing`）；`02`缺五节`execution-contract-incomplete`原样保留；六项历史`upstream-changed`精确waiver原样保留；末尾`✓ Dependency graph is healthy. No cycles or dangling references found.`。
- `bin/flowforge frontier --dir ff-wiki/proposals/cost-aware-verification-planning`：实际exit 0。`05`在`=== READY (Unblocked & Executable) ===`，`02`仍`⛔ #02 ... (Waiting on: [05])`。frontier由当前文件推导，没有持久化ready状态。

### 设计分支→落地位置

| design分支 | 落地位置 |
| --- | --- |
| #allocation 完整逻辑与拆分理由 | `flowforge-plan/SKILL.md:2. Draft tracer increments`、`3. Publish execution contracts`的新票schedule条；`flowforge-solution-design/SKILL.md:Verification strategy check`（必要性归design、分组归Plan） |
| #schedule 紧凑一组、提前关口、旧票/残缺路径、Fix优先级 | `_shared/ARTIFACT-CONTRACT.md:Verification schedule`、`Execution and review loop`；`flowforge-implement/SKILL.md`验时入口预检；`agents/issue-tracker.md:Publishing`/`Execution and review loop`；`AGENTS.md`（Execution unit policy新增`Verification timing`一段） |
| #evidence 一次执行确认、五节精确引用、四元组与共享运行 | `flowforge-implement/SKILL.md:Phase 0`（0b并入）、`3a`/`3c`/`3d`；`flowforge-refine-ticket/SKILL.md:4. Fill the five execution-contract sections`；`_shared/ARTIFACT-CONTRACT.md`循环第2步与`Shared runs and reused evidence` |
| #owners / #fix-contract 主会话owner与只读leaf | `flowforge-review/SKILL.md:1. Pin the fixed point and scope`、`4. Dispatch both axes in parallel`、`6. Fix planning`、`If zero findings`；四个subagent源正文；`AGENTS.md:Workflow-state delegation`；`agents/issue-tracker.md:Execution and review loop`；`_shared/ARTIFACT-CONTRACT.md:Ownership` |
| 04已修正门禁/repair/Fix入口保留 | 未删除：`Stage proof before dispatch`、`Effective verification entry`、`needs-repair`/`Repair of:`、Round 0先于双轴、2次/5轮上限、失败原文 |

### 测试冻结状态

G0冻结的两份测试在本轮开始前、修改后、提交前三次读取均为同一SHA256；实施期间未写入：

- internal/command/cost_aware_assets_test.go：17eca4b481e495fdb57cf5ad6206da39753582b9f7cf2e17825fba03b1b90033
- internal/command/assets_deploy_test.go：e60de2b8f618405479c03e88fab29efb33c5120918ce96e6aec17b4d735c687c

### 例外与偏离（实际发生过的事项）

- `Phase 0b`按Change3并入`Phase 0`，票`Touch points`仍写作`Phase 0/0b/3a–3e`；这是Change3明确要求的结果，不是语义偏离，票未改写原文。
- `assets/AGENTS.md`新增`Verification timing`一段：AGENTS原本没有新旧票验证时机口径，而design#allocation/#schedule要求入口一致；其余入口各自已有同类语句，此段未在其他文件重复。
- `internal/command/assets/**`与`bin/flowforge`由`make dev`生成且被`.gitignore`忽略（`.gitignore:30`、`.gitignore:17`），因此已重新生成但不在提交内；它们的正确性由同固定点的`diff -qr`与`TestCostAwareRulePackageSourceEmbeddedDeployed`/`TestCostAwareRolesSourceEmbeddedDeployed`实测证明。
- 本轮无失败重试、无`## Blocked evidence`、无写集外修改；未运行本项目的`flowforge init`/`agents deploy`，未改权限profile、模型配置、PI保护或tracker算法。
- 完整内部Go在本票勾选与本节落盘前运行；其后只改了本票与本节文档。唯一读取真实`ff-wiki/proposals`的Go断言是`TestRepositoryAndTangramProposalLayouts`的FlowForge阶段（只断言没有`spec.md`进入执行投影，且本次未增删文件），因此未重跑全套；该判断请独立Round 0核对。

### 剩余风险与未执行项

- 条款测试只证明条款存在，不证明语义一致；owner/leaf分支、提前关口理由、旧票路径和入口一致性必须由独立Round 0/双轴对照design rev4判定。
- 02的两条真实行为链（B1、B2）与发布门禁未执行；本轮不声称净提速或产品完成。
- 新规则需新会话消费；本次只改源，未部署本项目`.agents`。

## G1 固定点

本条由G1实施者提交后补记，见紧邻提交`G1 固定点记录`行/主会话返回的commit。

