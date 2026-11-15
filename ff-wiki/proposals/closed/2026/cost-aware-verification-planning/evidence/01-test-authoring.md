---
flowforge:
  schema: 1
  role: evidence
  id: cost-aware-verification-planning-01-test-authoring
  revision: 1
---

# 01 G0 测试作者执行记录

本记录覆盖 design rev 3 / Change 8 的合法测试作者执行与 Refine 前置事实，不是 Changes 1–7、9–10 的实施证据，也不是模型受控遵从证明。执行起点为 `fc2dd9a7ae13e44709e678b22290bdef4599d0ac`；开始时 `git status --short` 无输出。

## 实例与授权

- sessionId：`01a1abb0-1838-7290-90a1-09521081d41e`
- instanceId：`033414d8-3cb6-470b-96d2-a9987adcc348`
- projectRoot：`/Users/qiangbi/develop/projects/Syl/tangram/flowforge`
- grantPath：`/Users/qiangbi/develop/projects/Syl/tangram/flowforge/.flowforge/test-author-grants/01a1abb0-1838-7290-90a1-09521081d41e.json`

首次正常 `flowforge_test_authority` 返回 `eligible=true`，`grant.status=missing`，`role=""`、`paths=[]`、`expiresAt=""`，reason 为 `no grant record for this session`。随后通过 `contact_supervisor(reason=need_decision)` 发送完整状态和精确三路径，请求主会话只为此实例签发 2h `role=plan` 授权；本实例未自行写 grant。

主会话签发后再次正常读取：`grant.status=valid`、`role=plan`、`expiresAt=2026-11-05T14:31:07.010Z`，paths 精确为：

1. `internal/command/cost_aware_assets_test.go`
2. `internal/tracker/cost_aware_evidence_test.go`
3. `internal/command/assets_deploy_test.go`

G0 后再次 `contact_supervisor` 请求撤销。主会话回复已删除此精确 grant 文件；随后正常 `flowforge_test_authority` 再读为 `eligible=true`、`grant.status=missing`、空 role/paths/expiresAt，reason 为 `no grant record for this session`。授权已撤销，本实例未自行删除 grant，撤销后未再修改任何测试。

## 精确应用与内容对照

只使用正常 `write` 创建两个新测试、正常 `edit` 改变旧测试的 `docs` 路径段为 `ff-wiki`。没有委派、关闭 guard、使用 bash 写测试或 push，没有增加或改变验收断言。`TestIssueTrackerDescribesRepairWorkflow` 的 `strings.Contains(string(data), "repair")` 和 `t.Fatal("issue-tracker.md must describe the repair workflow")` 原样保留。

原补丁保持未改，SHA-256 为 `511d9ec3b773408634864dc30a51e2dc23d4b6d0341310d9baf38ebdf616169f`。落地后运行 `git apply --check --reverse ff-wiki/proposals/cost-aware-verification-planning/evidence/test-authoring.patch`，exit 0，证明当前三路径相对原基线精确匹配补丁，没有尚未应用的补丁内容。`git diff --check` exit 0。

| 路径 | 落地 SHA-256 | Git blob hash |
| --- | --- | --- |
| `internal/command/cost_aware_assets_test.go` | `7c3fcd13322cb89213898690d16bd6f831be4424f1c1d87f21f25076f2d52f43` | `c5b4f2bd04d3abbd1958acd80cceeebeaade57bd` |
| `internal/tracker/cost_aware_evidence_test.go` | `101e4219d06dadfd243ec1a005e436b86135cbcec02d00a5a0a4a6f8f6da68de` | `29dcdb22f6e4b7e99b031a921cf43d23a90f0b02` |
| `internal/command/assets_deploy_test.go` | `503b574e22fec31ef10ec908c3fb3c2c6c1be6a19d04ea2591c85b0b5ccd2511` | `90c183ff77e8943a13d7e6d1d8d1a71f359478d6` |

`gofmt -d internal/command/cost_aware_assets_test.go internal/tracker/cost_aware_evidence_test.go internal/command/assets_deploy_test.go` exit 0 且无输出：原补丁已经是标准格式，无需格式化写入，落地内容与补丁无 gofmt 差异。

## 实际测试

以下命令均在本实例真实运行，exit 均为 0。完整日志由工具会话保存；本文件是各命令的持久结果工件。

| 命令 | 实际结果 |
| --- | --- |
| `go test ./internal/command -run '^TestCostAwareRulePackageSourceEmbeddedDeployed$' -count=1 -v` | PASS；顶层用例与七个源→嵌入→隔离部署子项全部通过，`ok flowforge/internal/command 0.571s` |
| `go test ./internal/tracker -run '^TestCostAwareSharedRunPreservesEvidenceGate$' -count=1 -v` | PASS；schedule=false/true × fixMissing=false/true 四组合全部通过，`ok flowforge/internal/tracker 0.532s` |
| `go test ./internal/command -run '^TestIssueTrackerDescribesRepairWorkflow$' -count=1 -v` | PASS；修正后的 `ff-wiki` 路径通过，`ok flowforge/internal/command 0.498s` |
| `go test ./internal/command -run 'Test(CostAwareRulePackageSourceEmbeddedDeployed\|PackagedSkillPointersResolve\|PlanPublishesExecutionDetailSkeleton\|PlanSkillDeclaresModeLine\|ArtifactContractDocumentsExecutionContractSections\|AgentRulesDescribeGenericCapabilityDispatch\|ImplementSkillCarriesWeakExecutorContract\|ReviewSkillCarriesRound0Audit\|ReviewSkillDescribesRepairEscalation\|IssueTrackerDescribesRepairWorkflow)$' -count=1` | PASS；`ok flowforge/internal/command 0.692s` |
| `go test ./internal/tracker -run 'Test(CostAwareSharedRunPreservesEvidenceGate\|ExecutionContractCompletenessAppliesOnlyToManagedExecutableTickets\|EvidenceQuadrupleParsing\|NeedsRepairIsNonExecutableAndRepairProvenanceIsParsed)$' -count=1` | PASS；`ok flowforge/internal/tracker 0.554s` |
| `GOPROXY=https://goproxy.cn,direct go test -v ./internal/...` | exit 0；所有执行用例通过、0 failures。三项既有环境性 SKIP：`TestRepositoryAndTangramProposalLayouts` 缺少兼容 checkout；`TestE2E_SelfupdateApply_ReplacesBinaryCorrectly` 与 `TestFullUpgradePipeline` 因 Darwin 无 `/proc/self/exe` 跳过。不能把这些 SKIP 称为对应行为已实测。 |

表中 regex 的反斜杠仅为 Markdown 表格转义；实际 shell 命令使用票据 Done and verify 的原始 `|`，未改变 regex。两个新增兼容测试在落地后的首次运行即 green，没有落地前的新测试 red 记录，也未故意破坏实现制造 red。旧路径维护测试在规划基线中已有真实 red：`open ../../docs/agents/issue-tracker.md: no such file or directory`；本轮只修路径后 green，不能把这项旧 red 归因于两个新增测试。

00/03 已消除资产嵌入漂移，所以资产交付兼容回归 baseline green 是允许且实际观察到的结果。这些测试只保护现有源/嵌入/部署一致性和 Catalog 证据格式/门禁兼容；它们不证明尚未编写的 Changes 1–7 源规则已语义合规，也不替代 02 的真实模型行为验收。

## lint 与工具边界

在 lint 前以 `git add` 将三个精确 Go 路径纳入 Git 跟踪/暂存，避免未跟踪文件被 rev 门禁遗漏。实际运行：

```text
/tmp/flowforge-verification-tools/golangci-lint-2.15.0/golangci-lint run --timeout 2m --new-from-rev 16cce8e --whole-files --max-issues-per-linter 0 --max-same-issues 0 ./cmd/... ./internal/...
```

exit 0，输出 `0 issues.`。这是固定工具 2.15.0 对指定 rev、cmd/internal 修改文件的整文件门禁结果，不是全仓 lint 通过证明；[既有全仓与默认产品包失败基线](01-lint-baseline.md)保持有效。未使用 `--disable`、`--fix`、豁免或改变 issues exit code。

本轮没有运行 `make dev`，没有修改源规则、嵌入资产、产品实现、DAG/Catalog/evidence 算法或其他测试。Changes 1–7、9–10 未执行。

## Refine 图状态

消费Blocked evidence并更新五节后，正常 `flowforge_check(feature=cost-aware-verification-planning)` 返回依赖图健康、无循环或悬空引用；本票的 `blocked-evidence-present` 已消失，本票不再有 `execution-contract-incomplete`。仍报告的唯一execution-contract gap属于02，其他全局warnings保留，不据此声称全仓诊断清零。

正常 `flowforge_frontier(feature=cost-aware-verification-planning)` 将01列入 `ready`（`ready_with_warnings=null`），02仍等待01。这里的ready只表示01可进入后续轻量实施，不表示01已完成或模型行为已验收。
