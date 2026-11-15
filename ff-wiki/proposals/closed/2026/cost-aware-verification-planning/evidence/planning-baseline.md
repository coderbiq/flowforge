---
flowforge:
  schema: 1
  role: evidence
  id: cost-aware-planning-baseline
  revision: 1
---

# 规划基线（实施前）

本记录仅证明规划时仓库事实，不是规则行为验收或 Completion evidence。固定点 `10c1f667492a81976037b3ff82a7e9ce10a64ad0`；采集前 `git status --short` 为空。

## 源与嵌入漂移

逐文件读取 `assets/**` 与 `internal/command/assets/**`，以下列分别为源路径、源 SHA-256、嵌入 SHA-256；其他源文件无差异。

```text
HEAD: 10c1f667492a81976037b3ff82a7e9ce10a64ad0
assets/pi/flowforge.ts f4f1b1fd09625e3016c702f9f2d1e432599316d0f8fe25b608b078b67d85ed6e 57ce634a369828847ff83094b68e3ad9f88b2fe411707ddf743a676611e7bfd2
assets/skills/_shared/ARTIFACT-CONTRACT.md 9f3f7ca2efffa459959ab3f6fec1f8bc4bebf3d6c42f6694dc3d072fbdfed1bc 1fec741700a455821facfd43a681475f55ae8cb20cf64aba5eef8681a967cf07
assets/skills/flowforge-implement/SKILL.md 516f3cfafeaf76babe226c529aebb01884d4c4e3d9b5f086f23c17373293a40e c480f427ea3035744af5fb5a2708a96550fd32723f04c7b92d1cac90fddbb78b
assets/skills/flowforge-review/SKILL.md db988915946cfa952929556f56c0981a2899212c3d94c273d0f1af9a4ca4eeb7 2fc5faa58804f96a1ca46caf07d9c748fe3430acfb414025d464706e90bcbfa6
```

`git diff --no-index assets internal/command/assets` exit=1（表示存在差异，不是比较工具失败）；方向为源 → 旧嵌入。共享契约/implement/review 的差异涉及关闭责任，PI 差异涉及已存在的 feature 参数。01 的 `make dev` 会同步整个资产树，PI 的这项既有漂移允许作为生成同步带入，不能另改 PI 源。

差分展示将制表符展开为空格并去除行尾空白；精确内容由上述固定 commit 与 SHA-256 绑定。

```diff
diff --git a/assets/pi/flowforge.ts b/internal/command/assets/pi/flowforge.ts
index 941298d..5a175f3 100644
--- a/assets/pi/flowforge.ts
+++ b/internal/command/assets/pi/flowforge.ts
@@ -257,16 +257,8 @@ export default function flowforgeExtension(pi: ExtensionAPI): void {
         label: "FlowForge Frontier",
         description: "Compute unblocked, ready-to-execute FlowForge tickets as JSON (wraps `flowforge frontier --json`).",
         promptSnippet: "List ready FlowForge tickets",
-        parameters: Type.Object({
-            feature: Type.Optional(Type.String({ description: "Proposal slug to scope output to (graph stays global; pass only when working one proposal)" })),
-        }),
-        execute: async (args) => {
-            const cliArgs = ["frontier", "--json"];
-            if (args?.feature) {
-                cliArgs.push("--feature", args.feature);
-            }
-            return runCli(binary, cliArgs);
-        },
+        parameters: Type.Object({}),
+        execute: async () => runCli(binary, ["frontier", "--json"]),
     });

     pi.registerTool({
@@ -274,15 +266,7 @@ export default function flowforgeExtension(pi: ExtensionAPI): void {
         label: "FlowForge Check",
         description: `Validate the FlowForge proposal graph in ${proposalsDir} (wraps \`flowforge check --dir ${proposalsDir}\`).`,
         promptSnippet: "Validate the FlowForge proposal graph",
-        parameters: Type.Object({
-            feature: Type.Optional(Type.String({ description: "Proposal slug to scope reported diagnostics to (validity stays global)" })),
-        }),
-        execute: async (args) => {
-            const cliArgs = ["check", "--dir", proposalsDir];
-            if (args?.feature) {
-                cliArgs.push("--feature", args.feature);
-            }
-            return runCli(binary, cliArgs);
-        },
+        parameters: Type.Object({}),
+        execute: async () => runCli(binary, ["check", "--dir", proposalsDir]),
     });
 }
diff --git a/assets/skills/_shared/ARTIFACT-CONTRACT.md b/internal/command/assets/skills/_shared/ARTIFACT-CONTRACT.md
index a347db5..c6a1c4b 100644
--- a/assets/skills/_shared/ARTIFACT-CONTRACT.md
+++ b/internal/command/assets/skills/_shared/ARTIFACT-CONTRACT.md
@@ -92,7 +92,7 @@ When the implementer is a lightweight model and review runs as a separate sessio
 2. **Implement** (lightweight mode) executes unchecked Changes in order, runs mechanical self-checks (build, focused tests, full verify), checks off completed Changes, writes an Implementation note, and stops. It does not review, close the ticket, or make design decisions.
 3. **Review** runs the dual-axis review on the fixed change set. If findings exist, it translates each fixable finding into a new unchecked `Fix:` Change appended to the ticket. If a finding requires an architecture/seam/interface change, it marks the finding as a design return instead. It records the round in Review rounds.
 4. **Implement** re-executes the new fix Changes. Steps 3-4 repeat.
-5. When a review round produces zero findings, the closeout owner writes Completion evidence and closes the ticket: Review in the separate lightweight loop, or Implement in full mode after Review returns its findings.
+5. When a review round produces zero findings, **Review** writes Completion evidence and closes the ticket.

 **Substantive finding → repair ticket:** If a finding is High/Critical, changes an execution contract or acceptance criterion, crosses the write set, or alters a boundary downstream work depends on, Review does not append it to the original ticket. Instead, it creates a new repair ticket (`**Repair of:** <original-id>`, `**Status:** open`), sets the original to `needs-repair`, rewires downstream `Blocked by` to the repair, and records the repair reference in the original's Review rounds. The repair ticket enters the frontier independently; the original stays non-executable until the repair passes its own clean review, after which the closeout owner links the repair's evidence back to the original and closes both once their delivery conditions are satisfied. `Repair of:` is provenance, not a DAG edge.

diff --git a/assets/skills/flowforge-implement/SKILL.md b/internal/command/assets/skills/flowforge-implement/SKILL.md
index 536d171..aa0ea62 100644
--- a/assets/skills/flowforge-implement/SKILL.md
+++ b/internal/command/assets/skills/flowforge-implement/SKILL.md
@@ -16,7 +16,6 @@ Read project invariants, the ticket, its linked requirement/design authorities a
 - A blocker stops dependent work.
 - A gap stops only affected work unless the caller explicitly includes gaps; preserve the diagnostic in the report.
 - A warning remains visible; strict policy may exclude it.
-- Scope tool output when working a single proposal: pass `feature: <slug>` to `flowforge_frontier` / `flowforge_check` (projects that proposal's queue and diagnostics; validity stays global).
 - An upstream revision mismatch is reviewed against current meaning before execution.

 Do not advance status to express readiness. Begin only when the ticket's Delivery, design context, decided Changes, Constraints, and feasible Done and verify method form an executable contract.
@@ -102,7 +101,7 @@ Classify discoveries immediately:

 Invoke `flowforge-review` with a resolvable fixed point, an explicit committed or working-tree diff scope, and the effective specification links/revisions. Keep Standards and Specification findings separate.

-Resolve every finding by correction or an authority-owned disposition. A Specification finding that means required delivery is absent keeps the affected ticket open until corrected, or until the requirement/design owner explicitly rescopes the effective specification or grants an exact reasoned waiver. Creating a follow-up alone is not resolution; it can permit close only when its recorded blocking relationship prevents premature close, or the applicable authority explicitly accepts the finding as nonblocking. In full mode, Review returns its findings and the Implement owner writes Completion evidence, closes the ticket, and publishes the next frontier; an independent Review session owns those closeout actions only when it is reviewing a separate lightweight implementation.
+Resolve every finding by correction or an authority-owned disposition. A Specification finding that means required delivery is absent keeps the affected ticket open until corrected, or until the requirement/design owner explicitly rescopes the effective specification or grants an exact reasoned waiver. Creating a follow-up alone is not resolution; it can permit close only when its recorded blocking relationship prevents premature close, or the applicable authority explicitly accepts the finding as nonblocking. Review does not close the ticket.

 #### 4b. Record evidence, then close

diff --git a/assets/skills/flowforge-review/SKILL.md b/internal/command/assets/skills/flowforge-review/SKILL.md
index 641e437..3063c4b 100644
--- a/assets/skills/flowforge-review/SKILL.md
+++ b/internal/command/assets/skills/flowforge-review/SKILL.md
@@ -163,7 +163,7 @@ Append fix Changes to the ticket's Changes section and record the round. The lig

 For full-mode Implement, return the clean review and all finding dispositions to the implementation owner. Implement writes Completion evidence, closes the ticket after its own closeout checks, and publishes the next frontier. Review continues to steps 7–8 without writing Completion evidence or closing the ticket.

-For a separate review after lightweight implementation, Review writes Completion evidence, closes the ticket, and publishes the next frontier only after the clean review and its closeout checks. If findings exist, neither mode closes the ticket until the findings are corrected or authoritatively dispositioned. Use the verified Implementation note results plus review dispositions from all rounds to write Completion evidence:
+For a separate review after lightweight implementation, write Completion evidence in the ticket using the verified Implementation note results plus review dispositions from all rounds:

 - delivered behavior;
 - commands run and observed results (from Implementation note);
```

## 已知流程债

- `TestIssueTrackerDescribesRepairWorkflow` 读取已迁移的 `docs/agents/issue-tracker.md`；01 明确 Change 仅改为 `ff-wiki`，保留 `repair` 断言。
- 继承 journal 中 Beads 无活动数据库、旧 flowforge-progress 与 INDEX 缺失的事实；本轮不初始化，不编造旧流程执行记录。
- `command -v golangci-lint` exit=1，无输出：当前 PATH 未提供该工具；需要运行 lint 时原样报告，不冒充通过。

## 规划验证

- 权威初检：`flowforge check --dir ff-wiki/proposals/cost-aware-verification-planning --strict` exit=0，0 issues，无诊断（发布票前）。
- 现有源/部署结构接缝：`go test ./internal/command -run 'Test(PackagedSkillPointersResolve|PlanPublishesExecutionDetailSkeleton|PlanSkillDeclaresModeLine|ArtifactContractDocumentsExecutionContractSections|AgentRulesDescribeGenericCapabilityDispatch|IssueTrackerDescribesRepairWorkflow)$' -count=1 -v` exit=1；前五个测试通过，仅旧路径用例失败。

```text
=== RUN   TestPackagedSkillPointersResolve
--- PASS: TestPackagedSkillPointersResolve (0.05s)
=== RUN   TestAgentRulesDescribeGenericCapabilityDispatch
--- PASS: TestAgentRulesDescribeGenericCapabilityDispatch (0.00s)
=== RUN   TestPlanPublishesExecutionDetailSkeleton
--- PASS: TestPlanPublishesExecutionDetailSkeleton (0.00s)
=== RUN   TestArtifactContractDocumentsExecutionContractSections
--- PASS: TestArtifactContractDocumentsExecutionContractSections (0.00s)
=== RUN   TestPlanSkillDeclaresModeLine
--- PASS: TestPlanSkillDeclaresModeLine (0.00s)
=== RUN   TestIssueTrackerDescribesRepairWorkflow
    assets_deploy_test.go:610: open ../../docs/agents/issue-tracker.md: no such file or directory
--- FAIL: TestIssueTrackerDescribesRepairWorkflow (0.00s)
FAIL
FAIL	flowforge/internal/command	0.472s
FAIL
```

- Tracker 既有保护：`go test ./internal/tracker -run 'Test(ExecutionContractCompletenessAppliesOnlyToManagedExecutableTickets|EvidenceQuadrupleParsing|NeedsRepairIsNonExecutableAndRepairProvenanceIsParsed)$' -count=1 -v` exit=0，3个顶层测试及其子例通过；不代表新拟测试已运行。
- 要求的完整内包回归：`GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` exit=1；唯一失败测试仍为 `TestIssueTrackerDescribesRepairWorkflow`，其他内包通过或无测试。update 两项依赖 `/proc/self/exe` 的用例按既有条件在本机跳过，不能称这些分支已实测。原始本地日志 `/tmp/cost-aware-plan-internal-tests.log`，关键输出如下。

```text
--- SKIP: TestE2E_SelfupdateApply_ReplacesBinaryCorrectly (0.00s)
=== RUN   TestUpgradeToVersion_AlreadyUpToDate
--- PASS: TestUpgradeToVersion_AlreadyUpToDate (0.00s)
=== RUN   TestUpgradeToVersion_RejectsDowngradeDistinctly
--- PASS: TestUpgradeToVersion_RejectsDowngradeDistinctly (0.00s)
=== RUN   TestFullUpgradePipeline
    upgrade_test.go:347: cannot read /proc/self/exe: open /proc/self/exe: no such file or directory
--- SKIP: TestFullUpgradePipeline (0.00s)
PASS
ok  	flowforge/internal/update	0.442s
?   	flowforge/internal/version	[no test files]
FAIL
385:    assets_deploy_test.go:610: open ../../docs/agents/issue-tracker.md: no such file or directory
386:--- FAIL: TestIssueTrackerDescribesRepairWorkflow (0.00s)
661:FAIL
662:FAIL	flowforge/internal/command	9.313s
1016:    upgrade_test.go:219: cannot read /proc/self/exe: open /proc/self/exe: no such file or directory
1023:    upgrade_test.go:347: cannot read /proc/self/exe: open /proc/self/exe: no such file or directory
1028:FAIL
```

- 新测试未落地：正常 `write` 被测试保护拒绝；未运行新测试、不记录预计 red 为实际 red。已保存 [未应用补丁](test-authoring.patch) 与 [作者交接](test-author-handoff.md)，`git apply --check ff-wiki/proposals/cost-aware-verification-planning/evidence/test-authoring.patch` exit=0，仅证明可适用。
- 未运行 `make dev`、新包部署或受控模型情景；本轮没有修改产品技能、嵌入、测试、保护或模型配置。
